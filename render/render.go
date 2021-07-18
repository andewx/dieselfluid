package render

/*
The rendering module takes the opaque datatype GLTF loaded from GLTF2.0 Schema
and stores their referent properties into the DSLTF objects, with pointers managing
scene data changes and the corresponding marhsalling and unmarshalling of the GLTF
datas.

GLTF file formats while highly transmissible lack certain cues, such as objects
and nodes are typically only referenced by their indexes, while names, ids and maps
make for easier and more understandable linking.

Additionally for the purposes of this project, the rendering system will need to be able
to store fluid descriptors and manage fluid geometry without the use of .bin / buffer
descriptions for rendered geometry.

This decouples somewhat the fluid "rendering" from the object renders where fluids
will need to be able to reference collision data etc... ...
*/

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
	ShaderUniforms      map[string]int32

	//Vertex Buffers / Indice Buffers
	VBO []uint32
	VAO []uint32
}

//Renderer API Call Context
type Renderer interface {
	Setup() error
	Init() error
	Draw() error
	Update() error
	Close() error
}
