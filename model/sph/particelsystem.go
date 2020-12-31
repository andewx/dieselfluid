package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
V  "github.com/andewx/dieselfluid/math/math64" //Diesel Vector Library - Simple Vec
)

//Define Particle Attribute Types
type Position V.Vec
type Density float64
type Velocity V.Vec
type Force  V.Vec
type Pressure float64


//Particle System Interface
type ParticleSystem interface{
  func (p ParticleSystem)Positions() []Position
  func (p ParticleSystem)Velocities() []Velocity
  func (p ParticleSystem)Densities() []Density
  func (p ParticleSystem)Forces() []Force
  func (p ParticleSystem)Pressures() []Pressure
  func (p ParticleSystem)TimeStep() float64
  func (p ParticleSystem)UpdateTime()float64
  func (p ParticleSystem)Density(positions []V.Vec, density_field []float64)
  func (p ParticleSystem)Gradient(scalar []float64,vector_gradient []V.Vec)
  func (p ParticleSystem)Div(vector_field []V.Vec, div_field []float64)
  func (p ParticleSystem)Laplacian(vector_field []V.Vec, force_field []V.Vec)
  func (p ParticleSystem)Curl(vector_field []V.Vec, scalar_field []float64)
  func (p ParticleSystem)Length()int
}


//Defines SPH Particle System - Implements ParticleSystem Interface Which Contains some Generalized
//functionality for field computation. Which allows reuse by implementations of SPH. For example
//A DFSPH / WCSPH / IISPH / PCISPH Can utilize the SPH Particle System Layout and compute the Gradients
//and field derivatives consistently. Also Particle In Cell Methods Benefit
type SPHParticleSystem struct{
  Positions []Position
  Densities []Density
  Velocities []Velocity
  Forces  []Force
  Pressures []Pressure
  Particle SPHParticle
  TimeStep float64
  MaxVel   float64
}


//-----------Implements SPHParticleSystem -----------------------------
func (p SPHParticleSystem)Positions() []Position{
  return p.Positions
}
func (p SPHParticleSystem)Velocities() []Velocity{
  return p.Velocities
}
func (p SPHParticleSystem)Densities() []Density{
  return p.Densities
}
func (p SPHParticleSystem)Forces() []Force{
  return p.Forces
}
func (p SPHParticleSystem)Pressures() []Pressure{
  return p.Pressures
}
func (p SPHParticleSystem)TimeStep() float64{
  return p.TimeStep
}
func (p SPHParticleSystem)UpdateTime()float64{
  K := 0.4 * p.Particle.KernelRad()/p.Particle.Sound()
  A := p.Particle.KernelRad()/p.MaxVel

  if A < K{
    p.TimeStep = A
    return A
  }
  p.TimeStep = K
  return K
}

//Need the Kernel Function Before Implementation
func (p SPHParticleSystem)Gradient(scalar []float64,vector_gradient []V.Vec){

}
func (p SPHParticleSystem)Div(vector_field []V.Vec, div_field []float64){

}
func (p SPHParticleSystem)Laplacian(vector_field []V.Vec, force_field []V.Vec){

}
func (p SPHarticleSystem)Curl(vector_field []V.Vec, scalar_field []float64){

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
