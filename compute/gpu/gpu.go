// +build !darwin

package gpu

import "github.com/andewx/dieselfluid/compute"
import "github.com/go-gl/gl/v4.3-core/gl"
import "github.com/andewx/dieselfluid/common"
import "github.com/andewx/dieslfluid/shader"
import "math/rand"
import "ioutil"
import "fmt"
import "log"

const (
  OK = 0
  ACK = 1
  WAIT  = 2
  RUN = 3
)

type ComputeGPU struct {
	desc    compute.Descriptor
	sources map[string]string
  layouts map[string]int
  buffers map[string]int
  shaders map[string]*shader.Shader
  program map[string]*shader.Program
  active_program *shader.Program
	log     string
}

/* -------------------------------------------
    Compute Manager System Variable Access
------------------------------------------- */


func New_ComputeGPU(descriptor compute.Descriptor)(ComputeGPU){
  mCompute := ComputeGPU{}
  mCompute.desc = descriptor
  mCompute.sources = make(map[string]string, 10)
  mCompute.layouts = make(map[string]int, 20)
  mCompute.buffers = make(map[string]int,20)
  mCompute.shaders = make(map[string]*shader.Shader, 10)
  mCompute.program = make(map[string]*shader.Program,10)
  mCompute.log += "Compute GPU Created"
}

func (cp ComputeGPU) Setup(hasContext bool) bool {
  if !hasContext{
    if err := gl.Init(); err != nil {
  		return false
  	}
  }

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL Version: %s\n", version)
  fmt.Printf("Please ensure your system is compatible with OpenGL4.3 or greater\n")
}

/* ----------------------------------------
Kernels() - Maps a list of compiled Kernel functions to their referent integer
IDs for recall and identification @return map[string]int the list of compiled
kernel functions and their unique refrence IDs - Kernels are shader sources in OpenGL
----------------------------------------*/
func (cp ComputeGPU) RegisterKernel(name string) bool {
	sh, err := shader.NewShaderFromSource(cp.sources[name], name, gl.COMPUTE_SHADER)
	if err != nil {
		log.Fatalf("Create Shader failed: %+v", err)
	}
  cp.log += "Registered Kernel Shader" +name"\n"
  cp.active_program = shader.NewProgram(name)
  cp.programs[name] = cp.active_program
  cp.active_program.Use()
  cp.active_program.AddShader(sh)
  cp.active_program.Link()
  //Register with a currently bound porgram
	return true
}

/*Files must be located in the data/shaders directory*/
func (cp ComputeGPU) AddSourceFile(filename string) bool {
	valid := false
	kern_bytes, err := ioutil.ReadFile(common.ProjectRelativePath("data/shaders/" + filename))
	if err != nil {
		cp.log += "Unable to open file: " + file + "please enure correct path used\n"
		cp.log += "Kernel function was not registered"
	} else {
		kern_src := string(kern_bytes) + "\x00"
		cp.sources[filename] = kern_src
		valid = true
	}

	return valid
}

func (cp ComputeGPU)AddSourceString(source string) bool{
  cp.log += "This method is not valid for non-darwin builds\n"
  cp.log += "Use (ComputGPU) AddSourceFile(filename string) bool instead \n"
  return false
}


/* -------------------------------------------
    GPU Contextual Commands - Execution Path
------------------------------------------- */

//Executes an individual compute program
func (cp ComputeGPU) Queue(name string){
  cp.SetActive(name)
  cp.program[name].DispatchCompute(cp.desc.Work.Size[0],cp.desc.Work.Size[1],cp.desc.Work.Size[2])
}

func (cp ComputeGPU) Set(d Descriptor)  {
  cp.desc = d
}
func (cp ComputeGPU) Get() Descriptor {
  return cp.desc
}
func (cp ComputeGPU) Pre(x chan int){
  x <- OK
}
func (cp ComputeGPU) Post(x chan int){
  x <- OK
}

//Buffer Routines first argument is the GPU Mapped Buffer ID which also maps to a Buffer Address
//Copy Functions copy from GPU memory to CPU Memory. Pass copies from CPU -> GPU
func (cp ComputeGPU) RegisterBuffer(buffer_length int, layout int, name string){
  cp.active_program.AddUniform(name)
  cp.layouts[name] = layout
  cp.buffers[name] = buffer_length
  gl.GenBuffers(1, cp.active_program[name].Address(name))
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, cp.active_program[name].Location(name))
}

func (cp ComputeGPU) ReadFloatBuffer(cpu_buffer []float32, name string) {
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER,cp.active_program.Location(name))
  gl.GetBufferSubData(cp_active_program.Location(name),0,len(cpu_buffer), gl.Ptr(cpu_buffer))
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

}
func (cp ComputeGPU) PassFloatBuffer(cpu_buffer []float32, name string) {
  gl.BufferData(gl.SHADER_STORAGE_BUFFER, cp.buffers[name], gl.Ptr(cpu_buffer), gl.STREAM_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, cp.layouts[name], cp.active_program.Location(name))
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
func (cp ComputeGPU) ReadIntBuffer(cpu_buffer []int32, name string){
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER,cp.active_program.Location(name))
  gl.GetBufferSubData(cp_active_program.Location(name),0,len(cpu_buffer), gl.Ptr(cpu_buffer))
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

}
func (cp ComputeGPU) PassIntBuffer(cpu_buffer []int32, name string){
  gl.BufferData(gl.SHADER_STORAGE_BUFFER, cp.buffers[name], gl.Ptr(cpu_buffer), gl.STREAM_DRAW)
  gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, cp.layouts[name], cp.active_program.Location(name))
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}

func (cp ComputeGPU) PassLayoutBuffer(data interface{}, name string){
  gl.BufferData(gl.SHADER_STORAGE_BUFFER, cp.buffers[name], gl.Ptr(data), gl.STREAM_DRAW)
  gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, cp.layouts[name], cp.active_program.Location(name))
  gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}

//Pseudo Valid State
func (cp ComputeGPU) ValidState() bool       {
  if cp.context != nil && cp.program != nil && cp.device !=nil{
    return true
  }
  return false
}

func (cp ComputeGPU) Log()string{
  return cp.log
}

func (cp ComputeGPU) SetActive(name string){
  cp.active_program = cp.program[name]
  cp.active_program.Use()
}
