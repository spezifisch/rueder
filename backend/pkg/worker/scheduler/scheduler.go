package scheduler

import (
	"time"

	"github.com/apex/log"
	"github.com/gofrs/uuid"
)

// Scheduler dispatches to the workers which feeds should be fetched
type Scheduler struct {
	repository Repository
	queue      *FeedQueue

	feedRehash          chan bool
	feedRehashRequested bool
	feedAdded           chan uuid.UUID

	workerPool      WorkerPool
	workerCount     int
	workerFeeds     chan Feed
	workerDoneFeeds chan Feed
	jobsInProgress  int

	minimumFetchDelay time.Duration
	retryDelay        time.Duration
}

// NewScheduler creates a new scheduler with the given count of workers
func NewScheduler(repository Repository, workerPool WorkerPool, workerCount int) *Scheduler {
	return &Scheduler{
		repository: repository,
		queue:      nil,

		feedRehash:          make(chan bool),                     // a feed was deleted or changed
		feedRehashRequested: false,                               // true if the signal came and until all active jobs are completed
		feedAdded:           make(chan uuid.UUID, workerCount+1), // a new feed was added externally

		workerPool:      workerPool,
		workerCount:     workerCount,
		workerFeeds:     make(chan Feed),                // block when no worker is ready
		workerDoneFeeds: make(chan Feed, workerCount+1), // we need 1 more queue element than workers because the scheduler can queue an additional feed while all workers are busy
		jobsInProgress:  0,                              // counts how many workers are currently processing a feed

		minimumFetchDelay: 10 * time.Minute, // don't fetch feeds faster than this
		retryDelay:        30 * time.Second,
	}
}

// Run starts the workers and the scheduling loop
func (s *Scheduler) Run() {
	// start workers
	log.Infof("Starting %d workers ...", s.workerCount)
	for w := 1; w <= s.workerCount; w++ {
		go s.workerPool.StartWorker(w, s.workerFeeds, s.workerDoneFeeds)
	}

	// start listening for newly added feeds
	err := s.repository.RunFeedChangeListener(s.feedAdded, s.feedRehash)
	if err != nil {
		return
	}

	// dispatch jobs
	log.Info("Starting job dispatcher loop")
	for {
		// init queue if it's not yet initialized
		if s.queue == nil {
			log.Info("init queue")
			if err := s.initQueue(); err != nil {
				log.WithError(err).WithFields(log.Fields{"wait": s.retryDelay}).Error("failed initializing FeedQueue")
				time.Sleep(s.retryDelay)
				continue
			}
		}

		// wait until a feed needs to be fetched.
		// this also processes the worker results queue and schedules rehashes
		s.sleepUntilNextJob()

		// this can be a:
		// timeout => start sleeping again in next iteration
		// rehash requested => wait until all jobs are done and then re-init queue in next iteration
		if s.queue == nil || s.queue.Len() == 0 {
			continue
		}

		job := s.queue.Pop()
		log.WithField("pop", job.feed).Info("scheduler sends")

		// send it to a worker to fetch it; blocks until a worker takes it.
		// a rehash can get requested while we're waiting here. this is handled in the next loop iteration in sleepUntilNextJob()
		s.workerFeeds <- job.feed
		s.jobsInProgress++

		log.WithField("pop", job.feed).WithField("inProgress", s.jobsInProgress).Info("scheduler sent")
	}
}

func (s *Scheduler) initQueue() (err error) {
	// get all feeds
	feeds, err := s.repository.Feeds()
	if err != nil {
		return
	}

	// add all existing feeds to queue
	s.queue = NewFeedQueue()
	for _, feed := range feeds {
		s.queueFeed(&feed)
	}
	return
}

func (s *Scheduler) rehashIfRequestedAndPossible() {
	if !s.feedRehashRequested {
		return
	}
	if s.jobsInProgress > 0 {
		log.WithField("jobsInProgress", s.jobsInProgress).Info("can't rehash yet. jobs active")
		return
	}
	s.queue = nil // this causes Run() to reinit the queue in the next loop
	s.feedRehashRequested = false
}

func (s *Scheduler) sleepUntilNextJob() {
	timer := s.refreshSleepTimer(nil)

	for {
		select {
		case <-s.feedRehash:
			// this signals that we need to refetch all feed ids.
			// but we need to wait until all currently active workers are done to avoid race conditions.
			log.Info("got rehash signal, queueing rehash")
			s.feedRehashRequested = true
		case newFeedID := <-s.feedAdded:
			logFeed := log.WithField("feed", newFeedID)

			// put new feed into the queue
			addedFeed, err := s.repository.GetFeed(newFeedID)
			if err != nil {
				logFeed.Error("failed adding new feed")
				return
			}

			if !s.feedRehashRequested {
				logFeed.Info("adding new feed")
				s.queueFeed(&addedFeed)
			} else {
				logFeed.WithField("jobsInProgress", s.jobsInProgress).Info("not adding new feed because rehash requested")
			}
		case doneFeed := <-s.workerDoneFeeds:
			logFeed := log.WithField("feed", doneFeed.ID)

			s.jobsInProgress--
			if !s.feedRehashRequested {
				// put fetched feeds back into the job queue with their updated deadline
				logFeed.Info("readding feed")
				s.queueFeed(&doneFeed)
			} else {
				logFeed.WithField("jobsInProgress", s.jobsInProgress).Info("not readding feed because rehash requested")
			}
		case <-*timer.C:
			// we've got a job!
			s.rehashIfRequestedAndPossible()
			return
		}

		// safe rehash point reached?
		s.rehashIfRequestedAndPossible()
		if s.queue == nil {
			return
		}

		// reinit timer, our deadline might have changed
		timer = s.refreshSleepTimer(timer)
	}
}

// queueFeed adds the given feed to the queue, to be fetched once at its deadline
func (s *Scheduler) queueFeed(feed *Feed) {
	delay := s.getFetchDelay(feed)
	deadline := feed.FetcherState.FetchedAt.Add(delay) // when to fetch
	log.WithField("feed", feed.ID).WithField("deadline", deadline).Debug("queued feed")

	s.queue.Push(&FeedQueueItem{
		feed:     *feed,
		deadline: deadline,
	})
}

// refreshSleepTimer returns a timer that triggers when the next feed is due
func (s *Scheduler) refreshSleepTimer(oldTimer *DeadlineTimer) *DeadlineTimer {
	if s.queue == nil {
		return oldTimer
	}

	if s.queue.Len() == 0 {
		// the queue is empty
		if oldTimer != nil {
			// let the old timer keep running, we don't want the overhead of creating a new one
			return oldTimer
		}
		// we choose a big value that doesn't wake us up too often
		return NewDeadlineTimer(time.Now().Add(24 * time.Hour))
	}

	// wake up for the next job
	job := s.queue.Peek()

	if oldTimer != nil {
		// see if we can reuse the old timer
		if oldTimer.Deadline == job.deadline {
			return oldTimer
		}

		// we can't, stop it
		oldTimer.Stop()
		oldTimer = nil
	}

	timer := NewDeadlineTimer(job.deadline)
	if timer.InFuture() {
		log.WithField("until", job.deadline).Info("sleeping")
	}
	return timer
}

func (s Scheduler) getFetchDelay(f *Feed) time.Duration {
	delay := time.Duration(f.FetcherState.FetchDelayS) * time.Second
	if delay < s.minimumFetchDelay {
		log.Warnf("feed with lower than minimum fetch delay: %v", f.ID)
		delay = s.minimumFetchDelay
	}
	return delay
}
