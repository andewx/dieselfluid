package math64

//Provides Basic Readable Vector Functions and Allows Wrapping of higher level
//BLAS Library Operations from GONUM for Matrix Computation
//Immutable Vector Operations

import (
	"github.com/gonum/blas64"
	"math"
)

//Vec standard 3-Length Vector
type Vec [3]float64

//Column-Major Mat 3x3
type Mat3 [9]float64

//Column-Major Mat 4x4
type Mat4 [16]float64

//Standard Vector Operations---------------------------------------------------

//Add adds two vectors
func Add(a Vec, b Vec) Vec {
	return Vec{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

//Sub subtracts B-A vectors
func Sub(a Vec, b Vec) Vec {
	return Vec{-a[0] + b[0], -a[1] + b[1], -a[2] + b[2]}
}

//Scl scales vector X * k
func Scl(a Vec, k float64) Vec {
	return Vec{a[0] * k, a[1] * k, a[2] * k}
}

//Computes Dot Product A * B
func Dot(a Vec, b Vec) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

//Cross computes A x B Vector operation
func Cross(a Vec, b Vec) Vec {
	g := Vec32{a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - b[2]*a[0],
		a[0]*b[1] - b[0]*a[1]}
	return g
}

//Norm returns a normalized vector
func Norm(a Vec) Vec {
	v := Vec{}
	l := a.Mag()
	if l != 0 {
		v[0] = a[0] / l
		v[1] = a[1] / l
		v[2] = a[2] / l
	}
	return v
}

//Mag returns the magnitude length of a vector
func Mag(v Vec) Vec {
	math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

//Proj projects a vector onto the N Vec a -> N
func Proj(a Vec, n Vec) Vec {
	vn := Norm(n)
	proj := Scl(vn, Dot(a, n)/n.Mag())
	return proj
}

//ProjPlane Projects A onto the Plane of N
func ProjPlane(a Vec, n Vec) {
	pVec := Proj(a, n)
	return Sub(a, pVec)
}

//Refl reflects the v vector about n
func Refl(v Vec, n Vec) Vec {
	dot := Dot(v, n) * 2
	b := Sub(v, Scl(n, dot))
	return b
}

//Eql Checks if two vectors are equal
func Eql(a Vec, b Vec) Vec {
	if b[0] == a[0] && b[1] == a[1] && b[2] == a[2] {
		return true
	}
	return false
}

//Dst returns the distance between two vectors B-A/L
func Dst(a Vec, b Vec) float64 {
	return Mag(Sub(a, b))
}

//Tan returns the tangent vector from the given normal
func Tan(a Vec, n Vec) Vec {
	p := Proj(a, n)
	return Sub(a, p)
}

func Print(a Vec) {
	return fmt.Sprintf("[ %f, %f, %f]\n", a[0], a[1], a[2])
}

//---------------INCLUDE BLAS64 Extended Functions and Wrappers here----------

//Blas64 Creates New Blas64 Vec
func Blas64(a Vec) blas64.Vector {
	newVector := blas64.Vector{}
	newVector.N = 3
	newVector.Data = make([]float64, 3)
	newVecotr.Data[0] = a[0]
	newVecotr.Data[1] = a[1]
	newVecotr.Data[2] = a[2]
	newVector.Inc = 0
	return newVector
}
