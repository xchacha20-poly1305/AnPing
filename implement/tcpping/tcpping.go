package tcpping

import (
	"context"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/state"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "tcp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type TcpPinger struct {
	*implement.AnPingerWrapper
}

func New(logWriter io.Writer) anping.AnPinger {
	t := &TcpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),
	}
	t.SetLogger(&state.DefaultLogger{Writer: logWriter})

	return t
}

func (t *TcpPinger) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go t.start(ctx, done)
	return done
}

func (t *TcpPinger) start(ctx context.Context, done chan struct{}) {
	defer implement.TryCloseDone(done)
	t.OnStart()
	defer t.OnFinish()

	host, port, _ := net.SplitHostPort(t.Opt.Address())

	timer := time.NewTimer(t.Opt.Interval)
	for i := t.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-timer.C:
			timer.Reset(t.Opt.Interval)
		}

		latency, err := libping.TcpPing(host, port, t.Opt.Timeout)
		t.Add(uint64(latency.Milliseconds()), err == nil)
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

func (t *TcpPinger) SetAddress(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
		port = anping.Port
	}

	if M.IsDomainName(host) {
		ip, err := anping.LookupSingleIP(host, t.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return t.Opt.SetAddress(net.JoinHostPort(ip.String(), port))
	}

	return t.Opt.SetAddress(address)
}
