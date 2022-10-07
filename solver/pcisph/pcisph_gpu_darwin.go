package pcisph

import (
	"fmt"

	"github.com/andewx/dieselfluid/compute"
	"github.com/andewx/dieselfluid/compute/gpu"
	"github.com/andewx/dieselfluid/model/sph"
)

const LOCAL_GROUP_SIZE = 4
const DIM = 64
const EXIT = -1
const PROCEED = 1

type TempPCI struct {
	vel  [3]float32
	pos  [3]float32
	pres []float32
}

type GPUPredictorCorrector struct {
	system             sph.SPH
	gpu_compute        *gpu.ComputeGPU
	log                string
	temp_particles     []TempPCI
	gl_position_buffer bool
}

/*
 Please note that a GPU Compute Evolution uses shared buffer object in OpenGL for
 memory transfer. These use binding points in {1-6} so other rendering shaders will
 need to place buffers later in memory. If you are calling with a gl_position_buffer argument set to
 true then you must have initialized an obtained a valid open gl buffer uint id for the fluid positions
*/
func New_GPUPredictorCorrector(computeGPU *gpu.ComputeGPU, sph sph.SPH, opencl *gpu.OpenCL, gl_position_buffer bool) (GPUPredictorCorrector, error) {

	mGPU := GPUPredictorCorrector{}
	mGPU.system = sph
	mGPU.gl_position_buffer = gl_position_buffer
	mGPU.gpu_compute = computeGPU

	//Compute Description validity - fails when parameters are not set correctly
	m_n := DIM % LOCAL_GROUP_SIZE
	if m_n != 0 {
		err := fmt.Errorf("Invalid Local Group")
		return mGPU, err
	}

	//Compute Group Description (X,Y,Z) compute parameters
	work_dim := int(DIM / LOCAL_GROUP_SIZE)
	work_group := []int{work_dim, work_dim, work_dim}
	local_group := []int{LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE}
	size := sph.N()

	descriptor := compute.Descriptor{Work: work_group, Local: local_group, Size: size}
	mGPU.temp_particles = make([]TempPCI, size)

	//Setup compute worload definitions
	mGPU.gpu_compute.SetDescriptor(descriptor)
	mGPU.log += "Initialized New_ComputeGPU()"

	//Pre-Arrange Buffers
	field := mGPU.system.Field()

	ints := []int{field.Particles.N(), field.Particles.Total() - field.Particles.N(), field.GetSampler().GetBuckets(), field.GetSampler().BucketSize()}
	floats := []float32{mGPU.system.CFL(), field.Mass(), mGPU.system.Delta(), mGPU.system.MaxV(), field.GetKernelLength()}
	hash_buffer := field.GetSampler().GetData1D()
	hash_buffer_len := len(hash_buffer)
	random_project_vectors := field.GetSampler().GetVectors()

	hash_bytes := hash_buffer_len * 4
	vector_bytes := len(random_project_vectors) * 4
	temp_bytes := 2*3*4 + 4

	/*Commit buffers*/
	var buffer_err error

	buffer_err = mGPU.gpu_compute.PassFloatBuffer(field.Particles.Positions(), "positions")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(field.Particles.Velocities(), "velocities")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(field.Particles.Forces(), "forces")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(field.Particles.Densities(), "densities")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(field.Particles.Pressures(), "pressures")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassIntBuffer(ints, "sizes")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(floats, "fluid_data")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(hash_buffer, hash_bytes, "hash")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(random_project_vectors, vector_bytes, "vecs")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(mGPU.temp_particles, temp_bytes, "temp")
	if buffer_err != nil {
		return mGPU, buffer_err
	}

	//Pass the kernel arguments
	mContext := mGPU.gpu_compute.Context()

	//Get Arument Info
	for i := 0; i < 5; i++ {
		name, err := mContext.Kernels["compute_density"].ArgName(i)
		if err != nil {
			fmt.Printf("Kernel Info for name failed: %v", err)
		} else {
			fmt.Printf("Kernel arg %d, %s\n", i, name)
		}
	}

	//Get Arument Info
	for i := 0; i < 6; i++ {
		name, err := mContext.Kernels["predict_correct"].ArgName(i)
		if err != nil {
			fmt.Printf("Kernel Info for name failed: %v", err)
		} else {
			fmt.Printf("Kernel arg %d, %s\n", i, name)
		}
	}

	k1 := mContext.Kernels["compute_density"]
	if err := k1.SetArgBuffer(0, mContext.Buffers["positions"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(1, mContext.Buffers["velocities"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(2, mContext.Buffers["forces"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(3, mContext.Buffers["densities"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(4, mContext.Buffers["pressures"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(5, mContext.Buffers["sizes"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(6, mContext.Buffers["floats"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(7, mContext.Buffers["sampler_data"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(8, mContext.Buffers["sampler_vecs"]); err != nil {
		return mGPU, err
	}

	k2 := mContext.Kernels["predict_correct"]
	if err := k2.SetArgBuffer(0, mContext.Buffers["positions"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(1, mContext.Buffers["velocities"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(2, mContext.Buffers["forces"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(3, mContext.Buffers["densities"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(4, mContext.Buffers["pressures"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(5, mContext.Buffers["sizes"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(6, mContext.Buffers["floats"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(7, mContext.Buffers["sampler_data"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(8, mContext.Buffers["sampler_vecs"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(9, mContext.Buffers["temps"]); err != nil {
		return mGPU, err
	}

	local, err := k2.WorkGroupSize(mContext.Device)
	if err != nil {
		return mGPU, fmt.Errorf("WorkGroupSize failed: %+v\n", err)
	}

	fmt.Printf("Work group size: %d\n", local)
	sizeb, _ := k2.PreferredWorkGroupSizeMultiple(nil)
	fmt.Printf("Preferred Work Group Size Multiple: %d\n", sizeb)

	return mGPU, nil

}

/* Executes one full compute cycle for PCI PSH compute shader which uses 2 shader kernels */
func (m GPUPredictorCorrector) Run(message chan string) error {
	var err error

	done := false
	for done != true {
		ack := <-message
		if ack == "QUIT" {
			message <- "QUIT"
			return nil
		}
		if ack != "PROCEED" {
			return fmt.Errorf("Error unrecognized message")
		}

		m.system.CFL()
		m.system.CacheIncr()
		err = m.gpu_compute.Queue("compute_density")
		if err != nil {
			err = fmt.Errorf("Error adding kernel to execution path. For work group size errors users may need to augment the size of the work group dimensions to match the preferred sizes \n%v", err)
			return err
		}
		err = m.gpu_compute.Queue("predict_correct")
		if err != nil {
			return err
		}

		message <- "PROCEED"
	}

	fmt.Printf("Executed Kernels\n")
	return nil
}
