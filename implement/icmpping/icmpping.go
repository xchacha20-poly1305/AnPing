package icmpping

import (
	"context"
	"crypto/rand"
	"io"
	"os"
	"time"

	M "github.com/sagernet/sing/common/metadata"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
	"github.com/xchacha20-poly1305/anping/statistics"
	"github.com/xchacha20-poly1305/libping"
	"golang.org/x/sys/unix"
)

const Protocol = "icmp"

var _ anping.AnPinger = (*IcmpPinger)(nil)

func init() {
	anping.AnPingerCreator[Protocol] = New
}

type IcmpPinger struct {
	*implement.AnPingerWrapper
	PayloadLength int
	done          chan struct{}
}

func New(logWriter io.Writer) anping.AnPinger {
	i := &IcmpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),
		PayloadLength:   anping.PayloadLength,
	}
	i.SetLogger(&statistics.DefaultLogger{Writer: logWriter})
	return i
}

func (i *IcmpPinger) Start(ctx context.Context) {
	i.done = make(chan struct{})
	go i.start(ctx)
}

func (i *IcmpPinger) start(ctx context.Context) {
	i.OnStart()
	context.AfterFunc(ctx, i.OnFinish)

	bindAddr := resolveBindIP(i.Opt.BindIP)
	timer := time.NewTimer(i.Opt.Interval)
	for j := i.Opt.Count; j != 0; j-- {
		select {
		case <-ctx.Done():
			return
		case <-i.done:
			return
		case <-timer.C:
			timer.Reset(i.Opt.Interval)
		}

		payload := make([]byte, i.PayloadLength)
		_, _ = rand.Read(payload)

		sendCtx, cancel := context.WithTimeout(ctx, i.Opt.Timeout)
		t, err := libping.IcmpPing(sendCtx, i.Opt.Address(), payload, bindAddr)
		cancel()
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

func (i *IcmpPinger) Close() error {
	select {
	case <-i.done:
		return os.ErrClosed
	default:
		close(i.done)
		return nil
	}
}

func (i *IcmpPinger) Done() <-chan struct{} {
	return i.done
}

func (i *IcmpPinger) Clone() anping.AnPinger {
	return &IcmpPinger{
		AnPingerWrapper: i.AnPingerWrapper.Clone().(*implement.AnPingerWrapper),
		PayloadLength:   i.PayloadLength,
	}
}

func resolveBindIP(bindIP string) unix.Sockaddr {
	if bindIP == "" {
		return nil
	}

	addr := M.ParseSocksaddr(bindIP)
	// Turn port to 0 to prevent use same port.
	if addr.IsIPv6() {
		return &unix.SockaddrInet6{
			Port: 0,
			Addr: [16]byte(addr.IPAddr().IP),
		}
	}
	if addr.IsIPv4() {
		return &unix.SockaddrInet4{
			Port: 0,
			Addr: [4]byte(addr.IPAddr().IP),
		}
	}

	return nil
}
