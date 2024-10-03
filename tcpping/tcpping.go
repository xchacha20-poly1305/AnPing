package tcpping

import (
	"context"
	"os"
	"time"

	"github.com/sagernet/sing/common/control"

	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/libping"
)

const Type = "tcp"

var _ anping.AnPing = (*TcpPinger)(nil)

type TcpPinger struct {
	TcpSettings
	done chan struct{}
}

func New(settings TcpSettings) (t *TcpPinger) {
	t = &TcpPinger{
		TcpSettings: settings,
	}
	return t
}

func (t *TcpPinger) Start(ctx context.Context) error {
	var controlFunc control.Func
	if d := t.Interface(); d != "" {
		controlFunc = control.Append(controlFunc, control.BindToInterface(t.InterfaceFinder, d, -1))
	}
	if t.done != nil {
		select {
		case <-t.done:
		default:
			return os.ErrExist
		}
	}
	t.done = make(chan struct{})
	go t.loop(ctx, controlFunc)
	return nil
}

func (t *TcpPinger) loop(ctx context.Context, controlFunc control.Func) {
	defer func() {
		t.Logger().OnStop()
		close(t.done)
	}()

	var shouldStop func(int64) bool
	count := t.Count()
	if count <= 0 {
		shouldStop = func(_ int64) bool {
			return true
		}
	} else {
		shouldStop = func(i int64) bool {
			return i != 0
		}
	}

	t.Logger().OnStart()
	for i := count; shouldStop(i); i-- {
		select {
		case <-ctx.Done():
			return
		case <-time.After(t.Interval()):
		}
		ctx, cancel := context.WithTimeout(ctx, t.Timeout())
		latency, err := libping.TcpPing(ctx, controlFunc, t.Target())
		cancel()
		if err == nil {
			t.Logger().OnReceive(latency)
		} else {
			t.Logger().OnFail(err)
		}
	}
}

func (t *TcpPinger) Settings() anping.Settings {
	return &t.TcpSettings
}

func (t *TcpPinger) Done() <-chan struct{} {
	return t.done
}
