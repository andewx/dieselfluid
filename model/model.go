package model

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

//Particle Interface
type ParticleInterface interface {
	Position() []float32
	Velocity() []float32
	Density() float32
	Pressure() float32
	Force() []float32
	SetPosition([]float32)
	SetVelocity([]float32)
	SetForce([]float32)
	AddForce([]float32)
	SetDensity(float32)
	SetPressure(float32)
	GetParticle() Particle
}

//FluidParticle returns new FluidParticle with water parameters given particle radius
func FluidParticle() Particle {
	return Particle{}
}
