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
	Positions() []V.Vec
	Velocities() []V.Vec
	Densities() []float64
	Forces() []V.Vec
	Pressures() []V.Vec
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
	Pos      []V.Vec
	Dens     []float64
	Vels     []V.Vec
	Fs       []V.Vec
	Ps       []float64
	Particle SPHParticle
	Time     float64
	MaxVel   float64
	Kern     kernel.Kernel
	Smplr    sampler.Sampler
}

//-----------Implements SPHParticleSystem -----------------------------
func (p SPHParticleSystem) Positions() []V.Vec {
	return p.Pos
}
func (p SPHParticleSystem) Velocities() []V.Vec {
	return p.Vels
}
func (p SPHParticleSystem) Densities() []float64 {
	return p.Dens
}
func (p SPHParticleSystem) Forces() []V.Vec {
	return p.Fs
}
func (p SPHParticleSystem) Pressures() []float64 {
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
				dist := V.Dst(positions[i], positions[pIndex]) //Change to dist
				weight += p.Kern.F(dist)
			}
		}
		density := p.Particle.Mass() * weight
		density_field[i] = density
	}
}

//Gradient computes SPH Particle scalar gradient dependent on density field
func (p SPHParticleSystem) Gradient(scalar []float64, densities []V.Vec, vector_gradient []V.Vec) {

	for i := 0; i < len(scalar); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		F := float32(0.0)
		p.Kern.Adjust(dens / p.Particle.D0())
		mass := p.Particle.Mass()
		mq2 := -mass * mass
		accumGrad := V.Vec{}
		//Conduct inner loop
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := V.Sub(positions[samples[j]], positions[i])
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := fluid.Kern.Gradient(dist, dir)
				F = ((scalar[i] / (dens * dens)) + (scalar[samples[j]] / (jDensity * jDensity)))
				accumGrad = V.Add(accumGrad, V.Scale(grad, -F))
			}
			//End Inner Loop
		}
		V.Add(vector_gradient[i], V.Scl(accumGrad, mq2))
	} //End Particle Loop

}

//Div computes vector field divergence
func (p SPHParticleSystem) Div(vector_field []V.Vec, densities []float64, div_field []float64) {
	for i := 0; i < len(scalar); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		F := float32(0.0)
		p.Kern.Adjust(dens / p.Particle.D0())
		mass := p.Particle.Mass()
		mq2 := -mass * mass
		accumGrad := V.Vec{}
		//Conduct inner loop
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := V.Sub(positions[samples[j]], positions[i])
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := fluid.Kern.Gradient(dist, dir)
				F = ((scalar[i] / (dens * dens)) + (scalar[samples[j]] / (jDensity * jDensity)))
				accumGrad = V.Add(accumGrad, V.Scale(grad, -F))
			}
			//End Inner Loop
		}
		V.Add(vector_field[i], V.Scl(accumGrad, mq2))
	} //End Particle Loop

}

//Laplacian computes vector field laplacian
func (p SPHParticleSystem) Laplacian(vector_field []V.Vec, densities []float64, force_field []V.Vec) {
	for i := 0; i < len(scalar); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		F := float32(0.0)
		p.Kern.Adjust(dens / p.Particle.D0())
		mass := p.Particle.Mass()
		mq2 := -mass * mass
		accumGrad := V.Vec{}
		//Conduct inner loop
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := V.Sub(positions[samples[j]], positions[i])
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := fluid.Kern.Gradient(dist, dir)
				F = ((scalar[i] / (dens * dens)) + (scalar[samples[j]] / (jDensity * jDensity)))
				accumGrad = V.Add(accumGrad, V.Scale(grad, -F))
			}
			//End Inner Loop
		}
		V.Add(vector_field[i], V.Scl(accumGrad, mq2))
	} //End Particle Loop
}

func (p SPHParticleSystem) Curl(vector_field []V.Vec, scalar_field []float64) {

}
