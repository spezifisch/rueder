package worker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/mmcdole/gofeed"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/spezifisch/rueder3/backend/pkg/worker/scheduler"
	"github.com/sym01/htmlsanitizer"
)

// FeedWorkerPool implements a scheduler.WorkerPool
type FeedWorkerPool struct {
	config     FeedWorkerConfig
	repository scheduler.Repository
}

// FeedWorkerConfig configures fetching parameters
type FeedWorkerConfig struct {
	HTTPTimeout       time.Duration
	UserAgent         string
	MinimumFetchDelay time.Duration
	MaximumFetchDelay time.Duration
	FetchJitterS      int
}

// DefaultFeedWorkerConfig has usable default values
var DefaultFeedWorkerConfig = FeedWorkerConfig{
	HTTPTimeout: 60 * time.Second,
	// https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
	UserAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
	MinimumFetchDelay: 15 * time.Minute,
	MaximumFetchDelay: 12 * time.Hour,
	FetchJitterS:      30,
}

// NewFeedWorkerPool creates a worker pool with the given Repo backend
func NewFeedWorkerPool(repository scheduler.Repository) *FeedWorkerPool {
	return &FeedWorkerPool{
		config:     DefaultFeedWorkerConfig,
		repository: repository,
	}
}

// StartWorker is launches as a goroutine that fetches feeds
func (p FeedWorkerPool) StartWorker(id int, feeds <-chan scheduler.Feed, doneFeeds chan<- scheduler.Feed) {
	workerLog := log.WithField("worker", id)
	workerLog.Info("FeedWorker running")

	for feed := range feeds {
		workerLog.WithField("id", feed.ID).Info("fetching feed")

		// fetch feed
		if err := p.fetchFeed(&feed); err != nil {
			workerLog.WithError(err).WithField("id", feed.ID).Error("failed fetching feed")
		}

		workerLog.WithField("id", feed.ID).Info("done fetching feed")

		// report back that it's done
		doneFeeds <- feed
	}
}

func (p FeedWorkerPool) fetchFeedURL(url string) (feed *gofeed.Feed, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.HTTPTimeout)
	defer cancel()

	fp := gofeed.NewParser()
	fp.UserAgent = p.config.UserAgent

	feed, err = fp.ParseURLWithContext(url, ctx)
	return
}

func (p FeedWorkerPool) fetchFeedTryingHTTPS(f *scheduler.Feed) (feed *gofeed.Feed, err error) {
	// if it's an HTTP URL try using HTTPS first
	if helpers.IsHTTPURL(f.FeedURL) {
		feedLog := log.WithField("feed_id", f.ID)

		httpsURL := helpers.RewriteToHTTPS(f.FeedURL)
		feed, err = p.fetchFeedURL(httpsURL)
		if err != nil {
			feedLog.Info("feed was not reachable via https")
		} else {
			if len(feed.Items) > 0 {
				feedLog.Info("feed was reachable via https, updating feed info")
				// the change gets saved at the end of fetchFeed
				f.FeedURL = httpsURL
				return
			}

			feedLog.Info("feed was reachable via https but no articles found, falling back to http")
			feed = nil
		}
	}

	feed, err = p.fetchFeedURL(f.FeedURL)
	return
}

func (p FeedWorkerPool) fetchFeed(f *scheduler.Feed) (err error) {
	fetchedAt := time.Now()
	f.FetcherState.FetchedAt = fetchedAt // ensure it's updated to avoid endless immediate requeueing of the job

	if !helpers.IsURL(f.FeedURL) {
		f.FetcherState.Working = false
		f.FetcherState.LastError = time.Now().Round(time.Second)
		f.FetcherState.Message = "Invalid URL"

		if e := p.repository.UpdateFeedInfo(f.ID, f); e != nil {
			log.WithError(e).Error("couldn't update feed info after url error")
		}

		err = errors.New("feed has invalid url")
		return
	}

	// fetch feed
	parsedFeed, err := p.fetchFeedTryingHTTPS(f)
	if err != nil {
		f.FetcherState.Working = false
		f.FetcherState.LastError = time.Now().Round(time.Second)
		f.FetcherState.Message = fmt.Sprintf("Fetcher Error: %s", err)

		if e := p.repository.UpdateFeedInfo(f.ID, f); e != nil {
			log.WithError(e).Error("couldn't update feed info after fetcher error")
		}

		return // with original err
	}

	// parse articles
	p.processArticles(f, parsedFeed)

	// write updated feed info to repository
	p.updateFeedFields(f, parsedFeed)

	err = p.repository.UpdateFeedInfo(f.ID, f)
	return
}

