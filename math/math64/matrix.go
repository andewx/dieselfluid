package math64

import "fmt"

//Non Generalized Matrices - for R3 - R4 Matrices and Vector Computations
//Matrix functions sometimes take the object function pointer for operations
type Mat2 [4]float64
type Mat3 [9]float64
type Mat4 [16]float64

const MAT2 = 2
const MAT3 = 3
const MAT4 = 4

//Matrix2 Constructs Custom 2x2 Matrix from construction vectors v0,v1
func Matrix2(v0 Vec, v1 Vec) *Mat2 {
	nMat := Mat2{}
	vecSwitch := 0
	for i := 0; i < MAT2; i++ {
		if vecSwitch == 0 {
			nMat[i] = v0[i%MAT2]
		}
		if vecSwitch == 1 {
			nMat[i] = v1[i%MAT2]
		}

		if i%2 == 0 {
			vecSwitch++
		}
	}
	return &nMat
}

//Matrix3 Constructs Custom 3x3 Matrix from construction vectors v0,v1, v2
func Matrix3(v0 Vec, v1 Vec, v2 Vec) Mat3 {
	nMat := Mat3{}
	vecSwitch := 0
	for i := 0; i < MAT3*MAT3; i++ {
		if vecSwitch == 0 {
			nMat[i] = v0[i%MAT3]
		}
		if vecSwitch == 1 {
			nMat[i] = v1[i%MAT3]
		}
		if vecSwitch == 2 {
			nMat[i] = v2[i%MAT3]
		}
		if i%3 == 0 {
			vecSwitch++
		}
	}
	return nMat
}

//Matrix4 Constructs Custom 4x4 Matrix from construction vectors v0,v1, v2
func Matrix4(v0 Vec, v1 Vec, v2 Vec, v3 Vec) Mat4 {
	nMat := Mat4{}
	vecSwitch := 0
	for i := 0; i < MAT4*MAT4; i++ {
		if vecSwitch == 0 {
			if i < 3 {
				nMat[i] = v0[i%MAT3]
			}
		}
		if vecSwitch == 1 {
			if i < 7 {
				nMat[i] = v1[i%MAT3]
			}
		}
		if vecSwitch == 2 {
			if i < 11 {
				nMat[i] = v2[i%MAT3]
			}
		}
		if vecSwitch == 3 {
			if i < 15 {
				nMat[i] = v3[i%MAT3]
			}
		}
		if i%4 == 0 {
			vecSwitch++
		}
	}
	return nMat
}

//Identity2 constructs the 2x2 identity Matrice
func Identity2() Mat2 {
	nMat := Mat2{}
	nMat[0] = 1.0
	nMat[3] = 1.0
	return nMat
}

//Identity3 constructs the 3 NxN Matrix Identity
func Identity3() Mat3 {
	nMat := Mat3{}
	nMat[0] = 1.0
	nMat[4] = 1.0
	nMat[8] = 1.0
	return nMat
}

//Identity4 constructs the 4 NxN Matrix Identity
func Identity4() Mat4 {
	nMat := Mat4{}
	nMat[0] = 1.0
	nMat[5] = 1.0
	nMat[10] = 1.0
	nMat[15] = 1.0
	return nMat
}

//Maps given indices to single linear array matrix
//gives error for ill defined map. Does not check bounds
func Map(i int, j int, mat_size int) (int, error) {

	if i < 0 || j < 0 {
		return 0, fmt.Errorf("Error Negative Indice")
	}

	if i >= mat_size || j >= mat_size {
		return 0, fmt.Errorf("Error Indices Out of Bounds %d, %d\n", i, j)
	}

	return (i * mat_size) + j, nil
}

//Cross Product Functions Between Mat x Mat / Mat x Vec Computations Immutable
func (a *Mat2) Det() float64 {
	return a[0]*a[3] - a[1]*a[2]
}

