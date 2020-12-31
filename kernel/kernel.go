package kernel

import V "github.com/andewx/dieselfluid/math/math64"

const PI = 3.141592653589
const PI2 = PI * 2
const SQRT2PI = 2.50662827463
const SQRPI = 5.5860525258
const A = (4*PI*PI - 30)
const A0 = (2*PI*PI - 15)

type Kernel interface{
  func (k Kernel)F(x float64)float64
  func (k Kernel)O1D(x float64)float64
  func (k Kernel)O2D(x float64)float64
  func (k Kernel)H()float64 //Smoothing Length
  func (k Kernel)H0()float64 //Adaptive Smoothing Length
  func (k Kernel)Adjust(ratio float64)float64
  func (k Kernel)Grad(x float64, dir V.Vec)
  func (k Kernel)W0()float64
}
