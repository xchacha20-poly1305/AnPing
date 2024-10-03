package anping

import (
	"context"
)

// AnPing is the abstract interface for pinger implementations.
type AnPing interface {
	// Start starts AnPing service.
	// You can cancel it by cancel ctx.
	Start(ctx context.Context) error

	// Settings returns settings for AnPing.
	Settings() Settings

	// Done marks AnPing stop.
	Done() <-chan struct{}
}
