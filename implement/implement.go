// Package implement provides simple wrapper of anping.AnPinger
package implement

import (
	"context"
	"io"
	"time"

	F "github.com/sagernet/sing/common/format"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/statistics"
)

const Protocol = "Unknown"

var _ anping.AnPinger = (*AnPingerWrapper)(nil)

// AnPingerWrapper is a simple wrapper of anping.AnPinger. It does nothing when running.
type AnPingerWrapper struct {
	Opt *anping.Options
	Sta *statistics.Statistics

	logger statistics.Logger
}

func New(logWriter io.Writer) anping.AnPinger {
	return &AnPingerWrapper{
		Opt:    anping.NewOptions(),
		Sta:    statistics.NewStatistics(),
		logger: &statistics.DefaultLogger{Writer: logWriter},
	}
}

func (a *AnPingerWrapper) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go a.start(ctx, done)
	return done
}

func (a *AnPingerWrapper) start(ctx context.Context, done chan struct{}) {
	defer TryCloseDone(done)
	a.OnStart()
	defer a.OnFinish()

	timer := time.NewTimer(a.Opt.Interval)
	defer timer.Stop()
	for i := a.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-timer.C:
			timer.Reset(a.Opt.Interval)
		}
	}
}

func (a *AnPingerWrapper) Protocol() string {
	return Protocol
}

func (a *AnPingerWrapper) Statistics() statistics.StatisticsGetter {
	// Make user can't edit it.
	return a.Sta.Getter()
}

func (a *AnPingerWrapper) ResetStatistics() {
	a.Sta = statistics.NewStatistics()
}

func (a *AnPingerWrapper) SetAddress(address M.Socksaddr) error {
	return a.Opt.SetAddress(address)
}

func (a *AnPingerWrapper) Address() M.Socksaddr {
	return a.Opt.Address()
}

func (a *AnPingerWrapper) SetLogger(logger statistics.Logger) {
	a.logger = logger
}

func (a *AnPingerWrapper) Options() *anping.Options {
	return a.Opt
}

func (a *AnPingerWrapper) OnStart() {
	if a.logger != nil {
		a.logger.OnStart(a.Opt.Address(), a.Sta)
	}
}

func (a *AnPingerWrapper) OnRecv(t time.Duration) {
	if a.logger != nil {
		a.logger.OnRecv(a.Opt.Address(), a.Sta, t)
	}
}

func (a *AnPingerWrapper) OnLost(errMsg ...any) {
	if a.logger != nil {
		a.logger.OnLost(a.Opt.Address(), a.Sta, F.ToString(errMsg...), -1)
	}
}

func (a *AnPingerWrapper) OnFinish() {
	if a.logger != nil {
		a.logger.OnFinish(a.Opt.Address(), a.Sta)
	}
}

func TryCloseDone(dones ...chan struct{}) (closed int) {
	for _, done := range dones {
		select {
		case <-done:
			continue
		default:
			close(done)
			closed++
		}
	}

	return
}
