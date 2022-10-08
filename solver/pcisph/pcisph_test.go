package pcisph

import (
	"testing"

	"github.com/andewx/dieselfluid/compute"
	"github.com/andewx/dieselfluid/compute/gpu"
	"github.com/andewx/dieselfluid/model/sph"
)

/*OpenCL State Must Remain*/
func TestOpenCompute(t *testing.T) {

	//Sph fluid generates boundary particles for GL Buffers if needed
	sph := sph.Init(float32(1.0), []float32{0, 0, 0}, nil, 64, true)

	//OpenCL GPU Compute Setup - Not on a shared context
	opencl := &gpu.OpenCL{}

	work_dim := int(DIM / LOCAL_GROUP_SIZE)
	work_group := []int{work_dim, work_dim, work_dim}
	local_group := []int{LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE}
	descriptor := compute.Descriptor{Work: work_group, Local: local_group, Size: DIM}

	gpu.InitOpenCL(opencl)
	compute_gpu := &gpu.ComputeGPU{}
	compute_gpu = gpu.New_ComputeGPU(compute_gpu, &descriptor, opencl)
	solver, err := New_GPUPredictorCorrector(compute_gpu, sph, opencl, 0)
	if err != nil {
		t.Errorf("Failed pcisph creation %v", err)
		return
	}

	messenger := make(chan string)
	go solver.Run(messenger)
	messenger <- "QUIT"
}
