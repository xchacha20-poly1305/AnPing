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
	"github.com/xchacha20-poly1305/anping/state"
)

const Protocol = "udp"

func init() {
	anping.AnPingerCreator[Protocol] = New
}

var _ anping.AnPinger = (*UdpPinger)(nil)

type UdpPinger struct {
	Opt *anping.Options
	*state.State

	PayloadLength int

	logger state.Logger
}

func New(logWriter io.Writer) anping.AnPinger {
	return &UdpPinger{
		Opt:           anping.NewOptions(),
		State:         state.NewState(),
		PayloadLength: anping.PayloadLength,
		logger:        &state.DefaultLogger{Writer: logWriter},
	}
}

func (u *UdpPinger) Run() {
	u.RunContext(context.Background())
}

func (u *UdpPinger) RunContext(ctx context.Context) {
	if u.logger != nil {
		u.logger.OnStart(u.Opt.Address())
	}

	payload := make([]byte, u.PayloadLength)
	_, _ = rand.Read(payload)

	addr, err := net.ResolveUDPAddr("udp", u.Opt.Address())
	if err != nil {
		if writer, isWriter := u.logger.(io.Writer); isWriter {
			_, _ = io.WriteString(writer, err.Error())
		}
		return
	}

	defer func() {
		if u.logger != nil {
			u.logger.OnFinish(u.Opt.Address(), u.Probed(), u.Lost(), u.Succeed(), u.Min(), u.Max(), u.Avg(), u.Mdev())
		}
	}()

	for i := u.Opt.Count; i != 0; i-- {
		select {
		case <-ctx.Done():
			return
		default:
		}

		latency, err := Ping(addr, time.Duration(u.Opt.Timeout)*time.Millisecond, payload)
		u.Add(int(latency.Milliseconds()), err == nil)
		if !u.Opt.Quite && u.logger != nil {
			if err != nil {
				u.logger.OnLost(u.Opt.Address(), err.Error())
			} else {
				u.logger.OnRecv(u.Opt.Address(), int(latency.Milliseconds()))
			}
		}
		time.Sleep(u.Opt.Interval)
	}
}

func (u *UdpPinger) Protocol() string {
	return Protocol
}

func (u *UdpPinger) SetAddress(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
		port = anping.Port
	}

	if M.IsDomainName(host) {
		ip, err := anping.LookupSingleIP(host, u.Opt.DomainStrategy)
		if err != nil {
			return err
		}
		return u.Opt.SetAddress(net.JoinHostPort(ip.String(), port))
	}

	return u.Opt.SetAddress(address)
}

func (u *UdpPinger) SetLogger(logger state.Logger) {
	u.logger = logger
}

func (u *UdpPinger) Options() *anping.Options {
	return u.Opt
}

func Ping(addr net.Addr, timeout time.Duration, payload []byte) (time.Duration, error) {
	udpConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return -1, err
	}
	defer udpConn.Close()

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
