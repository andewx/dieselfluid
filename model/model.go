package model

import "math"

//Particle Const Enum Defs
const PARTICLE_SPH = 0
const PARTICLE_BOUND = 1
const PARTICLE_GHOST = 2

//Go Routine Enums:
const THREAD_WAIT = 100
const THREAD_GO = 101
const THREAD_ERR = 102
const THREAD_DONE = 103
const SPH_THREAD_WAITING = 104

//Collision Handling
const SPH_MESH_COLLIS = 1
const SPH_PARTICLE_COLLIS = 2

//SPH IMPLEMENTATION
const USE_STD = 0
const USE_WCSPH = 1
const USE_PCISPH = 2
const USE_DFSPH = 3
const USE_GRID = 4
const USE_FLIP = 5

//SAMPLER Enums
const VOXEL_SAMPLER = 0
const VOXEL_CACHE_SAMPLER = 1
const SAMPLER_HEURISTIC_NEIGHBORS = 2
const SAMPLER_ALL_NEIGHBORS = 3

//SPH Method Return Values
const SPH_VALID = 0
const SPH_EOS_INF = -1

//Water Model
const FLUID_MASS = 0.1
const FLUID_DENSITY = 1.0
const FLUID_STIFF = 6.1
const FLUID_KERN_RAD = 0.2
const FLUID_SOS = float32(1480.0)
const FLUID_WATER = 0

//Particle Field Interface for interacting with both buffer based and struct based particle fields
type ParticleField interface {
	Positions() []float32
	Velocities() []float32
	Densities() []float32
	Forces() []float32
	Pressure(x int, d0 float32, p0 float32) float32
	Pressures() []float32
	Density(x int) float32
	Position(x int) []float32
	Velocity(x int) []float32
	Force(x int) []float32
	Mass() float32
	Set(x int, particle Particle)
	Get(x int) Particle
	D0() float32
	AddBoundaryParticles([]float32)
	N() int
	Total() int //Total particles
}

//EOSGamma() - Full Tait Equation of State for Water like incrompressible fluids where
//gamma maps the stiffess parameter of the fluid with suggested values in the 6.0-7.0 range
//@param x - density input
//@param c0 - Speed of sound & reference pressure relation (2.15)gpa
//@param d0 - Reference density (1000)kg/m^3 or (1.0g/cm^3)
//@param gamma - Stifness parameter (7.15)
//@param p0 - reference pressure for the system (101,325)
//@notes The reference speed of sound c0 is sometimes taken to be approximately 10 times the maximum
//expected velocity for the system

func EosGamma(x float32, c0 float32, d0 float32, gamma float32, p0 float32) float32 {
	return (c0/gamma)*float32(math.Pow(float64(x/d0), float64(gamma))-1) + p0
}

//EOS() - Tait Equation of State for Water like incrompressible fluids where
//gamma maps the stiffess parameter of the fluid with suggested values in the 6.0-7.0 range
//@param x - density input
//@param c0 - Speed of sound & reference pressure relation (2.15)gpa
//@param d0 - Reference density (1000)kg/m^3 or (1.0g/cm^3)
//@param gamma - Stifness parameter
//@param p0 - reference pressure for the system (101325 pa) or 1013.25hPA
//@notes The reference speed of sound c0 is taken to be approximately 10 times the maximum
//expected velocity for the system if
func TaitEos(x float32, d0 float32, p0 float32) float32 {
	g := float32(7.16)
	w := float32(2.15)
	y := (w/g)*float32(math.Pow(float64(x/d0), float64(g))-1) + p0
	if y <= p0 {
		return p0
	}
	return y
}

//Set the a buffer size >= 3 with the size 3 b float32 array
func Float3_buffer_set(x int, buffer []float32, b [3]float32) {

	if x+2 > len(buffer) {
		return
	}

	if len(buffer) < 3 || len(b) < 3 {
		return
	}
	buffer[x] = b[0]
	buffer[x+1] = b[1]
	buffer[x+2] = b[2]
}

func Float3_set(x int, a [3]float32, buffer []float32) {

	if x+2 > len(buffer) {
		return
	}

	if len(buffer) < 3 || len(a) < 3 {
		return
	}

	a[0] = buffer[x]
	a[1] = buffer[x+1]
	a[2] = buffer[x+2]
}
