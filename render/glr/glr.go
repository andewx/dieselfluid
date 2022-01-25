package glr

/*

Global Render Interface - OpenGL

Holds render attributes and maps to windowing, event, CPU/GPU interface

Holds the GLTF Scene Description and Renders those objects
*/

import (
	"dslfluid.com/dsl/math/mgl"
	"dslfluid.com/dsl/render/camera"
	"dslfluid.com/dsl/render/defs"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"os"
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

//VAO Index Array Positions
const (
	VERT_NORM_HALF            = 0
	VERT_NORM_COLOR_HALF      = 1
	VERT_NORM_TEX_HALF        = 2
	VERT_NORM_TEX_COLOR_HALF  = 3
	VERT_NORM_TEX_BINORM_HALF = 4
	VERT_HALF                 = 5
	VERT_COLOR_HALF           = 6
	MAX_VAO_BINDING           = 7
)

const (
	CAMERA_SPEED = 0.05
)

//OGLRenderer interface type
type GLRenderer struct {
	MVPMatrix  []float32
	MVPMat     []float32
	VAO        []uint32
	Tex        []uint32
	VBO        []uint32
	ElementVBO []uint32
	Indices    []uint16
	Window     *glfw.Window
	Camera     *camera.Camera
}

//MouseState holds mouse state for GLFW polling
type MouseState struct {
	Hold     bool
	Position [2]int
	PosX     float64
	PosY     float64
	XRot     float32
	YRot     float32
	Dx       float32
	Dy       float32
	Time     time.Time
}

//KeyState holds Key pressed state for GLFW polling
type KeyState struct {
	Pressed int
	Time    time.Time
	Vec     mgl.Vec
	Scale   float32
	Select  int
}

type InputState struct {
	Keys      KeyState
	Mouse     MouseState
	Time      time.Time
	DebugPass bool
}

//Exposes static OpenGL Input as static global variable for GLFW threaded calls
var Input *InputState = new(InputState)

func InitGLFW(width int, height int, title string) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("glr | Setup() - Failed glfw.Init()\n")
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	window.SetMouseButtonCallback(ProcessMouse)
	window.SetCursorPosCallback(ProcessCursor)
	window.SetKeyCallback(ProcessInput)
	return window, nil
}

func InitOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL Version: %s\n", version)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.DEPTH)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.15, 0.15, 0.15, 1.0)
	return nil
}

func Renderer() *GLRenderer {

	mCamera := camera.NewCamera(mgl.Vec{0, 0, 10})
	mRenderer := new(GLRenderer)
	mRenderer.Camera = &mCamera
	return mRenderer
}

//Setup() - Setup GL Renderer
func (renderer *GLRenderer) Setup(width int, height int, title string) error {
	var err error

	if height == 0 || width == 0 {
		return fmt.Errorf("Invalid height/width parameter")
	}

	if renderer.Window, err = InitGLFW(width, height, title); err != nil {
		return err
	}

	InitOpenGL()

	aspect := float32(width) / float32(height)
	renderer.MVPMat = make([]float32, 16)
	renderer.MVPMat = mgl.ProjectionMatF(45.0, aspect, 1.0, 1000)
	Input.Time = time.Now()

	return nil

}

func (renderer *GLRenderer) AddShader(path string, gl_shader_type uint32) (uint32, error) {
	var shaderSource string

	if (gl_shader_type != gl.VERTEX_SHADER) && (gl_shader_type != gl.FRAGMENT_SHADER) {
		return 0, fmt.Errorf("glr | AddShader() - Shader type unspecified\ngl.VERTEX_SHADER or gl.FRAGMENT_SHADER int")
	}

	cSourceString, err := ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	shaderSource = string(cSourceString) + "\x00"

	sho, err := compileShader(shaderSource, gl_shader_type)

	if err != nil {
		return 0, err
	}

	return sho, err
}

func (renderer *GLRenderer) LinkShaders(vertexGLID uint32, fragmentGLID uint32) (uint32, error) {

	prog1 := gl.CreateProgram()
	gl.AttachShader(prog1, vertexGLID)
	gl.AttachShader(prog1, fragmentGLID)
	gl.LinkProgram(prog1)
	if prog1 == gl.INVALID_VALUE || prog1 == gl.INVALID_OPERATION {
		err := fmt.Errorf("Invalid Linking[vert %d,frag %d]\n", vertexGLID, fragmentGLID)
		return prog1, err
	}
	return prog1, nil
}

func (renderer *GLRenderer) GetUniformLocation(program uint32, name string) (int32, error) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if loc == gl.INVALID_VALUE || loc == gl.INVALID_OPERATION {
		err := fmt.Errorf("Uniform Location %s Not Found\n", name)
		return loc, err
	}
	return loc, nil
}

