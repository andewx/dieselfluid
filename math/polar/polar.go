package polar

import (
	"math"

	"github.com/andewx/dieselfluid/math/common"
	"github.com/andewx/dieselfluid/math/vector"
)

//(radial, azimuth, polar angle)
type Polar struct {
	Sphere vector.Vec
	Origin vector.Vec
}

func NewPolar(rad float32) Polar {
	return Polar{vector.Vec{rad, 0, 0}, vector.Vec{0, 0, 0}}
}

//Converts Vector to Spherical Coordinates Atan2 errors with az parameter on 0
func Vec2Sphere(x vector.Vec) (Polar, error) {
	var err error
	r := x.Mag()
	az := PolarAzimuth(x)
	incl := float32(math.Acos(float64(x[2] / r)))

	sph := Polar{vector.Vec{r, float32(az), incl}, vector.Vec{0, 0, 0}}
	return sph, err
}

//Spherical Coordinate to Vector with adjusted Azimuthal Vectors
func Sphere2Vec(x Polar) (vector.Vec, error) {
	var err error
	var a, b, c float64
	a = float64(x.Sphere[0])
	b = float64(x.Sphere[1])
	c = float64(x.Sphere[2])

	x0 := float32(a * math.Sin(c) * math.Cos(b))
	x1 := float32(a * math.Sin(c) * math.Sin(b))
	x2 := float32(a * math.Cos(c))
	return vector.Vec{x0, x1, x2}, err

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

//Add azimuth in degrees
func (s Polar) AddAzimuthDegrees(b float32) Polar {
	c := b * common.DEG2RAD
	s.Sphere[1] = s.Sphere[1] + c
	return s
}

//Add Azimuth Angle
func (s Polar) AddAzimuth(b float32) Polar {
	s.Sphere[1] = s.Sphere[1] + b
	return s
}

//Add inclination angle in degrees
func (s Polar) AddPolarDegrees(b float32) Polar {
	c := b * common.DEG2RAD
	s.Sphere[2] = s.Sphere[2] + c
	return s
}

//Add inclination angle
func (s Polar) AddPolar(b float32) Polar {
	s.Sphere[2] = s.Sphere[2] + b
	return s
}

func (s Polar) Copy() Polar {
	return Polar{s.Sphere, s.Origin}
}

//-------------------------------RAY SPHERE --------------------//

type Intersection struct {
	T float32
}

func Priority(t []*Intersection) float32 {
	if len(t) == 0 {
		return 0.0
	}
	min := t[0].T
	for i := 1; i < len(t); i++ {
		next := t[i].T
		abs_next := next * next
		abs_min := min * min
		if abs_next < abs_min {
			min = next
		}
	}
	return min
}

//Calculates azimuth from vector
func PolarAzimuth(vec vector.Vec) float64 {
	x := vec[0]
	y := vec[1]
	res := math.Atan2(float64(x), float64(y))
	return res
}

func RaySphereIntersect(r vector.Vec, o vector.Vec, s Polar) []*Intersection {

	// The vector from the s origin to the r origin.
	sphereToRayVec := vector.Sub(o, s.Origin)

	// Compute the discriminant to tell whether the r intersects with the s at all.
	a := vector.Dot(r, r)
	b := 2 * vector.Dot(r, sphereToRayVec)
	c := vector.Dot(sphereToRayVec, sphereToRayVec) - s.Radius()*s.Radius()
	discriminant := math.Pow(float64(b), 2) - float64(4*a*c)

	// If the discriminant is negative, then the r misses the s and no intersections occur.
	if discriminant < 0 {
		return []*Intersection{}
	}

	// Compute the t values.
	t1 := ((-1 * b) - float32(math.Sqrt(discriminant))) / (2 * a)
	t2 := ((-1 * b) + float32(math.Sqrt(discriminant))) / (2 * a)

	// Return the intersection t values and object in increasing order
	return []*Intersection{{T: t1}, {T: t2}}
}
