package mgl

import (
	"fmt"
	"math"
	"testing"
)

//Vector module testing
func TestVecAdd(t *testing.T) {
	var x = Vec{1.0, 1.0, 1.0}
	var y = Vec{1, 1, 1}
	var eq = Vec{2, 2, 2}

	if Eql(Add(x, y), eq) {

	} else {
		t.Errorf("Vector Addition failed %f", x[0])
	}

}

func TestMM(t *testing.T) {

	a := []float32{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}
	b := []float32{4, 5, 6, 7, 3, 4, 5, 6, 2, 3, 4, 5, 1, 2, 3, 4}
	c := []float32{20, 30, 40, 50, 60, 86, 112, 138, 70, 96, 122, 148, 30, 40, 50, 60}
	d := MulM(a, b)
	if !MatEqual(d, c) {
		t.Errorf("Matrix Multplication MultMM() failed!!\n")
	} else {
		fmt.Printf("Matrix MulMM() PASS\n")
	}
}

//Vector module testing
func TestVecDot(t *testing.T) {
	var x = Vec{1, 2, 3}
	var y = Vec{1, 1, 1}
	var eq = float32(6.0)

	if Dot(x, y) == eq && Dot(x, y) == eq {

	} else {
		t.Errorf("Vector failed %f", x[0])
	}

}

func TestVector(t *testing.T) {

	fmt.Printf("Executing TestVector(t *Testing.T)-------\n\nDo not use Vec.ProjPlane() \nas this function is suspected to be unusable\n")
	x := Vec{2.0, 2.0, 2.0}
	y := Vec{}

	a := Vec{2, 2, 2}
	b := Vec{0, 0, 0}

	if !Eql(x, a) && !Eql(y, b) {
		t.Error()
	}

	if !Eql(Scale(a, 2.0), Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}
	if !Eql(Add(a, Vec{2.0, 2.0, 2.0}), Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}

	if !Eql(Cross(Vec{-2, -2, -2}, Vec{1, 2, 1}), Vec{2, 0, -2}) {
		r := Cross(Vec{-2, -2, -2}, Vec{1, 2, 1})
		t.Errorf("Cross %f,%f,%f", r[0], r[1], r[2])
	}

	a = Vec{2, 2, 2}

	if Mag(a) != float32(math.Sqrt(12)) {
		t.Errorf("Error Mag")
	}

	if Mag(a) != float32(math.Sqrt(12)) {
		t.Errorf("Error Mag")
	}

	a = Vec{2, 2, 0}
	p := Vec{0, 2, 0}
	r := Proj(a, p)
	//Plane projection should take 1 input vector and 2 plane param vectors
	//This needs to be looked at
	//	h := ProjPlane(a, p)

	if !Eql(r, Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f %f %f\n", r[0], r[1], r[2])
	} else {
		fmt.Printf("Vec.Proj() PASS\n")
	}

	/*
		if !Eql(h, Vec{2, 0, 0}) {
			t.Errorf("Error Proj Plane  %f %f %f", h[0], h[1], h[2])
		}
	*/

	if !Eql(Proj(a, p), Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f, %f, %f\n", a[0], a[1], a[2])
	}

	p = Vec{1, -1, 0}
	o := Vec{0, 1, 0}
	re := Refl(p, o)
	if !Eql(re, Vec{1, 1, 0}) {
		t.Errorf("Error Reflection %f, %f, %f\n", re[0], re[1], re[2])
	}

}

func TestMatrix(t *testing.T) {
	fmt.Printf("Executing TestMatrix(t *Testing.T)-------\n\n")
	A := Mat{-2, 2, -3, -1, 1, 3, 2, 0, -1}
	B := Mat3(1.0)
	C := Mat{3, 2, 0, 1, 4, 0, 1, 2, 3, 0, 2, 1, 9, 2, 3, 1}
	D := Mat4(1.0)
	a := Vec{0, 0, 0}
	b := Vec{0, 0, 0, 0}

	id3 := Mat3(1.0)
	id4 := Mat4(1.0)

	var err error

	//Construct Mat 3
	for i := 0; i < B.Dim(); i++ {
		for j := 0; j < B.Dim(); j++ {
			index := Map(i, j, B.Dim())
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
			index := Map(i, j, D.Dim())
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
	p := Vec{1, -1, 0}
	o := Vec{0, 1, 0}

	for i := 0; i < b.N; i++ {
		r := Add(p, o)
		Cross(r, p)
		r = Proj(r, o)
		r = Add(r, o)
	}
}
