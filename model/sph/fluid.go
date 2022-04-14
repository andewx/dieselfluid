package sph

//Provides particle Interface - Implements basic particles and particle system
//for reuse with application specific SPH methods

import (
	"github.com/andewx/dieselfluid/geom"
	"github.com/andewx/dieselfluid/geom/grid"
	V "github.com/andewx/dieselfluid/math/mgl"
	"github.com/andewx/dieselfluid/model/field"
	"math"
)

//SPHSystem represents the core SPH iteration
type SPHSystem interface {
	GetPos() []V.Vec     //Gette
	GetVel() []V.Vec     //Getter
	GetDens() []float32  //Getter
	GetForce() []V.Vec   //Getter
	GetPress() []V.Vec   //Getter
	TimeStep() float32   //Getter
	UpdateTime() float32 //Update Time Delta
	Length() int         //Num Particles
}

//SPH Standard SPH Particle System - Implements SPHSystem Interface
type SPHCore struct {
	Pos       []V.Vec         //Positions
	Dens      []float32       //Densities
	Vels      []V.Vec         //Velocities
	Fs        []V.Vec         //Forces
	Ps        []float32       //Pressures
	Time      float64         //Time Step
	MaxVel    float32         //Max Vel - Courant Condition
	Field     field.SPHField  //Gradient Differential Methods
	Colliders []geom.Collider //List of collidables
}

//Map 3D Index - maps a 3D I,J,K position to 1D flattened array
func Map3DIndex(i int, j int, k int, i_w int, j_w int) int {
	return (k * i_w * j_w) + (k * i_w) + i
}

//---------- Build SPH Core Structure On Grid -----------//
func BuildSPHCube(Scale V.Vec, TransOrigin V.Vec, f field.SPHField, colliders []geom.Collider) SPHCore {

	//Build The Kernel Grid Structure
	grid, dim := grid.BuildKernGrid(Scale, TransOrigin, f.Kern.H()) //Builds Grid Based On Kernel Size (dimenionality of grid cube depends on kernel)
	core := SPHCore{}
	num := dim * dim * dim

	//Make Vectors
	core.Pos = make([]V.Vec, num)
	core.Dens = make([]float32, num)
	core.Vels = make([]V.Vec, num)
	core.Fs = make([]V.Vec, num)
	core.Ps = make([]float32, num)

	//Set Vars
	core.Field = f
	core.UpdateTime()
	core.Colliders = colliders

	//Update Positions -- Spacing Based on Kernel
	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			for k := 0; k < num; k++ {
				nPos := grid.GridPosition(i, j, k)
				index := Map3DIndex(i, j, k, dim, dim)
				core.Pos[index] = nPos
			}
		}
	}

	return core
}

//-----------Implements SPHCore -----------------------------
func (p SPHCore) GetPos() []V.Vec {
	return p.Pos
}
func (p SPHCore) GetVel() []V.Vec {
	return p.Vels
}
func (p SPHCore) GetDens() []float32 {
	return p.Dens
}
func (p SPHCore) GetForce() []V.Vec {
	return p.Fs
}
func (p SPHCore) GetPress() []float32 {
	return p.Ps
}
func (p SPHCore) TimeStep() float64 {
	return p.Time
}
func (p SPHCore) UpdateTime() float64 {
	K := float64(0.4 * p.Field.Part.KernelRad() / p.Field.Part.Sound())
	A := float64(p.Field.Part.KernelRad() / p.MaxVel)

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
		if math.IsNaN(float64(psi)) {
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
	k_stiff := float32(-0.25) //Restitution Coefficient. Further research req'd
	friction := float32(0.01)
	velN := V.Scl(norm, V.Dot(norm, vel))
	velTan := V.Sub(vel, velN)
	dtVN := V.Scl(velN, (k_stiff - 1.0))
	velN = V.Scl(velN, k_stiff)

	//Compute friction coefficients
	if V.Mag(velTan) > 0.0 {
		fcomp := float32(1.0 - friction*V.Mag(dtVN)/V.Mag(velTan))
		frictionScale := float32(math.Max(float64(fcomp), 0.0))
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
		vels[i] = V.Add(vels[i], V.Scl(a, float32(ts)))     //Apply Acceleration
		pos[i] = V.Add(pos[i], V.Scl(vels[i], float32(ts))) //Apply Velocity
		forces[i][0] = 0.0
		forces[i][1] = 0.0
		forces[i][2] = 0.0
		if V.Mag(vels[i]) > p.MaxVel {
			p.MaxVel = V.Mag(vels[i])
		}
	}
}
