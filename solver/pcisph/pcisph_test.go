package pcisph

import (
	"testing"

	"github.com/andewx/dieselfluid/model/sph"
)

/*OpenCL State Must Remain*/
func TestOpenCompute(t *testing.T) {
	//Sph fluid generates boundary particles for GL Buffers if needed
	sph := sph.Init(float32(1.0), []float32{0, 0, 0}, nil, 16, true)
	solver := NewPCIMethod(&sph, 0, 0)
	message := make(chan string)
	go solver.Run(message, false, nil)
	message <- "QUIT"
}
