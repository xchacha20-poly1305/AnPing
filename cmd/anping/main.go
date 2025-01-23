package main

import (
	"context"
	crand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/control"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/xchacha20-poly1305/anping"
	"github.com/xchacha20-poly1305/anping/icmpping"
	"github.com/xchacha20-poly1305/anping/tcpping"
	"github.com/xchacha20-poly1305/anping/udpping"
)

var (
	count              int64
	interfaceName      string
	interval           time.Duration
	quite              bool
	timeout            time.Duration
	justIPv4, justIPv6 bool

	payloadLength  int
	payloadContent string

	showVersion bool
)

func init() {
	flag.Int64Var(&count, "c", 0, "stop after <count> replies")
	flag.StringVar(&interfaceName, "I", "", "bind interface name")
	flag.DurationVar(&interval, "i", time.Second, "duration between sending each packet")
	flag.BoolVar(&quite, "q", false, "quiet mode")
	flag.DurationVar(&timeout, "W", 500*time.Millisecond, "time to wait for response")
	flag.BoolVar(&justIPv4, "4", false, "use IPv4")
	flag.BoolVar(&justIPv6, "6", false, "use IPv4")
	flag.BoolVar(&showVersion, "V", false, "print version and exit")

	flag.IntVar(&payloadLength, "pl", 40, "[UDP/ICMP] payload length")
	flag.StringVar(&payloadContent, "p", "", "[UDP/ICMP] contents of padding byte")

	flag.Parse()

	if showVersion {
		printVersion()
	}
	if justIPv4 && justIPv6 {
		fatal("Refer IPv4 & IPv6 at the same time!\n")
	}
}

func main() {
	args := flag.Args()
	switch len(args) {
	case 0:
		fatal("Not provided any address!\n")
	case 1:
		switch strings.ToLower(args[0]) {
		case "v", "version":
			printVersion()
		case "h", "help":
			flag.Usage()
			os.Exit(0)
		}
		args = append([]string{DefaultType}, args...)
	case 2:
	default:
		fatal("To many args: %v\n", args)
	}

	switch args[0] {
	case icmpping.Type:
		if runtime.GOOS == "windows" {
			fatal("Not support ICMP ping on Windows")
		}
	case tcpping.Type, udpping.Type:
	default:
		fatal("Unknown type: %s\n", args[0])
	}

	addr := M.ParseSocksaddr(args[1])
	if !addr.IsValid() {
		fatal("Invalid address: %s", args[1])
	}
	if justIPv4 && addr.IsIPv6() {
		fatal("Provided IPv6 address but use IPv4\n")
	}
	if justIPv6 && addr.IsIPv4() {
		fatal("Provided IPv4 address but use IPv6\n")
	}
	if addr.IsFqdn() {
		const resolveTimeout = 3 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
		var network string
		if justIPv4 {
			network = N.NetworkIP + "4"
		}
		if justIPv6 {
			network = N.NetworkIP + "6"
		}
		ips, err := net.DefaultResolver.LookupIP(ctx, network, addr.Fqdn)
		cancel()
		if err != nil {
			fatal("Look up IP for %s: $v\n", addr.Fqdn, err)
		}
		if justIPv4 {
			ips = common.Filter(ips, func(it net.IP) bool {
				return it.To16() == nil
			})
		}
		if justIPv6 {
			ips = common.Filter(ips, func(it net.IP) bool {
				return it.To16() != nil
			})
		}
		if len(ips) == 0 {
			fatal("Not found any IP of %s\n", addr.Fqdn)
		}
		addr.Fqdn = ""
		addr.Addr = M.AddrFromIP(ips[rand.IntN(len(ips))])
	}

	var logger anping.Logger
	if quite {
		logger = anping.NewLogger("", "", io.Discard)
	} else {
		var target string
		switch args[0] {
		case icmpping.Type:
			target = addr.AddrString()
		default:
			target = addr.String()
		}
		logger = anping.NewLogger(args[0], target, os.Stdout)
	}
	commonSettings := anping.NewSettings(control.NewDefaultInterfaceFinder())
	commonSettings.SetLogger(logger)
	commonSettings.SetTarget(addr)
	commonSettings.SetInterval(interval)
	commonSettings.SetTimeout(timeout)
	commonSettings.SetCount(count)
	commonSettings.SetInterface(interfaceName)

	var instance anping.AnPing
	switch args[0] {
	case icmpping.Type:
		settings := icmpping.NewSettings(commonSettings, payloadFunc())
		instance = icmpping.New(settings)
	case tcpping.Type:
		settings := tcpping.NewSettings(commonSettings)
		instance = tcpping.New(settings)
	case udpping.Type:
		settings := udpping.NewSettings(commonSettings, payloadFunc())
		instance = udpping.New(settings)
	default:
		panic("Unexpected type: " + args[0])
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	err := instance.Start(ctx)
	if err != nil {
		fatal("Run instance: %v", err)
	}
	select {
	case <-ctx.Done():
		<-instance.Done()
		return
	case <-instance.Done():
		cancel()
	}
}

func fatal(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func payloadFunc() func() []byte {
	if payloadContent != "" && payloadLength > 0 {
		fatal("Set payload content and length")
	}
	if payloadLength > 0 {
		payload := make([]byte, payloadLength)
		_, _ = crand.Read(payload)
		return func() []byte {
			return payload
		}
	}
	if payloadContent != "" {
		payload := []byte(payloadContent)
		return func() []byte {
			return payload
		}
	}
	// Use default
	return nil
}
