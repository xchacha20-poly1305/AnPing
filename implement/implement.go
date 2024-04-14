// Package implement provides simple wrapper of anping.AnPinger
package implement

import (
	"context"
	"io"
	"time"

	F "github.com/sagernet/sing/common/format"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/state"
)

const Protocol = "Unknown"

var _ anping.AnPinger = (*AnPingerWrapper)(nil)

// AnPingerWrapper is a simple wrapper of anping.AnPinger. It do nothing when running.
type AnPingerWrapper struct {
	Opt *anping.Options
	*state.State

	logger state.Logger
}

func New(logWriter io.Writer) anping.AnPinger {
	return &AnPingerWrapper{
		Opt:    anping.NewOptions(),
		State:  state.NewState(),
		logger: &state.DefaultLogger{Writer: logWriter},
	}
}

func (a *AnPingerWrapper) Run() {
	a.RunContext(context.Background())
}

func (a *AnPingerWrapper) RunContext(ctx context.Context) {
	a.OnStart()

	defer a.OnFinish()

	for i := a.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(a.Opt.Interval)
	}
}

func (a *AnPingerWrapper) Protocol() string {
	return Protocol
}

func (a *AnPingerWrapper) SetAddress(address string) error {
	return a.Opt.SetAddress(address)
}

func (a *AnPingerWrapper) SetLogger(logger state.Logger) {
	a.logger = logger
}

func (a *AnPingerWrapper) Options() *anping.Options {
	return a.Opt
}
func (a *AnPingerWrapper) OnStart() {
	if a.logger != nil {
		a.logger.OnStart(a.Opt.Address(), a.State)
	}
}

func (a *AnPingerWrapper) OnRecv(t time.Duration) {
	if a.logger != nil {
		a.logger.OnRecv(a.Opt.Address(), a.State, t)
	}
}

func (a *AnPingerWrapper) OnLost(errMsg ...any) {
	if a.logger != nil {
		a.logger.OnLost(a.Opt.Address(), a.State, F.ToString(errMsg...))
	}
}

func (a *AnPingerWrapper) OnFinish() {
	if a.logger != nil {
		a.logger.OnFinish(a.Opt.Address(), a.State)
	}
}
