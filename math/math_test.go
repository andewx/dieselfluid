package math

import (
	"fmt"
	"math"
	"testing"

	"github.com/andewx/dieselfluid/math/matrix"
	"github.com/andewx/dieselfluid/math/vector"
)

//Vector module testing
func TestVecAdd(t *testing.T) {
	var x = vector.Vec{1.0, 1.0, 1.0}
	var y = vector.Vec{1, 1, 1}
	var eq = vector.Vec{2, 2, 2}

	if vector.Eql(vector.Add(x, y), eq) {

	} else {
		t.Errorf("Vector Addition failed %f", x[0])
	}

}

func TestMM(t *testing.T) {

	a := []float32{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}
	b := []float32{4, 5, 6, 7, 3, 4, 5, 6, 2, 3, 4, 5, 1, 2, 3, 4}
	c := []float32{20, 30, 40, 50, 60, 86, 112, 138, 70, 96, 122, 148, 30, 40, 50, 60}
	d := matrix.MulM(a, b)
	if !matrix.MatEqual(d, c) {
		t.Errorf("Matrix Multplication MultMM() failed!!\n")
	} else {
		fmt.Printf("Matrix MulMM() PASS\n")
	}
}

//Vector module testing
func TestVecDot(t *testing.T) {
	var x = vector.Vec{1, 2, 3}
	var y = vector.Vec{1, 1, 1}
	var eq = float32(6.0)

	if vector.Dot(x, y) == eq && vector.Dot(x, y) == eq {

	} else {
		t.Errorf("Vector failed %f", x[0])
	}

}

func TestVector(t *testing.T) {

	fmt.Printf("Executing TestVector(t *Testing.T)-------\n\nDo not use Vec.ProjPlane() \nas this function is suspected to be unusable\n")
	x := vector.Vec{2.0, 2.0, 2.0}
	y := vector.Vec{}

	a := vector.Vec{2, 2, 2}
	b := vector.Vec{0, 0, 0}

	if !vector.Eql(x, a) && !vector.Eql(y, b) {
		t.Error()
	}

	if !vector.Eql(vector.Scale(a, 2.0), vector.Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}
	if !vector.Eql(vector.Add(a, vector.Vec{2.0, 2.0, 2.0}), vector.Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}

	if !vector.Eql(vector.Cross(vector.Vec{-2, -2, -2}, vector.Vec{1, 2, 1}), vector.Vec{2, 0, -2}) {
		r := vector.Cross(vector.Vec{-2, -2, -2}, vector.Vec{1, 2, 1})
		t.Errorf("Cross %f,%f,%f", r[0], r[1], r[2])
	}

	a = vector.Vec{2, 2, 2}

	if vector.Mag(a) != float32(math.Sqrt(12)) {
		t.Errorf("Error Mag")
	}

	if vector.Mag(a) != float32(math.Sqrt(12)) {
		t.Errorf("Error Mag")
	}

	a = vector.Vec{2, 2, 0}
	p := vector.Vec{0, 2, 0}
	r := vector.Proj(a, p)

	if !vector.Eql(r, vector.Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f %f %f\n", r[0], r[1], r[2])
	} else {
		fmt.Printf("Vec.Proj() PASS\n")
	}

	if !vector.Eql(vector.Proj(a, p), vector.Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f, %f, %f\n", a[0], a[1], a[2])
	}

	p = vector.Vec{1, -1, 0}
	o := vector.Vec{0, 1, 0}
	re := vector.Refl(p, o)
	if !vector.Eql(re, vector.Vec{1, 1, 0}) {
		t.Errorf("Error Reflection %f, %f, %f\n", re[0], re[1], re[2])
	}

}

func TestMatrix(t *testing.T) {
	fmt.Printf("Executing TestMatrix(t *Testing.T)-------\n\n")
	A := matrix.Mat{-2, 2, -3, -1, 1, 3, 2, 0, -1}
	B := matrix.Mat3(1.0)
	C := matrix.Mat{3, 2, 0, 1, 4, 0, 1, 2, 3, 0, 2, 1, 9, 2, 3, 1}
	D := matrix.Mat4(1.0)
	a := vector.Vec{0, 0, 0}
	b := vector.Vec{0, 0, 0, 0}

	id3 := matrix.Mat3(1.0)
	id4 := matrix.Mat4(1.0)

	var err error

	//Construct Mat 3
	for i := 0; i < B.Dim(); i++ {
		for j := 0; j < B.Dim(); j++ {
			index := matrix.Map(i, j, B.Dim())
			B[index] = float32(index)
			if i == j {
				id3[index] = 1.0
			}
		}
		a[i] = float32(i)

	}

	//Construct Mat 4
	for i := 0; i < D.Dim(); i++ {
		for j := 0; j < D.Dim(); j++ {
			index := matrix.Map(i, j, D.Dim())
			D[index] = float32(index)

			if i == j {
				id4[index] = 1.0
			}
		}
		b[i] = float32(i)
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	det3 := A.Det()
	det4 := C.Det()
	A.CrossVec(a)
	C.CrossVec(b)

	fmt.Printf("Determinant Mat: %f\n", det3)
	fmt.Printf("Determinant Mat: %f\n", det4)
	//Compute Matrix Inverses
	_ = A.Inv()
	fmt.Printf("Inverse A Matrix 3x3\n")

}

func BenchmarkVecOp(b *testing.B) {

	fmt.Printf("Executing BenchmarkVecOp(t *Testing.T)-------\n\n")
	p := vector.Vec{1, -1, 0}
	o := vector.Vec{0, 1, 0}

	for i := 0; i < b.N; i++ {
		r := vector.Add(p, o)
		vector.Cross(r, p)
		r = vector.Proj(r, o)
		r = vector.Add(r, o)
	}
}
