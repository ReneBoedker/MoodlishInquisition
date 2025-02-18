// Package unif provides functions to simplify generation of random values.
package unif

import (
	"fmt"
	"math/rand"
)

// IntInInterval generates a uniformly random integer in [a,b].
// If a>b, the function panics
func IntInInterval(a, b int, allowZero bool) int {
	if a > b {
		panic(fmt.Errorf("Cannot generate integer in empty interval [%d, %d]", a, b))
	}
	for {
		tmp := a + rand.Intn(b-a+1)
		if allowZero || tmp != 0 {
			return tmp
		}
	}
}

// BoundedInt generates a uniformly random integer in [-a,a].
func BoundedInt(a int, allowZero bool) int {
	if a < 0 {
		return BoundedInt(-a, allowZero)
	}
	return IntInInterval(-a, a, allowZero)
}

// Bool generates a uniformly random boolean value.
func Bool() bool {
	return rand.Uint32()&1 == 1
}
