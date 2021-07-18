package glr

import (
	V "dslfluid.com/dsl/math/math32"
	"dslfluid.com/dsl/render"
	"dslfluid.com/dsl/render/scene"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"strings"
	"time"
)

//Implements  RenderAPIContext & Renderer Interface in float32 Env
const (
	CtxVarSize = 10
)

var GlobalTrans *V.Vec
var RotateTime time.Time
var RotateTimeLast time.Time
var state_hold_mouse bool
var xT float64
var yT float64
var RotAngleX float32
var RotAngleY float32
var Fps float64
var lastTime time.Time
var animationTime float64

//GLRenderer Holds RenderAPIContext Vars
type GLRenderer struct {
	GLCTX    render.Context
	GLHandle *glfw.Window
	GLScene  scene.DSLScene
}

func InitRenderer() GLRenderer {
	return GLRenderer{render.Context{}, nil, scene.DSLScene{}}
}

//Defines Application Global Parameters
type GLAppState struct {
	LocalMat        [16]float32
	RotX            float32
	RotY            float32
	FPS             float64
	TimeLast        time.Time
	ApplicationTime float64
	MouseStateHold  bool
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
		var logLength = int32(1000)
		fmt.Printf("Log length %d\n", logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Printf("%s", log)
		return 0, fmt.Errorf("GLSL Shader failed to compile\n: %v", log)
	}
	free()
	return shader, nil
}

//Setup Calls Window Context -- After setting up GLRenderer you should pass in relevant
//OpenGL Draw Data initialization to include Shader Paths, Vertex, Indice, and Color Buffers
//Then Call Init and Draw
func (glRender *GLRenderer) Setup(filepath string, width float32, height float32) {

	//Initialize the render contexts
	newContext := render.Context{}
	runtime.LockOSThread()

	fmt.Printf("glfw.Init()\n")
	err := glfw.Init()
	if err != nil {
		fmt.Printf("Error glfw.Init() %s\n", err)
		panic(err)
	}

	//Still Need Camera Location and Matrix Identities as well as animation params
	newContext.VertexShaderPaths = make(map[string]string, CtxVarSize)
	newContext.FragmentShaderPaths = make(map[string]string, CtxVarSize)
	newContext.ShaderPrograms = make(map[string]uint32, CtxVarSize)
	newContext.FragShaderID = make(map[string]uint32, CtxVarSize)
	newContext.VertShaderID = make(map[string]uint32, CtxVarSize)
	newContext.VBO = nil
	newContext.VAO = nil
	newContext.ProgramID = make(map[string]uint32, CtxVarSize)
	newContext.ShaderUniforms = make(map[string]int32, CtxVarSize)
	newContext.VertexShaderPaths["default"] = "/Users/briananderson/go/src/github.com/andewx/dslfluid/shader/glsl/geom.vert"
	newContext.FragmentShaderPaths["default"] = "/Users/briananderson/go/src/github.com/andewx/dslfluid/shader/glsl/geom.frag"

	glRender.GLCTX = newContext
	fmt.Printf("Initiating DSL Scene\n")
	glRender.GLScene = scene.InitDSLScene(filepath, width, height)
	glRender.GLScene.ImportGLTF()
	fmt.Printf("Scene Loaded\n")
	fmt.Printf("Accessors(%d)\n", len(glRender.GLScene.GetAccessors()))

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	fmt.Printf("Creating Window\n...\n")
	window, err := glfw.CreateWindow(640, 400, "DSL Fluid", nil, nil)
	fmt.Printf("Window Created...\n")
	glRender.GLHandle = window
	checkError(err)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	window.SetKeyCallback(ProcessInput)
	window.SetMouseButtonCallback(ProcessMouse)
	window.SetCursorPosCallback(ProcessCursor)
}

