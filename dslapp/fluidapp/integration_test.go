package fluidapp

import (
	"runtime"
	"testing"

	"unsafe"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/compute"
	"github.com/andewx/dieselfluid/compute/gpu"
	"github.com/andewx/dieselfluid/geom/mesh"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model/sph"
	"github.com/andewx/dieselfluid/render"
	"github.com/andewx/dieselfluid/solver/pcisph"
)

const (
	FLUID_DIM   = 64
	N_PARTICLES = FLUID_DIM * FLUID_DIM * FLUID_DIM
	LOCAL_GROUP = 4
)

func byteSliceToFloat32Slice(src []byte) []float32 {
	if len(src) == 0 {
		return nil
	}

	l := len(src) / 4
	ptr := unsafe.Pointer(&src[0])
	// It is important to keep in mind that the Go garbage collector
	// will not interact with this data, and that if src if freed,
	// the behavior of any Go code using the slice is nondeterministic.
	// Reference: https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	return (*[1 << 26]float32)((*[1 << 26]float32)(ptr))[:l:l]
}

func TestFluid(t *testing.T) {

	if FLUID_DIM%LOCAL_GROUP != 0 {
		t.Errorf("Fluid dimensions must be multiple of the local group size\n")
		return
	}

	runtime.LockOSThread()

	//OpenGL Setup
	Sys, _ := render.Init(common.ProjectRelativePath("data/meshes/plane/plane.gltf"))
	if err := Sys.Init(1024, 720, common.ProjectRelativePath("render"), false); err != nil {
		t.Error(err)
	}

	colliderMeshes := make([]mesh.Mesh, len(Sys.MeshEntities))
	for i, meshy := range Sys.MeshEntities {
		vertices := byteSliceToFloat32Slice(meshy.Mesh.Vertices)
		tris := make([]vector.Vec, len(vertices)/3)
		for i := 0; i < len(tris); i++ {
			x := i * 3
			tris[i] = vector.Vec{vertices[x], vertices[x+1], vertices[x+2]}
		}
		m_mesh := mesh.InitMesh(tris, []float32{0, 0, 0})
		colliderMeshes[i] = m_mesh
	}

	//Sph fluid generates boundary particles for GL Buffers if needed
	sph := sph.Init(float32(1.0), []float32{0, 0, 0}, colliderMeshes, FLUID_DIM, true)

	if err := Sys.CompileLink(); err != nil {
		t.Error(err)
	}
	if err := Sys.Meshes(); err != nil {
		t.Error(err)
	}

	buffer_id, _ := Sys.RegisterParticleSystem(sph.Field().Particles.Positions(), 8)

	//OpenCL GPU Compute Setup
	opencl := &gpu.OpenCL{}
	descriptor := compute.Descriptor{Work: []int{}, Local: []int{}, Size: LOCAL_GROUP}
	gpu.InitOpenCL(opencl)
	compute_gpu := &gpu.ComputeGPU{}
	compute_gpu = gpu.New_ComputeGPU(compute_gpu, descriptor, opencl)

	sizes := []int{sph.Particles().N(), sph.Particles().N() - sph.Particles().Total(), sph.Field().GetSampler().GetBuckets(), sph.Field().GetSampler().BucketSize()}
	fluid_data := []float32{sph.CFL(), sph.Delta(), sph.Field().Mass(), sph.MaxV(), sph.Field().Kernel().H0()}
	sampler_data := sph.Field().GetSampler().GetData()
	sampler_vecs := sph.Field().GetSampler().GetVectors()

	//Initiate Buffers
	if err := compute_gpu.RegisterGLBuffer(buffer_id, sph.Particles().Total()*3*4, "positions"); err != nil {
		t.Error(err)
		return
	}
	compute_gpu.RegisterBuffer(sph.Particles().N()*3*4, 0, "velocities")
	compute_gpu.RegisterBuffer(sph.Particles().N()*3*4, 0, "forces")
	compute_gpu.RegisterBuffer(sph.Particles().N()*4, 0, "densities")
	compute_gpu.RegisterBuffer(sph.Particles().N()*4, 0, "pressures")
	compute_gpu.RegisterBuffer(len(sizes)*4, 0, "sizes")
	compute_gpu.RegisterBuffer(len(fluid_data)*4, 0, "fluid_data")
	compute_gpu.RegisterBuffer(len(sampler_data)*4, 0, "sampler")
	compute_gpu.RegisterBuffer(len(sampler_vecs)*4, 0, "vecs")
	compute_gpu.RegisterBuffer(sph.Particles().N()*7*4, 0, "temps")

	//Initiate GPU Solver Context
	solver, err := pcisph.New_GPUPredictorCorrector(compute_gpu, sph, opencl, true)
	messager := make(chan string)

	go solver.Run(messager)

	if err != nil {
		t.Errorf("Failed pcisph creation %v", err)
		return
	}

	if err := Sys.Run(messager); err != nil {
		t.Error(err)
	}

}
