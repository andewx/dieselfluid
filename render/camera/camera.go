package render

//Camera defines GL camera orientation
type Camera struct {
	Pos   [3]float32
	Front [3]float32
	Up    [3]float32
	Info  string
}
