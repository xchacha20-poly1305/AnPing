package icmpping

import (
	"context"
	"os"
	"time"

	"github.com/sagernet/sing/common/control"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/libping"
)

const Type = "icmp"

var _ anping.AnPing = (*IcmpPinger)(nil)

type IcmpPinger struct {
	IcmpSettings
	done chan struct{}
}

func New(settings IcmpSettings) (t *IcmpPinger) {
	t = &IcmpPinger{
		IcmpSettings: settings,
	}
	return t
}

func (p *IcmpPinger) Start(ctx context.Context) error {
	var controlFunc control.Func
	if d := p.Interface(); d != "" {
		controlFunc = control.Append(controlFunc, control.BindToInterface(p.InterfaceFinder, d, -1))
	}
	if p.done != nil {
		select {
		case <-p.done:
		default:
			return os.ErrExist
		}
	}
	p.done = make(chan struct{})
	go p.loop(ctx, controlFunc)
	return nil
}

func (p *IcmpPinger) loop(ctx context.Context, controlFunc control.Func) {
	defer func() {
		p.Logger().OnStop()
		close(p.done)
	}()

	var shouldStop func(int64) bool
	count := p.Count()
	if count <= 0 {
		shouldStop = func(_ int64) bool {
			return true
		}
	} else {
		shouldStop = func(i int64) bool {
			return i != 0
		}
	}

	payload := p.payload()

	p.Logger().OnStart()
	for i := count; shouldStop(i); i-- {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval()):
		}
		ctx, cancel := context.WithTimeout(ctx, p.Timeout())
		latency, err := libping.IcmpPing(ctx, p.Target(), payload, controlFunc)
		cancel()
		if err == nil {
			p.Logger().OnReceive(latency)
		} else {
			p.Logger().OnFail(err)
		}
	}
}

func (p *IcmpPinger) Settings() anping.Settings {
	return &p.IcmpSettings
}

func (p *IcmpPinger) Done() <-chan struct{} {
	return p.done
}