func (m *Mat2) Inverse() (*Mat2, error) {
	inverse := Mat2{}
	det := m.Det()

	if det == 0 {
		return &inverse, fmt.Errorf("Error Determinant 0\n")
	}
	inverse[0] = m[3] / det
	inverse[1] = -(m[1] / det)
	inverse[2] = -(m[2] / det)
	inverse[3] = m[0] / det

	return &inverse, nil

}

//Det calcs the 3x3 Determinant
// return m[0][0] *(m[1][1]*m[2][2] - m[1][2]*m[2][1]) - m[0][1]*(m[1][0]*m[2][2] - m[1][2]*m[2][0]) - m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
func (m *Mat3) Det() float64 {
	g := Mat2{m[4], m[5], m[7], m[8]}
	h := Mat2{m[3], m[5], m[6], m[8]}
	j := Mat2{m[3], m[4], m[6], m[7]}

	a := m[0]
	b := m[1]
	c := m[2]

	return a*g.Det() - b*h.Det() + c*j.Det()
}

//Inverse solves the 3x3 Matrix through determinant
func (m *Mat3) Inverse() *Mat3 {
	inv := Mat3{}
	det := m.Det()
	minor := make([]Mat2, 9)
	for i := 0; i < MAT3; i++ {
		for j := 0; j < MAT3; j++ {
			indx, _ := Map(i, j, MAT3)
			//Get Diagional Indexes and apply modulus
			d0, _ := Map((i+1)%3, (j+1)%3, MAT3)
			d1, _ := Map((i+1)%3, (j+2)%3, MAT3)
			d2, _ := Map((i+2)%3, (j+1)%3, MAT3)
			d3, _ := Map((i+2)%3, (j+2)%3, MAT3)

			minor[indx][0] = m[d0]
			minor[indx][1] = m[d1]
			minor[indx][2] = m[d2]
			minor[indx][3] = m[d3]

			inv[indx] = minor[indx].Det() / det
		}
	}
	return &inv
}

//Det Solves the 4x4 Matrix Determinant
func (m *Mat4) Det() float64 {
	a := m[0]
	b := m[1]
	c := m[2]
	d := m[3]

	i := Mat3{m[5], m[6], m[7], m[9], m[10], m[11], m[13], m[14], m[15]}
	j := Mat3{m[4], m[6], m[7], m[8], m[10], m[11], m[12], m[14], m[15]}
	k := Mat3{m[4], m[5], m[7], m[8], m[9], m[11], m[12], m[13], m[15]}
	l := Mat3{m[4], m[5], m[6], m[8], m[9], m[10], m[12], m[13], m[14]}

	return a*i.Det() - b*j.Det() + c*k.Det() - d*l.Det()
}

//Dot Product Rows and Columns
func (a *Mat3) CrossMat(b *Mat3) (*Mat3, error) {
	//Initialize Empty Matrix
	r := Mat3{}
	var Err error
	var index int
	var entry float64
	var rv, cv int
	for i := 0; i < MAT3; i++ {
		for j := 0; j < MAT3; j++ {
			index, Err = Map(i, j, MAT3)
			entry = float64(0.0)
			for k := 0; k < MAT3; k++ {
				rv, _ = Map(i, k, MAT3)
				cv, Err = Map(k, j, MAT3)
				entry += a[rv] * b[cv]
			}
			r[index] = entry
		}
		//End Loop
	}
	return &r, Err
}

//Cross multiplies 4 x 4 Square Matrix - Non Generalized return error for out of bound invalid indexes
func (a *Mat4) CrossMat(b *Mat4) (*Mat4, error) {
	//Initialize Empty Matrix
	r := Mat4{}
	var Err error
	var index int
	var entry float64
	var rv, cv int
	for i := 0; i < MAT4; i++ {
		for j := 0; j < MAT4; j++ {
			index, Err = Map(i, j, MAT4)
			entry = float64(0.0)
			for k := 0; k < MAT4; k++ {
				rv, _ = Map(i, k, MAT4)
				cv, Err = Map(k, j, MAT4)
				entry += a[rv] * b[cv]
			}
			r[index] = entry
		}
		//End Loop
	}
	return &r, Err
}

