package icmpping

import (
	"context"
	"io"
	"math/rand/v2"
	"net"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "icmp"

var _ anping.AnPinger = (*IcmpPinger)(nil)

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type IcmpPinger struct {
	*anping.Options

	logger anping.LoggerNotNil
}

func New(logWriter io.Writer) anping.AnPinger {
	return &IcmpPinger{
		Options: anping.NewOptions(),
		logger: anping.LoggerNotNil{
			L: &anping.DefaultLogger{
				Writer: logWriter,
			},
		},
	}
}

func (i *IcmpPinger) Run() {
	i.RunContext(context.Background())
}

func (i *IcmpPinger) RunContext(ctx context.Context) {
	i.logger.OnStart(i.Options)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		t, err := libping.IcmpPing(i.Address(), i.Timeout())
		i.Add(int(t), err == nil)
		if err != nil {
			i.logger.OnLost(i.Options, err.Error())
			time.Sleep(i.Interval())
			continue
		}

		i.logger.OnRecv(i.Options, int(t))
		time.Sleep(i.Interval())
	}
}

func (i *IcmpPinger) Clean() error {
	i.logger.OnFinish(i.Options)
	return nil
}

func (i *IcmpPinger) Protocol() string {
	return Protocol
}

func (i *IcmpPinger) SetLogger(logger anping.Logger) {
	i.logger = anping.LoggerNotNil{L: logger}
}

func (i *IcmpPinger) SetAddress(address string) error {
	_, _, err := net.SplitHostPort(address)
	if err == nil {
		return E.New("ICMP shouldn't has port")
	}

	if M.IsDomainName(address) {
		ip, err := net.LookupIP(address)
		if err != nil {
			return E.Cause(err, "look up ip")
		}
		i.Options.SetAddress(ip[rand.IntN(len(ip))].String())
		return nil
	}

	i.Options.SetAddress(address)
	return nil
}
