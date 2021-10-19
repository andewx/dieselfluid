package glr

/*
Package glr

Render implementation handles Global Render Routines, Setup, Draw, CreateGLTFRenderObjects,
Update, Clear etc...

Consider Refactoring Imperative Function Names and Map Concrete GlRenderer -> Typed Interface
Renderer

Note that GLR Renderer handles scene, camera, matrices, transforms, etc adding some dependencies

For GO purposes however working the Core GL Exported C Functions will probably comprise the scope
of available renderers. (Except for Realtime GI RT)
*/
import (
	"dslfluid.com/dsl/math/mgl"
	"dslfluid.com/dsl/render"
	"dslfluid.com/dsl/render/camera"
	"dslfluid.com/dsl/render/scene"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"strings"
	"time"
)

//Implements  RenderAPIContext & Renderer Interface in float32 Env
const (
	MIN_PARAM_SIZE      = 10
	COORDS_PER_VERTEX   = 3
	BYTES_PER_FLOAT     = 4
	BYTES_PER_SHORT     = 2
	VERT_POSITION_DATA  = 0
	VERT_ATTR_NORM_DATA = 3
	VERT_ATTR_RGB_DATA  = 6
	VERT_ATTR_UV_DATA   = 9
)

//GLRenderer Holds RenderAPIContext Vars
type GLRenderer struct {
	GLCTX         render.Context      //Maps paramaters and State
	GLHandle      *glfw.Window        //Holds GLFW Window
	GLScene       scene.DSLScene      //GLTF JSON
	RenderObjects []*GLTFRenderObject //Render Description
	MVPMat        []float32
	VAO           []uint32
	VBO           []uint32
	Indices       []uint16
}

type GLTFRenderObject struct {
	Indices         []byte //Short 2 Bytes * 1
	IndexByteOffset int
	IndexByteLength int
	IndexBufferId   uint32

	Vertices         []byte //Float 4 Bytes * 3
	VertexByteOffset int
	VertexByteLength int
	VertexBufferId   uint32

	Normals           []byte //Float 4 Bytes * 3
	NormalsByteOffset int
	NormalsByteLength int
	NormalsBufferId   uint32

	//Model Matrix
	Model []float32
}

//Render State Global State Structure
type RenderGlobal struct {
	MVPMatrix      []float32
	Camera         camera.Camera
	Translate      mgl.Vec
	MousePosition  [2]int
	TimeLast       time.Time
	MouseStateHold bool
	KeyPressed     int
	delT           float32
	XT             float64
	YT             float64
	CamRotY        float32
	CamRotX        float32
	RotateTime     time.Time
	RotateTimeLast time.Time
}

var RenderState *RenderGlobal = new(RenderGlobal)

func InitRenderer() GLRenderer {
	RenderState.Camera = camera.NewCamera(mgl.Vec{0, 0, 10})
	return GLRenderer{}

}

//Setup Calls Window Context -- After setting up GLRenderer you should pass in relevant
//OpenGL Draw Data initialization to include Shader Paths, Vertex, Indice, and Color Buffers
//Then Call Init and Draw
func (glRender *GLRenderer) Setup(filepath string, width float32, height float32) error {

	//Initialize the render contexts
	newContext := render.Context{}

	if len(glRender.MVPMat) == 0 {
		glRender.MVPMat = make([]float32, 16)
		glRender.MVPMat = mgl.ProjectionMatF(100.0, 1.0, 1000)
	}

	//Still Need Camera Location and Matrix Identities as well as animation params
	newContext.VertexShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	newContext.FragmentShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	newContext.FragShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	newContext.VertShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	newContext.ProgramID = make(map[string]uint32, MIN_PARAM_SIZE)
	newContext.ShaderUniforms = make(map[string]int32, MIN_PARAM_SIZE)
	newContext.VertexShaderPaths["default"] = "../../shader/glsl/geom.vert"
	newContext.FragmentShaderPaths["default"] = "../../shader/glsl/geom.frag"

	glRender.GLCTX = newContext
	glRender.GLScene = scene.InitDSLScene("../../resources/", filepath, width, height)
	if err := glRender.GLScene.ImportGLTF(); err != nil {
		return err
	}
	glRender.GLHandle.SetMouseButtonCallback(ProcessMouse)
	glRender.GLHandle.SetCursorPosCallback(ProcessCursor)
	glRender.GLHandle.SetKeyCallback(ProcessInput)
	return nil

}

