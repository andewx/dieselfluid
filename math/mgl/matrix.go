package mgl

import "fmt"
import "math"

//-----------------------DSL FLUID Mat LIBRARY ----------------------------
//DSL FLUID Mat Library implementation is non mutating (unless) the function
//had a tag @MUTATE and typically a suffix V. 3D Mat Libraries are typically
//Square so by convention we will only allow mxm square defined Mates. This
//Allows us to wrangle the matrices into one data type slice array of floats

type Mat []float32

const MAT2 = 2
const MAT2_SIZE = 4
const MAT3 = 3
const MAT3_SIZE = 9
const MAT4 = 4
const MAT4_SIZE = 16

//--------------Mat Utilities-----------------------------------------------

//Mat2 Creates new float slice type Mat of length 2 set to identity val
func Mat2(v float32) Mat {
	a := make([]float32, MAT2_SIZE)
	for i := 0; i < MAT2_SIZE; i += (MAT2 + 1) {
		a[i] = v
	}
	return a
}

//Mat2 Creates new float slice type Mat of length 2 set to identity val
func Mat3(v float32) Mat {
	a := make([]float32, MAT3_SIZE)
	for i := 0; i < MAT3_SIZE; i += (MAT3 + 1) {
		a[i] = v
	}
	return a
}

//Mat2 Creates new float slice type Mat of length 2 set to identity val
func Mat4(v float32) Mat {
	a := make([]float32, MAT4_SIZE)
	for i := 0; i < MAT4_SIZE; i += (MAT4 + 1) {
		a[i] = v
	}
	return a
}

//MatN Creates new float slice type Mat of length 2 set to identity val
func MatN(dim int, v float32) Mat {
	a := make([]float32, dim*dim)
	for i := 0; i < dim*dim; i += (dim + 1) {
		a[i] = v
	}
	return a
}

//Mat2 Constructs Custom 2x2 Mat from construction vectors v0,v1
func Mat2V(v0 Vec, v1 Vec) Mat {
	nMat := make([]float32, 4)
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
	return nMat
}

