package render

//Realtime Renderer

//Works with

import V "dslfluid.com/dsl/math/math32"

type TexCoord struct {
	coord [2]float32
}

//0.0 - 1.0 RGB
type RGB32 struct {
	color [3]float32
}

type RGBA32 struct {
	color [4]float32
}

//Normalizes management of Render API Access
type RenderAPIContext interface {
	//Initialize Graphics
	GetContext() bool

	//Shader Access
	GetVertexShaderPaths() []string
	GetFragmentShaderPaths() []string
	GetShaderPrograms() []int
	GetFrags() []int
	GetVerts() []int
	GetVarMap() []map[string]int

	//Vertex Buffers
	GetVertData() []V.Vec
	GetIndiceData() []int
	GetVBO() []int
	GetVAO() []int
	GetIBO() []int

	//Colors
	GetRGB32Data() []RGB32
	GetRGBA32Data() []RGBA32

	//Textures
	GetTexCoords() []TexCoord
	GetTexObjs() []int
	GetTexturePaths() []string
	GetTexData_RGBA32(textureIDs map[string]int) []RGBA32
	GetTexData_RGB32(textureIDs map[string]int) []RGBA32

	//FrameBuffers
	GetFrameData(FBO int) []RGB32
}

type Renderer interface {
	Setup() bool
	Init(context RenderAPIContext) bool
	Draw(context RenderAPIContext, t chan int) bool
	Update(context RenderAPIContext) bool
	Close(context RenderAPIContext) bool
}
