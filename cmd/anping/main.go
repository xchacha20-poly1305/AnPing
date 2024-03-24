package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/xchacha20-poly1305/anping"
	_ "github.com/xchacha20-poly1305/anping/include"
)

var version = "Unknown"

var (
	showVersion bool

	interval time.Duration
	quite    bool
	number   int
	timeout  int
)

func init() {
	flag.BoolVar(&showVersion, "V", false, "Show AnPing version")

	flag.DurationVar(&interval, "i", anping.Interval, "Ping interval")
	flag.BoolVar(&quite, "q", false, "Quite mode")
	flag.IntVar(&number, "c", anping.Number, "Ping count")
	flag.IntVar(&timeout, "W", anping.Timeout, "Ping timeout")

	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
		return
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
		args = append(args, "icmp")
		// turn to: "icmp <address>"
		slices.Reverse(args)
	}

	creator, ok := anping.AnPingerCreator[args[0]]
	if !ok {
		log.Fatalf("Not found protocol: %s\n", args[0])
	}

	var writer io.Writer
	if !quite {
		writer = os.Stdout
	}
	pinger := creator(writer)
	err := pinger.SetAddress(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	pinger.SetNumber(number)
	pinger.SetTimeout(int32(timeout))
	pinger.SetInterval(interval)

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
	_ = pinger.Clean()
}

func printVersion() {
	fmt.Printf("AnPing: %s\n", version)
}
