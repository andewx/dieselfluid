package model

//Fixed Particle structure array
type Particle struct {
	Position [3]float32
	Velocity [3]float32
	Force    [3]float32
	Density  float32
	Hpa      float32
}

//Add vector b to a. Mutates a
func add_float3(a [3]float32, b [3]float32) {
	a[0] += b[0]
	a[1] += b[1]
	a[2] += b[2]
}

func (p Particle) AddPosition(x [3]float32) {
	add_float3(p.Position, x)
}

func (p Particle) AddForce(x [3]float32) {
	add_float3(p.Force, x)
}

func (p Particle) AddVelocity(x [3]float32) {
	add_float3(p.Velocity, x)
}

func (p Particle) Pressure(d0 float32, p0 float32) float32 {
	p.Hpa = TaitEos(p.Density, d0, p0)
	return p.Hpa
}
