//go:build !linux

package include

import (
	"github.com/xchacha20-poly1305/anping/implement/tcpping"
)

const DefaultProtocol = tcpping.Protocol
