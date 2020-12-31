package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	"github.com/andewx/dieselfluid/kernel"
	V "github.com/andewx/dieselfluid/math/math64" //Diesel Vector Library - Simple Vec
	"github.com/andewx/dieselfluid/sampler"
)

//Define Particle Attribute Types
type Position V.Vec
type Density float64
type Velocity V.Vec
type Force V.Vec
type Pressure float64

//Particle System Interface
type ParticleSystem interface {
	Positions() []Position
	Velocities() []Velocity
	Densities() []Density
	Forces() []Force
	Pressures() []Pressure
	TimeStep() float64
	UpdateTime() float64
	Density(positions []V.Vec, density_field []float64)
	Gradient(scalar []float64, vector_gradient []V.Vec)
	Div(vector_field []V.Vec, div_field []float64)
	Laplacian(vector_field []V.Vec, force_field []V.Vec)
	Curl(vector_field []V.Vec, scalar_field []float64)
	Length() int
}

//Defines SPH Particle System - Implements ParticleSystem Interface Which Contains some Generalized
//functionality for field computation. Which allows reuse by implementations of SPH. For example
//A DFSPH / WCSPH / IISPH / PCISPH Can utilize the SPH Particle System Layout and compute the Gradients
//and field derivatives consistently. Also Particle In Cell Methods Benefit
type SPHParticleSystem struct {
	Pos      []Position
	Dens     []Density
	Vels     []Velocity
	Fs       []Force
	Ps       []Pressure
	Particle SPHParticle
	Time     float64
	MaxVel   float64
	Kern     kernel.Kernel
	Smplr    sampler.Sampler
}

//-----------Implements SPHParticleSystem -----------------------------
func (p SPHParticleSystem) Positions() []Position {
	return p.Pos
}
func (p SPHParticleSystem) Velocities() []Velocity {
	return p.Vels
}
func (p SPHParticleSystem) Densities() []Density {
	return p.Dens
}
func (p SPHParticleSystem) Forces() []Force {
	return p.Fs
}
func (p SPHParticleSystem) Pressures() []Pressure {
	return p.Ps
}
func (p SPHParticleSystem) TimeStep() float64 {
	return p.Time
}
func (p SPHParticleSystem) UpdateTime() float64 {
	K := 0.4 * p.Particle.KernelRad() / p.Particle.Sound()
	A := p.Particle.KernelRad() / p.MaxVel

	if A < K {
		p.Time = A
		return A
	}
	p.Time = K
	return K
}

//Density -- Computes density field for SPH Particle System
func (p SPHParticleSystem) Density(positions []V.Vec, density_field []float64) {
	N := len(positions)
	for i := 0; i < N; i++ {
		sampleList := p.Smplr.GetSamples(i)
		weight := p.Kern.W0()

		for j := 0; j < len(sampleList); j++ {
			pIndex := sampleList[j]
			if i != j {
				dist := V.Dist(positions[i], positions[j])
				weight += p.Kern(dist)
			}
		}
		density := p.Particle.Mass() * weight
		density_field[i] = density
	}
}

func (p SPHParticleSystem) Gradient(scalar []float64, vector_gradient []V.Vec) {

}
func (p SPHParticleSystem) Div(vector_field []V.Vec, div_field []float64) {

}
func (p SPHParticleSystem) Laplacian(vector_field []V.Vec, force_field []V.Vec) {

}
func (p SPHParticleSystem) Curl(vector_field []V.Vec, scalar_field []float64) {

}

//---------------------MOVE to Parallel and Threading Class -----------------//
/*
//Vec_IterParallel Spawns System Threads for an iteration loop over a list
//With implementing callback functions on that list
func (p SPHParticleSystem)Vec_IterParallel(callback func(vector_list []V.Vec)V.Vec){
        //Spawn Thread Functions For Each List Section N/(System_Processors)
}

//Vec_IterParallel Spawns System Threads for an iteration loop over a list
//With implementing callback functions on that list
func (p SPHParticleSystem)Scalar_IterParallel(callback func(scalar_list[]float64)float64){
    //Spawn Thread Functions For Each List Section N/(System_Processors)
}
*/