func (glRender *GLRenderer) CompileShaders() error {
	vtxSources := make(map[string]string, MIN_PARAM_SIZE)
	frgSources := make(map[string]string, MIN_PARAM_SIZE)

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
		frgSHO, err := compileShader(source, gl.FRAGMENT_SHADER)
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

	var logLength = int32(1000)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(prog1, logLength, nil, gl.Str(log))
	fmt.Printf("%s", log)

	active := int32(0)
	gl.GetProgramiv(prog1, gl.ACTIVE_UNIFORMS, &active)
	fmt.Printf("SHADER UNIFORMS[%d]\n", active)
	glRender.GLCTX.ProgramID["default"] = prog1

	//Initialize Scene and Assign Camera projection and var
	glRender.GLCTX.ShaderUniforms["mvp"] = gl.GetUniformLocation(prog1, gl.Str("mvp\x00"))
	glRender.GLCTX.ShaderUniforms["model"] = gl.GetUniformLocation(prog1, gl.Str("model\x00"))
	glRender.GLCTX.ShaderUniforms["viewMat"] = gl.GetUniformLocation(prog1, gl.Str("viewMat\x00"))

	glRender.CreateGLTFRenderObjects()
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.DEPTH)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	return nil
}

//Prepares render data for each GLTF Mesh Primitive -- move this to scene import
func (glRender *GLRenderer) CreateGLTFRenderObjects() {

	meshes := len(glRender.GLScene.GetMeshes())
	glRender.RenderObjects = make([]*GLTFRenderObject, len(glRender.GLScene.GetMeshes()))
	glRender.VAO = make([]uint32, meshes)
	gl.GenVertexArrays(int32(meshes), &glRender.VAO[0])

	if len(glRender.GLScene.GetMeshes()) == 0 {
		fmt.Printf("No GLTF Meshes\n")
		glRender.GLScene.Info()
	}
	//Note in real glTF nodes have children we're only loading top level siblings
	for i := 0; i < len(glRender.GLScene.GetMeshes()); i++ {
		gl.BindVertexArray(glRender.VAO[i])
		mesh, err := glRender.GLScene.GetMeshIx(i)
		checkError(err)
		//Add each primitive into the render object list.
		for j := 0; j < len(mesh.Primitives); j++ {
			renderObject := new(GLTFRenderObject)
			primitive := mesh.Primitives[j]
			posAccessorIdx := primitive.Attributes["POSITION"]
			normAccessorIdx := primitive.Attributes["NORMAL"]
			indicesAccessorIdx := primitive.Indices

			//Setup Positional data
			PosAccessor, PosBufferView, Err := glRender.GLScene.GetAccessorBufferView(posAccessorIdx)
			checkError(Err)

			//Position RenderObject Data
			if PosAccessor.ComponentType == gl.FLOAT {
				renderObject.Vertices = glRender.GLScene.Buffers[PosBufferView.Buffer]
				renderObject.VertexByteLength = PosBufferView.ByteLength
				renderObject.VertexByteOffset = PosBufferView.ByteOffset

			} else {
				fmt.Errorf("TAG Not Implemented")
			}

			//Index Data for Object
			//Setup Positional data
			_, IdxBufferView, _ := glRender.GLScene.GetAccessorBufferView(indicesAccessorIdx)

			//Assumes Target
			renderObject.Indices = glRender.GLScene.Buffers[IdxBufferView.Buffer]
			renderObject.IndexByteLength = IdxBufferView.ByteLength
			renderObject.IndexByteOffset = IdxBufferView.ByteOffset

			//Normal Data For Object
			//Setup Positional data
			normAccessor, normBufferView, NErr := glRender.GLScene.GetAccessorBufferView(normAccessorIdx)
			checkError(NErr)

			//Position RenderObject Data
			if normAccessor.ComponentType == gl.FLOAT {
				renderObject.Normals = glRender.GLScene.Buffers[normBufferView.Buffer]
				renderObject.NormalsByteLength = normBufferView.ByteLength
				renderObject.NormalsByteOffset = normBufferView.ByteOffset

			} else {
				fmt.Errorf("TAG Not Implemented")
			}

			//Prepare to upload buffers to GPU
			vbos := make([]uint32, 3)
			gl.GenBuffers(3, &vbos[0])

			renderObject.VertexBufferId = vbos[0]
			renderObject.NormalsBufferId = vbos[1]
			renderObject.IndexBufferId = vbos[2]

			//Upload Vertex Buffer to GPU
			gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.VertexBufferId)
			gl.BufferData(gl.ARRAY_BUFFER, renderObject.VertexByteLength, gl.Ptr(&renderObject.Vertices[renderObject.VertexByteOffset]), gl.STATIC_DRAW)
			gl.VertexAttribPointer(0, COORDS_PER_VERTEX, gl.FLOAT, false, 0, gl.PtrOffset(0))
			gl.EnableVertexAttribArray(0)
			//Upload Normals Buffer to GPU

			gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.NormalsBufferId)
			gl.BufferData(gl.ARRAY_BUFFER, renderObject.NormalsByteLength, gl.Ptr(&renderObject.Normals[renderObject.NormalsByteOffset]), gl.STATIC_DRAW)
			gl.VertexAttribPointer(3, COORDS_PER_VERTEX, gl.FLOAT, false, 0, gl.PtrOffset(0))
			gl.EnableVertexAttribArray(3)
			//Upload Vertex Buffer to GPU
			numElements := renderObject.IndexByteLength / (BYTES_PER_SHORT)
			kb := renderObject.IndexByteLength / 1024
			fmt.Printf("Indices: %d [%dKB]\n", numElements, kb)

			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderObject.IndexBufferId)
			gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, renderObject.IndexByteLength, gl.Ptr(&renderObject.Indices[renderObject.IndexByteOffset]), gl.STATIC_DRAW)

			glRender.RenderObjects[i] = renderObject
			gl.BindVertexArray(0)
		}

	}

	nodes := glRender.GLScene.GetNodes()

	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		meshIdx := node.Mesh      //Check if field is empty
		trans := node.Translation //3 float 65
		if len(trans) == 0 {
			trans = make([]float32, 3)
		}
		scale := node.Scale //3 float64
		if len(scale) == 0 {
			scale = mgl.Mat3(1.0)
		}
		rot := node.Rotation //4 float64
		if len(rot) == 0 {
			rot = make([]float32, 4)
		}
		M := MatrixTRS(trans, rot, scale)
		if meshIdx >= 0 && meshIdx < len(glRender.RenderObjects) {
			glRender.RenderObjects[meshIdx].Model = M
		}
	}

}

