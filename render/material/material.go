package render

type MaterialObject struct {
	ObjectKey  string     //Map name should correspond to all assignments (Vertex, Indice, Tex Buffers)
	ShaderKey  string     //GLSL Shader Key
	ColorRGB   [3]float32 //Object Color
	TextureKey string     //Texture Object Key
	DrawType   uint32     //User defined draw type - 0 - 3DPoint, 1 - 3DTriangle ... ...
}
