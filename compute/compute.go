package compute

/* MacOSX does not allow OpenGL compute so we are implementing an abstract interface
and API into compute functionality. When building compute tasks build with either build
tags (osx, darwin). osx render functionality will need to be replaced with a metal/vulkan
implementation in the near future*/

//Kernel Descriptor
type Descriptor struct {
	Work  []int
	Local []int
	Size  int
}

//Simplistic Compute Interface when looking for run simple commands with Go Routines
type Compute interface {
	Run(x chan int)
	Commit()
	Pre()
	Post()
	Set(Descriptor)
	Get() Descriptor
}

//GPU Computing Interface - Exposes a general purpose API for linearity between platforms
type GPUCompute interface {

	//Setup
	Setup(bool) bool

	//Storage Access
	Queue(string)
	Run(x chan int)
	Set(Descriptor)
	Get() Descriptor
	//Buffer Routines first argument is the GPU Mapped Buffer ID which also maps to a Buffer Address
	//Copy Functions copy from GPU memory to CPU Memory. Pass copies from CPU -> GPU
	RegisterBuffer(int, int, string)
	ReadFloatBuffer([]float32, string)
	PassFloatBuffer([]float32, string)
	PassIntBuffer([]int, string)
	ReadIntBuffer([]int, string)

	//Registers a Kernel and adds the filename as a listing identifier
	RegisterKernel(string) bool
	AddSourceFile(string) bool
	AddSourceString(string) bool

	//State Booleans
	HasDeviceContext() bool
	ValidState() bool
	Log() string
}
