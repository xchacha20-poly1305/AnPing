package anping

import (
	"fmt"
	
	"github.com/xchacha20-poly1305/AnPing/libtcping"
)

func (a *AnPinger) Tcpping() {
	pinger := libtcping.New(a.Addr.Host)
	pinger.Timeout = a.Timeout
	pinger.Count = a.Count
	
	fmt.Printf("TCPING %s:\n", a.Addr.Host)
	
	pinger.Run()
}
