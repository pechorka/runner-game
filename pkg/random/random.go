package random

import "math/rand"

func Float32(from, to float32) float32 {
	return from + (to-from)*rand.Float32()
}