//Constructs a Matrix from Translation scale rotation quat
func MatrixTRS(t []float32, r []float32, s []float32) []float32 {
	M := mgl.Mat4(1.0)

	//Trans Matrix Affine
	T := mgl.Mat4(1.0)
	T[12] = t[0]
	T[13] = t[1]
	T[14] = t[2]

	S := mgl.Mat4(1.0)
	S[0] = s[0]
	S[5] = s[1]
	S[10] = s[2]

	M = T.MulM(S)

	return M
}

func CheckGlError(op string) {
	error := gl.GetError()
	if error == gl.NO_ERROR {
		return
	}
	fmt.Printf(op+"GL Error %d: ", error)
}

//Draw rendering routine fundamentally recognizes that a VBO draw VAO draw routine
//Is attached by key map name to a specific shader program all vertex buffers are drawn
//With associated shader routines - should iteratively draw meshes
func (glRender *GLRenderer) Draw() error {

	mglView := RenderState.Camera.Update()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(glRender.GLCTX.ProgramID["default"])

	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["mvp"], 1, false, &glRender.MVPMat[0])
	gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["viewMat"], 1, false, &mglView[0])

	for i := 0; i < len(glRender.RenderObjects); i++ {

		gl.UniformMatrix4fv(glRender.GLCTX.ShaderUniforms["model"], 1, false, &glRender.RenderObjects[i].Model[0])
		gl.BindVertexArray(glRender.VAO[i])
		gl.DrawElements(gl.TRIANGLES, int32(glRender.RenderObjects[i].IndexByteLength/BYTES_PER_SHORT), gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.BindVertexArray(0)

	}

	glRender.GLHandle.SwapBuffers()
	glfw.PollEvents()
	return nil

}

