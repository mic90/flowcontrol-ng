package graph

import (
	"context"
	"github.com/desertbit/timer"
	"time"
)

type ExitReason uint8

const (
	ExitReasonNone ExitReason = iota
	ExitReasonTimeout
	ExitReasonContext
)

type Timer interface {
	WaitFor(duration time.Duration)
	ScheduleAt(cronSchedule string)
	ExitReason() ExitReason
}

type ProcessTimer struct {
	timer      timer.Timer
	exitReason ExitReason
	ctx        context.Context
}

func NewProcessTimer(ctx context.Context) *ProcessTimer {
	return &ProcessTimer{*timer.NewStoppedTimer(), ExitReasonNone, ctx}
}

func (pt *ProcessTimer) WaitFor(duration time.Duration) {
	pt.exitReason = ExitReasonNone
	pt.timer.Reset(duration)
	select {
	case <-pt.ctx.Done():
		pt.timer.Stop()
		pt.exitReason = ExitReasonContext
	case <-pt.timer.C:
		pt.exitReason = ExitReasonTimeout
	}
}

//TODO: Add functionality to trigger actions at certain time intervals based on cron-like string input
func (pt *ProcessTimer) ScheduleAt(cronSchedule string) {
	panic("ScheduleAt is not implemented")
}

func (pt *ProcessTimer) ExitReason() ExitReason {
	return pt.exitReason
}
