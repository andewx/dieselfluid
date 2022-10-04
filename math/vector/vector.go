package vector

//Implements Vector Data structure and Condenses Vector format to the
//sliced array. This format representation of vectors allows for generalized,
//parameterization

//NonPointer Methods with Explicit Parameters are considered immutable while,
//Functions with associated type pointers are mutating. For the sake of completness
//the documentation will indiciate if the function mutates or not with keywords
// @MUTATES / @DELEGATES will be added if there is any question

//Vector operations for vectors with 0 Dimensions or unequal dimensions will fail
//quietly by simply logging errros

import (
	"fmt"
	"math"

	"github.com/andewx/dieselfluid/math/common"
)

//Vector slice array
type Vec []float32

//------------Vector Utility Functions----------------------//

func Cast(a [3]float32) []float32 {
	return []float32{a[0], a[1], a[2]}
}

func CastFixed(a Vec) [3]float32 {
	return [3]float32{a[0], a[1], a[2]}
}

//Vec2 creates new vector as 2-Length slice
func Vec2() Vec {
	return []float32{0, 0}
}

//Vec3 creates new vector as 3-Length slice
func Vec3() Vec {
	return []float32{0, 0, 0}
}

//Vec4 creates new vector as 3-Length slice
func Vec4() Vec {
	return []float32{0, 0, 0, 0}
}

func VecN(n int) Vec {
	return make([]float32, n)
}

//Vec Dim returns vector dimensions
func (v Vec) Dim() int {
	return len(v)
}

//Copy copies the reference vector into the store vec
//@MUTATES store and updates with copies of reference Vec
//Copy will attempt to copy even if vectors are different sizes
func Copy(store Vec, reference Vec) {
	r := 0
	if len(store) >= len(reference) {
		r = len(reference)
	} else {
		r = len(store)
	}
	for i := 0; i < r; i++ {
		store[i] = reference[i]
	}
}

//Copy copies the reference vector into the store vec
//@MUTATES store and updates with copies of reference Vec
//Copy will attempt to copy even if vectors are different sizes
func (store Vec) Copy(reference Vec) {
	r := 0
	if len(store) >= len(reference) {
		r = len(reference)
	} else {
		r = len(store)
	}
	for i := 0; i < r; i++ {
		store[i] = reference[i]
	}
}

//DimEq(a, b) checks if two vectors have equal sizes
func DimEq(a Vec, b Vec) bool {
	if len(a) == len(b) {
		return true
	}
	return false
}

//DimEq(a, b) checks if two vectors have equal sizes
func (a Vec) DimEq(b Vec) bool {
	if len(a) == len(b) {
		return true
	}
	return false
}

//Equal checks if two vectors are logically equal their
//subtraction and the resulting mag() of that vector gives
//a value less than EPSIOLN if two vectors are equal
func (a Vec) Equal(b Vec) bool {
	dist := Sub(b, a)
	if dist.Mag() < common.EPSILON {
		return true
	}
	return false
}

//Equal checks if two vectors are logically equal their
//subtraction and the resulting mag() of that vector gives
//a value less than EPSIOLN if two vectors are equal
func Equal(a Vec, b Vec) bool {
	dist := Sub(b, a)
	if dist.Mag() < common.EPSILON {
		return true
	}
	return false
}

//Negates the current vector by -1.0. But stores the result in a new value
//@NONMUTATE function see NegV for mutating
func Neg(a Vec) Vec {
	return Scale(a, -1.0)
}

//Negates the current vector by -1.0. But stores the result in a new value
//@NONMUTATE function see NegV for mutating
func (a Vec) Neg() Vec {
	return Scale(a, -1.0)
}

//Negates the current vector by -1.0. But stores the result in a new value
//@NONMUTATE function see NegV for mutating
func (a Vec) NegV() Vec {
	a = a.Scale(-1.0)
	return a
}

func (v Vec) ToString() string {
	mstring := "[ "
	for i := 0; i < len(v); i++ {
		mstring += fmt.Sprintf(" %.2f ", v[i])
	}
	mstring += "]\n"
	return mstring
}

func LogVecError(vec Vec, op string) Vec {
	fmt.Printf(op+": FAIL | Vector %s", vec.ToString())
	return vec
}

//--------------Vector Math Vector Operations-----------------------------
// Basic Vector math operations, most vector operations shuold not be mutating
// Vector state. With the exception of certain functions where the intended use
// is most often a method with a side effect such as scale

//Add sums each element index against the other vector a stores the result
//in a new vector @DELEGATES
func Add(a Vec, b Vec) Vec {
	l := len(a)
	g := len(b)
	if len(b) < len(a) {
		l = len(b)
		g = len(a)
	}
	c := make([]float32, g)
	for i := 0; i < l; i++ {
		c[i] = a[i] + b[i]
	}
	return c
}

