// Package include is used to register ping protocol.
// Default protocol is ICMP.
//
// Tags:
//   - `without_icmp`: Not include ICMP. On this time, default protocol is TCP.
package include

import (
	_ "github.com/xchacha20-poly1305/anping/tcpping"
)
