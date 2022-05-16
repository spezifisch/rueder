package scheduler

import "time"

// DeadlineTimer is a deadline-aware timer
type DeadlineTimer struct {
	C        *<-chan time.Time
	Deadline time.Time

	timer    *time.Timer
	duration *time.Duration
}

// NewDeadlineTimer starts a timer triggers at the given deadline
func NewDeadlineTimer(deadline time.Time) *DeadlineTimer {
	duration := time.Until(deadline)
	d := &DeadlineTimer{
		Deadline: deadline,
		timer:    time.NewTimer(duration),
		duration: &duration,
	}
	d.C = &d.timer.C
	return d
}

// Stop the timer
func (d *DeadlineTimer) Stop() {
	d.timer.Stop()
}

// InFuture returns true if the deadline was in the future at the time of creating the timer
func (d *DeadlineTimer) InFuture() bool {
	return *d.duration > 0
}
