package math

import (
	"dslfluid.com/dsl/math/math32"
	"dslfluid.com/dsl/math/math64"
)

func Vec64(a math32.Vec) math64.Vec {
	v64 := math64.Vec{float64(a[0]), float64(a[1]), float64(a[2])}
	return v64
}

func Vec32(a math64.Vec) math32.Vec {
	v32 := math32.Vec{float32(a[0]), float32(a[1]), float32(a[2])}
	return v32
}
