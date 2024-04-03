//go:build !without_icmp && linux

package include

import (
	"github.com/xchacha20-poly1305/anping/implement/icmpping"
)

const DefaultProtocol = icmpping.Protocol
