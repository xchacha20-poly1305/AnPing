package anping

import (
	"context"
	"io"
	"math"

	"github.com/sagernet/sing/common/atomic"
)

type InitPinger func(logWriter io.Writer) AnPinger

var AnPingerCreator = map[string]InitPinger{}

// AnPinger is an abstract interface to ping.
// You can absorb it by canal context.
type AnPinger interface {
	// Run runs AnPinger.
	// If you want to control to stop it, please use RunContext.
	Run()

	// RunContext runs AnPinger with context.
	RunContext(ctx context.Context)

	// Protocol returns the protocol of AnPinger.
	Protocol() string

	SetLogger(logger Logger)

	// SetAddress sets the target address.
	SetAddress(address string) error

	// SetNumber sets the number of runs.
	// If number <= 0, AnPinger will run forever.
	SetNumber(number int)

	SetTimeout(timeout int32)
}

type Options struct {
	Number  int
	Address string
	Timeout int32

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
	opts := &Options{}
	opts.min.Store(math.MaxUint64)
	return opts
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
