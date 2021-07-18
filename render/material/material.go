package render

type TexMaterial struct {
	Name      string
	Type      string
	UVIndex   int
	FilePath  string
	TextureID string
}
type Material struct {
	Name       string
	ShaderID   string
	Textures   []TexMaterial
	Diffuse    [3]float32
	Specular   [3]float32
	Glossiness float32
	Opacity    float32
}
