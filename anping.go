package anping

import (
	"context"
	"io"

	"github.com/xchacha20-poly1305/anping/statistics"
)

type InitPinger func(logWriter io.Writer) AnPinger

var AnPingerCreator = map[string]InitPinger{}

// AnPinger is an abstract interface to ping.
// You can absorb it by canal context.
type AnPinger interface {
	// Start starts AnPinger until context cancel or other reason.
	Start(ctx context.Context) (done <-chan struct{})

	// Protocol returns the protocol of AnPinger.
	Protocol() string

	// Statistics provides read-time statistics of ping. It will refresh after use Start.
	Statistics() statistics.StatisticsGetter

	ResetStatistics()

	// SetAddress set the address of pinger. Use it instead of set it in options.
	SetAddress(address string) error

	SetLogger(logger statistics.Logger)

	Options() *Options
}
