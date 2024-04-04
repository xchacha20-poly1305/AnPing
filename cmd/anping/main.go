package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/include"
)

var version = "Unknown"

var (
	showVersion bool

	interval time.Duration
	quite    bool
	count    int
	timeout  int

	prefer6        bool
	prefer4        bool
	domainStrategy = anping.PreferNone
)

func init() {
	flag.BoolVar(&showVersion, "V", false, "Show AnPing version")

	flag.DurationVar(&interval, "i", anping.Interval, "Ping interval")
	flag.BoolVar(&quite, "q", false, "Quite mode")
	flag.IntVar(&count, "c", anping.Count, "Ping count")
	flag.IntVar(&timeout, "W", anping.Timeout, "Ping timeout")

	flag.BoolVar(&prefer6, "6", false, "Prefer to IPv6")
	flag.BoolVar(&prefer4, "4", false, "Prefer to IPv4")

	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
		return
	}

	if prefer4 {
		domainStrategy = anping.PreferIpv4
	}
	// Final prefer IPv6.
	if prefer6 {
		domainStrategy = anping.PreferIpv6
	}
}

func main() {
	args := flag.Args()
	switch len(args) {
	case 0:
		log.Fatalln("Didn't set any args!")
	case 1:
		switch strings.ToLower(args[0]) {
		case "v", "version":
			printVersion()
			return
		}
		args = append(args, include.DefaultProtocol)
		// turn to: "icmp <address>"
		slices.Reverse(args)
	}

	creator, ok := anping.AnPingerCreator[args[0]]
	if !ok {
		log.Fatalf("Not found protocol: %s\n", args[0])
	}

	pinger := creator(os.Stdout)
	err := pinger.SetAddress(args[1])
	if err != nil {
		log.Fatalln(err)
	}

	opts := pinger.Options()
	opts.Count = count
	opts.Timeout = int32(timeout)
	opts.Interval = interval
	opts.Quite = quite
	opts.DomainStrategy = domainStrategy

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			cancel()
			close(osSignals)
		}()
		pinger.RunContext(ctx)
	}()
	<-osSignals
	cancel()

	// Wait for finish print.
	time.Sleep(interval)
}

func printVersion() {
	fmt.Printf("AnPing: %s\n", version)
}
