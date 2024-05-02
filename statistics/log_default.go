package statistics

import (
	"fmt"
	"io"
	"strings"
	"time"

	M "github.com/sagernet/sing/common/metadata"
)

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	io.Writer
}

func (d *DefaultLogger) OnStart(address M.Socksaddr, statistic StatisticsGetter) {
	if d.Writer == nil {
		return
	}

	// PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
	_, _ = fmt.Fprintf(d.Writer, fmt.Sprintf("PING %s:\n", address.String()))
}

func (d *DefaultLogger) OnRecv(address M.Socksaddr, statistic StatisticsGetter, t time.Duration) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(
		d.Writer,
		fmt.Sprintf("From %s: time=%d ms\n", address.String(), t.Milliseconds()),
	)
}

func (d *DefaultLogger) OnLost(address M.Socksaddr, statistic StatisticsGetter, errMessage string,
	t time.Duration) {
	if d.Writer == nil {
		return
	}

	_, _ = fmt.Fprintf(d.Writer, "Failed to ping %s: %s\n", address.String(), errMessage)
}

func (d *DefaultLogger) OnFinish(address M.Socksaddr, statistics StatisticsGetter) {
	if d.Writer == nil {
		return
	}

	var b strings.Builder

	_, _ = fmt.Fprintf(&b, "\n--- %s ping statistics ---\n", address.String())

	_, _ = fmt.Fprintf(&b, "%d packets transmitted, %d packets received, %s packet loss\n",
		statistics.Probed(), statistics.Succeed(), percent(statistics.Lost(), statistics.Probed()))

	_, _ = fmt.Fprintf(&b, "round-trip min/avg/max/stddev = %d/%d/%d/%d ms\n",
		statistics.Min(), statistics.Avg(), statistics.Max(), statistics.Mdev())

	_, _ = io.WriteString(d.Writer, b.String())
}
