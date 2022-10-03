package kernel

import V "github.com/andewx/dieselfluid/math/vector"

const PI = 3.141592653589

type Cubic struct {
	A  float32
	B  float32 //01D
	C  float32 //O2D
	H1 float32 //Ref Smooth
	H_ float32 //Ref Smooth
	H2 float32
	H3 float32
	H4 float32
	H5 float32
}

//Spiky Cubic Kernel with Std Guassian in the Normal Function
func Build_Kernel(h float32) Cubic {
	cubic := Cubic{0, 0, 0, h, h, 0, 0, 0, 0}
	cubic.H2 = h * h
	cubic.H3 = h * h * h
	cubic.H4 = h * h * h * h
	cubic.H5 = h * h * h * h * h
	cubic.A = (315 / (64 * PI * cubic.H3))
	cubic.B = (-45.0 / (PI * cubic.H4))
	cubic.C = (90.0) / (PI * cubic.H5)
	return cubic
}

//**------------------Standard----------------------//
func (K Cubic) F(x float32) float32 {
	if x >= K.H_ {
		return float32(0.0)
	}
	q := 1.0 - (x*x)/(K.H_*K.H_)
	return K.A * q * q
}

func (K Cubic) W0() float32 {
	return K.F(0)
}

func (K Cubic) H() float32 {
	return K.H_
}

func (K Cubic) H0() float32 {
	return K.H_
}

//O1D - 1st order Differential
func (K Cubic) O1D(x float32) float32 {
	if x >= K.H_ {
		return 0.0
	}
	q := 1.0 - x/K.H_
	return K.B * q * q
}

//O2D - Returns 2nd Order Differential
func (K Cubic) O2D(x float32) float32 {
	if x > K.H_ {
		return 0.0
	}

	q := 1.0 - x/K.H_
	return K.C * q

}

//Grad finds the kernel gradient
func (K Cubic) Grad(x float32, dir V.Vec) V.Vec {
	return V.Scale(dir, -K.O1D(x))
}
