package model

//Fixed Particle structure array
type Particle struct {
	Position [3]float32
	Velocity [3]float32
	Force    [3]float32
	Density  float32
	Press    float32
}

//Add vector b to a. Mutates a
func add_float3(a [3]float32, b [3]float32) [3]float32 {
	a[0] += b[0]
	a[1] += b[1]
	a[2] += b[2]
	return a
}

func (p *Particle) AddPosition(x [3]float32) {
	p.Position = add_float3(p.Position, x)
}

func (p *Particle) AddForce(x [3]float32) {
	p.Force = add_float3(p.Force, x)
}

func (p *Particle) AddVelocity(x [3]float32) {
	p.Velocity = add_float3(p.Velocity, x)
}

func (p *Particle) Pressure(d0 float32, p0 float32) float32 {
	p.Press = TaitEos(p.Density, d0, p0)
	return p.Press
}

func (p *Particle) SetPosition(x []float32) {
	p.Position = [3]float32{x[0], x[1], x[2]}
}

func (p *Particle) SetVelocity(x []float32) {
	p.Velocity = [3]float32{x[0], x[1], x[2]}
}

func (p *Particle) SetForce(x []float32) {
	p.Force = [3]float32{x[0], x[1], x[2]}
}
