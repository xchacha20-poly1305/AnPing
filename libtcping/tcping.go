package libtcping

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Pinger struct {
	Dest    string
	Timeout int // ms
	Count   int

	OnRecv   func()
	OnFail   func()
	OnFinish func()

	resault
	statistics
}

func New(dest string) *Pinger {
	p := &Pinger{
		Dest:    dest,
		Timeout: 1000,
		Count:   -1,
	}
	p.resault = *newResault()
	p.statistics = *newStatistics()

	p.OnRecv = func() {
		fmt.Printf("From %s: time=%v ms\n",
			p.Dest, p.UsedTime.Milliseconds())
	}
	p.OnFail = func() {
		fmt.Printf("Failed to connect %s: %s\n",
			p.Dest, p.Err.Error())
	}

	p.OnFinish = func() {
		fmt.Printf("\n--- %s tcping statistics ---\n", p.Dest)
		fmt.Printf("%v packets transmitted, %v packets received, %v packet loss\n",
			p.GetSent(), p.GetRecv(), p.GetLoss())
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v ms\n",
			p.GetMinRtt(), p.GetAvgRtt(), p.GetMaxRtt(), p.GetStddevRtt())
	}
	return p
}

func (p *Pinger) Run() {
	// p.OnRecvPrinter()
	// fmt.Println(p.Count)
	if p.Count == 0 {
		p.Count = -1
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {

		<-sigChan

		p.OnFinish()

		os.Exit(0)
	}()

	for i := p.Count; i > 0 || i < 0; i-- {
		// fmt.Println(i)
		p.ping()
		if p.Err == nil {
			p.OnRecv()
			p.freshRtt(int(p.resault.UsedTime / time.Millisecond))
		} else {
			p.OnFail()
		}
		time.Sleep(1 * time.Second)
	}

	p.OnFinish()

	// select {}
}
