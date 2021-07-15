package random

import "math/rand"

func Get() float32 {
	Seed()
	return rand.Float32()
}
