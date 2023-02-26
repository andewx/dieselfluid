package pcisph

import (
	"fmt"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/compute/gpu"
	"github.com/andewx/dieselfluid/model/sph"
)

const LOCAL_GROUP_SIZE = 4
const DIM = 16
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
	gl_position_buffer uint32
}

/*
 Please note that a GPU Compute Evolution uses shared buffer object in OpenGL for
 memory transfer. These use binding points in {1-6} so other rendering shaders will
 need to place buffers later in memory. If you are calling with a gl_position_buffer argument set to
 true then you must have initialized an obtained a valid open gl buffer uint id for the fluid positions
*/
func New_GPUPredictorCorrector(computeGPU *gpu.ComputeGPU, sph sph.SPH, opencl *gpu.OpenCL, gl_position_buffer uint32) (GPUPredictorCorrector, error) {

	mGPU := GPUPredictorCorrector{}
	mGPU.system = sph
	mGPU.gl_position_buffer = gl_position_buffer
	mGPU.gpu_compute = computeGPU
	mGPU.gl_position_buffer = gl_position_buffer

	//Compute Description validity - fails when parameters are not set correctly
	m_n := DIM % LOCAL_GROUP_SIZE
	if m_n != 0 {
		err := fmt.Errorf("Invalid Local Group")
		return mGPU, err
	}

	size := sph.N()
	mGPU.temp_particles = make([]TempPCI, size)

	//Setup compute worload definitions
	mGPU.log += "Initialized New_ComputeGPU()"

	//Pre-Arrange Buffers
	field := mGPU.system.Field()

	ints := []int{field.Particles.N(), field.Particles.Total() - field.Particles.N(), field.GetSampler().GetBuckets(), field.GetSampler().BucketSize()}
	floats := []float32{mGPU.system.CFL(), field.Mass(), mGPU.system.Delta(), mGPU.system.MaxV(), field.GetKernelLength()}
	hash_buffer := field.GetSampler().GetData1D()
	random_project_vectors := field.GetSampler().GetVectors()
	temp_bytes := 2*3*4 + 4
	parray := sph.Particles()

	mGPU.gpu_compute.RegisterBuffer(parray.Total()*3*4, 0, "positions")
	mGPU.gpu_compute.RegisterBuffer(parray.N()*3*4, 0, "velocities")
	mGPU.gpu_compute.RegisterBuffer(parray.N()*3*4, 0, "forces")
	mGPU.gpu_compute.RegisterBuffer(parray.N()*4, 0, "densities")
	mGPU.gpu_compute.RegisterBuffer(parray.N()*4, 0, "pressures")
	mGPU.gpu_compute.RegisterBuffer(len(ints)*4, 0, "sizes")
	mGPU.gpu_compute.RegisterBuffer(len(floats)*4, 0, "floats")
	mGPU.gpu_compute.RegisterBuffer(len(hash_buffer)*4, 0, "sampler")
	mGPU.gpu_compute.RegisterBuffer(len(random_project_vectors)*4, 0, "vecs")
	mGPU.gpu_compute.RegisterBuffer(parray.N()*7*4, 0, "temps")

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
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(ints, len(ints)*4, "sizes")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(floats, "floats")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(hash_buffer, len(hash_buffer)*4, "sampler")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassFloatBuffer(random_project_vectors, "vecs")
	if buffer_err != nil {
		return mGPU, buffer_err
	}
	buffer_err = mGPU.gpu_compute.PassLayoutBuffer(mGPU.temp_particles, temp_bytes, "temps")
	if buffer_err != nil {
		return mGPU, buffer_err
	}

	//Pass the kernel arguments
	if merr := mGPU.gpu_compute.AddSourceFile(common.ProjectRelativePath("data/shaders/opencl/pcisph/pci_density.c")); merr != nil {
		return mGPU, merr
	}
	if merr := mGPU.gpu_compute.AddSourceFile(common.ProjectRelativePath("data/shaders/opencl/pcisph/pci_predict.c")); merr != nil {
		return mGPU, merr
	}
	if merr := mGPU.gpu_compute.BuildProgram(common.ProjectRelativePath("data/shaders/opencl/include")); merr != nil {
		return mGPU, merr
	}

	if merr := mGPU.gpu_compute.RegisterKernel("compute_density"); merr != true {
		return mGPU, fmt.Errorf("Register kernel compute density failed\n")
	}

	if merr := mGPU.gpu_compute.RegisterKernel("predict_correct"); merr != true {
		return mGPU, fmt.Errorf("Register kernel predict_correct failed\n")
	}

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
	if err := k1.SetArgBuffer(7, mContext.Buffers["sampler"]); err != nil {
		return mGPU, err
	}
	if err := k1.SetArgBuffer(8, mContext.Buffers["vecs"]); err != nil {
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
	if err := k2.SetArgBuffer(7, mContext.Buffers["sampler"]); err != nil {
		return mGPU, err
	}
	if err := k2.SetArgBuffer(8, mContext.Buffers["vecs"]); err != nil {
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
	fmt.Printf(mGPU.MemoryRequirements())

	return mGPU, nil

}

func (m GPUPredictorCorrector) MemoryRequirements() string {
	particles := m.system.Field().Particles.N()
	boundary := m.system.Field().Particles.Total() - particles
	total := m.system.Field().Particles.Total()
	totalSizeKB := float64(((total * 4 * 3) + (particles * 4 * 3 * 2) + (2*particles*4)/1024))
	sizeMB := totalSizeKB * 0.001
	return fmt.Sprintf("Fluid GPU PCI (Fluid Particles[%d]  Boundary[%d]\nAllocated %.2fkB (%.2fMB)\n\n", particles, boundary, totalSizeKB, sizeMB)

}

/* Executes one full compute cycle for PCI PSH compute shader which uses 2 shader kernels */
func (m GPUPredictorCorrector) Run(message *chan string) error {
	mContext := m.gpu_compute.Context()
	done := false
	for !done {
		m.system.CFL()
		m.system.CacheIncr()
		description := m.gpu_compute.Descriptor()
		if _, err := mContext.Queue.EnqueueNDRangeKernel(mContext.Kernels["compute_density"], nil, description.Work, description.Local, nil); err != nil {
			fmt.Printf("Enqueue Compute desnity failed: %s\n", err.Error())
			return err
		}

		if err := mContext.Queue.Finish(); err != nil {
			fmt.Printf("Finish failed: %s\n", err.Error())
			return err
		}

		if _, err := mContext.Queue.EnqueueNDRangeKernel(mContext.Kernels["predict_correct"], nil, m.gpu_compute.Descriptor().Work, m.gpu_compute.Descriptor().Local, nil); err != nil {
			fmt.Printf("Enqueue Predict Correct failed: %s\n", err.Error())
			return err
		}

		if err := mContext.Queue.Finish(); err != nil {
			fmt.Printf("Finish PC failed failed: %s\n", err.Error())
			return err
		}

		positions := m.system.Field().Particles.Positions()
		_, err := mContext.Queue.EnqueueReadBufferFloat32(mContext.Buffers["positions"], true, 0, positions, nil)
		if err == nil {
			*message <- "CL_REFRESH"
		} else {
			fmt.Printf("Readbuffer error\n")
		}
		done = false
	}
	return nil
}
