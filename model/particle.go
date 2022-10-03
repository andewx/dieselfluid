package model

import "github.com/andewx/dieselfluid/math/vector"

//Water Model
const FLUID_MASS = 0.1
const FLUID_DENSITY = 1.0
const FLUID_STIFF = 6.1
const FLUID_KERN_RAD = 0.2
const FLUID_SOS = float32(1480.0)
const FLUID_WATER = 0

//SPH Fluid Incrompressible Water Particle with sized arrays
type Particle struct {
	Pos  [3]float32
	Vel  [3]float32
	F    [3]float32
	P    float32
	Dens float32
}

//-----------Implements SPH PARTICLE ------------------------//
func (p Particle) Position() []float32 {
	return vector.Cast(p.Pos)
}
func (p Particle) Velocity() []float32 {
	return vector.Cast(p.Vel)
}
func (p Particle) Density() float32 {
	return p.Dens
}
func (p Particle) Pressure() float32 {
	return p.P
}
func (p Particle) Force() []float32 {
	return vector.Cast(p.F)
}

func (p Particle) AddForce(x []float32) {
	p.F = vector.CastFixed(vector.Add(x, vector.Cast(p.F)))
}

func (p Particle) AddPosition(x []float32) {
	p.Pos = vector.CastFixed(vector.Add(x, vector.Cast(p.Pos)))
}

func (p Particle) AddVelocity(x []float32) {
	p.Vel = vector.CastFixed(vector.Add(x, vector.Cast(p.Vel)))
}

func (p Particle) SetPosition(x []float32) {
	p.Pos = vector.CastFixed(x)
}
func (p Particle) SetVelocity(x []float32) {
	p.Vel = vector.CastFixed(x)
}
func (p Particle) SetDensity(x float32) {
	p.Dens = x
}
func (p Particle) SetPressure(x float32) {
	p.P = x
}
func (p Particle) SetForce(x []float32) {
	p.F = vector.CastFixed(x)
}

func (p Particle) GetParticle() Particle {
	return p
}