func (renderer *GLRenderer) ShaderLog(programGLID uint32) {
	var logLength = int32(1000)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(programGLID, logLength, nil, gl.Str(log))
	fmt.Printf("%s", log)

	active := int32(0)
	gl.GetProgramiv(programGLID, gl.ACTIVE_UNIFORMS, &active)
	fmt.Printf("SHADER UNIFORMS[%d]\n", active)
}

/* Layout(num_vbo, num_tex)( error) - initializes as many VBOs and Textures as
requested by the RenderSystem caller. The Raw OpenGL handles are then copied into corresponding
VAO/VBO/Tex GLRenderer state objects. Initialize VAO State layouts for position, color, tex data
descriptions.
*/
func (renderer *GLRenderer) Layout(num_vao int, num_vbo int, num_tex int) error {

	if num_vbo < 0 || num_tex < 0 || num_vao < 0 {
		r := fmt.Errorf("glr | Layout() - Rendersystem tried to initiate GPU memory layout with zero size param\n")
		return r
	}

	renderer.Tex = make([]uint32, num_tex+1)
	renderer.VAO = make([]uint32, num_vao+1)
	renderer.VBO = make([]uint32, num_vbo+1)
	renderer.ElementVBO = make([]uint32, num_vao+1)

	gl.GenVertexArrays(int32(num_vao), &renderer.VAO[0])
	gl.GenTextures(int32(num_tex), &renderer.Tex[0])
	gl.GenBuffers(int32(num_vbo), &renderer.VBO[0])
	gl.GenBuffers(int32(num_vao), &renderer.ElementVBO[0])

	return nil
}

func (renderer *GLRenderer) BindVertexArray(vao_index int) {
	if vao_index == -1 {
		gl.BindVertexArray(0)
	} else {
		gl.BindVertexArray(renderer.VAO[vao_index])
	}
}

func (renderer *GLRenderer) VertexArrayAttr(index int, components int, gltype uint32, offset int) {
	gl.VertexAttribPointer(uint32(index), int32(components), gltype, false, 0, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(uint32(index))
}

/*
BufferArrayData(vboID, width, byteSize, ref) - binds a VBO and fills its byte data
using the simplest terms possible. The rendersystem will need to pass in the unique buffer ID
it wants to fill. This VBO ID is not the GL generated VBO id but the reference index. The ref array
will need to have the precomputed offset already computed so that ref[0] is the valid buffer data position.
*/
func (renderer *GLRenderer) BufferArrayData(vboID int, width int, offset int, ref []byte) error {
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO[vboID])
	gl.BufferData(gl.ARRAY_BUFFER, width, gl.Ptr(&ref[offset]), gl.STATIC_DRAW)
	return nil
}

/*
BufferIndexData(vboID, width, byteSize, ref) - binds a VBO and fills its byte data
using the simplest terms possible. The rendersystem will need to pass in the unique buffer ID
it wants to fill. This VBO ID is not the GL generated VBO id but the reference index. The ref array
will need to have the precomputed offset already computed so that ref[0] is the valid buffer data position.
*/
func (renderer *GLRenderer) BufferIndexData(vboID int, width int, offset int, ref []byte) error {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderer.ElementVBO[vboID])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, width, gl.Ptr(&ref[offset]), gl.STATIC_DRAW)
	return nil
}

func (renderer *GLRenderer) InvalidateBuffer(vboID int, width int) error {
	gl.MapBufferRange(gl.MAP_INVALIDATE_BUFFER_BIT, 0, width, renderer.VBO[vboID])
	return nil
}

func (renderer *GLRenderer) BindArrayBuffer(vbo_index int) {
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO[vbo_index])
}

func (r *GLRenderer) SwapBuffers() {
	r.Window.SwapBuffers()
}

func (r *GLRenderer) LoadTexture(path string, index int) (uint32, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", path, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("Unsupported Stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return texture, nil
}

func (r *GLRenderer) LoadEnvironment(path string, index int) (uint32, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", path, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("Unsupported Stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	for i := 0; i < 6; i++ { //Load Faces - Properly load with proper cubic offsets
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0,
			gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

	return texture, nil
}

func (r *GLRenderer) SetActiveTexture(glid uint32, index uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + index)
	gl.BindTexture(gl.TEXTURE_2D, glid)

}

