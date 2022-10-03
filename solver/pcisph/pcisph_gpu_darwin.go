package pcisph

import "github.com/andewx/dieselfluid/model/sph"
import "github.com/andewx/dieselfluid/compute/gpu"
import "github.com/andewx/dieselfluid/compute"
import "github.com/andewx/dieselfluid/geom"
import "github.com/andewx/dieselfluid/math/vector"
import "log"
import "fmt"

const LOCAL_GROUP_SIZE = 4

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
func New_GPUPredictorCorrector(n int, collider []geom.Collider) (GPUPredictorCorrector, error) {
	var err error
	mGPU := GPUPredictorCorrector{}
	pci := sph.Init(float32(1.0), vector.Vec{}, collider, n, true)

	//Compute Description validity - fails when parameters are not set correctly
	m_n := n % LOCAL_GROUP_SIZE
	if m_n != 0 {
		err = fmt.Errorf("Invalid Local Group")
		fmt.Printf("New_GPUPredictorCorrector() - Method Failed\n")
		err = nil
		log.Fatalf("Parameter n int (%d) must be a factor of the LOCAL_GROUP_SIZE (%d) this ensures that the work groups allocated to the GPU Kernel will completely be able to cover the particle index space\nFuture updates may alleviate this issue however for now the user must be careful when selecting the particle cubic dimensions for GPU based parallel workloads\n", n, LOCAL_GROUP_SIZE)
	}

	//Compute Group Description (X,Y,Z) compute parameters
	work_dim := int(n / LOCAL_GROUP_SIZE)
	work_group := []int{work_dim, work_dim, work_dim}
	local_group := []int{LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE, LOCAL_GROUP_SIZE}
	size := pci.N()
	descriptor := compute.Descriptor{work_group, local_group, size}
	mGPU.temp_particles = make([]TempPCI, size)

	//Setup compute worload definitions
	mGPU.system = pci
	mGPU.gpu_compute = gpu.New_ComputeGPU(descriptor)
	mGPU.log += "Initialized New_ComputeGPU()"

	//Setup and initialize the gpu comute unit
	mGPU.gpu_compute.Setup(false)

	//Add the GLSL Source Files
	mGPU.gpu_compute.AddSourceFile("pcisph_kern0.ocl")
	mGPU.gpu_compute.AddSourceFile("pcisph_kern1.ocl")

	err = mGPU.gpu_compute.BuildProgram()

	//Link the programs and kernels
	mGPU.gpu_compute.RegisterKernel("compute_density")
	mGPU.gpu_compute.RegisterKernel("predict_correct")

	//Pre-Arrange
	ints := []int{mGPU.system.N(), len(mGPU.system.Particles()) - mGPU.system.N(), mGPU.system.Field().GetSampler().GetBuckets(), mGPU.system.Field().GetSampler().BucketSize()}
	floats := []float32{mGPU.system.CFL(), mGPU.system.Field().Mass(), mGPU.system.Delta(), mGPU.system.MaxV(), mGPU.system.Field().GetKernelLength()}
	hash_buffer := mGPU.system.Field().GetSampler().GetData()
	hash_buffer_len := ints[2] * ints[3]
	random_project_vectors := mGPU.system.Field().GetSampler().GetVectors()
	particle_bytes := len(mGPU.system.Field().Particles())*(3*(3*4)) + (2 * 4)
	hash_bytes := hash_buffer_len * 4
	vector_bytes := 8 * 3 * 4
	temp_bytes := 2*3*4 + 4

	//Creates
	mGPU.gpu_compute.RegisterBuffer(len(mGPU.system.Field().Particles()), 1, "particles")
	mGPU.gpu_compute.RegisterBuffer(len(ints), 2, "intdata")
	mGPU.gpu_compute.RegisterBuffer(len(floats), 3, "floatdata")
	mGPU.gpu_compute.RegisterBuffer(hash_buffer_len, 4, "nn_hash")
	mGPU.gpu_compute.RegisterBuffer(len(random_project_vectors), 5, "hash_vectors")
	mGPU.gpu_compute.RegisterBuffer(len(mGPU.temp_particles), 6, "temp_particles")

	mGPU.gpu_compute.PassLayoutBuffer(mGPU.system.Field().Particles(), particle_bytes, "particles")
	mGPU.gpu_compute.PassIntBuffer(ints, "intdata")
	mGPU.gpu_compute.PassFloatBuffer(floats, "floatdata")
	mGPU.gpu_compute.PassLayoutBuffer(hash_buffer, hash_bytes, "nn_hash")
	mGPU.gpu_compute.PassLayoutBuffer(random_project_vectors, vector_bytes, "hash_vectors")
	mGPU.gpu_compute.PassLayoutBuffer(mGPU.temp_particles, temp_bytes, "temp_particles")

	return mGPU, err

}

/* Executes one full compute cycle for PCI PSH compute shader which uses 2 shader kernels */
func (m GPUPredictorCorrector) Run() {
	m.system.CFL()
	m.system.CacheIncr()
	m.gpu_compute.Queue("compute_density")
	m.gpu_compute.Queue("predict_correct")
	fmt.Printf("Executed Kernels\n")
}
