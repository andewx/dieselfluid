package geom

import (
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/render/transform"
)

type Ray struct {
	Ray       *vector.Vec
	Origin    *vector.Vec
	Transform *transform.Transform
}

type Sphere struct {
	Radius    float32
	Transform *transform.Transform
}

type Intersection struct {
	T []float32
}
