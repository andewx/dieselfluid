package sph

import (
	"github.com/andewx/dieselfluid/geom"
	"github.com/andewx/dieselfluid/geom/grid"
	"github.com/andewx/dieselfluid/kernel"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model"
	"github.com/andewx/dieselfluid/model/field"
	"github.com/andewx/dieselfluid/sampler/lsh"
	"math"
)

const (
	VISCOSITY_WATER = 1.3059
	CACHE_L         = 0.8
)

//SPH Standard SPH Particle System - Implements SPHSystem Interface
type SPH struct {
	time       float32         //Time Step
	maxVel     float32         //Max Vel - Courant Condition
	field      field.SPHField  //SPH Field Methods
	colliders  []geom.Collider //List of collidables
	particles  int             //Number Particles
	cache_life float32         //Cache Extinction Coefficient
	mu         float32         //viscosity coefficient
	delta      float32         //pcisph delta computation
}

/*
InitSPH() Creates SPH particle grid using where n3 is the cubic root of the number of particles desired
so that N = n3*n3*n3 and the kernel smoothing lengthing is taken to be the the cubic average scale vector
which is defaulted to 1.0. So that h = (||s||/(N3)). To ensure that the GPU shader is well formed n3 must
be a multiple of the local gpu group size which is 4. So n3 = 4 * X.
*/
func Init(scl float32, origin vector.Vec, colliders []geom.Collider, n3 int, pci bool) SPH {

	//Build The Kernel Grid Structure using a cubic dimension of the particles
	scale := vector.Vec{scl, scl, scl}
	h := (scale.Mag()) / float32(n3)
	num := n3 * n3 * n3
	dim_vec := vector.Vec{float32(n3), float32(n3), float32(n3)}
	kern := kernel.Build_Kernel(h)
	grid := grid.BuildGrid(scale, origin, dim_vec) //Builds Grid Based On Kernel Size (dimenionality of grid cube depends on kernel)
	core := SPH{}

	//Instantiates and allocates the fielded particle lists which includes the collider implicit particle fields
	particles := make([]model.Particle, num)
	sampler := lsh.Allocate(num, 255, 8, particles)
	core.field = field.InitSPH(particles, sampler, kern, num)
	core.particles = num
	core.cache_life = CACHE_L
	core.CFL()
	core.mu = VISCOSITY_WATER
	core.field.BoundaryParticles(colliders)
	core.field.AlignWithGrid(grid)

	if pci {
		core.pcidelta()
	}

	return core
}

//Get the field particles list, note that boundary particles are appended
func (p SPH) Particles() []model.Particle {
	return p.field.Particles()
}

func (p SPH) Field() field.SPHField {
	return p.Field()
}

//Update Nearest Neighbors
func (p SPH) NN() {
	p.field.NN()
}

//Get the number of live particles , does not include boundary particles in the
//field particle list
func (p SPH) N() int {
	return p.particles
}

//CFL Time Step Condition - Ensure the GPU Forumlas match this constraint
func (p SPH) CFL() float32 {
	if p.maxVel >= 5.0 {
		p.time = 1 / p.maxVel
	} else {
		p.time = 0.2
	}
	return p.time
}

func (p SPH) Viscosity() float32 {
	return p.mu
}

func (p SPH) SetViscosity(x float32) {
	p.mu = x
}

//-----------SPH Core Methods Utilized the SPHField Methods-----------//

//Computes all particle densities
func (p SPH) DensityAll() {
	for i := 0; i < p.particles; i++ {
		p.field.Density(i)
	}
}

//Iterates over density field and calculates the particle pressures using tait EOS mapping
func (p SPH) PressureAll() int {
	retVal := 0                           //SPH VALID
	standard_pressure := float32(1013.25) //1013.25 when workimg with air/particle medium
	for i := 0; i < p.particles; i++ {
		psi := p.field.TaitEOS(p.Particles()[i].Dens, standard_pressure) //hpa
		p.field.Particles()[i].SetPressure(psi)
		if math.IsNaN(float64(psi)) {
			retVal = -1
			p.field.Particles()[i].SetPressure(standard_pressure)
		}
	}
	return retVal
}

//(N)Applys an artificial viscosity force by calculating the laplacian of the velocity field
//And maps the force to the particle force fields
func (p SPH) ViscousAll() {
	for i := 0; i < p.particles; i++ {
		p.field.Particles()[i].AddForce(vector.Scale(p.field.LaplacianForce(i, p.field.GetTensorFields()["velocity"]), p.mu))
	}
}

//External add an external force to all particles
func (p SPH) ExternalAll(force vector.Vec) {
	for i := 0; i < p.particles; i++ {
		p.field.Particles()[i].AddForce(force)
	}
}

//Update updates all particle positions -- non-blocking mutex locked
func (p SPH) Update() {

	m := 1 / p.field.Mass()
	ts := p.CFL()

	//Calculate Velocities Update Position / Clear Force To Gravity only
	for i := 0; i < p.particles; i++ {
		a := vector.Scale(p.field.Particles()[i].Force(), m)
		p.field.Particles()[i].AddVelocity(vector.Scale(a, float32(ts)))
		p.field.Particles()[i].AddPosition(vector.Scale(p.field.Particles()[i].Velocity(), float32(ts))) //Apply Velocity
		p.field.Particles()[i].SetForce(vector.Vec{0, -9.81 * p.field.Mass(), 0})
		if vector.Mag(p.field.Particles()[i].Velocity()) > p.maxVel {
			p.maxVel = vector.Mag(p.field.Particles()[i].Velocity())
		}
	}
}

func (p SPH) Delta() float32 { return p.delta }

func (p SPH) MaxV() float32 { return p.maxVel }

func (p SPH) CacheIncr() float32 {
	p.cache_life *= p.cache_life
	if p.cache_life < 0.1 {
		p.cache_life = CACHE_L
		p.NN()
	}
	return p.cache_life
}

/* PCISPHDelta() - Computes delta scalar for PCISPH Pressure Correction term
which is based on a default initialized grid with "full" neighborhood. Kernel
length is Size_Grid/Dim so 0.5. And grid contains 8 particles in the grid.
*/
func (p SPH) pcidelta() float32 {
	sample_sph := Init(1.0, vector.Vec{}, nil, 4, false)
	//d0 := float32(0.0)
	d1 := float32(0.0)
	d2 := vector.Vec{0, 0, 0}
	for i := 0; i < sample_sph.particles; i++ {
		point := sample_sph.field.Particles()[i].Position()
		dist := vector.Mag(point) * vector.Mag(point)
		h0 := sample_sph.field.GetKernelLength()
		if dist < h0*h0 {
			dist = vector.Mag(point)
			dir := vector.Vec{0, 0, 0}
			if dist > 0.0 {
				vector.Scale(point, dist)
			}
			gradW := sample_sph.field.Kernel().Grad(dist, dir)
			d2 = d2.Add(gradW)
			d1 += vector.Dot(gradW, gradW)
		}
	}
	res := vector.Dot(d2, d2) - d1
	if res <= float32(0.0) {
		res = -res
	}
	if res > float32(0.0) {
		q := sample_sph.field.Mass() * sample_sph.time / sample_sph.field.D0()
		beta := float32(2.0) * q * q
		res = (-1 / (beta * res))
	} else {
		return float32(0.0)
	}
	p.delta = res
	return res
}
