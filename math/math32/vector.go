package math32

//Provides Basic Readable Vector Functions and Allows Wrapping of higher level
//BLAS Library Operations from GONUM for Matrix Computation
//Immutable Vector Operations

import (
	"fmt"
	"math"
)

const EPSILON = 0.0000000019

//Vec standard 3-Length Vectors
type Vec [3]float32
type Vec2 [2]float32
type Vec4 [4]float32

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
func Scl(a Vec, k float32) Vec {
	return Vec{a[0] * k, a[1] * k, a[2] * k}
}

func VecScl(a Vec, b Vec) Vec {
	return Vec{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

//Computes Dot Product A * B
func Dot(a Vec, b Vec) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

//Cross computes A x B Vector operation
func Cross(a Vec, b Vec) Vec {
	g := Vec{a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - b[2]*a[0],
		a[0]*b[1] - b[0]*a[1]}
	return g
}

//Mag returns the magnitude length of a vector
func Mag(v Vec) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

//Norm returns a normalized vector
func Norm(a Vec) Vec {
	v := Vec{}
	l := Mag(a)
	if l != 0.0 {
		v[0] = a[0] / l
		v[1] = a[1] / l
		v[2] = a[2] / l
	}
	return v
}

//Proj projects a vector onto the N Vec a -> N
func Proj(a Vec, n Vec) Vec {
	vn := Norm(n)
	proj := Scl(vn, Dot(a, n)/Mag(n))
	return proj
}

//ProjPlane Projects A onto the Plane of N
func ProjPlane(a Vec, n Vec) Vec {
	pVec := Proj(a, n)
	return Sub(pVec, a)
}

//Refl reflects the v vector about n
func Refl(v Vec, n Vec) Vec {
	b := Scl(n, Dot(v, n)*2.0)
	return Sub(b, v)
}

//Eql Checks if two vectors are equal
func Eql(a Vec, b Vec) bool {
	if b[0] == a[0] && b[1] == a[1] && b[2] == a[2] {
		return true
	}
	return false
}

func isEpsilon(a float32) bool {
	a = float32(math.Abs(float64(a)))
	if a < EPSILON {
		return true
	}
	return false
}

//Dst returns the distance between two vectors B-A/L
func Dist(a Vec, b Vec) float32 {
	return Mag(Sub(a, b))
}

//Tan returns the tangent vector from the given normal
func Tan(a Vec, n Vec) Vec {
	p := Proj(a, n)
	return Sub(a, p)
}

func Print(a Vec) string {
	return fmt.Sprintf("[ %f, %f, %f]\n", a[0], a[1], a[2])
}

//---------------INCLUDE BLAS64 Extended Functions and Wrappers here----------

/*
func Blas64(a Vec) blas64.Vector {
	newVector := blas64.Vector{}
	newVector.Data = make([]float32, 3)
	newVector.Data[0] = a[0]
	newVector.Data[1] = a[1]
	newVector.Data[2] = a[2]
	newVector.Inc = 0
	return newVector
}
*/
