package state

import (
	"fmt"
	"io"
	"strings"
)

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	Writer io.Writer
}

func (d *DefaultLogger) OnStart(address string) {
	if d.Writer == nil {
		return
	}

	// PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
	_, _ = fmt.Fprintf(d.Writer, fmt.Sprintf("PING %s:\n", address))
}

func (d *DefaultLogger) OnRecv(address string, t int) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(d.Writer, fmt.Sprintf("From %s: time=%d ms\n", address, t))
}

func (d *DefaultLogger) OnLost(address, errMessage string) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(d.Writer, "Failed to ping %s: %s\n", address, errMessage)
}

func (d *DefaultLogger) OnFinish(address string, probed, lost, succeed, min, avg, max, mdev uint64) {
	if d.Writer == nil {
		return
	}

	var b strings.Builder

	_, _ = b.WriteString(fmt.Sprintf("\n--- %s ping statistics ---\n", address))

	_, _ = b.WriteString(
		fmt.Sprintf("%d packets transmitted, %d packets received, %s packet loss\n",
			probed, succeed, percent(lost, probed)),
	)

	_, _ = b.WriteString(
		fmt.Sprintf("round-trip min/avg/max/stddev = %d/%d/%d/%d ms\n",
			min, avg, max, mdev),
	)

	_, _ = io.WriteString(d.Writer, b.String())
}
