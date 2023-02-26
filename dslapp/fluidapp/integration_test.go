package fluidapp

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/model/sph"
	"github.com/andewx/dieselfluid/render"
	"github.com/andewx/dieselfluid/solver/pcisph"
)

const (
	FLUID_DIM   = 16
	N_PARTICLES = FLUID_DIM * FLUID_DIM * FLUID_DIM
	LOCAL_GROUP = 4
)

func TestFluid(t *testing.T) {

	//OpenGL Setup
	runtime.LockOSThread()
	Sys, _ := render.Init(common.ProjectRelativePath("data/meshes/materialcube/materialcube.gltf"))

	if err := Sys.Init(1920, 1080, common.ProjectRelativePath("render"), true); err != nil {
		t.Error(err)
	}

	if err := Sys.Meshes(); err != nil {
		t.Error(err)
	}

	if err := Sys.CompileLink(); err != nil {
		t.Error(err)
	}

	//Sph fluid generates boundary particles for GL Buffers if needed
	sph := sph.Init(float32(1.0), []float32{0, 0, 0}, nil, FLUID_DIM, true)
	pos := sph.Field().Particles.Positions()
	vbo, vao := Sys.RegisterParticleSystem(pos, 13)
	Sys.MyRenderer.HasParticleSystem = true
	Sys.MyRenderer.NumParticles = int32(sph.Field().Particles.N())
	solver := pcisph.NewPCIMethod(&sph, vbo, vao)
	message := make(chan string)
	go sph.Field().GetSampler().Run(message)
	go solver.Run(message, true, Sys)
	if err := Sys.Run(message, &sph); err != nil {
		t.Error(err)
	}

	fmt.Printf("Executed Process\n")

}