func (p FeedWorkerPool) updateFeedFields(storedFeed *scheduler.Feed, parsedFeed *gofeed.Feed) {
	storedFeed.FetcherState.Working = true
	storedFeed.FetcherState.LastSuccess = time.Now().Round(time.Second)

	if parsedFeed.Title != "" {
		storedFeed.Title = parsedFeed.Title
	}
	// TODO description maybe
	if parsedFeed.Link != "" {
		storedFeed.SiteURL = parsedFeed.Link
	}
	if parsedFeed.Image != nil && parsedFeed.Image.URL != "" {
		storedFeed.Icon = parsedFeed.Image.URL
	} else if parsedFeed.Items != nil && len(parsedFeed.Items) > 0 && parsedFeed.Items[0].Image != nil && parsedFeed.Items[0].Image.URL != "" {
		// take feed icon from newest article
		storedFeed.Icon = parsedFeed.Items[0].Image.URL
	}
	// also old.FeedURL might have been updated in fetchFeedTryingHTTPS.
	// we don't need to do more in this case.

	// update fetch delay
	p.updateFetchDelay(storedFeed)
}

func (p FeedWorkerPool) updateFetchDelay(feed *scheduler.Feed) {
	// get up to date article count
	f, err := p.repository.GetFeed(feed.ID)
	if err != nil {
		log.WithField("feed_id", feed.ID).WithError(err).Error("failed getting feed for article count")
		return
	}

	articleCount := f.ArticleCount
	durationSubscribed := time.Since(f.CreatedAt)
	fetchDelay := 4 * p.config.MinimumFetchDelay

	if durationSubscribed < 24*time.Hour {
		return
	} else if articleCount == 0 {
		fetchDelay = p.config.MaximumFetchDelay
	} else {
		articlesPerDay := int(math.Round(float64(articleCount) / math.Round(durationSubscribed.Hours()/24.0)))
		if articlesPerDay >= 15 {
			fetchDelay = p.config.MinimumFetchDelay
		} else if articlesPerDay >= 8 {
			fetchDelay = 2 * p.config.MinimumFetchDelay
		} else if articlesPerDay >= 5 {
			fetchDelay = 3 * p.config.MinimumFetchDelay
		}
	}

	// add random component to the fetches drift apart so we don't always fetch them all at the same time
	fetchDelay += time.Duration(rand.Intn(p.config.FetchJitterS)) * time.Second
	log.WithFields(log.Fields{"articleCount": articleCount, "durationSubscribed": durationSubscribed, "fetchDelay": fetchDelay}).Info("calculated new fetchDelay")

	feed.FetcherState.FetchDelayS = int(math.Round(fetchDelay.Seconds()))
}

// custom sort that orders articles from oldest to newest
type feedCustomSort struct {
	feed *gofeed.Feed
	sort.Interface

	hasArticleDates bool
}

func (s *feedCustomSort) Less(i, j int) bool {
	if s.feed.Items[i].PublishedParsed != nil && s.feed.Items[j].PublishedParsed != nil {
		s.hasArticleDates = true
		return s.Interface.Less(i, j)
	}
	return false
}

func reverseFeedItems(feed *gofeed.Feed) {
	for i := 0; i < len(feed.Items)/2; i++ {
		j := len(feed.Items) - i - 1
		feed.Swap(i, j)
	}
}

// these regexes are for removing tags whose content would otherwise weird leftovers
var regexpStyle = regexp.MustCompile(`<style[\S\s]+?<\/style>*`)
var regexpScript = regexp.MustCompile(`<script[\S\s]+?<\/script>`)

