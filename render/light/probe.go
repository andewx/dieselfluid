package light

import "dslfluid.com/dsl/math/mgl"

type Probe struct { //Use for light probes
	Lights   LightRig
	Size     float32
	Dims     [3]int
	Position mgl.Vec
	Samples  int
}
