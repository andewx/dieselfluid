package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	V "github.com/andewx/dieselfluid/math/math64" //Diesel Vector Library - Simple Vec
	F "github.com/andewx/dieselfluid/model/field"
)

//Particle System Interface
type ParticleSystem interface {
	GetPos() []V.Vec
	GetVel() []V.Vec
	GetDens() []float64
	GetForce() []V.Vec
	GetPress() []V.Vec
	TimeStep() float64
	UpdateTime() float64
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
	Field    F.SPHField
}

//-----------Implements SPHParticleSystem -----------------------------
func (p SPHParticleSystem) GetPos() []V.Vec {
	return p.Pos
}
func (p SPHParticleSystem) GetVel() []V.Vec {
	return p.Vels
}
func (p SPHParticleSystem) GetDens() []float64 {
	return p.Dens
}
func (p SPHParticleSystem) GetForce() []V.Vec {
	return p.Fs
}
func (p SPHParticleSystem) GetPress() []float64 {
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
