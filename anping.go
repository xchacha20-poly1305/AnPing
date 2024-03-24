package anping

import (
	"context"
	"io"
	"time"
)

type InitPinger func(logWriter io.Writer) AnPinger

var AnPingerCreator = map[string]InitPinger{}

// AnPinger is an abstract interface to ping.
// You can absorb it by canal context.
type AnPinger interface {
	// Run runs AnPinger.
	// If you want to control to stop it, please use RunContext.
	Run()

	// RunContext runs AnPinger with context. It will block the thread.
	RunContext(ctx context.Context)

	// Clean used to do some chores when finished. Such as print log.
	Clean() error

	// Protocol returns the protocol of AnPinger.
	Protocol() string

	SetLogger(logger Logger)

	// Address returns the target address.
	Address() string

	// SetAddress sets the target address.
	SetAddress(address string) error

	// Number returns the number of runs.
	Number() int

	// SetNumber sets the number of runs.
	// If number <= 0, AnPinger will run forever.
	SetNumber(number int)

	Timeout() int32

	SetTimeout(timeout int32)

	Interval() time.Duration
	SetInterval(i time.Duration)

	// Quite returns true if enabled quite mode
	Quite() bool
	SetQuite(yes bool)
}
