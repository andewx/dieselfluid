package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	"dslfluid.com/dsl/geom"
	V "dslfluid.com/dsl/math/math64"
	"dslfluid.com/dsl/model/field"
	"math"
)

//SPHSystem represents the core SPH iteration
type SPHSystem interface {
	GetPos() []V.Vec     //Gette
	GetVel() []V.Vec     //Getter
	GetDens() []float64  //Getter
	GetForce() []V.Vec   //Getter
	GetPress() []V.Vec   //Getter
	TimeStep() float64   //Getter
	UpdateTime() float64 //Update Time Delta
	Length() int         //Num Particles
}

//SPH Standard SPH Particle System - Implements SPHSystem Interface
type SPHCore struct {
	Pos       []V.Vec         //Positions
	Dens      []float64       //Densities
	Vels      []V.Vec         //Velocities
	Fs        []V.Vec         //Forces
	Ps        []float64       //Pressures
	Time      float64         //Time Step
	MaxVel    float64         //Max Vel - Courant Condition
	Field     field.SPHField  //Gradient Differential Methods
	Colliders []geom.Collider //List of collidables
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
	for i := 0; i < len(p.Colliders); i++ {
		for j := 0; j < len(p.GetPos()); j++ {

			//Get Collisions for each point and calculate force
			norm, _, point, collis := p.Colliders[i].Collision(p.GetPos()[j], p.GetVel()[j], p.TimeStep(), p.Field.Part.Rad())
			//Sets new position and velocity
			if collis {
				v, _ := p.CalcCollision(j, norm)
				p.GetPos()[j] = point
				p.GetVel()[j] = v
			}
		}
	}
}

//Collision Calculations returns Velocity (Vec32), Force (Vec32) Momentum Vector
func (p SPHCore) CalcCollision(index int, norm V.Vec) (V.Vec, V.Vec) {
	vel := p.GetVel()[index]
	k_stiff := -0.25 //Restitution Coefficient. Further research req'd
	friction := 0.01
	velN := V.Scl(norm, V.Dot(norm, vel))
	velTan := V.Sub(vel, velN)
	dtVN := V.Scl(velN, (k_stiff - 1.0))
	velN = V.Scl(velN, k_stiff)

	//Compute friction coefficients
	if V.Mag(velTan) > 0.0 {
		fcomp := float64(1.0 - friction*V.Mag(dtVN)/V.Mag(velTan))
		frictionScale := math.Max(fcomp, 0.0)
		velTan = V.Scl(velTan, frictionScale)
	}

	nVelocity := V.Add(velN, velTan)

	//Scale Force
	forceNormal := V.Scl(velN, -p.Field.Part.Mass())
	p.GetForce()[index] = forceNormal

	return nVelocity, forceNormal
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
