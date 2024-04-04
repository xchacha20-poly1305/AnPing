package icmpping

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/state"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "icmp"

var _ anping.AnPinger = (*IcmpPinger)(nil)

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type IcmpPinger struct {
	implement.AnPingerWrapper
	PayloadLength int
}

func New(logWriter io.Writer) anping.AnPinger {
	i := &IcmpPinger{
		AnPingerWrapper: implement.AnPingerWrapper{
			Opt:   anping.NewOptions(),
			State: state.NewState(),
		},
		PayloadLength: anping.PayloadLength,
	}
	i.SetLogger(&state.DefaultLogger{Writer: logWriter})
	return i
}

func (i *IcmpPinger) RunContext(ctx context.Context) {
	i.OnStart()

	defer i.OnFinish()

	for j := i.Opt.Count; j != 0; j-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		payload := make([]byte, i.PayloadLength)
		_, _ = rand.Read(payload)

		t, err := libping.IcmpPing(i.Opt.Address(), time.Duration(i.Opt.Timeout), payload)
		i.Add(int(t), err == nil)
		if !i.Opt.Quite {
			if err != nil {
				i.OnLost(err)
			} else {
				i.OnRecv(t)
			}
		}
		time.Sleep(i.Opt.Interval)
	}
}

func (i *IcmpPinger) Protocol() string {
	return Protocol
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
