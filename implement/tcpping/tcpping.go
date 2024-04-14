package tcpping

import (
	"context"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/state"
)

const Protocol = "tcp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type TcpPinger struct {
	implement.AnPingerWrapper
}

func New(logWriter io.Writer) anping.AnPinger {
	t := &TcpPinger{
		AnPingerWrapper: implement.AnPingerWrapper{
			Opt:   anping.NewOptions(),
			State: state.NewState(),
		},
	}
	t.SetLogger(&state.DefaultLogger{Writer: logWriter})

	return t
}

func (t *TcpPinger) RunContext(ctx context.Context) {
	t.OnStart()

	go context.AfterFunc(ctx, t.OnFinish)

	for i := t.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		latency, err := Ping(t.Opt.Address(), t.Opt.Timeout)
		t.Add(uint64(latency.Milliseconds()), err == nil)
		if !t.Opt.Quite {
			if err != nil {
				t.OnLost(err)
			} else {
				t.OnRecv(latency)
			}
		}
		time.Sleep(t.Opt.Interval)
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

func Ping(address string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout(N.NetworkTCP, address, timeout)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = conn.Close()
		conn = nil
	}()

	return time.Since(start), nil
}
