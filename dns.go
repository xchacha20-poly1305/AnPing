package anping

import (
	"math/rand/v2"
	"net"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
)

// DomainStrategy is the behavior of dual stack.
type DomainStrategy uint8

const (
	PreferNone DomainStrategy = iota
	PreferIpv6
	PreferIpv4
)

func LookupSingleIP(address M.Socksaddr, prefer DomainStrategy) (ip net.IP, err error) {
	ips, err := net.LookupIP(address.Fqdn)
	if err != nil {
		return nil, E.Cause(err, "look up ip")
	}

	switch prefer {
	case PreferNone:
		return ips[rand.IntN(len(ips))], nil
	case PreferIpv6:
		ipv6s := common.Filter(ips, func(it net.IP) bool {
			return it.To4() == nil
		})

		if len(ipv6s) == 0 {
			return nil, E.New("not found IPv6 address of: ", address)
		}

		return ipv6s[rand.IntN(len(ipv6s))], nil
	case PreferIpv4:
		ipv4s := common.Filter(ips, func(it net.IP) bool {
			return it.To4() != nil
		})

		if len(ipv4s) == 0 {
			return nil, E.New("not found IPv4 address of: ", address.Fqdn)
		}

		return ipv4s[rand.IntN(len(ipv4s))], nil
	default:
		return nil, E.New("invalid domain strategy: ", prefer)
	}
}
