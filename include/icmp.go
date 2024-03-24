//go:build !without_icmp

package include

import (
	_ "github.com/xchacha20-poly1305/anping/icmpping"
)
