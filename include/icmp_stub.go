//go:build without_icmp || !linux

package include

import (
	"github.com/xchacha20-poly1305/anping/tcpping"
)

const DefaultProtocol = tcpping.Protocol
