package kernel

import V "github.com/andewx/dieselfluid/math/vector"

type Kernel interface {
	F(x float32) float32
	O1D(x float32) float32
	O2D(x float32) float32
	H() float32  //Adpative Smoothing Length
	H0() float32 //Adaptive Smoothing Length
	Grad(x float32, dir V.Vec) V.Vec
	W0() float32
}
