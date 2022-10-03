package gpu

import "github.com/andewx/dieselfluid/compute"
import "github.com/andewx/dieselfluid/common"
import "github.com/jgillich/go-opencl/cl"
import "io/ioutil"
import "fmt"
import "log"

const (
	OK   = 0
	ACK  = 1
	WAIT = 2
	RUN  = 3
)

type ComputeGPU struct {
	desc    compute.Descriptor
	sources []string
	kerns   map[string]*cl.Kernel
	buffers map[string]*cl.MemObject
	log     string
	devices []*cl.Device
	context *cl.Context
	program *cl.Program
	queue   *cl.CommandQueue
}

/* -------------------------------------------
    Compute Manager System Variable Access
------------------------------------------- */

func New_ComputeGPU(descriptor compute.Descriptor) ComputeGPU {
	mCompute := ComputeGPU{}
	mCompute.desc = descriptor
	mCompute.kerns = make(map[string]*cl.Kernel, 5)
	mCompute.buffers = make(map[string]*cl.MemObject, 20)
	mCompute.log = ""
	mCompute.devices = make([]*cl.Device, 10)
	mCompute.sources = make([]string, 10)
	return mCompute
}

func (cp ComputeGPU) Setup(verbose bool) bool {

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

		if verbose {
			fmt.Printf("  Address Bits: %d\n", d.AddressBits())
			fmt.Printf("  Available: %+v\n", d.Available())
			fmt.Printf("  Compiler Available: %+v\n", d.CompilerAvailable())
			fmt.Printf("  Double FP Config: %s\n", d.DoubleFPConfig())
			fmt.Printf("  Driver Version: %s\n", d.DriverVersion())
			fmt.Printf("  Error Correction Supported: %+v\n", d.ErrorCorrectionSupport())
			fmt.Printf("  Execution Capabilities: %s\n", d.ExecutionCapabilities())
			fmt.Printf("  Extensions: %s\n", d.Extensions())
			fmt.Printf("  Global Memory Cache Type: %s\n", d.GlobalMemCacheType())
			fmt.Printf("  Global Memory Cacheline Size: %d KB\n", d.GlobalMemCachelineSize()/1024)
			fmt.Printf("  Half FP Config: %s\n", d.HalfFPConfig())
			fmt.Printf("  Host Unified Memory: %+v\n", d.HostUnifiedMemory())
			fmt.Printf("  Image Support: %+v\n", d.ImageSupport())
			fmt.Printf("  Image2D Max Dimensions: %d x %d\n", d.Image2DMaxWidth(), d.Image2DMaxHeight())
			fmt.Printf("  Image3D Max Dimenionns: %d x %d x %d\n", d.Image3DMaxWidth(), d.Image3DMaxHeight(), d.Image3DMaxDepth())
			fmt.Printf("  Little Endian: %+v\n", d.EndianLittle())
			fmt.Printf("  Local Mem Size Size: %d KB\n", d.LocalMemSize()/1024)
			fmt.Printf("  Local Mem Type: %s\n", d.LocalMemType())
			fmt.Printf("  Max Clock Frequency: %d\n", d.MaxClockFrequency())
			fmt.Printf("  Max Constant Args: %d\n", d.MaxConstantArgs())
			fmt.Printf("  Max Constant Buffer Size: %d KB\n", d.MaxConstantBufferSize()/1024)
			fmt.Printf("  Max Mem Alloc Size: %d KB\n", d.MaxMemAllocSize()/1024)
			fmt.Printf("  Max Parameter Size: %d\n", d.MaxParameterSize())
			fmt.Printf("  Max Read-Image Args: %d\n", d.MaxReadImageArgs())
			fmt.Printf("  Max Write-Image Args: %d\n", d.MaxWriteImageArgs())
			fmt.Printf("  Memory Base Address Alignment: %d\n", d.MemBaseAddrAlign())
			fmt.Printf("  Native Vector Width Char: %d\n", d.NativeVectorWidthChar())
			fmt.Printf("  Native Vector Width Short: %d\n", d.NativeVectorWidthShort())
			fmt.Printf("  Native Vector Width Int: %d\n", d.NativeVectorWidthInt())
			fmt.Printf("  Native Vector Width Long: %d\n", d.NativeVectorWidthLong())
			fmt.Printf("  Native Vector Width Float: %d\n", d.NativeVectorWidthFloat())
			fmt.Printf("  Native Vector Width Double: %d\n", d.NativeVectorWidthDouble())
			fmt.Printf("  Native Vector Width Half: %d\n", d.NativeVectorWidthHalf())
			fmt.Printf("  Profiling Timer Resolution: %d\n", d.ProfilingTimerResolution())
		}

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
	cp.log += "Created Device\n"

	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		log.Fatalf("CreateCommandQueue failed %+v\n", err)
	}

	cp.log += "Created Command Queue\n"

	cp.log += "Sources Built Compiled and Linked\n"

	cp.devices = devices
	cp.context = context

	cp.queue = queue
	return true
}

