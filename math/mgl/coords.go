package mgl

import (
	"fmt"
	"math"
)

//(radial, azimuth, polar angle)
type Polar [3]float32

const (
	PI      = 3.14159265359
	DEG2RAD = 0.01745329252
	RAD2DEG = 57.2958
)

//Converts vector into sperical coordinates
func Vec2Sphere(x Vec) (Polar, error) {

	r := x.Mag()
	az := float32(math.Atan2(float64(x[1]), float64(x[0])))
	incl := float32(math.Acos(float64(x[2] / r)))
	sph := Polar{r, az, incl}
	if math.IsNaN(float64(az)) {
		return sph, fmt.Errorf("Invalid Operation Vec2Sphere x coord = 0 arctan is NaN")
	}
	return sph, nil
}

func Sphere2Vec(x Polar) (Vec, error) {
	if math.IsNaN(float64(x[1])) {
		return Vec{}, fmt.Errorf("Az Is NaN")
	}
	var a, b, c float64
	a = float64(x[0])
	b = float64(x[1])
	c = float64(x[2])
	x0 := float32(a * math.Sin(c) * math.Cos(b))
	x1 := float32(a * math.Sin(c) * math.Sin(b))
	x2 := float32(a * math.Cos(c))
	return Vec{x0, x1, x2}, nil

}

func (s Polar) Radius() float32 {
	return s[0]
}

func (s Polar) Azimuth() float32 {
	return s[1]
}

func (s Polar) Polar() float32 {
	return s[2]
}

func (s Polar) Add(b Polar) Polar {
	s[1] = s[1] + b[1]
	s[2] = s[2] + b[2]
	return s
}

func (s Polar) AddAzimuthDegrees(b float32) Polar {
	c := b * DEG2RAD
	s[1] = s[1] + c
	return s
}

func (s Polar) AddAzimuth(b float32) Polar {
	s[1] = s[1] + b
	return s
}

func (s Polar) AddPolarDegrees(b float32) Polar {
	c := b * DEG2RAD
	s[2] = s[2] + c
	return s
}

func (s Polar) AddPolar(b float32) Polar {
	s[2] = s[2] + b
	return s
}

func (s Polar) Copy() Polar {
	return Polar{s[0], s[1], s[2]}
}
