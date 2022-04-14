package light

import "github.com/andewx/dieselfluid/math/mgl"

type Probe struct { //Use for light probes
	//Lights   LightRig
	Size     float32
	Dims     [3]int
	Position mgl.Vec
	Samples  int
}
