package tcpping

import (
	"context"
	"io"
	"net"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
)

const Protocol = "tcp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type TcpPinger struct {
	*anping.Options

	logger anping.LoggerNotNil
}

func New(logWriter io.Writer) anping.AnPinger {
	return &TcpPinger{
		Options: anping.NewOptions(),
		logger: anping.LoggerNotNil{
			L: &anping.DefaultLogger{
				Writer: logWriter,
			},
		},
	}
}

func (t *TcpPinger) Run() {
	t.RunContext(context.Background())
}

func (t *TcpPinger) RunContext(ctx context.Context) {
	t.logger.OnStart(t.Options)

	for i := t.Number(); i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		latency, err := Ping(t.Address(), time.Millisecond*time.Duration(t.Timeout()))
		t.Add(int(latency.Milliseconds()), err == nil)
		if err != nil {
			t.logger.OnLost(t.Options, err.Error())
			time.Sleep(t.Interval())
			continue
		}

		t.logger.OnRecv(t.Options, int(latency.Milliseconds()))
		time.Sleep(t.Interval())
	}
}

func (t *TcpPinger) Clean() error {
	t.Options.PrintedLogOnce.Do(
		func() {
			t.logger.OnFinish(t.Options)
		},
	)
	return nil
}

func (t *TcpPinger) Protocol() string {
	return Protocol
}

func (t *TcpPinger) SetLogger(logger anping.Logger) {
	t.logger = anping.LoggerNotNil{L: logger}
}

func (t *TcpPinger) SetAddress(address string) error {
	domain, port, err := net.SplitHostPort(address)
	if err != nil {
		return E.Cause(err, "parse address")
	}

	if M.IsDomainName(domain) {
		ip, err := anping.LookupSingleIP(domain, t.DomainStrategy())
		if err != nil {
			return err
		}
		_ = t.Options.SetAddress(net.JoinHostPort(ip.String(), port))
		return nil
	}

	_ = t.Options.SetAddress(address)
	return nil
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
