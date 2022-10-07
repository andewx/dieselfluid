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
	kerns   map[string]*cl.Kernel
	buffers map[string]*cl.MemObject
	log     string
	devices []*cl.Device
	device  *cl.Device
	context *cl.Context
	program *cl.Program
	queue   *cl.CommandQueue
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

/*Go will release the referenced C pointer device objects so high level code
  must retrieve the context and pass down contextually to the GPU class */
func New_ComputeGPU(descriptor compute.Descriptor, context OpenCL) ComputeGPU {
	mCompute := ComputeGPU{}
	mCompute.desc = descriptor
	mCompute.kerns = context.Kernels
	mCompute.buffers = context.Buffers
	mCompute.log = ""
	mCompute.devices = context.Devices
	mCompute.sources = make(map[string]string, 0)
	mCompute.context = context.Context
	mCompute.device = context.Device
	mCompute.program = context.Program
	mCompute.queue = context.Queue
	return mCompute
}

func (cp ComputeGPU) BuildProgram() error {
	var program *cl.Program
	var err error

	values := []string{}
	for _, value := range cp.sources {
		values = append(values, value)
	}
	if len(values) > 0 {
		program, err = cp.context.CreateProgramWithSource(values)
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

func (cp ComputeGPU) Context() OpenCL {
	return OpenCL{cp.context, cp.device, cp.program, cp.buffers, cp.kerns, cp.devices, cp.queue}
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

	if cp.kerns[name] == nil {
		return fmt.Errorf("Kernel %s has nil reference\n", name)
	}

	if _, err := cp.queue.EnqueueNDRangeKernel(cp.kerns[name], nil, cp.desc.Work, cp.desc.Local, nil); err != nil {
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
	if cp.buffers[name] == nil {
		return fmt.Errorf("%s is not a regisered buffer or the reference has been released\n", name)
	}
	return nil
}

//Buffer Routines first argument is the GPU Mapped Buffer ID which also maps to a Buffer Address
//Copy Functions copy from GPU memory to CPU Memory. Pass copies from CPU
func (cp ComputeGPU) RegisterBuffer(bytes_size int, t int, name string) error {
	buffer, err := cp.context.CreateEmptyBuffer(cl.MemReadWrite, bytes_size)
	if err == nil {
		cp.log += "RegisterBuffer() - Created Buffer " + name + "\n"
		cp.buffers[name] = buffer
	}
	return err
}

func (cp ComputeGPU) ReadFloatBuffer(cpu_buffer []float32, name string) error {
	var err error
	err = cp.isregistered(name)
	if err != nil {
		return err
	}
	if _, err = cp.queue.EnqueueReadBufferFloat32(cp.buffers[name], true, 0, cpu_buffer, nil); err != nil {
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
	if _, err = cp.queue.EnqueueWriteBufferFloat32(cp.buffers[name], true, 0, cpu_buffer, nil); err != nil {
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
	if _, err = cp.queue.EnqueueReadBuffer(cp.buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
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
	if _, err = cp.queue.EnqueueWriteBuffer(cp.buffers[name], true, 0, len(cpu_buffer), common.Ptr(cpu_buffer), nil); err != nil {
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
	if _, err = cp.queue.EnqueueWriteBuffer(cp.buffers[name], true, 0, bytes, common.Ptr(data), nil); err != nil {
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
	if &cp.context != nil && cp.program != nil {
		return true
	}
	return false
}

func (cp ComputeGPU) SetArgs(name string, args ...interface{}) error {
	kern := cp.kerns[name]
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
