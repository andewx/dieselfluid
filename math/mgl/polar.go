package mgl

import (
	"fmt"
	"math"
)

//(radial, azimuth, polar angle)
type Polar struct {
	Sphere Vec
	Origin Vec
}

const (
	PI      = 3.14159265359
	DEG2RAD = 0.01745329252
	RAD2DEG = 57.2958
)

func NewPolar(rad float32) Polar {
	return Polar{Vec{rad, 0, 0}, Vec{0, 0, 0}}
}

//Converts Vector to Spherical Coordinates Atan2 errors with az parameter on 0
func Vec2Sphere(x Vec) (Polar, error) {
	var err error
	r := x.Mag()
	az := float32(math.Atan2(float64(x[1]), float64(x[0])))
	incl := float32(math.Acos(float64(x[2] / r)))

	if math.IsNaN(float64(az)) {
		x[0] = 0.1
		az = float32(math.Atan2(float64(x[1]), float64(x[0])))
		err = fmt.Errorf("Atan2 failed -- vec(x) set to 0.0001")
	}

	sph := Polar{Vec{r, az, incl}, Vec{0, 0, 0}}
	return sph, err
}

//Spherical Coordinate to Vector with adjusted Azimuthal Vectors
func Sphere2Vec(x Polar) (Vec, error) {
	var err error
	var a, b, c float64
	a = float64(x.Sphere[0])
	b = float64(x.Sphere[1])
	c = float64(x.Sphere[2])

	if math.IsNaN(float64(x.Sphere[1])) {
		b = 0.1
		err = fmt.Errorf("Azimuth is NaN --Setting Azimuth to 0.0001")
	}

	if b == 0 {
		b = 0.1
	}

	x0 := float32(a * math.Sin(c) * math.Cos(b))
	x1 := float32(a * math.Sin(c) * math.Sin(b))
	x2 := float32(a * math.Cos(c))
	return Vec{x0, x1, x2}, err

}

func (s Polar) Radius() float32 {
	return s.Sphere[0]
}

func (s Polar) Azimuth() float32 {
	return s.Sphere[1]
}

func (s Polar) Polar() float32 {
	return s.Sphere[2]
}

func (s Polar) Add(b Polar) Polar {
	s.Sphere[1] = s.Sphere[1] + b.Sphere[1]
	s.Sphere[2] = s.Sphere[2] + b.Sphere[2]
	return s
}

func (s Polar) AddAzimuthDegrees(b float32) Polar {
	c := b * DEG2RAD
	s.Sphere[1] = s.Sphere[1] + c
	return s
}

func (s Polar) AddAzimuth(b float32) Polar {
	s.Sphere[1] = s.Sphere[1] + b
	return s
}

func (s Polar) AddPolarDegrees(b float32) Polar {
	c := b * DEG2RAD
	s.Sphere[2] = s.Sphere[2] + c
	return s
}

func (s Polar) AddPolar(b float32) Polar {
	s.Sphere[2] = s.Sphere[2] + b
	return s
}

func (s Polar) Copy() Polar {
	return Polar{s.Sphere, s.Origin}
}