func (renderer *GLRenderer) Draw(mesh_entities []*defs.MeshEntity, shaders *defs.ShadersMap, materials map[string]*defs.Material, mapID uint32) error {

	mglView := renderer.Camera.Update()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(shaders.ProgramID["default"])

	gl.UniformMatrix4fv(shaders.ShaderUniforms["mvp"], 1, false, &renderer.MVPMat[0])
	gl.UniformMatrix4fv(shaders.ShaderUniforms["view"], 1, false, &mglView[0])
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, mapID)

	for i := 0; i < len(mesh_entities); i++ {

		myMesh := mesh_entities[i]
		matname := myMesh.Mesh.MaterialComponent.Name
		material := materials[matname]
		metal := material.MetallicRoughMaterial.Metallic
		rough := material.MetallicRoughMaterial.Roughness

		gl.Uniform1f(shaders.ShaderUniforms["metallness"], metal)
		gl.Uniform1f(shaders.ShaderUniforms["roughness"], rough)

		gl.BindVertexArray(renderer.VAO[myMesh.VAO])
		gl.UniformMatrix4fv(shaders.ShaderUniforms["model"], 1, false, &myMesh.Mesh.TransformComponent.Model[0])
		gl.DrawElements(gl.TRIANGLES, int32(myMesh.Mesh.IndiceComponent.IndexByteLength/BYTES_PER_SHORT), gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.BindVertexArray(0)
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

	return nil

}

func (r *GLRenderer) Update(dt float64) error {
	elapsed := time.Now().Sub(Input.Time).Seconds()
	if elapsed > dt {
		if Input.Keys.Pressed == 1 {
			dir := r.Camera.Transform.Matrix.Get(Input.Keys.Select)
			scale := CAMERA_SPEED * Input.Keys.Scale
			r.MoveCamera(dir, scale)
		}
		if Input.Mouse.Hold == true {
			r.Camera.RotateFPS(mgl.Vec{float32(-Input.Mouse.Dy), float32(-Input.Mouse.Dx), 0})
		}
		Input.Time = time.Now()
	}
	return nil
}

func (renderer *GLRenderer) ShouldClose() bool {
	window := renderer.Window
	return window.ShouldClose()
}

func (renderer *GLRenderer) Status() error {

	if renderer.Window == nil {
		return fmt.Errorf("glr | Status() - MyInput.Window not available (GLFW not intialized properly)\n")
	}

	return nil

}

func (r *GLRenderer) MoveCamera(dir mgl.Vec, mag float32) {
	move := mgl.Scale(dir, mag)
	r.Camera.Transform.Translate(move)
}

func ProcessInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		Input.Keys.Pressed = 1
		Input.Keys.Time = time.Now()
	} else {
		Input.Keys.Pressed = 0
	}

	if action == glfw.Repeat {
		Input.Keys.Pressed = 1
	}

	if key == glfw.KeyW {
		Input.Keys.Select = 2
		Input.Keys.Scale = -1.0
	}
	if key == glfw.KeyS {
		Input.Keys.Select = 2
		Input.Keys.Scale = 1.0
	}
	if key == glfw.KeyA {
		Input.Keys.Select = 0
		Input.Keys.Scale = -1.0
	}
	if key == glfw.KeyD {
		Input.Keys.Select = 0
		Input.Keys.Scale = 1.0
	}
	if key == glfw.KeyUp {
		Input.Keys.Select = 1
		Input.Keys.Scale = 1.0
	}
	if key == glfw.KeyDown {
		Input.Keys.Select = 1
		Input.Keys.Scale = -1.0
	}
	if key == glfw.KeyTab {
		//Nothing
	}
}

//ProcessMouse sets the mouse state structure during click events
func ProcessMouse(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		if !Input.Mouse.Hold {
			Input.Mouse.Hold = true
			Input.Mouse.Time = time.Now()
		}
	}
	if button == glfw.MouseButtonLeft && action == glfw.Release {
		Input.Mouse.Hold = false
		Input.Mouse.PosX = 0
		Input.Mouse.PosY = 0
		Input.Mouse.Time = time.Now()
	}
}

//ProcessCursor - GLFW Callback. While mouse is in state hold calculate a X,Y derivative
func ProcessCursor(w *glfw.Window, xPos float64, yPos float64) {
	if Input.Mouse.Hold {
		dt := float64(Input.Mouse.Time.Sub(time.Now()).Seconds())
		if !(Input.Mouse.PosX == 0) && !(Input.Mouse.PosY == 0) {
			Input.Mouse.Dx = float32((Input.Mouse.PosX - xPos) / (dt * 500))
			Input.Mouse.Dy = float32((Input.Mouse.PosY - yPos) / (dt * 500))
			Input.Mouse.Time = time.Now()
		}
		Input.Mouse.PosX = xPos
		Input.Mouse.PosY = yPos
	} else {
		Input.Mouse.Dx = 0
		Input.Mouse.Dy = 0
	}

}
