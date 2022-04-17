package common

import "github.com/andewx/dieselfluid/math/vector"

//Generalized callback function interface with []float32 return parameter
type ComputeFunction struct {
	Evaluate func(vector.Vec) vector.Vec
}

//Integral Computation Function Map with object parameter
type Evaluatable interface {
	Evaluate(vector.Vec) vector.Vec
}
