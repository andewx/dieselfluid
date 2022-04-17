package integral

import (
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/sampler"
)

//General Purpose Integral Evaluation Interface which operates vector parameter
//return slice sets using functional queues. Functions and Integrators are available
//To this queue. The integrator
type Integral interface {
	Integrator() *Integrator
	Evaluate(vector.Vec) vector.Vec
}

type Integrator interface {
	Sampler() *sampler.Sampler
	Bounds() vector.Vec
}
