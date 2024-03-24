//go:build !without_icmp

package include

import (
	"github.com/xchacha20-poly1305/anping/icmpping"
)

const DefaultProtocol = icmpping.Protocol
