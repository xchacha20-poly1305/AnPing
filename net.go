package anping

import (
	"math/rand/v2"
	"net"

	E "github.com/sagernet/sing/common/exceptions"
)

// DomainStrategy is the behavior of dual stack.
type DomainStrategy uint8

const (
	PreferNone DomainStrategy = iota
	PreferIpv6
	PreferIpv4
)

func LookupSingleIP(domain string, prefer DomainStrategy) (ip net.IP, err error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, E.Cause(err, "look up ip")
	}

	switch prefer {
	case PreferNone:
		ip = ips[rand.IntN(len(ips))]
	case PreferIpv6:
		var ipv6s []net.IP

		for _, singleIP := range ips {
			if ip.To16() != nil {
				ipv6s = append(ipv6s, singleIP)
			}
		}

		l := len(ipv6s)
		if l == 0 {
			return nil, E.New("not found IPv6 address of ", domain)
		}

		ip = ipv6s[rand.IntN(l)]
	case PreferIpv4:
		var ipv4s []net.IP

		for _, singleIP := range ips {
			if ip.To4() != nil {
				ipv4s = append(ipv4s, singleIP)
			}
		}

		l := len(ipv4s)
		if l == 0 {
			return nil, E.New("not found IPv4 address of ", domain)
		}

		ip = ipv4s[rand.IntN(l)]
	}

	return
}
