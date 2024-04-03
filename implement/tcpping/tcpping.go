package tcpping

import (
	"context"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/state"
)

const Protocol = "tcp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type TcpPinger struct {
	Opt *anping.Options
	*state.State

	logger state.Logger
}

func New(logWriter io.Writer) anping.AnPinger {
	return &TcpPinger{
		Opt:    anping.NewOptions(),
		State:  state.NewState(),
		logger: &state.DefaultLogger{Writer: logWriter},
	}
}

func (t *TcpPinger) Run() {
	t.RunContext(context.Background())
}

func (t *TcpPinger) RunContext(ctx context.Context) {
	if t.logger != nil {
		t.logger.OnStart(t.Opt.Address())
	}

	for i := t.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		latency, err := Ping(t.Opt.Address(), time.Millisecond*time.Duration(t.Opt.Timeout))
		t.Add(int(latency.Milliseconds()), err == nil)
		if !t.Opt.Quite && t.logger != nil {
			if err != nil {
				t.logger.OnLost(t.Opt.Address(), err.Error())
			} else {
				t.logger.OnRecv(t.Opt.Address(), int(latency.Milliseconds()))
			}
		}
		time.Sleep(t.Opt.Interval)
	}
}

func (t *TcpPinger) Clean() error {
	t.State.FinishOnce.Do(
		func() {
			if t.logger != nil {
				t.logger.OnFinish(t.Opt.Address(), t.Probed(), t.Lost(), t.Succeed(), t.Min(), t.Max(), t.Avg(), t.Mdev())
			}
		},
	)
	return nil
}

func (t *TcpPinger) Protocol() string {
	return Protocol
}

func (t *TcpPinger) SetLogger(logger state.Logger) {
	t.logger = logger
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

func (t *TcpPinger) Options() *anping.Options {
	return t.Opt
}

func Ping(address string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout(N.NetworkTCP, address, timeout)
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	return time.Since(start), nil
}
