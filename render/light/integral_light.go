package light

import "math"

//Approximates area light for contributions.
type LightIntegral interface {
	Lights() []Light
	Luminance(point vector.Vec) float32
	Position() vector.Vec
}

//Rect lights sample the light boundaries luminance is distributed across the
//planar region. Lights can be handled as virtual area lights with 180 cutoff
//becaues the output energy is distributed evenly by the area the normal light
//point vectors estimate the contribution while output is conserved.
type RectLight struct {
	SampleLx   *Area
	Lum        float32 //Total output
	Pos        vector.Vec
	Width      float32
	Height     float32
	Num_w      float32
	Num_h      float32
	EdgeCutoff float32
	Plane      []vector.Vec //Triangle plane represents object
}

//Cylindrical lights are simply a line of oriented attenuated lights
type CylinderLight struct {
	SampleLx  *Area
	Luminance float32 //Total output
	Pos       vector.Vec
	Axis      vector.Vec
	Length    float32
}

type DiscLight struct {
	SampleLx  *Area
	Luminance float32
	Pos       vector.Vec
	Dir       vector.Vec
	Radius    vector.Vec
}

//-------------------RectLight----------------------------------//

func NewRectLight(color vector.Vec, lum float32, pos vector.Vec, width float32,
	height float32, num_w int, num_h int) *RectLight {
	plane := make([]vector.Vec, 3)
	x1 := (-width) / 2
	y1 := (-height) / 2
	plane[0][0] = x1
	plane[0][1] = y1
	plane[1][0] = x1
	plane[1][1] = -y1
	plane[2][0] = -x1
	plane[2][1] = y1
	myRect := RectLight{&Area{pos, vector.Vec{}, math.Pi / 2, Source{color, lum, 0}}, lum, pos, width, height, float32(num_w), float32(num_h),
		math.Pi / 4, plane}
	return &myRect
}

//Area light is defined by a planar orientation given by the planar matrix
//This calculates all light positions and returns an array for the given light
func (light *RectLight) Lights() []Area {
	num_lights := math.Floor(float64(light.Num_w * light.Num_h))
	if num_lights == 0 {
		return nil
	}

	a := light.Plane[1].Sub(light.Plane[0])
	b := light.Plane[2].Sub(light.Plane[1])
	n := vector.Cross(a, b)

	objMat := vector.Mat3V(a, b, n)

	x_step := float32(light.Width / light.Num_w)
	y_step := float32(light.Width / light.Num_h)

	lightsArray := make([]Area, int(num_lights))
	lum := light.Lum / float32(num_lights)
	for i := 0; i < int(light.Num_w); i++ {
		for j := 0; j < int(light.Num_h); j++ {

			cutoff := float32(math.Pi / 2)
			if j == 0 || i == 0 {
				cutoff = light.EdgeCutoff
			}
			x := (-light.Width / 2.0) + (float32(i) * x_step)
			y := (-light.Height / 2.0) + (float32(j) * y_step)
			pos := vector.Vec{x, y, 0}
			nPos := objMat.CrossVec(pos).Add(light.Pos)

			lgt := Area{nPos, n, cutoff, Source{light.SampleLx.Lx().RGB, lum, 0}}
			lightsArray[i*j+j] = lgt
		}
	}
	return lightsArray
}

//Placeholder body
func (light *RectLight) Luminance(point vector.Vec) float32 {
	return light.Lum
}
func (light *RectLight) Position() vector.Vec {
	return light.Pos
}