func (glRender *GLRenderer) Init() error {
	vtxSources := make(map[string]string, CtxVarSize)
	frgSources := make(map[string]string, CtxVarSize)

	for key, path := range glRender.GLCTX.VertexShaderPaths {
		srcVertex, err := ioutil.ReadFile(path)
		if checkError(err) {
			return err
		}
		vtxSources[key] = string(srcVertex) + "\x00"
	}

	for key, path := range glRender.GLCTX.FragmentShaderPaths {
		srcFrag, err := ioutil.ReadFile(path)
		if checkError(err) {
			return err
		}
		frgSources[key] = string(srcFrag) + "\x00"
	}

	//Compile Vertex Shaders

	//Read-In All Vertex Shaders
	fmt.Printf("COMPILING VERTEX SHADERS\n")
	for key, source := range vtxSources {

		vtxSHO, err := compileShader(source, gl.VERTEX_SHADER)
		if checkError(err) {
			fmt.Printf("Could not compile vertex\n")
		}
		glRender.GLCTX.VertShaderID[key] = vtxSHO
	}

	//Compile Fragment Shaders
	//Read In All Fragment Shaders
	fmt.Printf("COMPILING FRAG SHADERS\n")
	for key, source := range frgSources {
		frgSHO, err := compileShader(source, gl.VERTEX_SHADER)
		if checkError(err) {
			fmt.Printf("Could not compile frag\n")
		}
		glRender.GLCTX.FragShaderID[key] = frgSHO
	}

	//Link Default Program
	prog1 := gl.CreateProgram()
	gl.AttachShader(prog1, glRender.GLCTX.VertShaderID["default"])
	gl.AttachShader(prog1, glRender.GLCTX.FragShaderID["default"])
	gl.LinkProgram(prog1)
	gl.UseProgram(prog1)
	glRender.GLCTX.ProgramID["default"] = prog1

	//Initialize Scene and Assign Camera projection and vars
	glRender.GLCTX.ShaderUniforms["model"] = gl.GetUniformLocation(prog1, gl.Str("model\x00"))
	glRender.GLCTX.ShaderUniforms["view"] = gl.GetUniformLocation(prog1, gl.Str("view\x00"))
	glRender.GLCTX.ShaderUniforms["proj"] = gl.GetUniformLocation(prog1, gl.Str("projection\x00"))
	glRender.GLCTX.ShaderUniforms["rotx"] = gl.GetUniformLocation(prog1, gl.Str("rotX\x00"))
	glRender.GLCTX.ShaderUniforms["roty"] = gl.GetUniformLocation(prog1, gl.Str("rotY\x00"))
	glRender.GLCTX.ShaderUniforms["rot0"] = gl.GetUniformLocation(prog1, gl.Str("rotOriginTrans0\x00"))
	glRender.GLCTX.ShaderUniforms["rot1"] = gl.GetUniformLocation(prog1, gl.Str("rotOriginTrans1\x00"))
	glRender.GLCTX.ShaderUniforms["rotOrigin"] = gl.GetUniformLocation(prog1, gl.Str("rotOrigin\x00"))

	//Set the VAO
	glRender.MakeGLBuffers()
	GlobalTrans = &V.Vec{0, 0, 0}

	return nil
}