//Add sums each element index against the other vector a stores the result
//in a new vector @DELEGATES
func (a Vec) Add(b Vec) Vec {
	l := len(a)
	g := len(b)
	if len(b) < len(a) {
		l = len(b)
		g = len(a)
	}
	c := make([]float32, g)
	for i := 0; i < l; i++ {
		c[i] = a[i] + b[i]
	}
	return c
}

//Sub subtracts the vector b minus a utilizing the same add functions
//above. @NONMUTATE
func Sub(b Vec, a Vec) Vec {
	return Add(b, Scale(a, -1.0))
}

//Sub subtracts the vector b minus a utilizing the same add functions
//above. @NONMUTATE
func (b Vec) Sub(a Vec) Vec {
	return Add(b, Scale(a, -1.0))
}

func (b Vec) Mul(a Vec) Vec {
	return Vec{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

//Scl scales vector kX this function does not mutate the current vector
func Scale(a Vec, k float32) Vec {
	c := make([]float32, len(a))
	for i := 0; i < a.Dim(); i++ {
		c[i] = a[i] * k
	}
	return c
}

//Scl scales vector kX this function does not mutate the current vector
func ScaleVar(a Vec, b Vec) Vec {
	c := make([]float32, len(a))
	for i := 0; i < a.Dim(); i++ {
		c[i] = a[i] * b[i]
	}
	return c
}

//Scale scales vector kX this function does not mutate the current vector
//@NONMUTATE function Scale
func (a Vec) Scale(k float32) Vec {
	c := make([]float32, len(a))
	for i := 0; i < a.Dim(); i++ {
		c[i] = a[i] * k
	}
	return c
}

//---TRIPLE CHECK THIS LOL
//Dot computes the dot product as in a.Dot(b) where the resulting value should return
//a scalar value representative of the ||a|||b||Cos() theta angle between two non-parallel vectors
func Dot(a Vec, b Vec) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

//Cross cross multpiplies A x B. Be aware that cross multplication is not commutative
//Cross is a non-mutating function @NONMUTATE and represents a Vector that is perpendicular
//To the plane created from two vectors when those vectors are 3 Vecs satisfying that property
//In general cross products produce an orthonormal vector between two vectors.
//For now only works on 3 Vecs and  3 Like 4 Vecs =
func Cross(a Vec, b Vec) Vec {

	if !DimEq(a, b) {
		return LogVecError(a, "CROSS FAIL Dim !Equal")
	}

	if len(a) != 3 {
		return LogVecError(a, "CROSS FAIL !Dim 3 Vec")
	}

	g := Vec{a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - b[2]*a[0],
		a[0]*b[1] - b[0]*a[1]}
	return g
}

//Mag (v Vec) returns the magnitued of any vector > length size 0. Operation
//does not mutate vector state
func Mag(v Vec) float32 {

	size := float32(0)
	for i := 0; i < len(v); i++ {
		size += v[i] * v[i]
	}
	return float32(math.Sqrt(float64(size)))
}

//Mag (v Vec) returns the magnitued of any vector > length size 0. Operation
//does not mutate vector state @DELEGATES
func (v Vec) Mag() float32 {
	size := float32(0)
	for i := 0; i < len(v); i++ {
		size += v[i] * v[i]
	}
	return float32(math.Sqrt(float64(size)))
}

//Returns normalized vector A.i / Mag() function is non mutating and delegates
//result to a new vector
func Norm(a Vec) Vec {
	v := make([]float32, len(a))
	l := Mag(a)
	if l != 0.0 {
		for i := 0; i < len(a); i++ {
			v[i] = a[i] / l
		}
	}
	return v
}

//Returns normalized vector A.i / Mag() function is non mutating and delegates
//result to a new vector - fails quietly with 0 vector
func (a Vec) Norm() Vec {
	v := make([]float32, len(a))
	l := Mag(a)
	if l != 0.0 {
		for i := 0; i < len(a); i++ {
			v[i] = a[i] / l
		}
	}
	return v
}

//Proj projects a vector onto the N Vec a -> N. Functino is non mutating and delegates
//The output vector to a new vector. Error log checking handled by internal dot product
//call. N typically may be considered to be a Normal vector that A will be projected onto.
func Proj(a Vec, n Vec) Vec {
	vn := Norm(n)
	proj := Scale(vn, Dot(a, n)/Mag(n))
	return proj
}

//Proj projects a vector onto the N Vec a -> N. Functino is non mutating and delegates
//The output vector to a new vector. Error log checking handled by internal dot product
//call. N typically may be considered to be a Normal vector that A will be projected onto.
func (a Vec) Proj(n Vec) Vec {
	vn := Norm(n)
	proj := Scale(vn, Dot(a, n)/Mag(n))
	return proj
}

//ProjPlane projects a vector on to the the plane defined by A & N. Functino is non mutating and delegates
//The output vector to a new vector. Error log checking handled by internal dot product
//call. N typically may be considered to be a Normal vector that A will be projected onto.
func ProjPlane(a Vec, n Vec) Vec {
	pVec := Proj(a, n)
	return Sub(pVec, a)
}

//ProjPlane projects a vector on to the the plane defined by A & N. Functino is non mutating and delegates
//The output vector to a new vector. Error log checking handled by internal dot product
//call. N typically may be considered to be a Normal vector that A will be projected onto.
func (a Vec) ProjPlane(n Vec) Vec {
	pVec := Proj(a, n)
	return Sub(pVec, a)
}

//Refl reflects the v vector about N this function is non mutating and delegates
//vector result into a new vector
func Refl(v Vec, n Vec) Vec {
	b := Scale(n, Dot(v, n)*2.0)
	return Sub(v, b)
}

//Refl reflects the v vector about N this function is non mutating and delegates
//vector result into a new vector
func (v Vec) Refl(n Vec) Vec {
	b := Scale(n, Dot(v, n)*2.0)
	return Sub(v, b)
}

//Eql Checks if two vectors are equal by lazy comparison
//If the user of this API wishes to check whether two vectors are within EPSILON of
//Eachother please use the Equal() function denoted at the top of this library
func Eql(a Vec, b Vec) bool {

	if !DimEq(a, b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//Eql Checks if two vectors are equal by lazy comparison
//If the user of this API wishes to check whether two vectors are within EPSILON of
//Eachother please use the Equal() function denoted at the top of this library
func (a Vec) Eql(b Vec) bool {

	if !DimEq(a, b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//isEpsilon is a non exported function, check whether the value a is within EPSILON
//for component vectors you may wish to see if Vec Mags are within Epsiolon
func isEpsilon(a float32) bool {
	a = float32(math.Abs(float64(a)))
	if a < common.EPSILON {
		return true
	}
	return false
}

//Dist returns the distance between two vectors B-A/L by non mutating ops
func Dist(a Vec, b Vec) float32 {
	return Mag(Sub(a, b))
}

//Dist returns the distance between two vectors B-A/L by non mutating ops
func (a Vec) Dist(b Vec) float32 {
	return Mag(Sub(a, b))
}

//Tan returns the tangent vector from the given normal profile - this may
//Just be a plane projection of A onto the the plane p....defined by projection
//of A onto N when A is not parallel to N.
func Tan(a Vec, n Vec) Vec {
	p := Proj(a, n)
	return Sub(a, p)
}

func (a Vec) Tan(n Vec) Vec {
	p := Proj(a, n)
	return Sub(a, p)
}

//RaySphereIntersection Calculates the Intersection Points for a Ray and Sphere
//r0 is ray origin , d0 direction , c is the sphere origin, r is the radius
func RaySphereIntersection(r0 Vec, d0 Vec, c Vec, r float32) (Vec, bool) {
	mVec := Vec{}
	vpc := Sub(c, r0)
	vmag := vpc.Mag()
	pc := Proj(c, d0)
	pcc2 := Mag(Sub(pc, c))
	pcc2 *= pcc2
	dist := float32(math.Sqrt(float64(r*r - (pcc2))))

	if Dot(vpc, d0) < 0 {
		if vmag > r {
			return vpc, false
		}

		if vmag == r {
			return r0, true
		}
		di1 := dist - Sub(pc, r0).Mag()
		inter := Add(r0, Scale(d0, di1))
		return inter, true
	} else {

		di1 := float32(0.0)
		if Mag(c.Sub(pc)) > r {
			return mVec, false
		} else {
			if vpc.Mag() > r {
				di1 = Mag(Sub(pc, r0)) - dist
			} else {
				di1 = Mag(Sub(pc, r0)) + dist
			}
		}
		inter := Add(r0, Scale(d0, di1))
		return inter, true
	}
}

//Clamps Vector between min and max values for each entry
//@mutate_selector / #utility / #mgl / #vector / #math
func (v Vec) Clamp(min float32, max float32) {
	v[0] = common.Clamp1f(v[0], min, max)
	v[1] = common.Clamp1f(v[1], min, max)
	v[2] = common.Clamp1f(v[2], min, max)
}

func SinDot(a Vec, b Vec) float32 {
	cosT := Dot(a, b) * (a.Mag() * b.Mag())
	return float32(math.Sqrt(1.0 - float64(cosT*cosT)))
}
