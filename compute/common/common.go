package common

import "github.com/andewx/dieselfluid/math/mgl"

//Generalized callback function interface with []float32 return parameter
type ComputeFunction struct {
	Evaluate func(mgl.Vec) mgl.Vec
}

//Integral Computation Function Map with object parameter
type Evaluatable interface {
	Evaluate(mgl.Vec) mgl.Vec
}
