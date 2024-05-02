package anping

import (
	"time"

	M "github.com/sagernet/sing/common/metadata"
)

const (
	Count         = -1
	Timeout       = 3000 * time.Millisecond
	Interval      = time.Second
	Port          = "443"
	PayloadLength = 20
)

type Options struct {
	Count          int
	address        M.Socksaddr
	Timeout        time.Duration
	Interval       time.Duration
	Quite          bool
	DomainStrategy DomainStrategy
}

func NewOptions() *Options {
	return &Options{
		Count:          Count,
		Interval:       Interval,
		Timeout:        Timeout,
		Quite:          false,
		DomainStrategy: PreferNone,
	}
}

func (o *Options) Address() M.Socksaddr {
	return o.address
}

func (o *Options) SetAddress(address M.Socksaddr) error {
	o.address = address
	return nil
}
