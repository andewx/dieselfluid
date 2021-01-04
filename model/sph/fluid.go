package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	V "dslfluid.com/dsl/math/math64" //Diesel Vector Library - Simple Vec
	F "dslfluid.com/dsl/model/field"
)

//SPHSystem represents the core SPH iteration
type SPHSystem interface {
	GetPos() []V.Vec               //Getter
	GetVel() []V.Vec               //Getter
	GetDens() []float64            //Getter
	GetForce() []V.Vec             //Getter
	GetPress() []V.Vec             //Getter
	TimeStep() float64             //Getter
	UpdateTime() float64           //Update Time Delta
	Length() int                   //Num Particles
	Run() error                    //Runs SPH Computation Loop
	Run_Threaded(t chan int) error //Runs SPH computation loop as thread
}

//SPH Standard SPH Particle System - Implements SPHSystem Interface
type SPHCore struct {
	Pos    []V.Vec    //Positions
	Dens   []float64  //Densities
	Vels   []V.Vec    //Velocities
	Fs     []V.Vec    //Forces
	Ps     []float64  //Pressures
	Time   float64    //Time Step
	MaxVel float64    //Max Vel - Courant Condition
	Field  F.SPHField //Gradient Differential Methods
}

//-----------Implements SPHCore -----------------------------
func (p SPHCore) GetPos() []V.Vec {
	return p.Pos
}
func (p SPHCore) GetVel() []V.Vec {
	return p.Vels
}
func (p SPHCore) GetDens() []float64 {
	return p.Dens
}
func (p SPHCore) GetForce() []V.Vec {
	return p.Fs
}
func (p SPHCore) GetPress() []float64 {
	return p.Ps
}
func (p SPHCore) TimeStep() float64 {
	return p.Time
}
func (p SPHCore) UpdateTime() float64 {
	K := 0.4 * p.Field.Part.KernelRad() / p.Field.Part.Sound()
	A := p.Field.Part.KernelRad() / p.MaxVel

	if A < K {
		p.Time = A
		return A
	}
	p.Time = K
	return K
}

//-----------SPH Core Private (From Interface) Methods --------------------//

//ComputeDensity calculates particle density -- Maps Particle Density to Pressure
func (p SPHCore) ComputeDensity() {

}

//AccumulatePressure -- Accumulates Particle Pressure Forces
func (p SPHCore) AccumulatePressure() int {
	return 0
}

//Accumulate_NonPressure -- accumulates non pressure forces (Viscosity)
func (p SPHCore) Accumulate_NonPressure() int {
	return 0
}

//External adds external forces -- (Gravity - Wind Resist Etc)
func (p SPHCore) External(force V.Vec) {

}

//Collide handles particle collision
func (p SPHCore) Collide() {

}

//Update updates all particle positions -- non-blocking mutex locked
func (p SPHCore) Update() {

}

//---------SPHCore Run Methods------------------------//
func (p SPHCore) Run() error {
	return nil
}

func (p SPHCore) Run_Threaded(t chan int) error {
	return nil
}
