package anping

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
	"unsafe"

	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/x/constraints"
)

type Logger interface {
	OnStart()
	OnReceive(latency time.Duration)
	OnFail(err error)
	OnStop()
}

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	Type   string
	Target string
	Writer io.Writer
	Start  time.Time

	mu sync.Mutex

	total    uint64
	received uint64
	failed   uint64

	min, avg, max, mdev time.Duration
}

func NewLogger(Type, target string, writer io.Writer) *DefaultLogger {
	if writer == nil {
		writer = io.Discard
	}
	return &DefaultLogger{Type: Type, Target: target, Writer: writer}
}

func (d *DefaultLogger) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Start = time.Time{}
	d.total = 0
	d.received = 0
	d.failed = 0
	d.min = 0
	d.avg = 0
	d.max = 0
	d.mdev = 0
}

func (d *DefaultLogger) OnStart() {
	d.Start = time.Now()
	message := F.ToString(d.Type, " ", d.Target, " :\n")
	_, _ = d.Writer.Write(toBytes(&message))
}

func (d *DefaultLogger) OnReceive(latency time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.updateStatistic(latency)
	message := F.ToString("From ", d.Target, ": used ", latency.Milliseconds(), "ms\n")
	_, _ = d.Writer.Write(toBytes(&message))
}

// make sure lock before invoke it
func (d *DefaultLogger) updateStatistic(latency time.Duration) {
	d.total += 1
	d.received += 1

	if d.min <= 0 || latency < d.min {
		d.min = latency
	}
	if d.max <= 0 || latency > d.max {
		d.max = latency
	}

	if d.avg <= 0 {
		d.avg = latency
	} else {
		d.avg = (d.avg + latency) / 2
	}

	if d.mdev < 0 {
		d.mdev = 0
	} else {
		d.mdev = (d.mdev + abs(latency-d.avg)) / 2
	}
}

func (d *DefaultLogger) OnFail(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.total += 1
	d.failed += 1
	if E.IsClosedOrCanceled(err) {
		return
	}
	message := F.ToString("Failed: ", err, "\n")
	_, _ = d.Writer.Write(toBytes(&message))
}

func (d *DefaultLogger) OnStop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	builder := &bytes.Buffer{}
	_, _ = builder.WriteString(F.ToString("\n--- ", d.Target, " ", d.Type, " statistic ---\n"))
	_, _ = builder.WriteString(F.ToString("Has sent ", d.total, " packet, "))
	_, _ = builder.WriteString(F.ToString("received ", d.received, " packet, "))
	_, _ = fmt.Fprintf(builder, "%.2f packet loss, ", float64(d.received)/float64(d.total))
	_, _ = builder.WriteString(F.ToString("time ", time.Since(d.Start).String(), "\n"))
	_, _ = builder.WriteString(
		F.ToString(
			"rtt min/avg/max/mdev = ",
			d.min.Milliseconds(), "/",
			d.avg.Milliseconds(), "/",
			d.max.Milliseconds(), "/",
			d.mdev.Milliseconds(),
			" ms\n",
		),
	)
	_, _ = d.Writer.Write(builder.Bytes())
}

func toBytes(str *string) []byte {
	return *(*[]byte)(unsafe.Pointer(str))
}

func abs[T constraints.Integer | constraints.Float](num T) T {
	if num < 0 {
		return -num
	}
	return num
}