func (glRender *GLRenderer) MakeGLBuffers() {
	AccessorArray := glRender.GLScene.GetAccessors()
	nAccessors := len(AccessorArray)
	fmt.Printf("Count(Accessors): %d\n", nAccessors)
	glRender.GLCTX.VBO = make([]uint32, nAccessors)
	glRender.GLCTX.VAO = make([]uint32, nAccessors)

	//Initialize GLBuffers
	gl.GenBuffers(int32(nAccessors), &glRender.GLCTX.VBO[0])
	gl.GenVertexArrays(int32(nAccessors), &glRender.GLCTX.VAO[0])

	for i := 0; i < nAccessors; i++ {
		vaoBind := glRender.GLCTX.VAO[i]
		accessor, buffer_view, err := glRender.GLScene.GetAccessorBufferView(i)

		if err == nil {

			componentType := accessor.ComponentType
			size := SizeGL(accessor.Type)
			bfr_offset := buffer_view.ByteOffset
			bfr_byteLen := buffer_view.ByteLength

			tgt := uint32(gl.ARRAY_BUFFER)
			if componentType == 5123 {
				tgt = uint32(gl.ELEMENT_ARRAY_BUFFER)
			}
			gl.BindVertexArray(vaoBind)
			gl.BindBuffer(tgt, glRender.GLCTX.VBO[i])
			gl.BufferSubData(tgt, bfr_offset, bfr_byteLen, gl.Ptr(&(*glRender.GLScene.Buffers[i][0]))
			gl.VertexAttribPointer(uint32(i), int32(size), uint32(componentType), false, 0, nil)
			gl.BindBuffer(tgt, 0)
		}

	}
	gl.BindVertexArray(0)
}

func (glRender *GLRenderer) CleanBuffers() {
	nAccessors := len(glRender.GLScene.GetAccessors())
	for i := 0; i < nAccessors; i++ {
		gl.DeleteVertexArrays(int32(i), &glRender.GLCTX.VAO[0])
		gl.DeleteBuffers(int32(i), &glRender.GLCTX.VBO[0])
	}
}

func SizeGL(typeID string) uint32 {
	if typeID == "SCALAR" {
		return 1
	}
	if typeID == "VEC3" {
		return 3
	}
	if typeID == "VEC2" {
		return 2
	}
	if typeID == "VEC4" {
		return 4
	}
	return 1
}

func SetRotMatrix(dsl *scene.DSLScene) {
	cosAX := float32(math.Cos(float64(RotAngleX)))
	sinAX := float32(math.Sin(float64(RotAngleX)))
	cosAY := float32(math.Cos(float64(RotAngleY)))
	sinAY := float32(math.Sin(float64(RotAngleY)))

	dsl.RotX[4] = cosAX
	dsl.RotX[5] = sinAX
	dsl.RotX[7] = -sinAX
	dsl.RotX[8] = cosAX

	dsl.RotY[0] = cosAY
	dsl.RotY[2] = -sinAY
	dsl.RotY[6] = sinAY
	dsl.RotY[8] = cosAY

}

//Draw rendering routine fundamentally recognizes that a VBO draw VAO draw routine
//Is attached by key map name to a specific shader program all vertex buffers are drawn
//With associated shader routines - should iteratively draw meshes
func (glRender *GLRenderer) Draw() error {

	//FPS
	elapse_s := lastTime.Sub(time.Now()).Seconds()
	Fps = 1 / elapse_s

	GlobalTrans[0] *= 0
	GlobalTrans[1] *= 0
	GlobalTrans[2] *= 0

	gl.UseProgram(glRender.GLCTX.ProgramID["default"])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["model"], 1, false, &glRender.GLScene.Model[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["view"], 1, false, &glRender.GLScene.View[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["proj"], 1, false, &glRender.GLScene.Cam.ProjMat[0])
	gl.UniformMatrix3fv(glRender.GLCTX.ShaderUniforms["rotx"], 1, false, &glRender.GLScene.RotX[0])
	gl.UniformMatrix3fv(glRender.GLCTX.ShaderUniforms["roty"], 1, false, &glRender.GLScene.RotY[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["rot0"], 1, false, &glRender.GLScene.Rot0[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["rot1"], 1, false, &glRender.GLScene.Rot1[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["rotOrign"], 1, false, &glRender.GLScene.RotOrigin[0])

	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	nAccessors := len(glRender.GLScene.GetAccessors())

	for i := 0; i < nAccessors; i++ {
		mAccessor, err := glRender.GLScene.GetAccessorIx(i)
		if err == nil {
			gl.BindVertexArray(glRender.GLCTX.VAO[i])
			gl.EnableVertexAttribArray(uint32(i))
			gl.DrawArrays(gl.TRIANGLES, 0, int32(mAccessor.Count))
			gl.DisableVertexAttribArray(uint32(i))
		}
	}
	return nil

}

func (glRender *GLRenderer) Update() error {
	glRender.CleanBuffers()
	return nil

}

func (glRender *GLRenderer) Close() error {

	return nil
}

func ProcessInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	cameraSpeed := float32(0.009) // adjust accordingly - Just use framerate

	trans := cameraSpeed

	if key == glfw.KeyW {
		//Accumulates the translation Vector
		p := V.Vec{0.0, 0.0, -1.0}
		GlobalTrans.Add(*p.Scl(trans))
	}
	if key == glfw.KeyS {
		//Accumulates the translation Vector
		p := V.Vec{0.0, 0.0, 1.0}
		GlobalTrans.Add(*p.Scl(trans))
	}
	if key == glfw.KeyA {
		//Accumulates the translation Vector
		//Accumulates the translation Vector
		p := V.Vec{-1.0, 0.0, 0.0}
		GlobalTrans.Add(*p.Scl(trans))
	}
	if key == glfw.KeyD {
		//Accumulates the translation Vector
		p := V.Vec{1.0, 0.0, 0.0}
		GlobalTrans.Add(*p.Scl(trans))
	}

	if key == glfw.KeyUp {
		p := V.Vec{0.0, -1.0, 0.0}
		GlobalTrans.Add(*p.Scl(trans))
	}
	if key == glfw.KeyDown {
		p := V.Vec{0.0, 1.0, 0.0}
		GlobalTrans.Add(*p.Scl(trans))
	}
	if key == glfw.KeyTab {
		fmt.Printf("Current Simulation Time: %f\n", animationTime)
	}
}

//Set Mouse Callba j
func ProcessMouse(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		//Save the raw xy postion intial and continuously update and poll
		if !state_hold_mouse {
			state_hold_mouse = true
			RotateTime = time.Now()
			RotateTimeLast = time.Now()

		} else {
			RotateTimeLast = time.Now()
		}

	}
	if button == glfw.MouseButtonLeft && action == glfw.Release {
		state_hold_mouse = false
		xT = 0
		yT = 0
	}
}

func ProcessCursor(w *glfw.Window, xPos float64, yPos float64) {

	if state_hold_mouse {
		dt := RotateTimeLast.Sub(time.Now()).Seconds()
		if !(xT == 0) && !(yT == 0) {
			xdt := (xT - xPos) / dt
			ydt := (yT - yPos) / dt
			RotAngleX += float32(ydt / 100)
			RotAngleY += float32(xdt / 100)
		}
		xT = xPos
		yT = yPos
	}
}
