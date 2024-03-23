package anping

import (
	"github.com/sagernet/sing/common/x/constraints"
)

func diffAbs[T constraints.Integer | constraints.Float](a, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}
