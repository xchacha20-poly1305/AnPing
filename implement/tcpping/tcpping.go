package tcpping

import (
	"context"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/statistics"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "tcp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type TcpPinger struct {
	*implement.AnPingerWrapper
	done chan struct{}
}

func New(logWriter io.Writer) anping.AnPinger {
	t := &TcpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),
	}
	t.SetLogger(&statistics.DefaultLogger{Writer: logWriter})

	return t
}

func (t *TcpPinger) Start(ctx context.Context) {
	t.done = make(chan struct{})
	go t.start(ctx)
}

func (t *TcpPinger) start(ctx context.Context) {
	t.OnStart()
	context.AfterFunc(ctx, t.OnFinish)

	timer := time.NewTimer(t.Opt.Interval)
	for i := t.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-t.done:
			return
		case <-timer.C:
			timer.Reset(t.Opt.Interval)
		}

		sendCtx, cancel := context.WithTimeout(ctx, t.Opt.Timeout)
		latency, err := libping.TcpPing(sendCtx, t.Opt.Address())
		cancel()
		t.Sta.Add(uint64(latency.Milliseconds()), err == nil)
		if !t.Opt.Quite {
			if err != nil {
				t.OnLost(err)
			} else {
				t.OnRecv(latency)
			}
		}
	}
}

func (t *TcpPinger) Protocol() string {
	return Protocol
}

func (t *TcpPinger) SetAddress(address M.Socksaddr) error {
	if !address.IsIP() {
		ip, err := anping.LookupSingleIP(address, t.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return t.Opt.SetAddress(M.ParseSocksaddrHostPort(ip.String(), address.Port))
	}

	return t.Opt.SetAddress(address)
}

func (t *TcpPinger) Close() error {
	select {
	case <-t.done:
		return net.ErrClosed
	default:
		close(t.done)
		return nil
	}
}

func (t *TcpPinger) Done() <-chan struct{} {
	return t.done
}

func (t *TcpPinger) Clone() anping.AnPinger {
	return &TcpPinger{
		AnPingerWrapper: t.AnPingerWrapper.Clone().(*implement.AnPingerWrapper),
	}
}
