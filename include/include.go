// Package include is used to register ping protocol.
//
// Default protocol is ICMP. But besides Linux,
// other system use TCP as default protocol and not include ICMP.
package include

import (
	_ "github.com/xchacha20-poly1305/anping/implement/icmpping"
	_ "github.com/xchacha20-poly1305/anping/implement/udpping"
)
