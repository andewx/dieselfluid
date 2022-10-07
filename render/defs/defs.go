package defs

import (
	"fmt"
	"unsafe"

	"github.com/EngoEngine/ecs"
)

//Render component defintions pacakge, decouples type definitions allows for global
//package use of definitions, exports effectively to top level module
//When render component definitions need low and high level module access please
//place module definition here

type VAOIndex []uint32
type VBOIndex []uint32
type TEXIndex []uint32

//TexCoord Render Structure
type TexCoord struct {
	coord [2]float32
}

//RenderContext Holds Graphics API Draw Parameters
type ShadersMap struct {
	//Shader
	VertexShaderPaths   map[string]string
	FragmentShaderPaths map[string]string
	ShaderPrograms      map[string]uint32
	FragShaderID        map[string]uint32
	VertShaderID        map[string]uint32
	ProgramID           map[string]uint32
	ShaderUniforms      map[string]int32
	ProgramLinks        map[uint32]uint32
}

type TexturesMap struct {
	TexID      []uint32
	TexIDMap   map[string]uint32
	TexUnit    []int32
	TexUnitMap map[string]uint32
	Names      []string
}

//OGLRenderer interface defines basic named operations for OpenGL renderer interaction
//via the RenderSystem component. Allows rendersystem to work with the heterogenous OGL
//API without a corrollary to the CommandBuffer type GPU APIs such as Vulkan, DX12, METAL
//Those renderer types are defined seperately in the CommandRenderer interace
type OGLRenderer interface {
	Setup(width int, height int, name string) error
	AddShader(path string, gl_shader_type uint32) (uint32, error)
	LinkShaders(vertexGLID uint32, fragmentGLID uint32) (uint32, error)
	GetUniformLocation(program uint32, name string) (int32, error)
	Layout(num_vao int, num_vbo int, num_tex int) error
	BufferArrayData(vboID int, width int, offset int, ref []byte) error
	BufferArrayFloat(vboID int, width int, offset int, ref []float32) error
	BufferArrayPointer(vboID int, width int, ref unsafe.Pointer) error
	BufferIndexData(vboID int, width int, offset int, ref []byte) error
	InvalidateBuffer(vboID int, width int) error
	Draw(mesh_entities []*MeshEntity, shaders *ShadersMap, materials map[string]*Material, mapid uint32) error
	VertexArrayAttr(index int, components int, gltype uint32, offset int)
	BindVertexArray(vao_index int)
	BindArrayBuffer(vbo_index int)
	LoadTexture(path string, index int) (uint32, error)
	LoadEnvironment(path string, index int) (uint32, error)
	SetActiveTexture(glid uint32, index uint32)
	Update(dt float64) error
	ShouldClose() bool
	SwapBuffers()
	Status() error
	GetVBO(int) uint32
	GetVAO(int) uint32
}

//PipelineRenderer interface deals with DX12, METAL, VULKAN like APIs that support
//more similar APIs for interactions with the GPU. This pipeline renderer remains
//abstract but the interface is declared for future support and testing purposes.
type PipelineRenderer interface {
	InitPipeline() error //thats all i got lol
}

type IndiceComponent struct {
	Indices         []byte
	IndexByteOffset int
	IndexByteLength int
	IndexBufferId   uint32
	IndicesVBO      int
}

type VertexComponent struct {
	Vertices         []byte //Float 4 Bytes * 3
	VertexByteOffset int
	VertexByteLength int
	VertexBufferId   uint32
	VerticesVBO      int
}

type NormalComponent struct {
	Normals           []byte //Float 4 Bytes * 3
	NormalsByteOffset int
	NormalsByteLength int
	NormalsBufferId   uint32
	NormalsVBO        int
}

type ColorComponent struct {
	Colors           []byte //Float 4 Bytes * 3
	ColorsByteOffset int
	ColorsByteLength int
	ColorsBufferId   uint32
	ColorsVBO        int
}

type TexComponent struct {
	TexCoords           []byte
	TexCoordsByteOffset int
	TexCoordsByteLength int
	TexCoordsBufferId   uint32
	TextureIDs          []uint32
	TexVBO              int
}

