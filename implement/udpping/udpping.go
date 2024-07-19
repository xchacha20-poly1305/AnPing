package udpping

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"os"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/implement"
)

const Protocol = "udp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

var _ anping.AnPinger = (*UdpPinger)(nil)

type UdpPinger struct {
	*implement.AnPingerWrapper

	PayloadLength int
	done          chan struct{}
}

func New(logWriter io.Writer) anping.AnPinger {
	u := &UdpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),

		PayloadLength: anping.PayloadLength,
	}
	return u
}

func (u *UdpPinger) Start(ctx context.Context) {
	u.done = make(chan struct{})
	go u.start(ctx)
}

func (u *UdpPinger) start(ctx context.Context) {
	u.OnStart()
	context.AfterFunc(ctx, u.OnFinish)

	payload := make([]byte, u.PayloadLength)
	_, _ = rand.Read(payload)
	dialer := &net.Dialer{}
	if udpAddr, _ := net.ResolveUDPAddr(N.NetworkUDP, u.Opt.BindIP); udpAddr != nil {
		udpAddr.Port = 0 // Not use same port
		dialer.LocalAddr = udpAddr
	}

	timer := time.NewTimer(u.Opt.Interval)
	defer timer.Stop()
	for i := u.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-u.done:
			return
		case <-timer.C:
			timer.Reset(u.Opt.Interval)
		}

		sendCtx, cancel := context.WithTimeout(ctx, u.Opt.Timeout)
		latency, err := Ping(
			sendCtx,
			dialer,
			u.Opt.Address(),
			payload,
		)
		cancel()
		u.Sta.Add(uint64(latency.Milliseconds()), err == nil)
		if !u.Opt.Quite {
			if err != nil {
				u.OnLost(err)
			} else {
				u.OnRecv(latency)
			}
		}
	}
}

func (u *UdpPinger) Protocol() string {
	return Protocol
}

func (u *UdpPinger) SetAddress(address M.Socksaddr) error {
	if !address.IsIP() {
		ip, err := anping.LookupSingleIP(address, u.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return u.Opt.SetAddress(M.ParseSocksaddrHostPort(ip.String(), address.Port))
	}

	return u.Opt.SetAddress(address)
}

func (u *UdpPinger) Close() error {
	select {
	case <-u.done:
		return os.ErrClosed
	default:
		close(u.done)
		return nil
	}
}

func (u *UdpPinger) Done() <-chan struct{} {
	return u.done
}

func (u *UdpPinger) Clone() anping.AnPinger {
	return &UdpPinger{
		AnPingerWrapper: u.AnPingerWrapper.Clone().(*implement.AnPingerWrapper),
		PayloadLength:   u.PayloadLength,
	}
}

func Ping(ctx context.Context, dialer *net.Dialer, addr M.Socksaddr, payload []byte) (time.Duration, error) {
	conn, err := dialer.DialContext(ctx, N.NetworkUDP, addr.String())
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	context.AfterFunc(ctx, func() {
		_ = conn.SetDeadline(time.Now())
	})

	start := time.Now()
	_, err = conn.Write(payload)
	if err != nil {
		return -1, E.Cause(err, "write to udpConn")
	}

	buf := make([]byte, 1)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return -1, E.Cause(err, "read udpConn")
	}
	return time.Since(start), nil
}
