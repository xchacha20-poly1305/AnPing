package anping

import (
	"fmt"
	"log"
	"runtime"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func (a *AnPinger) Icmpping() {
	pinger, err := probing.NewPinger(a.Addr.Hostname())
	if err != nil {
		log.Println(err)
		return
	}

	pinger.Timeout = time.Duration(a.Timeout*1000) * time.Millisecond
	pinger.Count = a.Count
	// log.Println(pinger.Count)
	// pinger.Count = 3

	switch runtime.GOOS {
	case "windows":
		pinger.SetPrivileged(true)
	default:
		err = getPermission()
		if err != nil {
			log.Println(err)
			log.Println("Because of permission is not enough, so now is UDP ping!")
		}
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v ms\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	err = pinger.Run()
	if err != nil {
		log.Println(err)
	}
}