//Matrix Vec Cross Multiplication - Gives a Row Vector
func (a *Mat3) CrossVec(b *Vec) (*Vec, error) {
	//Initialize Empty VecN
	r := Vec{}
	var Err error

	for i := 0; i < MAT3; i++ {
		r[i] = a[i]*b[0] + a[MAT3+i]*b[1] + a[(MAT3*2)+i]*b[2]
	}

	return &r, Err
}

//Matrix Vec Cross Multiplication - Gives a Row Vector
func (a *Mat2) CrossVec(b *Vec2) (*Vec2, error) {
	//Initialize Empty Vec ()
	r := Vec2{}
	var Err error

	r[0] = a[0]*b[0] + a[1]*b[1]
	r[1] = a[2]*b[0] + a[3]*b[1]

	return &r, Err
}

//Matrix Vec Cross Multiplication - Gives a Row Vector
func (a *Mat4) CrossVec(b *Vec4) (*Vec4, error) {
	//Initialize Empty Vec ()
	r := Vec4{}
	var Err error

	for i := 0; i < MAT4; i++ {
		r[i] = a[i]*b[0] + a[MAT4+i]*b[1] + a[(MAT4*2)+i]*b[2] + a[(MAT4*3)+i]*b[3]
	}

	return &r, Err
}

//Translation accumulates and adds in the given translation vector into the current matrix.
func (a *Mat4) Translation(b *Vec) {

	a[12] = a[12] + b[0]
	a[13] = a[13] + b[1]
	a[14] = a[14] + b[2]
	a[15] = 1.0

}

//Transpose alters and transposes the current matrix data
func (a *Mat4) Transpose() *Mat4 {
	for i := 0; i < MAT4; i++ {
		for j := 0; j < MAT4; j++ {
			id1 := i*MAT4 + j
			id2 := j*MAT4 + i
			tmp := a[id1]
			a[id1] = a[id2]
			a[id2] = tmp
		}
	}
	return a
}

//Projection Matrix creates a projection matrix
func ProjectionMatrix(l float64, r float64, t float64, b float64, n float64, f float64) Mat4 {
	proj := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, (-f / (f - n)), (-f * n) / (f - n), 0, 0, -1, 0} //scratch a pixel projection matrix
	return proj
}

func (m *Mat3) Dot(b *Vec) *Vec {
	//Initialize Empty Vec ()
	r := Vec{}

	for i := 0; i < MAT3; i++ {
		r[i] = m[i]*b[0] + m[MAT3+i]*b[1] + m[(MAT3*2)+i]*b[2]
	}
	return &r
}

//Mat2 string
func (a *Mat2) String() string {
	s := ""

	for i := 0; i < MAT2; i++ {
		s += fmt.Sprintf("| ")
		for j := 0; j < MAT2; j++ {
			id, _ := Map(i, j, MAT2)
			s += fmt.Sprintf(" %f ", a[id])
		}
		s += fmt.Sprintf(" |\n")
	}
	return s
}

//Matrix Print Operations
func (a *Mat3) String() string {
	s := ""

	for i := 0; i < MAT3; i++ {
		s += fmt.Sprintf("| ")
		for j := 0; j < MAT3; j++ {
			id, _ := Map(i, j, MAT3)
			s += fmt.Sprintf(" %f ", a[id])
		}
		s += fmt.Sprintf(" |\n")
	}
	return s
}

//Matrix Print Operations
func (a *Mat4) String() string {
	s := ""

	for i := 0; i < MAT4; i++ {
		s += fmt.Sprintf("| ")
		for j := 0; j < MAT4; j++ {
			id, _ := Map(i, j, MAT4)
			s += fmt.Sprintf(" %f ", a[id])
		}
		s += fmt.Sprintf(" |\n")
	}
	return s
}
