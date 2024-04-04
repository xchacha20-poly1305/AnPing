package icmpping

import (
	"context"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/state"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "icmp"

var _ anping.AnPinger = (*IcmpPinger)(nil)

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type IcmpPinger struct {
	Opt *anping.Options
	*state.State

	logger state.Logger
}

func New(logWriter io.Writer) anping.AnPinger {
	return &IcmpPinger{
		Opt:    anping.NewOptions(),
		State:  state.NewState(),
		logger: &state.DefaultLogger{Writer: logWriter},
	}
}

func (i *IcmpPinger) Run() {
	i.RunContext(context.Background())
}

func (i *IcmpPinger) RunContext(ctx context.Context) {
	if i.logger != nil {
		i.logger.OnStart(i.Opt.Address())
	}

	defer func() {
		if i.logger != nil {
			i.logger.OnFinish(i.Opt.Address(), i.Probed(), i.Lost(), i.Succeed(), i.Min(), i.Avg(), i.Max(), i.Mdev())
		}
	}()

	for j := i.Opt.Count; j != 0; j-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		t, err := libping.IcmpPing(i.Opt.Address(), i.Opt.Timeout)
		i.Add(int(t), err == nil)
		if !i.Opt.Quite && i.logger != nil {
			if err != nil {
				i.logger.OnLost(i.Opt.Address(), err.Error())
			} else {
				i.logger.OnRecv(i.Opt.Address(), int(t))
			}
		}
		time.Sleep(i.Opt.Interval)
	}
}

func (i *IcmpPinger) Protocol() string {
	return Protocol
}

func (i *IcmpPinger) SetLogger(logger state.Logger) {
	i.logger = logger
}

func (i *IcmpPinger) SetAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	if M.IsDomainName(host) {
		ip, err := anping.LookupSingleIP(host, i.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return i.Opt.SetAddress(ip.String())
	}

	return i.Opt.SetAddress(host)
}

func (i *IcmpPinger) Options() *anping.Options {
	return i.Opt
}
