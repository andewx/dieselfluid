package kernel

import V "dslfluid.com/dsl/math/math64"

const PI = 3.141592653589

type Cubic struct {
	A  float64
	H  float64
	H0 float64
}

//Build Cubic builds cubic kernel function
func Build_Cubic(h float64) Cubic {
	cubic := Cubic{0, h, h}
	cubic.A = 1 / (PI * h * h * h)
	return cubic
}

//**-------------------CUBIC KERNEL BSPLINE----------------------//
func (K Cubic) F(x float64) float64 {
	r := x / K.H0
	if r > 2.0 {
		return 0.0
	}

	s := (2 - r)

	if r < 1.0 {
		p := (1 - r)
		return K.A * ((0.25 * s * s * s) - (p * p * p))
	}

	return K.A * 0.25 * s * s * s
}

func (K Cubic) W0() float64 {
	return K.F(0)
}

//O1D - 1st order Differential
func (K Cubic) O1D(x float64) float64 {
	//Try the functional derivative
	r := x / K.H0
	q := (2 - r)
	p := (1 - r)
	if r > 2.0 {
		return 0.0
	}
	if r < 1.0 {
		return K.A * ((0.75 * q * q) - 3*(p*p))
	} else {
		return K.A * 0.75 * (q * q)
	}

}

//O2D - Returns 2nd Order Differential
func (K Cubic) O2D(x float64) float64 {

	//Try the functional derivative
	r := x / K.H0
	q := (2 - r)
	p := (1 - r)

	if r > 2.0 {
		return 0.0
	}

	if r < 1.0 {
		return K.A * ((1.5 * q) - 6*p)
	} else {
		return K.A * 1.5 * q
	}

}

//Adjust changes the kernel smoothing length based on density ratio
func (K Cubic) Adjust(densityRatio float64) float64 {

	K.H0 = K.H

	return densityRatio
}

//Grad finds the kernel gradient
func (K Cubic) Grad(x float64, dir V.Vec) V.Vec {
	return V.Scl(dir, -K.O1D(x))
}
