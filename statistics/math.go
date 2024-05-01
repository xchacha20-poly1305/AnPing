package statistics

import (
	"fmt"

	"github.com/sagernet/sing/common/x/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func diffAbs[T Number](a, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}

func percent[T Number](dividend, divisor T) string {
	switch divisor {
	case 0:
		// 0 cannot be used as a divisor.
		return "0.00%"
	}
	result := fmt.Sprintf("%.2f", float64(dividend)/float64(divisor)*100)
	/*
		switch result {
		case "+Inf":
			result = "0.00"
		}
	*/

	return result + "%"
}
