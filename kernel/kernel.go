package kernel

import V "github.com/andewx/dieselfluid/math/math64"

type Kernel interface {
	F(x float64) float64
	O1D(x float64) float64
	O2D(x float64) float64
	H() float64  //Smoothing Length
	H0() float64 //Adaptive Smoothing Length
	Adjust(ratio float64) float64
	Grad(x float64, dir V.Vec)
	W0() float64
}
