package geom

import (
	"github.com/andewx/dieselfluid/math/mgl"
	"github.com/andewx/dieselfluid/render/transform"
)

type Ray struct {
	Ray       *mgl.Vec
	Origin    *mgl.Vec
	Transform *transform.Transform
}

type Sphere struct {
	Radius    float32
	Transform *transform.Transform
}

type Intersection struct {
	T []float32
}
