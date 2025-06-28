package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
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
	flood              bool
	justIPv4, justIPv6 bool

	payloadLength  int
	payloadContent string

	showVersion bool
)

func init() {
	flag.Int64Var(&count, "c", 0, "Stop after <count> replies. If it <= 0, it will loop forever.")
	flag.StringVar(&interfaceName, "I", "", "Bind interface name.")
	flag.DurationVar(&interval, "i", time.Second, "Duration between sending each packet.")
	flag.BoolVar(&quite, "q", false, "Enabled quiet mode.")
	flag.DurationVar(&timeout, "W", 500*time.Millisecond, "Time to wait for response.")
	flag.BoolVar(&justIPv4, "4", false, "Use IPv4.")
	flag.BoolVar(&justIPv6, "6", false, "Use IPv4.")
	flag.BoolVar(&showVersion, "V", false, "Print version and exit.")
	flag.BoolVar(&flood, "f", false, "Flood ping.")

	flag.IntVar(&payloadLength, "pl", 40, "[UDP/ICMP] Payload length.")
	flag.StringVar(&payloadContent, "p", "", "[UDP/ICMP] Contents of padding byte.")

	flag.Parse()

	if showVersion {
		printVersion()
	}
	if justIPv4 && justIPv6 {
		fatalf("Refer IPv4 & IPv6 at the same time!\n")
	}
	if flood {
		const floodMinimumInterval = 2 * time.Millisecond
		if interval > floodMinimumInterval {
			fatalf("Provided interval(%s) is greater than 2 millisecond, try to use -i 2ms.\n", interval)
		}
		quite = true
	}
}

func main() {
	args := flag.Args()
	switch len(args) {
	case 0:
		fatalf("Not provided any address!\n")
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
		fatalf("To many args: %v\n", args)
	}

	switch args[0] {
	case icmpping.Type:
		if runtime.GOOS == "windows" {
			fatalf("Not support ICMP ping on Windows")
		}
	case tcpping.Type, udpping.Type:
	default:
		fatalf("Unknown type: %s\n", args[0])
	}

	addr := M.ParseSocksaddr(args[1])
	if !addr.IsValid() {
		fatalf("Invalid address: %s", args[1])
	}
	if justIPv4 && addr.IsIPv6() {
		fatalf("Provided IPv6 address but use IPv4\n")
	}
	if justIPv6 && addr.IsIPv4() {
		fatalf("Provided IPv4 address but use IPv6\n")
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
			fatalf("Look up IP for %s: %v\n", addr.Fqdn, err)
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
			fatalf("Not found any IP of %s\n", addr.Fqdn)
		}
		addr.Fqdn = ""
		addr.Addr = M.AddrFromIP(ips[0])
	}

	var target string
	switch args[0] {
	case icmpping.Type:
		target = addr.AddrString()
	default:
		target = addr.String()
	}
	logger := anping.NewLogger(args[0], target, os.Stdout, quite)
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
		fatalf("Run instance: %v", err)
	}
	select {
	case <-ctx.Done():
		<-instance.Done()
		return
	case <-instance.Done():
		cancel()
	}
}

func fatalf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func payloadFunc() func() []byte {
	if payloadContent != "" && payloadLength > 0 {
		fatalf("Set payload content and length in the same time")
	}
	if payloadLength > 0 {
		payload := make([]byte, payloadLength)
		_, _ = rand.Read(payload)
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
