package udpping

// Credit:
// https://github.com/wangyu-/UDPping

import (
	"context"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
)

const Type = "udp"

var listenPool = sync.Pool{
	New: func() any {
		return &net.ListenConfig{}
	},
}

func UdpPing(ctx context.Context, target M.Socksaddr, payload []byte, controlFunc control.Func) (latency time.Duration, err error) {
	listener := listenPool.Get().(*net.ListenConfig)
	defer func() {
		listener.Control = nil
		listenPool.Put(listener)
	}()
	listener.Control = controlFunc

	conn, err := listener.ListenPacket(ctx, network.NetworkUDP, "")
	if err != nil {
		return -1, E.Cause(err, "listen packet")
	}
	stop := context.AfterFunc(ctx, func() {
		_ = conn.Close()
	})
	defer func() {
		stop()
		_ = conn.Close()
	}()

	start := time.Now()
	_, err = conn.WriteTo(payload, target.UDPAddr())
	if err != nil {
		return -1, E.Cause(err, "send payload")
	}
	_, _, err = conn.ReadFrom(make([]byte, 1))
	if err != nil {
		return -1, E.Cause(err, "read receive")
	}

	return time.Since(start), nil
}

var _ anping.AnPing = (*UdpPinger)(nil)

type UdpPinger struct {
	UdpSettings
	done chan struct{}
}

func New(settings UdpSettings) (t *UdpPinger) {
	t = &UdpPinger{
		UdpSettings: settings,
	}
	return t
}

func (u *UdpPinger) Start(ctx context.Context) error {
	var controlFunc control.Func
	if d := u.Interface(); d != "" {
		controlFunc = control.Append(controlFunc, control.BindToInterface(u.InterfaceFinder, d, -1))
	}
	if u.done != nil {
		select {
		case <-u.done:
		default:
			return os.ErrExist
		}
	}
	u.done = make(chan struct{})
	go u.loop(ctx, controlFunc)
	return nil
}

func (u *UdpPinger) loop(ctx context.Context, controlFunc control.Func) {
	defer func() {
		u.Logger().OnStop()
		close(u.done)
	}()

	var shouldStop func(int64) bool
	count := u.Count()
	if count <= 0 {
		shouldStop = func(_ int64) bool {
			return true
		}
	} else {
		shouldStop = func(i int64) bool {
			return i != 0
		}
	}

	payload := u.payload()

	u.Logger().OnStart()
	for i := count; shouldStop(i); i-- {
		select {
		case <-ctx.Done():
			return
		case <-time.After(u.Interval()):
		}
		ctx, cancel := context.WithTimeout(ctx, u.Timeout())
		latency, err := UdpPing(ctx, u.Target(), payload, controlFunc)
		cancel()
		if err == nil {
			u.Logger().OnReceive(latency)
		} else {
			u.Logger().OnFail(err)
		}
	}
}

func (u *UdpPinger) Settings() anping.Settings {
	return &u.UdpSettings
}

func (u *UdpPinger) Done() <-chan struct{} {
	return u.done
}
