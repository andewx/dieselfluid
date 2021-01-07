package math32

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
	x := Vec{2.0, 2.0, 2.0}
	y := Vec{}

	a := Vec{2, 2, 2}
	b := Vec{0, 0, 0}

	if !Eql(x, a) && !Eql(y, b) {
		t.Error()
	}

	if !Eql(Scl(a, 2.0), Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}
	if !Eql(Add(a, Vec{2.0, 2.0, 2.0}), Vec{4.0, 4.0, 4.0}) {
		t.Error()
	}

	if !isEpsilon(Mag(Norm(x)) - 1.0) {
		t.Errorf("Normalized vector error: A Mag(): %f, %f, %f", x[0], x[1], x[2])
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
	h := ProjPlane(a, p)

	if !Eql(r, Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f %f %f", r[0], r[1], r[2])
	}

	if !Eql(h, Vec{2, 0, 0}) {
		t.Errorf("Error Proj Plane  %f %f %f", h[0], h[1], h[2])
	}

	if !Eql(Proj(a, p), Vec{0, 2, 0}) {
		t.Errorf("Error Projection %f, %f, %f", a[0], a[1], a[2])
	}

	p = Vec{1, -1, 0}
	o := Vec{0, 1, 0}
	re := Refl(p, o)
	if !Eql(re, Vec{1, 1, 0}) {
		t.Errorf("Error Reflection %f, %f, %f", re[0], re[1], re[2])
	}

}

func TestMatrix(t *testing.T) {
	A := Mat3{-2, 2, -3, -1, 1, 3, 2, 0, -1}
	B := Mat3{}
	C := Mat4{3, 2, 0, 1, 4, 0, 1, 2, 3, 0, 2, 1, 9, 2, 3, 1}
	D := Mat4{}

	a := Vec{}
	b := Vec4{}

	id3 := Mat3{}
	id4 := Mat4{}

	var err error

	//Construct Mat 3
	for i := 0; i < MAT3; i++ {
		for j := 0; j < MAT3; j++ {
			index, e := Map(i, j, MAT3)
			B[index] = float32(index)
			err = e
			if i == j {
				id3[index] = 1.0
			}
		}
		a[i] = float32(i)

	}

	//Construct Mat 4
	for i := 0; i < MAT4; i++ {
		for j := 0; j < MAT4; j++ {
			index, e := Map(i, j, MAT4)
			D[index] = float32(index)
			err = e
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
	_, _ = id3.CrossMat(&A)
	_, _ = id4.CrossMat(&C)
	_, _ = A.CrossVec(&a)
	_, _ = C.CrossVec(&b)

	fmt.Printf("Determinant Mat3: %f\n", det3)
	fmt.Printf("Determinant Mat4: %f\n", det4)
	//Compute Matrix Inverses
	_ = A.Inverse()
	fmt.Printf("Inverse A Matrix 3x3\n")

}

func BenchmarkVecOp(b *testing.B) {

	p := Vec{1, -1, 0}
	o := Vec{0, 1, 0}

	for i := 0; i < b.N; i++ {
		r := Add(p, o)
		Cross(r, p)
		r = Proj(r, o)
		r = Add(r, o)
	}
}
