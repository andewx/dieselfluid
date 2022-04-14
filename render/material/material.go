package render

type TextureID struct {
	Index    int32
	UVCoord  int32
	Scale    float32
	Strength float32
}

//PBRMaterial Type Represents a ready to use Metallic Roughness Based material
//For rendering
type PBRMaterial struct {
	Name                     string
	ShaderID                 string
	BaseColorTexture         TextureID
	NormalTexture            TextureID
	OcculusionTexture        TextureID
	EmmissiveTexture         TextureID
	AmbientOcclusionTexture  TextureID
	DisplacementTexture      TextureID
	MetallicRoughnessTexture TextureID //Blue Metalness/Green Roughnes
	BaseColor                [4]float32
	MetallicFactor           float32
	RoughnessFactor          float32
}