func (p FeedWorkerPool) processArticles(f *scheduler.Feed, feed *gofeed.Feed) {
	if feed == nil || feed.Items == nil || len(feed.Items) == 0 {
		log.WithField("id", f.ID).Info("no articles")
		return
	}

	now := time.Now().Round(time.Second)
	articleCount := len(feed.Items)
	newArticleCount := 0
	brokenArticleCount := 0
	failedArticleCount := 0

	// sort feed items so that the oldest article is first.
	// (the oldest article is then added to the db first and gets the lowest sequence number.)
	// if the items don't have timestamps we want to keep the order in this step,
	// therefore we need a stable sort.
	fcs := &feedCustomSort{
		feed:            feed,
		Interface:       feed,
		hasArticleDates: false,
	}
	sort.Stable(fcs)
	if !fcs.hasArticleDates {
		// assume the articles are ordered from new to old, so we need to reverse the order.
		reverseFeedItems(feed)
	}

	// first check which articles are new
	guids := make([]string, len(feed.Items))
	for i, item := range feed.Items {
		if item.GUID != "" {
			guids[i] = item.GUID
		} else if item.Link != "" {
			// use link instead for this borked feed
			guids[i] = item.Link
		} else {
			// this is a very broken feed, ignore the article as we have no way to identify it
			guids[i] = ""
			brokenArticleCount++
		}
	}
	exists, err := p.repository.CheckExistingArticles(f.ID, guids)
	if err != nil {
		log.WithError(err).Error("failed checking existing articles")
		return
	}

	// parse articles
	for i, item := range feed.Items {
		if exists[i] || guids[i] == "" {
			// ignore article
			continue
		}
		newArticleCount++

		article := scheduler.Article{
			SiteGUID:  guids[i],
			Link:      item.Link,
			Tags:      item.Categories,
			RawTitle:  item.Title,
			RawTeaser: item.Description,
			RawText:   item.Content,
		}

		// parse values that can fail
		if item.PublishedParsed != nil {
			article.Time = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			article.Time = *item.UpdatedParsed
		} else {
			article.Time = now
		}
		if item.Authors != nil && len(item.Authors) > 0 {
			article.Authors = make([]string, len(item.Authors))
			for i, author := range item.Authors {
				article.Authors[i] = author.Name
			}
		}
		if item.Image != nil {
			article.Image = item.Image.URL
			article.ImageTitle = item.Image.Title
		}
		if item.Enclosures != nil && len(item.Enclosures) > 0 {
			article.Enclosures = make([]scheduler.ArticleEnclosure, len(item.Enclosures))
			for i, enclosure := range item.Enclosures {
				article.Enclosures[i].Length = enclosure.Length
				article.Enclosures[i].Type = enclosure.Type
				article.Enclosures[i].URL = enclosure.URL
			}
		}

		// set content to teaser if content is empty
		if article.RawText == "" {
			article.RawText = article.RawTeaser
		}

		articleLog := log.WithFields(log.Fields{
			"feed": f.ID,
			"guid": article.SiteGUID})

		// strip all html from title and teaser
		{
			s := htmlsanitizer.NewHTMLSanitizer()
			s.AllowList = nil

			if title, err := s.SanitizeString(article.RawTitle); err != nil {
				articleLog.WithError(err).Warn("failed sanitizing title html")
			} else {
				article.Title = strings.TrimSpace(title)
			}
			if teaser, err := s.SanitizeString(article.RawTeaser); err != nil {
				articleLog.WithError(err).Warn("failed sanitizing teaser html")
			} else {
				article.Teaser = strings.TrimSpace(teaser)
			}
		}
		// sanitize content html
		{
			// remove style/script tags first because the sanitizer strips only the tags but leaves their content
			rawText := regexpStyle.ReplaceAllString(article.RawText, "")
			rawText = regexpScript.ReplaceAllString(rawText, "")

			s := htmlsanitizer.NewHTMLSanitizer()
			// use defaultAllowList (removing iframes and such) but forbid id and class.
			s.GlobalAttr = []string{}
			// we don't do additional work here like removing tracking pixels and links,
			// that's job of the frontend (also it might be different depending on user settings.)

			if content, err := s.SanitizeString(rawText); err != nil {
				articleLog.WithError(err).Warn("failed sanitizing content html")
			} else {
				article.Text = strings.TrimSpace(content)
			}
		}

		if err := p.repository.AddArticle(f.ID, &article); err != nil {
			failedArticleCount++
			articleLog.WithError(err).Error("failed adding article")
			continue
		}

		articleLog.WithField("title", article.Title).Debug("got article")
	}

	log.WithFields(log.Fields{
		"feed":   f.ID,
		"new":    newArticleCount,
		"failed": failedArticleCount,
		"broken": brokenArticleCount}).
		Infof("got %d articles", articleCount)
}
