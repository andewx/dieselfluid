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

//Particle Interface
type Particle interface {
	D0() float64              //Target float64
	MapEOS(d float64) float64 //Maps float64
	KernelVolume() float64    //Kernel Volume
	KernelRad() float64       //Kernel Radius
	Rad() float64             //Particle Radius
	Vol() float64             //Particle Volume
	Stiffness() float64       //Stifness Constant
	Gamma() float64           //Gamma Constant
	Sound() float64           //Speed of Sound
	Type() int                //Particle Type
	Mass() float64
}

//Defines General Particle Parameters - Implements Particle
type SPHParticle struct {
	Ms      float64
	Dens    float64
	Radius  float64
	KrnlRad float64
	EosK    float64
	EosG    float64
	Sos     float64
	Enum    int
}

//-----------Implements SPH PARTICLE ------------------------//

func (p SPHParticle) D0() float64 {
	return p.Dens
}

func (p SPHParticle) Mass() float64 {
	return p.Ms
}

//MapEOS maps a particle float64 to float64
func (p SPHParticle) MapEOS(d float64) float64 {
	if d < p.Dens {
		d = p.Dens
	}
	R := d / p.Dens
	return float64(p.Stiffness() * (math.Pow(float64(R), p.EosG) - 1.0))
}

//Kernel Volume Returns kernel Volume
func (p SPHParticle) KernelVolume() float64 {
	d := 2 * p.KrnlRad
	return d * d * d
}

//Kernel Rad returns Particle Kernel Radius
func (p SPHParticle) KernelRad() float64 {
	return p.KrnlRad
}

//Rad returns the particle radius
func (p SPHParticle) Rad() float64 {
	return p.Radius
}

//Vol returns particle volume
func (p SPHParticle) Vol() float64 {
	d := 2 * p.Radius
	return d * d * d
}

//Stiffness return Eos equation K Parameter
func (p SPHParticle) Stiffness() float64 {
	return (float64(p.Dens) * p.EosK) / p.EosG
}

//Gamma returns particle Eos equation gamma
func (p SPHParticle) Gamma() float64 {
	return p.EosG
}

//Sound returns speed of sound
func (p SPHParticle) Sound() float64 {
	return p.Sos
}

func (p SPHParticle) Type() int {
	return p.Enum
}

//SPHParticle returns new SPHParticle with water parameters given particle radius
//under the Weak EOS Equation K = 1 G = 1
func Build_SPHParticle(rad float64) SPHParticle {
	return SPHParticle{1.0, 1000, rad, 1.5 * rad, 1, 1, 4800, 0}
}

func Build_SPHParticle0(rad float64, stiff float64, gamma float64, sound float64) SPHParticle {
	return SPHParticle{1.0, 1000, rad, 1.5 * rad, stiff, gamma, sound, 0}
}