//Mat3 Constructs Custom 3x3 Mat from construction vectors v0,v1, v2
func Mat3V(v0 Vec, v1 Vec, v2 Vec) Mat {
	nMat := make([]float32, 9)

	vecSwitch := 0
	for i := 0; i < MAT3; i++ {
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

//Mat4 Constructs Custom 4x4 Mat from construction vectors v0,v1, v2
func Mat4V(v0 Vec, v1 Vec, v2 Vec, v3 Vec) Mat {
	nMat := make([]float32, 16)
	vecSwitch := 0
	for i := 0; i < MAT4; i++ {
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

//Copy returns a new matrix copy from the addressed Mat object
func (m Mat) Copy() Mat {
	newMat := make(Mat, len(m))
	for i := 0; i < len(m); i++ {
		newMat[i] = m[i]
	}

	return newMat
}

//MatEqual returns two Matrixes equal by epsilon value comparision
func MatEqual(a Mat, e Mat) bool {
	if len(a) != len(e) {
		return false
	}

	for k := 0; k < len(a); k++ {
		comp := math.Abs(float64(e[k] - a[k]))
		if comp > 0.0000001 {
			return false
		}
	}
	return true
}

//Dim Returns Dimensionality of Square Matrix Based on Size...If size not recognized
//Then Dim returns non nil error
func (M Mat) Dim() int {
	if len(M) == 4 {
		return 2
	}

	if len(M) == 9 {
		return 3
	}

	if len(M) == 16 {
		return 4
	}

	return 0
}

//Dim Returns Dimensionality of Square Matrix Based on Size...If size not recognized
//Then Dim returns non nil error
func Dim(M Mat) int {
	if len(M) == 4 {
		return 2
	}

	if len(M) == 9 {
		return 3
	}

	if len(M) == 16 {
		return 4
	}

	return 0
}

//MAPS any 2D Matrix Position to its linear ROW-Major order layout in memory with
//the returned index being 0 for invalid inputs
func Map(i int, j int, mat_size int) int {

	if i < 0 || j < 0 {
		return 0
	}

	if i >= mat_size || j >= mat_size {
		return 0
	}

	return (i * mat_size) + j
}

func (a Mat) Det() float32 {
	dim := a.Dim()
	if dim == 4 {
		return a.Det4()
	}

	if dim == 3 {
		return a.Det3()
	}

	if dim == 2 {
		return a.Det2()
	} else {
		return 0.0
	}
}

//Det2() Handles this matrix as size 2 matrix and returns the determinant
func Det2(a Mat) float32 {
	if len(a) != 4 {
		fmt.Printf("Invalid Matrix size for Det2()\n")
		return 0.0
	}
	a0 := a[0]
	a3 := a[3]
	a1 := a[1]
	a2 := a[2]
	return a0*a3 - a1*a2
}

//Det2() Handles this matrix as size 2 matrix and returns the determinant
func (a Mat) Det2() float32 {
	if len(a) != 4 {
		fmt.Printf("Invalid Matrix size for Det2()\n")
		return 0.0
	}
	a0 := a[0]
	a3 := a[3]
	a1 := a[1]
	a2 := a[2]
	return a0*a3 - a1*a2
}

//Inverse2() Computes the Matrix Inverse of a 2 Sized Matrix for now all inverse
//And determinant functinos will be named since they are inmplemented differently
//Inverse 2 is non mutating and delgates the value of the new matrix
func (m Mat) Inv2() Mat {
	inverse := Mat2(0.0)

	if m.Dim() != 2 {
		fmt.Printf("Invalid Matrix size for Inverse2() returning Mat2 Identity\n")
		return Mat2(1.0)
	}

	det := m.Det2()

	if det == 0 {
		return inverse
	}
	inverse[0] = m[3] / det
	inverse[1] = -(m[1] / det)
	inverse[2] = -(m[2] / det)
	inverse[3] = m[0] / det

	return inverse

}

//Det calcs the 3x3 Determinant
// return m[0][0] *(m[1][1]M[2][2] - m[1][2]M[2][1]) - m[0][1]*(m[1][0]M[2][2] - m[1][2]M[2][0]) - m[0][2]*(m[1][0]M[2][1]-m[1][1]M[2][0])
func (m Mat) Det3() float32 {
	g := Mat([]float32{m[4], m[5], m[7], m[8]})
	h := Mat([]float32{m[3], m[5], m[6], m[8]})
	j := Mat([]float32{m[3], m[4], m[6], m[7]})

	a := m[0]
	b := m[1]
	c := m[2]

	return a*g.Det2() - b*h.Det2() + c*j.Det2()
}

//Set(row int, v Vec) sets a matrix row to the specified vector
//WARNING this is a mutatiing @MUTATE operation that changes the internal
//Values of the referenced matrix M
func (m Mat) Set(row int, v Vec) {
	dim := m.Dim()
	if dim < len(v) {
		LogVecError(v, "Set(row int, v Vec) FAILED vector length less the matrix dim")
		return
	}
	idx := Map(row, 0, dim)
	for i := 0; i < 3; i++ {
		m[i+idx] = v[i]
	}
}

//Get obtains the Row Vector from the parameter row copies vector value into a
//new vector and returns.
func (m Mat) Get(row int) Vec {
	idx := Map(row, 0, m.Dim())
	v := VecN(m.Dim())
	for i := 0; i < m.Dim(); i++ {
		v[i] = m[i+idx]
	}
	return v
}

// Mul performs a scalar multiplcation of the Mat. This is equivalent to iterating
// over every element of the Mat and multiply it by c.
func (m1 Mat) Mul(c float32) Mat {
	m2 := MatN(m1.Dim(), 0.0)
	for i := 0; i < m1.Dim()*m1.Dim(); i++ {
		m2[i] = m1[i] * c
	}
	return m2
}

// Inv computes the inverse of a square Mat. An inverse is a square Mat such that when multiplied by the
// original, yields the identity.
//
// M_inv * M = M * M_inv = I
//
// In this library, the math is precomputed, and uses no loops, though the multiplications, additions, determinant calculation, and scaling
// are still done. This can still be (relatively) expensive for a 4x4.
//
// This function checks the determinant to see if the Mat is invertible.
// If the determinant is 0.0, this function returns the zero Mat. However, due to floating point errors, it is
// entirely plausible to get a false positive or negative.
// In the future, an alternate function may be written which takes in a pre-computed determinant.
//Tries to simplify main call for matrixes of different dimensions by delgating to other calls if the dimension aren't equal
func (m Mat) Inv() Mat {

	//Main attach point for other inverses
	if m.Dim() != 4 {

		if m.Dim() == 3 {
			return m.Inv3()
		}

		if m.Dim() == 2 {
			return m.Inv2()
		}
		return Mat4(1.0)
	}

	det := m.Det4()
	if det == 0.0 {
		return Mat4(1.0)
	}

	retMat := Mat{
		-m[7]*m[10]*m[13] + m[6]*m[11]*m[13] + m[7]*m[9]*m[14] - m[5]*m[11]*m[14] - m[6]*m[9]*m[15] + m[5]*m[10]*m[15],
		m[3]*m[10]*m[13] - m[2]*m[11]*m[13] - m[3]*m[9]*m[14] + m[1]*m[11]*m[14] + m[2]*m[9]*m[15] - m[1]*m[10]*m[15],
		-m[3]*m[6]*m[13] + m[2]*m[7]*m[13] + m[3]*m[5]*m[14] - m[1]*m[7]*m[14] - m[2]*m[5]*m[15] + m[1]*m[6]*m[15],
		m[3]*m[6]*m[9] - m[2]*m[7]*m[9] - m[3]*m[5]*m[10] + m[1]*m[7]*m[10] + m[2]*m[5]*m[11] - m[1]*m[6]*m[11],
		m[7]*m[10]*m[12] - m[6]*m[11]*m[12] - m[7]*m[8]*m[14] + m[4]*m[11]*m[14] + m[6]*m[8]*m[15] - m[4]*m[10]*m[15],
		-m[3]*m[10]*m[12] + m[2]*m[11]*m[12] + m[3]*m[8]*m[14] - m[0]*m[11]*m[14] - m[2]*m[8]*m[15] + m[0]*m[10]*m[15],
		m[3]*m[6]*m[12] - m[2]*m[7]*m[12] - m[3]*m[4]*m[14] + m[0]*m[7]*m[14] + m[2]*m[4]*m[15] - m[0]*m[6]*m[15],
		-m[3]*m[6]*m[8] + m[2]*m[7]*m[8] + m[3]*m[4]*m[10] - m[0]*m[7]*m[10] - m[2]*m[4]*m[11] + m[0]*m[6]*m[11],
		-m[7]*m[9]*m[12] + m[5]*m[11]*m[12] + m[7]*m[8]*m[13] - m[4]*m[11]*m[13] - m[5]*m[8]*m[15] + m[4]*m[9]*m[15],
		m[3]*m[9]*m[12] - m[1]*m[11]*m[12] - m[3]*m[8]*m[13] + m[0]*m[11]*m[13] + m[1]*m[8]*m[15] - m[0]*m[9]*m[15],
		-m[3]*m[5]*m[12] + m[1]*m[7]*m[12] + m[3]*m[4]*m[13] - m[0]*m[7]*m[13] - m[1]*m[4]*m[15] + m[0]*m[5]*m[15],
		m[3]*m[5]*m[8] - m[1]*m[7]*m[8] - m[3]*m[4]*m[9] + m[0]*m[7]*m[9] + m[1]*m[4]*m[11] - m[0]*m[5]*m[11],
		m[6]*m[9]*m[12] - m[5]*m[10]*m[12] - m[6]*m[8]*m[13] + m[4]*m[10]*m[13] + m[5]*m[8]*m[14] - m[4]*m[9]*m[14],
		-m[2]*m[9]*m[12] + m[1]*m[10]*m[12] + m[2]*m[8]*m[13] - m[0]*m[10]*m[13] - m[1]*m[8]*m[14] + m[0]*m[9]*m[14],
		m[2]*m[5]*m[12] - m[1]*m[6]*m[12] - m[2]*m[4]*m[13] + m[0]*m[6]*m[13] + m[1]*m[4]*m[14] - m[0]*m[5]*m[14],
		-m[2]*m[5]*m[8] + m[1]*m[6]*m[8] + m[2]*m[4]*m[9] - m[0]*m[6]*m[9] - m[1]*m[4]*m[10] + m[0]*m[5]*m[10],
	}

	return retMat.Mul(1 / det)
}

//Inverse solves the 3x3 Mat through determinant
func (m Mat) Inv3() Mat {

	if m.Dim() != 3 {
		return Mat3(1.0)
	}

	inv := make([]float32, 9)
	det := m.Det3()
	minor := make([]Mat, 9)

	for k := 0; k < 9; k++ {
		minor[k] = Mat2(0.0)
	}
	for i := 0; i < MAT3; i++ {
		for j := 0; j < MAT3; j++ {
			indx := Map(i, j, MAT3)
			//Get Diagional Indexes and apply modulus
			d0 := Map((i+1)%3, (j+1)%3, MAT3)
			d1 := Map((i+1)%3, (j+2)%3, MAT3)
			d2 := Map((i+2)%3, (j+1)%3, MAT3)
			d3 := Map((i+2)%3, (j+2)%3, MAT3)

			minor[indx][0] = m[d0]
			minor[indx][1] = m[d1]
			minor[indx][2] = m[d2]
			minor[indx][3] = m[d3]

			inv[indx] = minor[indx].Det2() / det
		}
	}
	return inv
}

//Det Solves the 4x4 Mat Determinant
func (m Mat) Det4() float32 {
	a := m[0]
	b := m[1]
	c := m[2]
	d := m[3]

	i := Mat{m[5], m[6], m[7], m[9], m[10], m[11], m[13], m[14], m[15]}
	j := Mat{m[4], m[6], m[7], m[8], m[10], m[11], m[12], m[14], m[15]}
	k := Mat{m[4], m[5], m[7], m[8], m[9], m[11], m[12], m[13], m[15]}
	l := Mat{m[4], m[5], m[6], m[8], m[9], m[10], m[12], m[13], m[14]}

	return a*i.Det() - b*j.Det() + c*k.Det() - d*l.Det()
}

//Multplies two matrices with mutation into left associated A
func MulM(a Mat, b Mat) Mat {
	if len(a) != len(b) {
		return a
	}
	c := MatN(a.Dim(), 0.0)
	var entry float32
	var rv, cv int
	dim := a.Dim()

	for i := 0; i < dim; i++ {
		entry = 0
		for j := 0; j < dim; j++ {
			idx := Map(i, j, dim)
			entry = 0
			for k := 0; k < dim; k++ {
				rv = Map(i, k, dim)
				cv = Map(k, j, dim)
				entry += a[rv] * b[cv] //Row Col Dot Product
			}
			c[idx] = entry
		}
	}

	return c
}

//Multplies two matrices with mutation into left associated A
func (a Mat) MulM(b Mat) Mat {
	if len(a) != len(b) {
		return a
	}
	c := MatN(a.Dim(), 0.0)
	var entry float32
	var rv, cv int
	dim := a.Dim()

	for i := 0; i < dim; i++ {
		entry = 0
		for j := 0; j < dim; j++ {
			idx := Map(i, j, dim)
			entry = 0
			for k := 0; k < dim; k++ {
				rv = Map(i, k, dim)
				cv = Map(k, j, dim)
				entry += a[rv] * b[cv] //Row Col Dot Product
			}
			c[idx] = entry
		}
	}

	return c
}

//CrossVec crosses a matrix with a vector by summing the Matrix Row (Row Major)
//With the associative column vector. Assuming matrices are all MxM in this LIBRARY
//Vector is returned length M
func (a Mat) CrossVec(b Vec) Vec {

	if a.Dim() != len(b) {
		LogVecError(b, "Cross Vec Matrix Dim != Vec Dim")
		return VecN(a.Dim())
	}
	r := VecN(a.Dim())

	for i := 0; i < a.Dim(); i++ {
		for j := 0; j < a.Dim(); j++ {
			mIdx := Map(i, j, a.Dim())
			r[i] += a[mIdx] * b[j]
		}
	}

	return r
}

//Transpose Matrix reverses matrix index positions for square matrixes only
//Transpose is a non mutating function
func (a Mat) Transpose() Mat {
	dim := a.Dim()
	retMat := MatN(dim, 1.0)
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			id1 := i*dim + j
			id2 := j*dim + i
			retMat[id2] = a[id1]
		}
	}
	return retMat
}

//Normalize normalizes the current matrix. This operation is a mutating operation
//WARNING THIS OPERATION IS A MUTATION OF THE CURRENT MATRIX
func (M Mat) Normalize() Mat {
	det := M.Det()
	if det == 0.0 {
		return M
	}
	for i := 0; i < 16; i++ {
		M[i] = M[i] / det
	}
	return M
}

//Projection Mat creates a projection Mat
func ProjectionMat(l float32, r float32, t float32, b float32, n float32, f float32) Mat {
	s := 1 / (float32(math.Tan(float64(45.0 / 2 * 3.141529 / 180))))
	proj := Mat{s, 0, 0, 0, 0, s, 0, 0, 0, 0, (-f / (f - n)), (-f * n) / (f - n), 0, 0, -1, 0}
	//	proj := Mat4{2 * n / (r - l), 0, (r + l) / (r - l), 0, 0, (2 * n) / (t - b), (t + b) / (t - b), 0, 0, 0, (-(f + n) / (f - n)), (-2 * f * n) / (f - n), 0, 0, -1, 0} //scratch a pixel projection Mat
	return proj
}

//Projection Mat creates a projection Mat
func ProjectionMatF(fov float32, n float32, f float32) Mat {
	s := 1 / (float32(math.Tan(float64((fov / 2) * (3.141529 / 180)))))
	proj := Mat{s, 0, 0, 0, 0, s, 0, 0, 0, 0, (-f / (f - n)), (-f * n) / (f - n), 0, 0, -1, 0}
	//	proj := Mat4{2 * n / (r - l), 0, (r + l) / (r - l), 0, 0, (2 * n) / (t - b), (t + b) / (t - b), 0, 0, 0, (-(f + n) / (f - n)), (-2 * f * n) / (f - n), 0, 0, -1, 0} //scratch a pixel projection Mat
	return proj
}

func Mat3T4(a Mat) Mat {
	newMat := Mat4(1.0)
	for i := 0; i < 9; i++ {
		newMat[i] = a[i]
	}
	return newMat
}

//Mat Print Operations
func (a Mat) String() string {
	s := ""
	dim := a.Dim()
	for i := 0; i < dim; i++ {
		s += fmt.Sprintf("| ")
		for j := 0; j < dim; j++ {
			id := Map(i, j, dim)
			s += fmt.Sprintf(" %f ", a[id])
		}
		s += fmt.Sprintf(" |\n")
	}
	return s
}
