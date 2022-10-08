package fluidapp

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/compute"
	"github.com/andewx/dieselfluid/compute/gpu"
	"github.com/andewx/dieselfluid/model/sph"
	"github.com/andewx/dieselfluid/render"
	"github.com/andewx/dieselfluid/solver/pcisph"
)

const (
	FLUID_DIM   = 32
	N_PARTICLES = FLUID_DIM * FLUID_DIM * FLUID_DIM
	LOCAL_GROUP = 4
)

func TestFluid(t *testing.T) {

	if FLUID_DIM%LOCAL_GROUP != 0 {
		t.Errorf("Fluid dimensions must be multiple of the local group size\n")
		return
	}

	//OpenGL Setup
	runtime.LockOSThread()
	Sys, _ := render.Init(common.ProjectRelativePath("data/meshes/materialsphere/MaterialSphere.gltf"))

	if err := Sys.Init(1024, 720, common.ProjectRelativePath("render"), true); err != nil {
		t.Error(err)
	}

	if err := Sys.Meshes(); err != nil {
		t.Error(err)
	}

	if err := Sys.CompileLink(); err != nil {
		t.Error(err)
	}

	//Sph fluid generates boundary particles for GL Buffers if needed
	sph := sph.Init(float32(1.0), []float32{0, 0, 0}, Sys.GetColliderMeshes(), FLUID_DIM, true)
	pos := sph.Field().Particles.Positions()
	particleVBO, _ := Sys.RegisterParticleSystem(pos, 13)
	Sys.MyRenderer.HasParticleSystem = true
	Sys.MyRenderer.NumParticles = int32(sph.Field().Particles.N())

	//OpenCL GPU Compute Setup - Not on a shared context
	opencl := &gpu.OpenCL{}

	work_dim := int(FLUID_DIM / LOCAL_GROUP)
	work_group := []int{work_dim, work_dim, work_dim}
	local_group := []int{LOCAL_GROUP, LOCAL_GROUP, LOCAL_GROUP}
	descriptor := compute.Descriptor{Work: work_group, Local: local_group, Size: FLUID_DIM}

	gpu.InitOpenCL(opencl)
	compute_gpu := &gpu.ComputeGPU{}
	compute_gpu = gpu.New_ComputeGPU(compute_gpu, &descriptor, opencl)
	solver, err := pcisph.New_GPUPredictorCorrector(compute_gpu, sph, opencl, particleVBO)
	if err != nil {
		t.Errorf("Failed pcisph creation %v", err)
		return
	}

	messenger := make(chan string)
	go solver.Run(messenger)
	if err := Sys.Run(messenger); err != nil {
		t.Error(err)
	}

	fmt.Printf("Executed Process\n")

}
