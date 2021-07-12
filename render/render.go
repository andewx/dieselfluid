package render

//Realtime Renderer

//TexCoord Render Structure
type TexCoord struct {
	coord [2]float32
}

//RenderContext Holds Graphics API Draw Parameters
type Context struct {
	//Shader
	VertexShaderPaths   map[string]string
	FragmentShaderPaths map[string]string
	ShaderPrograms      map[string]uint32
	FragShaderID        map[string]uint32
	VertShaderID        map[string]uint32
	ProgramID           map[string]uint32
	ShaderUniforms      map[string]uint32

	//Vertex Buffers / Indice Buffers
	Vertexes map[string][]float32
	Indices  map[string][]uint32
	VBO      map[string]uint32
	VAO      map[string]uint32
	IBO      map[string]uint32
}

//Renderer API Call Context
type Renderer interface {
	Setup() error
	Init() error
	Draw() error
	Update() error
	Close() error
}
