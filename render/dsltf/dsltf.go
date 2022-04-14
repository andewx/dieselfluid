package dsltf

/*
Corrollary to the gltf module with concrete underlying datatypes after
type assertions have been linted from the gltf asset. Modification of the asset
in memory is achieved with pointer representations with the marhsalling of the object
being handled by the encapsulated GLTF holding

This file handles the struct reqs' for DSLTF concrete datatypes (//no interface{} types)
Should be 1-1 mapping controlling the GLTF asset in memory.

Immutable GLTF properties will be handled without a pointer
*/
//Accessors and buffers

type Buffer struct {
	Uri        string
	ByteLength uint32
}

//Accessor should be immutable to the GLTF asset
type Accessor struct {
	BufferView    uint32
	ByteOffset    uint32
	Type          string
	ComponentType uint32
	Count         uint32
	Min           []float32
	Max           []float32
}

type BufferView struct {
	Buffer     uint32
	ByteOffset uint32
	ByteLength uint32
	ByteStride uint32
	Target     uint32
}

//Meshes

type Primitive struct {
	Indices  uint32
	Mode     uint32
	Material uint32
	Targets  map[string]uint32 //Map attributes to morph target accessors
	Weights  []float32         //Mesh Weights
}

type Mesh struct {
	Attributes map[string]uint32 // Accessor Index for Mesh attributes ("POSITION", "NORMAL", etc)
	Primitives []Primitive
}

//Scenes
type Node struct {
	Name        string
	Mesh        uint32
	Camera      uint32
	Skin        uint32
	Matrix      []float32
	Translation []float32
	Rotation    []float32
	Scale       []float32
	Children    []uint32
}

type Scene struct {
	Nodes []uint32
}

//PBR Materials

type ColorTexture struct {
	Scale    float32
	Strength float32
	Index    uint32
	TexCoord uint32
}

type PbrMetallicRough struct {
	BaseColorTexture         ColorTexture
	BaseColorFactor          []float32
	MetallicRoughnessTexture ColorTexture
	MetallicFactor           float32
	RoughnessFactor          float32
	NormalTexture            ColorTexture
	EmissiveTexture          ColorTexture
	OcclusionTexture         ColorTexture
	EmmissiveFactor          []float32
}

type Material struct {
	PbrMetallicRoughness PbrMetallicRough
}

//Cameras

type PersCamera struct {
	AspectRatio float32
	Yfov        float32
	Zfar        float32
	Znear       float32
}

type OrthCamera struct {
	Xmag  float32
	Ymag  float32
	Zfar  float32
	Znear float32
}

type Camera struct {
	Type         string
	Perpective   PersCamera
	Orthographic OrthCamera
}

//Textures
type TextureDescriptor struct {
	Source  uint32
	Sampler uint32
}

type Image struct {
	Uri        string
	BufferView uint32
	MimeType   string
}

type Sampler struct {
	MagFilter uint32
	MinFilter uint32
	WrapS     uint32
	WrapT     uint32
}

//SKINS and Animations
type Skin struct {
	InverseBindMatrices uint32 //Accessor whic references Joint Matrices buffers
	Joints              []uint32
}

type TargetAnim struct {
	Node uint32
	Path string
}

type AnimSampler struct {
	Input          uint32
	Interpoloation string
	Output         uint32
}
type Channel struct {
	Target []TargetAnim
}

type Animation struct {
	Channels []Channel
	Samplers []AnimSampler
}

//DSLTF Main Root Node
type DSLTF struct {
	Scenes       []Scene
	Nodes        []Node
	Cameras      []Camera
	Meshes       []Mesh
	Buffers      []Buffer
	BufferView   []BufferView
	Material     []Material
	Textures     []TextureDescriptor
	Images       []Image
	Sampler      []Sampler
	Skins        []Skin
	Animations   []Animation
	MasterBuffer []byte
}
