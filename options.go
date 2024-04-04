package anping

import (
	"time"
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
	address        string
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

func (o *Options) Address() string {
	return o.address
}

// SetAddress sets the target address. If you are external user, please use AnPinger.SetAddress instead it.
func (o *Options) SetAddress(address string) error {
	o.address = address
	return nil
}
