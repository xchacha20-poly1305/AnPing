package udpping

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
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
}

func New(logWriter io.Writer) anping.AnPinger {
	u := &UdpPinger{
		AnPingerWrapper: implement.New(logWriter).(*implement.AnPingerWrapper),

		PayloadLength: anping.PayloadLength,
	}
	return u
}

func (u *UdpPinger) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go u.start(ctx, done)
	return done
}

func (u *UdpPinger) start(ctx context.Context, done chan struct{}) {
	defer implement.TryCloseDone(done)
	u.OnStart()
	defer u.OnFinish()

	payload := make([]byte, u.PayloadLength)
	_, _ = rand.Read(payload)

	timer := time.NewTimer(u.Opt.Interval)
	defer timer.Stop()
	for i := u.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-timer.C:
			timer.Reset(u.Opt.Interval)
		}

		latency, err := Ping(u.Opt.Address(), u.Opt.Timeout, payload)
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

func Ping(addr M.Socksaddr, timeout time.Duration, payload []byte) (time.Duration, error) {
	udpConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = udpConn.Close()
		udpConn = nil
	}()

	start := time.Now()
	_, err = udpConn.WriteTo(payload, addr)
	if err != nil {
		return -1, E.Cause(err, "write to udpConn")
	}

	_ = udpConn.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1)
	_, err = udpConn.Read(buf)
	if err != nil {
		return -1, E.Cause(err, "read udpConn")
	}
	return time.Since(start), nil
}