func (cp ComputeGPU) BuildProgram() error {

	program, err := cp.context.CreateProgramWithSource(cp.sources)
	if err != nil {
		fmt.Errorf("CreateProgramWithSource failed: \n")
	}

	cp.log += "Sources Compiled\n"

	if err := program.BuildProgram(nil, ""); err != nil {
		fmt.Errorf("program.BuildProgram failed:\n ")
	}
	cp.log += "Build Successful\n"
	cp.program = program
	return err
}

/* ----------------------------------------
Kernels() - Maps a list of compiled Kernel functions to their referent integer
IDs for recall and identification @return map[string]int the list of compiled
kernel functions and their unique refrence IDs
----------------------------------------*/
func (cp ComputeGPU) RegisterKernel(name string) bool {
	kernel, err := cp.program.CreateKernel(name)
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
	cp.kerns[name] = kernel
	cp.log += "Registered Kernel " + name + "\n"
	return true
}

/*Creates source from file*/
func (cp ComputeGPU) AddSourceFile(filename string) bool {
	valid := false
	kern_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		cp.log += "Unable to open file: " + filename + "please enure correct path used\n"
		cp.log += "Kernel function was not registered"
	} else {
		kern_src := string(kern_bytes) + "\x00"
		cp.sources = append(cp.sources, kern_src)
		valid = true
	}

	return valid
}

func (cp ComputeGPU) AddSourceString(source string) bool {
	valid := false
	kern_src := source + "\x00"
	cp.sources = append(cp.sources, kern_src)
	valid = true
	return valid
}

/* -------------------------------------------
    GPU Contextual Commands - Execution Path
------------------------------------------- */
func (cp ComputeGPU) Queue(name string) {
	if _, err := cp.queue.EnqueueNDRangeKernel(cp.kerns[name], nil, cp.desc.Work, cp.desc.Local, nil); err != nil {
		log.Fatalf("EnqueueNDRangeKernel failed: %+v\n", err)
	}
	cp.log += "Queue() - Enqueued Kernel " + name + "for execution\n"
}

func (cp ComputeGPU) Set(d compute.Descriptor) {
	cp.desc = d
}
func (cp ComputeGPU) Get() compute.Descriptor {
	return cp.desc
}

//Buffer Routines first argument is the GPU Mapped Buffer ID which also maps to a Buffer Address
//Copy Functions copy from GPU memory to CPU Memory. Pass copies from CPU
func (cp ComputeGPU) RegisterBuffer(bytes_size int, t int, name string) {
	buffer, err := cp.context.CreateEmptyBuffer(cl.MemReadWrite, bytes_size)
	if err != nil {
		log.Fatalf("CreateBuffer failed for output: %+v\n", err)
	}
	cp.log += "RegisterBuffer() - Created Buffer " + name + "\n"
	cp.buffers[name] = buffer
}

func (cp ComputeGPU) ReadFloatBuffer(cpu_buffer []float32, name string) {
	if _, err := cp.queue.EnqueueReadBufferFloat32(cp.buffers[name], true, 0, cpu_buffer, nil); err != nil {
		log.Fatalf("EnqueueReadBufferFloat32 failed: %+v\n", err)
	}
}
func (cp ComputeGPU) PassFloatBuffer(cpu_buffer []float32, name string) {
	if _, err := cp.queue.EnqueueWriteBufferFloat32(cp.buffers[name], true, 0, cpu_buffer, nil); err != nil {
		log.Fatalf("EnqueueWriteBufferFloat32 failed: %+v\n", err)
	}
}
func (cp ComputeGPU) ReadIntBuffer(cpu_buffer []int, name string) {
	if _, err := cp.queue.EnqueueReadBuffer(cp.buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
		log.Fatalf("EnqueueReadBufferInt32 failed: %+v\n", err)
	}
}
func (cp ComputeGPU) PassIntBuffer(cpu_buffer []int, name string) {
	if _, err := cp.queue.EnqueueWriteBuffer(cp.buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
		log.Fatalf("EnqueueWriteBufferInt32 failed: %+v\n", err)
	}
	cp.log += "Passed Integer Buffer " + name + "\n"
}

func (cp ComputeGPU) PassLayoutBuffer(data interface{}, bytes int, name string) {
	if _, err := cp.queue.EnqueueWriteBuffer(cp.buffers[name], true, 0, bytes, common.Ptr(data), nil); err != nil {
		log.Fatalf("EnqueueWriteBuffer failed: %+v\n", err)
	}
	cp.log += "Passed Layout Buffer " + name + "\n"
}

//State Booleans
func (cp ComputeGPU) HasDeviceContext() bool {
	if cp.context != nil {
		return true
	}
	return true
}

//Pseudo Valid State
func (cp ComputeGPU) ValidState() bool {
	if cp.context != nil && cp.program != nil {
		return true
	}
	return false
}

func (cp ComputeGPU) Log() string {
	return cp.log
}
