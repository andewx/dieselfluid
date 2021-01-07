package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	V "dslfluid.com/dsl/math/math64" //Diesel Vector Library - Simple Vec
	"dslfluid.com/dsl/model"
	F "dslfluid.com/dsl/model/field"
	"math"
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

	p.MaxVel = 0.0

	return K
}

//-----------SPH Core Private (From Interface) Methods --------------------//

//ComputeDensity calculates particle density -- Maps Particle Density to Pressure
func (p SPHCore) ComputeDensity() {
	p.Field.Density(p.GetPos(), p.GetDens())
}

//AccumulatePressure -- Accumulates Particle Pressure Forces
func (p SPHCore) AccumulatePressure() int {
	retVal := 0 //SPH VALID
	pressures := p.GetPress()
	densities := p.GetDens()
	for i := 0; i < len(pressures); i++ {
		psi := p.Field.Part.MapEOS(densities[i])
		pressures[i] = psi
		if math.IsNaN(psi) {
			retVal = -1
			pressures[i] = 0
		}
	}
	return retVal
}

//Accumulate_NonPressure -- accumulates non pressure forces (Viscosity)
func (p SPHCore) Accumulate_NonPressure() int {
	p.Field.Laplacian(p.GetPos(), p.GetDens(), p.GetVel(), p.GetForce())
	return 0
}

//External adds external forces -- (Gravity - Wind Resist Etc)
func (p SPHCore) External(force V.Vec) {
	fs := p.GetForce()
	for i := 0; i < len(fs); i++ {
		fs[i] = V.Add(fs[i], force)
	}
}

//Collide handles particle collision
func (p SPHCore) Collide() {

}

//Update updates all particle positions -- non-blocking mutex locked
func (p SPHCore) Update() {

	m := 1 / p.Field.Part.Mass()
	ts := p.TimeStep()
	forces := p.GetForce()
	vels := p.GetVel()
	pos := p.GetPos()

	//Calculate Velocities Update Position / Clear Force
	for i := 0; i < len(p.GetPos()); i++ {
		a := V.Scl(forces[i], m)
		vels[i] = V.Add(vels[i], V.Scl(a, ts))     //Apply Acceleration
		pos[i] = V.Add(pos[i], V.Scl(vels[i], ts)) //Apply Velocity
		forces[i][0] = 0.0
		forces[i][1] = 0.0
		forces[i][2] = 0.0
		if V.Mag(vels[i]) > p.MaxVel {
			p.MaxVel = V.Mag(vels[i])
		}
	}
}

//---------SPHCore Run Methods------------------------//
func (p SPHCore) Run() error {
	done := false

	for !done {
		p.ComputeDensity()
		p.AccumulatePressure()
		p.Accumulate_NonPressure()
		p.Collide()
		p.Update()
		p.UpdateTime()
	}

	return nil
}

//Run Threaded Executes SPH Loop in Thread Blocking I/O Manner. If an application
//Needs exclusive resource access to SPHCore data structures they should pass
//model.THREAD_WAIT to block thread execution. When access is no longer required
//the method should pass model.THREAD_GO to the specified channel
//Buffer access to SPHCore go slices should be read only access, other wise for thread safe
//Execution THREAD_WAIT should be called if modifying buffers or relying on temporal coherence
//for volatile data buffers
func (p SPHCore) Run_Threaded(t chan int) error {
	done := false
	sync := true

	for !done {

		//Executes full in frame computation loop.
		if sync {
			p.ComputeDensity()
			p.AccumulatePressure()
			p.Accumulate_NonPressure()
			p.Collide()
			p.Update()
			p.UpdateTime()

			//Channel Monitor - Monitor Blocking I/O Request
			status := <-t
			if status == model.THREAD_WAIT {
				sync = false
				t <- model.SPH_THREAD_WAITING
				waitStatus := <-t
				if waitStatus == model.THREAD_GO {
					sync = true
				}
			}
			//End Thread Blocking
		}

		//Handle Thread Block Outside Initial Message
		waitStatus := <-t
		if waitStatus == model.THREAD_GO {
			sync = true
		}

	}

	return nil
}
