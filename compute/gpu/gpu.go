package gpu

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/compute"
	"github.com/andewx/go-opencl/cl"
)

const (
	OK   = 0
	ACK  = 1
	WAIT = 2
	RUN  = 3
)

type ComputeGPU struct {
	desc    compute.Descriptor
	sources map[string]string
	context *OpenCL
	log     string
}

/* -------------------------------------------
    Compute Manager System Variable Access
------------------------------------------- */

type OpenCL struct {
	Context *cl.Context
	Device  *cl.Device
	Program *cl.Program
	Buffers map[string]*cl.MemObject
	Kernels map[string]*cl.Kernel
	Devices []*cl.Device
	Queue   *cl.CommandQueue
}

func (p ComputeGPU) SetDescriptor(desc compute.Descriptor) {
	p.desc = desc
}

func InitOpenCL(opencl *OpenCL) *OpenCL {
	/*Initiate OpenCL State*/
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
	return opencl
}

/*Go will release the referenced C pointer device objects so high level code
  must retrieve the context and pass down contextually to the GPU class */
func New_ComputeGPU(mCompute *ComputeGPU, descriptor compute.Descriptor, context *OpenCL) *ComputeGPU {
	mCompute.desc = descriptor
	mCompute.context = context
	mCompute.log = ""
	return mCompute
}

/*
	sph := sph.Init(float32(1.0), vector.Vec{0, 0, 0}, nil, DIM, true)

	//Pre-Arrange
	field := sph.Field()
	ints := []int{sph.N(), field.Particles.Total() - field.Particles.N(), sph.Field().GetSampler().GetBuckets(), sph.Field().GetSampler().BucketSize()}
	floats := []float32{sph.CFL(), sph.Field().Mass(), sph.Delta(), sph.MaxV(), sph.Field().GetKernelLength()}
	hash_buffer_len := ints[2] * ints[3]
	particle_bytes := field.Particles.Total() * (3 * 4)
	vec3_bytes := field.Particles.N() * 3 * 4
	vector_bytes := 8 * 3 * 4
	temp_bytes := 2*3*4 + 4

	bf1, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, len(ints)*4)
	bf2, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, len(floats)*4)
	bf3, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, hash_buffer_len*4)
	bf4, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, vector_bytes)
	bf5, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, temp_bytes*field.Particles.N())

	bf6, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, particle_bytes)         //positions
	bf7, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, vec3_bytes)             //velocities
	bf8, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, vec3_bytes)             //forces
	bf9, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, field.Particles.N()*4)  //densities
	bf10, _ := opencl.Context.CreateEmptyBuffer(cl.MemReadWrite, field.Particles.N()*4) //pressures

	opencl.Buffers["ints"] = bf1
	opencl.Buffers["floats"] = bf2
	opencl.Buffers["hash"] = bf3
	opencl.Buffers["vecs"] = bf4
	opencl.Buffers["temp"] = bf5
	opencl.Buffers["positions"] = bf6
	opencl.Buffers["velocities"] = bf7
	opencl.Buffers["forces"] = bf8
	opencl.Buffers["densities"] = bf9
	opencl.Buffers["pressures"] = bf10


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

	gpu, err := New_GPUPredictorCorrector(sph, opencl, false)
	if err != nil {
		t.Errorf("Failed gpu PCISPH Implementation %v", err)
	}
*/

func (cp ComputeGPU) BuildProgram() error {
	var program *cl.Program
	var err error

	values := []string{}
	for _, value := range cp.sources {
		values = append(values, value)
	}
	if len(values) > 0 {
		program, err = cp.context.Context.CreateProgramWithSource(values)
		if err != nil {
			return fmt.Errorf("CreateProgramWithSource failed: \n%s\n", cp.log)
		}
	} else {
		return fmt.Errorf("No source files available\n%s\n", cp.log)
	}

	cp.log += "Sources Compiled\n"

	if err := program.BuildProgram(nil, ""); err != nil {
		fmt.Errorf("program.BuildProgram failed:\n %s", cp.log)
	}
	cp.log += "Build Successful\n"
	cp.context.Program = program
	return err
}

/* ----------------------------------------
Kernels() - Maps a list of compiled Kernel functions to their referent integer
IDs for recall and identification @return map[string]int the list of compiled
kernel functions and their unique refrence IDs
----------------------------------------*/
func (cp ComputeGPU) RegisterKernel(name string) bool {
	kernel, err := cp.context.Program.CreateKernel(name)
	if err != nil {
		log.Fatalf("CreateKernel failed: %+v\n", err)
	}
	for i := 0; i < 3; i++ {
		name, err := kernel.ArgName(i)
		if err == cl.ErrUnsupported {
			break
		} else if err != nil {
			log.Printf("GetKernelArgInfo for name failed: %+v\n", err)
			break
		} else {
			log.Printf("Kernel arg %d: %s\n", i, name)
		}
	}
	cp.context.Kernels[name] = kernel
	cp.log += "Registered Kernel " + name + "\n"
	return true
}

func (cp ComputeGPU) Context() *OpenCL {
	return cp.context
}

