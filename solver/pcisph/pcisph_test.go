package pcisph

import "testing"
import "github.com/jgillich/go-opencl/cl"
import "github.com/andewx/dieselfluid/compute/gpu"
import "github.com/andewx/dieselfluid/math/vector"
import "github.com/andewx/dieselfluid/common"
import "github.com/andewx/dieselfluid/model/sph"
import "io/ioutil"
import "fmt"
import "log"

/*OpenCL State Must Remain*/
func TestOpenCompute(t *testing.T) {

	/*Initiate OpenCL State*/
	opencl := gpu.OpenCL{}
	opencl.Buffers = make(map[string]*cl.MemObject, 10)
	opencl.Kernels = make(map[string]*cl.Kernel, 10)

	platforms, err := cl.GetPlatforms()
	if err != nil {
		log.Fatalf("Failed to get platforms: %+v\n", err)
	}
	for i, p := range platforms {
		fmt.Printf("Platform %d:\n", i)
		fmt.Printf("  Name: %s\n", p.Name())
		fmt.Printf("  Vendor: %s\n", p.Vendor())
		fmt.Printf("  Profile: %s\n", p.Profile())
		fmt.Printf("  Version: %s\n", p.Version())
		fmt.Printf("  Extensions: %s\n", p.Extensions())
	}
	platform := platforms[0]

	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		log.Fatalf("Failed to get devices: %+v\n", err)
	}
	if len(devices) == 0 {
		log.Fatalf("GetDevices returned no devices")
	}
	deviceIndex := -1
	for i, d := range devices {
		if deviceIndex < 0 && d.Type() == cl.DeviceTypeGPU {
			deviceIndex = i
			fmt.Printf("Device %d (%s): %s Selected!\n", i, d.Type(), d.Name())
		}
		fmt.Printf("Device %d (%s): %s\n", i, d.Type(), d.Name())
		fmt.Printf("  OpenCL C Version: %s\n", d.OpenCLCVersion())
		fmt.Printf("  Profile: %s\n", d.Profile())
		fmt.Printf("  Vendor: %s\n", d.Vendor())
		fmt.Printf("  Version: %s\n", d.Version())
		fmt.Printf("  Max Samplers: %d\n", d.MaxSamplers())
		fmt.Printf("  Max Work Group Size: %d\n", d.MaxWorkGroupSize())
		fmt.Printf("  Max Work Item Dimensions: %d\n", d.MaxWorkItemDimensions())
		fmt.Printf("  Max Work Item Sizes: %d\n", d.MaxWorkItemSizes())
		fmt.Printf("  Global Memory Size: %d MB\n", d.GlobalMemSize()/(1024*1024))
		fmt.Printf("  Max Compute Units: %d\n", d.MaxComputeUnits())
	}
	if deviceIndex < 0 {
		deviceIndex = 0
	}
	device := devices[deviceIndex]
	fmt.Printf("Using device %d\n", deviceIndex)
	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		log.Fatalf("CreateContext failed: %+v\n", err)
	}

	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		log.Fatalf("CreateCommandQueue failed %+v\n", err)
	}

	opencl.Devices = devices
	opencl.Device = device
	opencl.Context = context
	opencl.Queue = queue

	sph := sph.Init(float32(1.0), vector.Vec{0, 0, 0}, nil, DIM, true)

	//Pre-Arrange
	ints := []int{sph.N(), len(sph.Particles()) - sph.N(), sph.Field().GetSampler().GetBuckets(), sph.Field().GetSampler().BucketSize()}
	floats := []float32{sph.CFL(), sph.Field().Mass(), sph.Delta(), sph.MaxV(), sph.Field().GetKernelLength()}
	hash_buffer_len := ints[2] * ints[3]
	particle_bytes := len(sph.Field().Particles())*(3*(3*4)) + (2 * 4)
	vector_bytes := 8 * 3 * 4
	temp_bytes := 2*3*4 + 4

	bf0, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, particle_bytes)
	bf1, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, len(ints)*4)
	bf2, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, len(floats)*4)
	bf3, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, hash_buffer_len*4)
	bf4, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, vector_bytes)
	bf5, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, temp_bytes*sph.N())

	opencl.Buffers["particles"] = bf0
	opencl.Buffers["ints"] = bf1
	opencl.Buffers["floats"] = bf2
	opencl.Buffers["hash"] = bf3
	opencl.Buffers["vecs"] = bf4
	opencl.Buffers["temp"] = bf5

	/*----------------------------Create Kernels And Programs *-----------------*/
	s0, err := ioutil.ReadFile(common.ProjectRelativePath("data/shaders/opencl/pcisph/pci_density.c"))
	s1, err1 := ioutil.ReadFile(common.ProjectRelativePath("data/shaders/opencl/pcisph/pci_predict.c"))
	inclDir := common.ProjectRelativePath("data/shaders/opencl/include")
	sources := []string{string(s0), string(s1)}

	if err != nil || err1 != nil {
		log.Fatalf("File not found")
	}
	var k1, k2 *cl.Kernel
	opencl.Program, err = opencl.Context.CreateProgramWithSource(sources)
	if err != nil {
		t.Errorf("Failed to create program from sources %v", err)
	}

	buildDevices := []*cl.Device{opencl.Device}
	if err := opencl.Program.BuildProgram(buildDevices, "-cl-kernel-arg-info -I "+inclDir); err != nil {
		t.Fatalf("BuildProgram failed: %+v", err)
	}

	k1, err = opencl.Program.CreateKernel("compute_density")
	if err != nil {
		t.Errorf("Failed to create compute_density kernel %v", err)
	}
	k2, err = opencl.Program.CreateKernel("predict_correct")
	if err != nil {
		t.Errorf("Failed to create predict correct kernel %v", err)
	}
	opencl.Kernels["compute_density"] = k1
	opencl.Kernels["predict_correct"] = k2

	gpu, err := New_GPUPredictorCorrector(sph, opencl)
	if err != nil {
		t.Errorf("Failed gpu PCISPH Implementation %v", err)
	}
	err = gpu.Run()
	if err != nil {
		t.Errorf("Failed Kernel Execution %v", err)
	}

}