func (glRender *GLRenderer) Update() error {

	return nil

}

func (glRender *GLRenderer) Close() error {

	return nil
}

func ProcessInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	cameraSpeed := float32(0.5) // adjust accordingly - Just use framerate
	RenderState.KeyPressed = 1
	trans := cameraSpeed
	dir := RenderState.Camera.Transform.Matrix.Get(2)
	ldir := RenderState.Camera.Transform.Matrix.Get(0)
	up := RenderState.Camera.Transform.Matrix.Get(1)

	if key == glfw.KeyW {
		//Accumulates the translation Vector
		p := mgl.Scale(dir, -trans)
		RenderState.Camera.Transform.Translate(p)
	}
	if key == glfw.KeyS {
		//Accumulates the translation Vector
		p := mgl.Scale(dir, trans)
		RenderState.Camera.Transform.Translate(p)
	}
	if key == glfw.KeyA {
		//Accumulates the translation Vector
		//Accumulates the translation Vector
		p := mgl.Scale(ldir, -trans)
		RenderState.Camera.Transform.Translate(p)
	}
	if key == glfw.KeyD {
		//Accumulates the translation Vector
		p := mgl.Scale(ldir, trans)
		RenderState.Camera.Transform.Translate(p)
	}

	if key == glfw.KeyUp {
		p := mgl.Scale(up, trans)
		RenderState.Camera.Transform.Translate(p)
	}
	if key == glfw.KeyDown {
		p := mgl.Scale(up, -trans)
		RenderState.Camera.Transform.Translate(p)
	}
	if key == glfw.KeyTab {
		RenderState.Camera.Log()
	}
}

//Set Mouse Callba j
func ProcessMouse(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		//Save the raw xy postion intial and continuously update and poll
		if !RenderState.MouseStateHold {
			RenderState.MouseStateHold = true
			RenderState.RotateTime = time.Now()
			RenderState.RotateTimeLast = time.Now()

		} else {
			RenderState.RotateTimeLast = time.Now()
		}

	}
	if button == glfw.MouseButtonLeft && action == glfw.Release {
		RenderState.MouseStateHold = false
		RenderState.XT = 0
		RenderState.YT = 0
		RenderState.RotateTime = time.Now()
		RenderState.RotateTimeLast = time.Now()
	}
}

//ProcessCursor - GLFW Callback. While mouse is in state hold calculate a X,Y derivative
func ProcessCursor(w *glfw.Window, xPos float64, yPos float64) {
	if RenderState.MouseStateHold {
		dt := RenderState.RotateTimeLast.Sub(time.Now()).Seconds()
		if !(RenderState.XT == 0) && !(RenderState.YT == 0) {
			xdt := (RenderState.XT - xPos) / (dt * 100)
			ydt := (RenderState.YT - yPos) / (dt * 100)
			RenderState.Camera.RotateFPS(mgl.Vec{float32(-ydt), float32(-xdt), 0})
			RenderState.RotateTimeLast = time.Now()
		}
		RenderState.XT = xPos
		RenderState.YT = yPos
	}

}
