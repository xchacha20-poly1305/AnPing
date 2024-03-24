package anping

import (
	"fmt"
	"io"
	"strings"
)

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	Writer io.Writer
}

func (d *DefaultLogger) OnStart(opts *Options) {
	if d.Writer == nil {
		return
	}

	// PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
	_, _ = io.WriteString(d.Writer, fmt.Sprintf("PING %s:\n", opts.Address()))
}

func (d *DefaultLogger) OnRecv(opts *Options, t int) {
	if opts.Quite() || d.Writer == nil {
		return
	}

	_, _ = io.WriteString(d.Writer, fmt.Sprintf("From %s: time=%d ms\n", opts.Address(), t))
}

func (d *DefaultLogger) OnLost(opts *Options, errMessage string) {
	if opts.Quite() || d.Writer == nil {
		return
	}

	_, _ = io.WriteString(d.Writer,
		fmt.Sprintf("Failed to ping %s: %s\n", opts.Address(), errMessage),
	)
}

func (d *DefaultLogger) OnFinish(opts *Options) {
	if d.Writer == nil {
		return
	}

	var b strings.Builder

	_, _ = b.WriteString(fmt.Sprintf("\n--- %s ping statistics ---\n", opts.Address()))

	probed := opts.Probed()
	lost := opts.Lost()
	_, _ = b.WriteString(
		fmt.Sprintf("%d packets transmitted, %d packets received, %s packet loss\n",
			probed, opts.Succeed(), percent(lost, probed)),
	)

	_, _ = b.WriteString(
		fmt.Sprintf("round-trip min/avg/max/stddev = %d/%d/%d/%d ms\n",
			opts.Min(), opts.Avg(), opts.Max(), opts.Mdev()),
	)

	_, _ = io.WriteString(d.Writer, b.String())
}
