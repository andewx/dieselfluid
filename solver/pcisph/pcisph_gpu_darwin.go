package pcisph

import "github.com/andewx/dieselfluid/model/sph"
import "github.com/andewx/dieselfluid/compute/gpu"
import "github.com/andewx/dieselfluid/compute"
import "fmt"

const LOCAL_GROUP_SIZE = 4
const DIM = 64

type TempPCI struct {
	vel  [3]float32
	pos  [3]float32
	pres []float32
}

type GPUPredictorCorrector struct {
	system         sph.SPH
	gpu_compute    gpu.ComputeGPU
	log            string
	temp_particles []TempPCI
}

/*
 Please note that a GPU Compute Evolution uses shared buffer object in OpenGL for
 memory transfer. These use binding points in {1-6} so other rendering shaders will
 need to place buffers later in memory
*/
func New_GPUPredictorCorrector(sph sph.SPH, opencl gpu.OpenCL) (GPUPredictorCorrector, error) {

	mGPU := GPUPredictorCorrector{}
	mGPU.system = sph

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

	descriptor := compute.Descriptor{work_group, local_group, size}
	mGPU.temp_particles = make([]TempPCI, size)

	//Setup compute worload definitions
	mGPU.gpu_compute = gpu.New_ComputeGPU(descriptor, opencl)
	mGPU.log += "Initialized New_ComputeGPU()"

	//Pre-Arrange Buffers
	ints := []int{mGPU.system.N(), len(mGPU.system.Particles()) - mGPU.system.N(), mGPU.system.Field().GetSampler().GetBuckets(), mGPU.system.Field().GetSampler().BucketSize()}
	floats := []float32{mGPU.system.CFL(), mGPU.system.Field().Mass(), mGPU.system.Delta(), mGPU.system.MaxV(), mGPU.system.Field().GetKernelLength()}
	hash_buffer := mGPU.system.Field().GetSampler().GetData1D()
	hash_buffer_len := len(hash_buffer)
	random_project_vectors := mGPU.system.Field().GetSampler().GetVectors()
	particle_bytes := len(mGPU.system.Field().Particles())*(3*(3*4)) + (2 * 4)
	hash_bytes := hash_buffer_len * 4
	vector_bytes := len(random_project_vectors) * 4
	temp_bytes := 2*3*4 + 4

	/*Commit buffers*/
	var buffer_err error
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(mGPU.system.Field().Particles(), particle_bytes, "particles")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassIntBuffer(ints, "ints")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(floats, "floats")
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
	if err := k1.SetArgBuffer(0, mContext.Buffers["particles"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(1, mContext.Buffers["ints"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(2, mContext.Buffers["floats"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(3, mContext.Buffers["hash"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(4, mContext.Buffers["vecs"]); err != nil {
		return mGPU, err
	}

	k2 := mContext.Kernels["predict_correct"]
	if err := k2.SetArgBuffer(0, mContext.Buffers["particles"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(1, mContext.Buffers["ints"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(2, mContext.Buffers["floats"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(3, mContext.Buffers["hash"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(4, mContext.Buffers["vecs"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(5, mContext.Buffers["temp"]); err != nil {
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
func (m GPUPredictorCorrector) Run() error {
	var err error
	m.system.CFL()
	m.system.CacheIncr()
	err = m.gpu_compute.Queue("compute_density")
	if err != nil {
		fmt.Printf("Error adding kernel to execution path. For work group size errors users may need to augment the size of the work group dimensions to match the preferred sizes")
		return err
	}
	err = m.gpu_compute.Queue("predict_correct")
	if err != nil {
		return err
	}
	fmt.Printf("Executed Kernels\n")
	return nil
}