/*Creates source from file*/
func (cp ComputeGPU) AddSourceFile(filename string) error {
	var kern_src string
	kern_bytes, err := ioutil.ReadFile(common.ProjectRelativePath("data/shaders/") + filename)
	if err != nil {
		cp.log += "Unable to open file: " + filename + "please enure correct path used\n"
		cp.log += "Kernel function was not registered"
		return err
	} else {
		kern_src = string(kern_bytes)
		cp.log += "Files added"
		cp.log += "Appended source string length\n"
		cp.sources[filename] = kern_src
	}
	return nil
}

func (cp ComputeGPU) AddSourceString(source string, key string) error {

	cp.sources[key] = source
	return nil
}

/* -------------------------------------------
    GPU Contextual Commands - Execution Path
------------------------------------------- */
func (cp ComputeGPU) Queue(name string) error {

	if cp.context.Kernels[name] == nil {
		return fmt.Errorf("Kernel %s has nil reference\n", name)
	}

	if _, err := cp.context.Queue.EnqueueNDRangeKernel(cp.context.Kernels[name], nil, cp.desc.Work, cp.desc.Local, nil); err != nil {
		return err
	}
	return nil
}

func (cp ComputeGPU) Set(d compute.Descriptor) {
	cp.desc = d
}
func (cp ComputeGPU) Get() compute.Descriptor {
	return cp.desc
}

func (cp ComputeGPU) isregistered(name string) error {
	if cp.context.Buffers[name] == nil {
		return fmt.Errorf("%s is not a regisered buffer or the reference has been released\n", name)
	}
	return nil
}

//Buffer Routines first argument is the GPU Mapped Buffer ID which also maps to a Buffer Address
//Copy Functions copy from GPU memory to CPU Memory. Pass copies from CPU
func (cp ComputeGPU) RegisterBuffer(bytes_size int, t int, name string) error {
	buffer, err := cp.context.Context.CreateEmptyBuffer(cl.MemReadWrite, bytes_size)
	if err == nil {
		cp.log += "RegisterBuffer() - Created Buffer " + name + "\n"
		cp.context.Buffers[name] = buffer
	}
	return err
}

func (cp ComputeGPU) RegisterGLBuffer(gl_buffer_id uint32, size int, name string) error {
	buffer, err := cp.context.Context.CreateFromGLBuffer(cl.MemReadWrite, uint(gl_buffer_id), size)
	if err == nil {
		cp.log += "RegisterBuffer() - Created Buffer " + name + "\n"
		cp.context.Buffers[name] = buffer
	}
	return err
}

func (cp ComputeGPU) ReadFloatBuffer(cpu_buffer []float32, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.context.Queue.EnqueueReadBufferFloat32(cp.context.Buffers[name], true, 0, cpu_buffer, nil); err != nil {
		log.Fatalf("EnqueueReadBufferFloat32 failed: %+v\n", err)
	}
	return nil
}
func (cp ComputeGPU) PassFloatBuffer(cpu_buffer []float32, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.context.Queue.EnqueueWriteBufferFloat32(cp.context.Buffers[name], true, 0, cpu_buffer, nil); err != nil {
		log.Fatalf("EnqueueWriteBufferFloat32 failed: %+v\n", err)
	}
	return nil
}
func (cp ComputeGPU) ReadIntBuffer(cpu_buffer []int, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.context.Queue.EnqueueReadBuffer(cp.context.Buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
		log.Fatalf("EnqueueReadBufferInt32 failed: %+v\n", err)
	}
	return nil
}
func (cp ComputeGPU) PassIntBuffer(cpu_buffer []int, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.context.Queue.EnqueueWriteBuffer(cp.context.Buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
		log.Fatalf("EnqueueWriteBufferInt32 failed: %+v\n", err)
	}
	cp.log += "Passed Integer Buffer " + name + "\n"
	return nil
}

func (cp ComputeGPU) PassLayoutBuffer(data interface{}, bytes int, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.context.Queue.EnqueueWriteBuffer(cp.context.Buffers[name], true, 0, bytes, common.Ptr(data), nil); err != nil {
		log.Fatalf("EnqueueWriteBuffer failed: %+v\n", err)
	}
	cp.log += "Passed Layout Buffer " + name + "\n"
	return nil
}

//State Booleans
func (cp ComputeGPU) HasDeviceContext() bool {
	if &cp.context != nil {
		return true
	}
	return true
}

//Pseudo Valid State
func (cp ComputeGPU) ValidState() bool {
	if &cp.context != nil && cp.context.Program != nil {
		return true
	}
	return false
}

func (cp ComputeGPU) SetArgs(name string, args ...interface{}) error {
	kern := cp.context.Kernels[name]
	for index, arg := range args {
		if arg == nil {
			return fmt.Errorf("Invalid nil argument (%d) passed to\n", index)
		}
	}
	if kern != nil {
		return kern.SetArgs(args)
	} else {
		return fmt.Errorf("Non valid kernel passed to set args with name %s\n", name)
	}
}

func (cp ComputeGPU) Log() string {
	return cp.log
}
