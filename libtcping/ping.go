package libtcping

import (
	"net"
	"time"
)

type resault struct {
	UsedTime time.Duration
	Err      error
}

func newResault() *resault {
	return &resault{
		UsedTime: 0,
		Err:      nil,
	}
}

// Credits:
// https://github.com/i3h/tcping/blob/6f55d63cd789777706bd46000f851374f732b84b/pkg/tcping/tcping.go
func (p *Pinger) ping() {
	p.statistics.packetsSent++

	d := &net.Dialer{
		Timeout: time.Duration(p.Timeout*1000) * time.Millisecond,
	}

	r := newResault()

	startTime := time.Now()
	conn, err := d.Dial("tcp", p.Dest)
	endTime := time.Now() // Used time
	r.UsedTime = endTime.Sub(startTime)

	if err != nil {
		r.Err = err
		p.statistics.packetsLoss++
	} else {
		p.statistics.packetsRecv++
		defer conn.Close()
	}

	p.resault = *r
}
