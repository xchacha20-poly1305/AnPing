//go:build !unix

package main

import (
	"github.com/xchacha20-poly1305/anping/tcpping"
)

const (
	DefaultType = tcpping.Type
)
