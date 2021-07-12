package render

import (
	V "dslfluid.com/dsl/math/math32"
	"dslfluid.com/dsl/render"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"time"
)

//Implements  RenderAPIContext & Renderer Interface in float32 Env
const (
	CtxVarSize = 10
)

//GLRenderer Holds RenderAPIContext Vars
type GLRenderer struct {
	GLCTX    render.Context
	GLHandle *glfw.Window
}

//GLTransforms defines GL Transform state
type GLTransforms struct {
	Model     *V.Mat4
	View      *V.Mat4
	Proj      *V.Mat4
	RotX      *V.Mat3
	RotY      *V.Mat3
	RotTrans  *V.Mat4
	RotTrans0 *V.Mat4
	RotOrign  *V.Vec
}

//Defines Application Global Parameters
type GLAppState struct {
	Translation     *V.Vec
	GLCamera        *Camera
	RotX            float32
	RotY            float32
	FPS             float64
	TimeLast        time.Time
	ApplicationTime float64
	MouseStateHold  bool
}

func init() {
	runtime.LockOSThread()
}

func checkError(err error) bool {
	if err != nil {
		fmt.Printf(err.Error())
		return true
	}
	return false
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("GLSL Shader failed to compile\n: %v", log)
	}
	free()
	return shader, nil
}

//AddVertexBuffer - Adds Vertex Buffer to OpenGL Context under the mapped name
func (glRender *GLRenderer) AddVertexBuffer(name string, vertBuffer []float32, length int) {
	glRender.GLCTX.Vertexes[name] = vertBuffer
}

//AddIndiceBuffer - Adds Indice Draw Command Buffer - indice buffers are draw specific commands
//which tell opengl which order to access vertex buffers, neccessary for polygonal rendering
//Make sure if you intend to dereference a previously used buffer no references are held in state
func (glRender *GLRenderer) AddIndiceBuffer(name string, indiceBuffer []uint32, length int) {
	glRender.GLCTX.Indices[name] = indiceBuffer
}

//Setup Calls Window Context -- After setting up GLRenderer you should pass in relevant
//OpenGL Draw Data initialization to include Shader Paths, Vertex, Indice, and Color Buffers
//Then Call Init and Draw
func (glRender *GLRenderer) Setup() {

	//Initialize the render contexts
	newContext := render.Context{}

	//Still Need Camera Location and Matrix Identities as well as animation params
	newContext.VertexShaderPaths = make(map[string]string, CtxVarSize)
	newContext.FragmentShaderPaths = make(map[string]string, CtxVarSize)
	newContext.ShaderPrograms = make(map[string]uint32, CtxVarSize)
	newContext.FragShaderID = make(map[string]uint32, CtxVarSize)
	newContext.VertShaderID = make(map[string]uint32, CtxVarSize)
	newContext.Vertexes = make(map[string][]float32, CtxVarSize)
	newContext.Indices = make(map[string][]uint32, CtxVarSize)
	newContext.VBO = make(map[string]uint32, CtxVarSize)
	newContext.VAO = make(map[string]uint32, CtxVarSize)
	newContext.IBO = make(map[string]uint32, CtxVarSize)
	newContext.ProgramID = make(map[string]uint32, CtxVarSize)
	newContext.ShaderUniforms = make(map[string]uint32, CtxVarSize)
	newContext.VertexShaderPaths["default"] = "resources/shaders/vert.glsl"
	newContext.FragmentShaderPaths["default"] = "resources/shaders/frag.glsl"

	glRender.GLCTX = newContext

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "DSL Fluid", nil, nil)

	glRender.GLHandle = window

	if err != nil {
		panic(err)
	}
	glRender.GLHandle.MakeContextCurrent()
}

func (glRender *GLRenderer) Init() error {
	vtxSources := make(map[string]string, CtxVarSize)
	frgSources := make(map[string]string, CtxVarSize)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	//Read-In All Vertex Shaders
	for key, path := range glRender.GLCTX.VertexShaderPaths {
		srcVertex, err := ioutil.ReadFile(path)
		if checkError(err) {
			return err
		}
		vtxSources[key] = string(srcVertex) + "\x00"
	}

	//Read In All Fragment Shaders
	for key, path := range glRender.GLCTX.FragmentShaderPaths {
		srcFrag, err := ioutil.ReadFile(path)
		if checkError(err) {
			return err
		}
		frgSources[key] = string(srcFrag) + "\x00"
	}

	//Compile Vertex Shaders
	for key, source := range vtxSources {
		vtxSHO, err := compileShader(source, gl.VERTEX_SHADER)
		if checkError(err) {
			return err
		}
		glRender.GLCTX.VertShaderID[key] = vtxSHO
	}

	//Compile Fragment Shaders
	for key, source := range frgSources {
		frgSHO, err := compileShader(source, gl.VERTEX_SHADER)
		if checkError(err) {
			return err
		}
		glRender.GLCTX.FragShaderID[key] = frgSHO
	}

	//Link Default Program
	prog1 := gl.CreateProgram()
	gl.AttachShader(prog1, glRender.GLCTX.VertShaderID["default"])
	gl.AttachShader(prog1, glRender.GLCTX.FragShaderID["default"])
	gl.LinkProgram(prog1)
	glRender.GLCTX.ProgramID["default"] = prog1

	return nil
}

//Draw rendering routine fundamentally recognizes that a VBO draw VAO draw routine
//Is attached by key map name to a specific shader program all vertex buffers are drawn
//With associated shader routines
func (glRender *GLRenderer) Draw() error {

	return nil

}

func (glRender *GLRenderer) Update() error {

	return nil

}

func (glRender *GLRenderer) Close() error {

	return nil
}
