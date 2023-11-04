package anping

import (
	"errors"
	"net/url"
	"strings"
)

type AnPinger struct {
	Addr    *url.URL
	Timeout int // ms
	Count   int
}

func New(addr string) (*AnPinger, error) {
	if !strings.Contains(addr, "://") {
		// Default to use ICMP
		addr = "icmp://" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	// 检测是否包含端口
	switch u.Port() {
	case "":
		if strings.HasSuffix(u.Host, ":") {
			u.Host += "443"
		} else {
			u.Host += ":443"
		}
	}

	return &AnPinger{
		Addr:    u,
		Timeout: 1,
		Count:   -1,
	}, nil
}

func (a *AnPinger) Start() error {

	switch a.Addr.Scheme {
	case "", "icmp":
		a.Icmpping()
	case "tcp":
		a.Tcpping()
	default:
		return errors.New("Unknow protocol: " + a.Addr.Scheme)
	}

	return nil
}