type ShaderComponent struct {
	ProgramID uint32
	Uniforms  []string
}

//Transform stores the mesh model transformation along with the "Position"
//Position is once calc'd center of mass for the mesh and used for OCTREE
//Identification of the Mesh
type TransformComponent struct {
	Model    []float32
	Position []float32
}

//CollisionMeshComponent stores a Reduced low poly bounds mesh for standard type
//explicit collision handling
type CollisionMeshComponent struct {
	Indice    IndiceComponent
	Vertex    VertexComponent
	Normal    NormalComponent
	Transform TransformComponent
}

type Texture struct {
	Index        int    //Image Resources Index as loaded
	GLID         uint32 //GL Texture ID
	GLTEXUNIT    uint32
	TEXCORRDATTR int    //Tex Coord Attribute
	Key          string //Associated Mapping Key
	Status       int    //Texture Status

}

type RGB struct {
	R float32
	G float32
	B float32
}

type Material struct {
	AlphaCutoff    float32
	AlphaMode      int
	DoubleSided    int
	BaseColor      RGB
	SpecularColor  RGB
	EmissiveFactor RGB

	ColorTexture     *Texture
	NormalTexture    *Texture
	OcclusionTexture *Texture
	EmmissiveTexture *Texture

	MetallicRoughMaterial *PBRMaterial

	Name string
}

type MaterialComponent struct {
	Name string
}

//Light color should be treated as a wattage output per unit surface area
//Given the initial positional conditions.
type Light struct {
	Pos   []float32 //vec4
	Color []float32 //vec3
}

type PBRMaterial struct {
	MetallicRoughnessTexture *Texture
	Metallic                 float32
	Roughness                float32
	Anisoptric               float32
}

//MeshEntity is a non-exported Entity component for ECS Systems
type MeshEntity struct {
	Basic ecs.BasicEntity
	Mesh  MeshComponent
	VAO   int //VAO Index
}

type MeshComponent struct {
	*IndiceComponent
	*VertexComponent
	*NormalComponent
	*ColorComponent
	*TransformComponent
	*TexComponent
	*ShaderComponent
	*MaterialComponent
}

func NewPBRMaterial(color RGB, metalness float32, roughness float32) *PBRMaterial {
	newPBR := PBRMaterial{nil, metalness, roughness, 0.0}
	return &newPBR
}

func NewMaterial(name string, base RGB) *Material {
	mat := Material{0, 0, 0, base, base, base, nil, nil, nil, nil, nil, name}
	return &mat
}

//-----------Debug Struct Objects------------------------//
func (m *MeshEntity) PrintDebug() {
	fmt.Printf("\nMesh Entity(%d): Vao(%d)\n", m.Basic.ID(), m.VAO)
	fmt.Printf("Vertex Byte Offset: %d\n", m.Mesh.VertexComponent.VertexByteOffset)
	fmt.Printf("Vertex Byte Length: %d\n", m.Mesh.VertexComponent.VertexByteLength)
	fmt.Printf("Vertex VBO: %d\n\n", m.Mesh.VertexComponent.VerticesVBO)
	fmt.Printf("Normal Offset: %d\n", m.Mesh.NormalComponent.NormalsByteOffset)
	fmt.Printf("Normal Byte Length: %d\n", m.Mesh.NormalComponent.NormalsByteLength)
	fmt.Printf("Normal VBO: %d\n\n", m.Mesh.NormalComponent.NormalsVBO)
	fmt.Printf("TexCoord Byte Offset: %d\n", m.Mesh.TexComponent.TexCoordsByteOffset)
	fmt.Printf("TexCoord Byte Length: %d\n", m.Mesh.TexComponent.TexCoordsByteLength)
	fmt.Printf("TexCoord VBO: %d\n\n", m.Mesh.TexComponent.TexVBO)
	fmt.Printf("Indice Byte Offset: %d\n", m.Mesh.IndiceComponent.IndexByteOffset)
	fmt.Printf("Indice Byte Length: %d\n", m.Mesh.IndiceComponent.IndexByteLength)
	fmt.Printf("Indice VBO: %d\n\n", m.Mesh.IndiceComponent.IndicesVBO)
	fmt.Print("--------------------\n")

}
