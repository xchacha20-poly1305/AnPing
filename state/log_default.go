package state

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	io.Writer
}

func (d *DefaultLogger) OnStart(address string, state *State) {
	if d.Writer == nil {
		return
	}

	// PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
	_, _ = fmt.Fprintf(d.Writer, fmt.Sprintf("PING %s:\n", address))
}

func (d *DefaultLogger) OnRecv(address string, state *State, t time.Duration) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(d.Writer, fmt.Sprintf("From %s: time=%d ms\n", address, t.Milliseconds()))
}

func (d *DefaultLogger) OnLost(address string, state *State, errMessage string) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(d.Writer, "Failed to ping %s: %s\n", address, errMessage)
}

func (d *DefaultLogger) OnFinish(address string, state *State) {
	if d.Writer == nil {
		return
	}

	var b strings.Builder

	_, _ = fmt.Fprintf(&b, "\n--- %s ping statistics ---\n", address)

	_, _ = fmt.Fprintf(&b, "%d packets transmitted, %d packets received, %s packet loss\n",
		state.Probed(), state.Succeed(), percent(state.Lost(), state.Probed()))

	_, _ = fmt.Fprintf(&b, "round-trip min/avg/max/stddev = %d/%d/%d/%d ms\n",
		state.Min(), state.Avg(), state.Max(), state.Mdev())

	_, _ = io.WriteString(d.Writer, b.String())
}
