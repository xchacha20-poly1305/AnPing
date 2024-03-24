package anping

import (
	"math"
	"sync"
	"time"

	"github.com/sagernet/sing/common/atomic"
)

const (
	Number   = -1
	Timeout  = 3000
	Interval = time.Second
)

type Options struct {
	number         int
	address        string
	timeout        int32
	interval       time.Duration
	quite          bool
	domainStrategy DomainStrategy

	PrintedLogOnce sync.Once

	// probed is the time that probed.
	probed atomic.Uint64
	// succeed is the time that probing succeed.
	succeed atomic.Uint64
	// lost is the time that packet lost.
	lost atomic.Uint64
	// min is the minimum of probing time.
	min atomic.Uint64
	// max is the maximum of probing time.
	max atomic.Uint64

	avg  atomic.Uint64
	mdev atomic.Uint64
}

func NewOptions() *Options {
	opts := &Options{
		number:   Number,
		interval: Interval,
		timeout:  Timeout,
	}
	opts.min.Store(math.MaxUint64)
	return opts
}

func (o *Options) Address() string {
	return o.address
}

func (o *Options) SetAddress(address string) error {
	o.address = address
	return nil
}

func (o *Options) Number() int {
	return o.number
}

func (o *Options) SetNumber(number int) {
	if number == 0 {
		o.number = Number
		return
	}

	o.number = number
}

func (o *Options) Timeout() int32 {
	return o.timeout
}

func (o *Options) SetTimeout(timeout int32) {
	o.timeout = timeout
}

func (o *Options) Interval() time.Duration {
	return o.interval
}

func (o *Options) SetInterval(t time.Duration) {
	if t < 0 {
		o.interval = Interval
		return
	}

	o.interval = t
}

func (o *Options) Quite() bool {
	return o.quite
}

func (o *Options) SetQuite(yes bool) {
	o.quite = yes
}

func (o *Options) DomainStrategy() DomainStrategy {
	return o.domainStrategy
}

func (o *Options) SetDomainStrategy(d DomainStrategy) {
	o.domainStrategy = d
}

func (o *Options) Add(t int, success bool) {
	o.probed.Add(1)
	if !success {
		o.lost.Add(1)
		return
	}

	uintTime := uint64(t)
	o.succeed.Add(1)
	if o.min.Load() > uintTime {
		o.min.Store(uintTime)
	}
	if o.max.Load() < uintTime {
		o.max.Store(uintTime)
	}

	avg := o.avg.Load()
	if avg == 0 {
		o.avg.Store(uintTime)
	} else {
		o.avg.Store((avg + uintTime) / 2)
	}

	adev := o.mdev.Load()
	abs := diffAbs(adev, avg)
	if adev == 0 {
		o.mdev.Store(abs)
	} else {
		o.mdev.Store((adev + abs) / 2)
	}
}

func (o *Options) Probed() uint64 {
	return o.probed.Load()
}

func (o *Options) Succeed() uint64 {
	return o.succeed.Load()
}

func (o *Options) Lost() uint64 {
	return o.lost.Load()
}

func (o *Options) Min() uint64 {
	return o.max.Load()
}

func (o *Options) Max() uint64 {
	return o.max.Load()
}

func (o *Options) Avg() uint64 {
	return o.avg.Load()
}

func (o *Options) Mdev() uint64 {
	return o.mdev.Load()
}
