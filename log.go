package anping

import (
	"fmt"
	"io"
	"strings"
)

type Logger interface {
	OnFinish(opts *Options)
}

var _ Logger = (*LoggerNotNil)(nil)

type LoggerNotNil struct {
	L Logger
}

func (l *LoggerNotNil) OnFinish(opts *Options) {
	if l.L != nil {
		l.L.OnFinish(opts)
	}
}

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	Writer io.Writer
}

func (d *DefaultLogger) OnFinish(opts *Options) {
	if d.Writer == nil {
		return
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n--- %s ping statistics ---\n", opts.Address))

	probed := opts.Probed()
	lost := opts.Lost()
	b.WriteString(fmt.Sprintf("%d packets transmitted, %d packets received, %.2f%% packet loss\n",
		probed, opts.Succeed(), float64(probed)/float64(lost)*100))

	b.WriteString(fmt.Sprintf("round-trip min/avg/max/stddev = %d/%d/%d/%d\n",
		opts.Min(), opts.Avg(), opts.Max(), opts.Mdev()))

	_, _ = io.WriteString(d.Writer, b.String())
}
