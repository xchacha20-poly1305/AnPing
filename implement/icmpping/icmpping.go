package icmpping

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/statistics"
	"github.com/xchacha20-poly1305/libping"
)

const Protocol = "icmp"

var _ anping.AnPinger = (*IcmpPinger)(nil)

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type IcmpPinger struct {
	*implement.AnPingerWrapper
	PayloadLength int
}

func New(logWriter io.Writer) anping.AnPinger {
	i := &IcmpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),
		PayloadLength:   anping.PayloadLength,
	}
	i.SetLogger(&statistics.DefaultLogger{Writer: logWriter})
	return i
}

func (i *IcmpPinger) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go i.start(ctx, done)
	return done
}

func (i *IcmpPinger) start(ctx context.Context, done chan struct{}) {
	defer implement.TryCloseDone(done)
	i.OnStart()
	defer i.OnFinish()

	timer := time.NewTimer(i.Opt.Interval)
	for j := i.Opt.Count; j != 0; j-- {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-timer.C:
			timer.Reset(i.Opt.Interval)
		}

		payload := make([]byte, i.PayloadLength)
		_, _ = rand.Read(payload)

		t, err := libping.IcmpPing(i.Opt.Address().AddrString(), i.Opt.Timeout, payload)
		i.Sta.Add(uint64(t.Milliseconds()), err == nil)
		if !i.Opt.Quite {
			if err != nil {
				i.OnLost(err)
			} else {
				i.OnRecv(t)
			}
		}
	}
}

func (i *IcmpPinger) Protocol() string {
	return Protocol
}

func (i *IcmpPinger) SetAddress(address M.Socksaddr) error {
	if !address.IsIP() {
		ip, err := anping.LookupSingleIP(address, i.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return i.Opt.SetAddress(M.ParseSocksaddrHostPort(ip.String(), address.Port))
	}

	return i.Opt.SetAddress(address)
}
