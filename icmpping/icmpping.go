package icmpping

import (
	"context"
	"errors"
	"io"
	"net"

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
	defer i.logger.OnFinish(i.Options)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		t, err := libping.IcmpPing(i.Address, i.Timeout)
		i.Add(int(t), err == nil)
		// TODO: log
	}
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
		return errors.New("ICMP shouldn't has port")
	}
	i.Address = address
	return nil
}

func (i *IcmpPinger) SetNumber(number int) {
	i.Number = number
}

func (i *IcmpPinger) SetTimeout(timeout int32) {
	i.Timeout = timeout
}
