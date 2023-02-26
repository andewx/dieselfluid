package model

import "fmt"

type ParticleArray struct {
	positions        []float32
	velocities       []float32
	densities        []float32
	forces           []float32
	pressures        []float32
	n_particles      int
	n_boundary       int
	mass             float32
	ReferenceDensity float32
}

//Returns the particle array given n particles and n_boundary particles
func NewParticleArray(n_particles int, n_boundary int, kernel float32, density float32, mass float32) ParticleArray {
	parray := ParticleArray{}
	parray.positions = make([]float32, (n_particles+n_boundary)*3)
	parray.velocities = make([]float32, (n_particles)*3)
	parray.densities = make([]float32, (n_particles))
	parray.forces = make([]float32, (n_particles)*3)
	parray.pressures = make([]float32, (n_particles))
	parray.mass = mass
	parray.ReferenceDensity = density * parray.mass
	parray.n_particles = n_particles
	parray.n_boundary = n_boundary

	fmt.Printf("\nSPH Particle Profile:\nParticles[%d]\nBoundary Particles[%d]\nMass:[%.2f]\nDensity[%.2f]\n",
		parray.n_particles, parray.n_boundary, parray.mass, parray.ReferenceDensity)
	return parray
}

func (p *ParticleArray) SetReferenceDensity(d float32) {
	p.ReferenceDensity = d
}

func (p *ParticleArray) Positions() []float32 {
	return p.positions
}
func (p *ParticleArray) Velocities() []float32 {
	return p.velocities
}
func (p *ParticleArray) Densities() []float32 {
	return p.densities
}

func (p *ParticleArray) Pressures() []float32 {
	return p.pressures
}
func (p *ParticleArray) Forces() []float32 {
	return p.forces
}
func (p *ParticleArray) Pressure(index int, d0 float32, p0 float32) float32 {
	return TaitEos(p.densities[index], d0, p0)
}
func (p *ParticleArray) Density(index int) float32 {
	return p.densities[index]
}
func (p *ParticleArray) Position(index int) []float32 {
	x := index * 3
	if (x+2) > len(p.positions) || x < 0 {
		return []float32{0, 0, 0}
	}
	return []float32{p.positions[x], p.positions[x+1], p.positions[x+2]}
}
func (p *ParticleArray) Velocity(index int) []float32 {
	x := index * 3
	if (x+2) > len(p.velocities) || x < 0 {
		return []float32{0, 0, 0}
	}
	return []float32{p.velocities[x], p.velocities[x+1], p.velocities[x+2]}
}

func (p *ParticleArray) Force(index int) []float32 {
	x := index * 3
	if (x+2) > len(p.forces) || x < 0 {
		return []float32{0, 0, 0}
	}
	return []float32{p.forces[x], p.forces[x+1], p.forces[x+2]}
}
func (p *ParticleArray) Mass() float32 {
	return p.mass
}
func (p *ParticleArray) Set(index int, particle Particle) {
	x := index * 3
	Float3_buffer_set(x, p.positions, &particle.Position)
	Float3_buffer_set(x, p.velocities, &particle.Velocity)
	Float3_buffer_set(x, p.forces, &particle.Force)
	p.densities[index] = particle.Density
	p.pressures[index] = particle.Press
}
func (p *ParticleArray) Get(index int) Particle {
	x := index * 3
	particle := Particle{}

	if index > p.n_particles && index < p.Total() {
		Float3_set(x, &particle.Position, p.positions)
		particle.Press = 0
		particle.Density = 0
		particle.Velocity = [3]float32{0, 0, 0}
		particle.Force = [3]float32{0, 0, 0}
		return particle
	}

	if index < p.n_particles {
		Float3_set(x, &particle.Position, p.positions)
		Float3_set(x, &particle.Velocity, p.velocities)
		Float3_set(x, &particle.Force, p.forces)
		particle.Density = p.densities[index]
		particle.Press = p.pressures[index]
		return particle
	}

	return particle
}
func (p *ParticleArray) D0() float32 {
	return p.ReferenceDensity
}

//Adds in boundary particle buffer of positions
func (p *ParticleArray) AddBoundaryParticles(positions []float32) []float32 {
	p.n_boundary += len(positions) / 3
	p.positions = append(p.positions, positions...)
	return p.positions

}

func (p *ParticleArray) N() int {
	return p.n_particles
}

func (p *ParticleArray) Total() int {
	return p.n_particles + p.n_boundary
}
