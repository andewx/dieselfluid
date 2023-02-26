package model

import (
	"fmt"

	"github.com/andewx/dieselfluid/math/vector"
)

type ParticleStructField struct {
	Particles        []Particle
	Boundary         []float32
	n_particles      int
	n_boundary       int
	mass             float32
	ReferenceDensity float32
}

func NewParticleStructField(n int, kernel float32, density float32, mass float32) ParticleStructField {
	particles := ParticleStructField{}
	particles.Particles = make([]Particle, n)
	particles.n_particles = n
	particles.mass = mass
	particles.ReferenceDensity = density * mass
	return particles
}

func (p ParticleStructField) Positions() []float32 {
	fmt.Printf("Warning duplicating particle field property of length  %d KB", int((p.n_particles+p.n_boundary)*3*4/1024))
	alloc := make([]float32, (p.n_particles*p.n_boundary)*3)
	for index, particle := range p.Particles {
		x := index * 3
		alloc[x] = particle.Position[0]
		alloc[x+1] = particle.Position[1]
		alloc[x+2] = particle.Position[2]
	}

	for i := p.n_particles; i < p.n_particles+p.n_boundary; i++ {
		x := i * 3
		alloc[x] = p.Particles[i].Position[0]
		alloc[x+1] = p.Particles[i].Position[1]
		alloc[x+2] = p.Particles[i].Position[2]
	}
	return alloc
}

func (p ParticleStructField) Velocities() []float32 {
	fmt.Printf("Warning duplicating particle field property of length  %d KB", int(p.n_particles*3*4/1024))
	alloc := make([]float32, p.n_particles*3)
	for index, particle := range p.Particles {
		x := index * 3
		alloc[x] = particle.Velocity[0]
		alloc[x+1] = particle.Velocity[1]
		alloc[x+2] = particle.Velocity[2]
	}
	return alloc
}

func (p ParticleStructField) Densities() []float32 {
	fmt.Printf("Warning duplicating particle field property of length  %d KB", int(p.n_particles*3*4/1024))
	alloc := make([]float32, p.n_particles)
	for index, particle := range p.Particles {
		alloc[index] = particle.Density
	}
	return alloc
}

func (p ParticleStructField) Pressures() []float32 {
	fmt.Printf("Warning duplicating particle field property of length  %d KB", int(p.n_particles*3*4/1024))
	alloc := make([]float32, p.n_particles)
	for index, particle := range p.Particles {
		alloc[index] = particle.Press
	}
	return alloc
}

func (p ParticleStructField) Forces() []float32 {
	fmt.Printf("Warning duplicating particle field property of length  %d KB", int(p.n_particles*3*4/1024))
	alloc := make([]float32, p.n_particles*3)
	for index, particle := range p.Particles {
		alloc[index] = particle.Force[0]
		alloc[index+1] = particle.Force[1]
		alloc[index+2] = particle.Force[2]
	}
	return alloc
}

func (p ParticleStructField) Pressure(x int, d0 float32, p0 float32) float32 {
	return TaitEos(p.Particles[x].Density, p.ReferenceDensity, 0.0)
}

func (p ParticleStructField) Density(x int) float32 {
	return p.Particles[x].Density
}
func (p ParticleStructField) Position(x int) []float32 {
	return vector.Cast(p.Particles[x].Position)
}
func (p ParticleStructField) Velocity(x int) []float32 {
	return vector.Cast(p.Particles[x].Velocity)
}

func (p ParticleStructField) Force(x int) []float32 {
	return vector.Cast(p.Particles[x].Force)
}

func (p ParticleStructField) Mass() float32 {
	return p.mass
}

func (p ParticleStructField) Set(x int, particle Particle) {
	p.Particles[x] = particle
}
func (p ParticleStructField) Get(x int) Particle {
	return p.Particles[x]
}
func (p ParticleStructField) D0() float32 {
	return p.ReferenceDensity
}
func (p ParticleStructField) AddBoundaryParticles([]float32) {

}
func (p ParticleStructField) N() int {
	return p.n_particles
}

func (p ParticleStructField) Total() int {
	return p.n_particles + p.n_boundary
}
