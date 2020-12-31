package sph

import "math"

//Particle Const Enum Defs
const PARTICLE_SPH = 0
const PARTICLE_BOUND = 1
const PARTICLE_GHOST = 2

//Particle Interface
type Particle interface {
	D0() float64               //Target Density
	MapEOS(d Density) Pressure //Maps Pressure
	KernelVolume() float64     //Kernel Volume
	KernelRad() float64        //Kernel Radius
	Rad() float64              //Particle Radius
	Vol() float64              //Particle Volume
	Stiffness() float64        //Stifness Constant
	Gamma() float64            //Gamma Constant
	Sound() float64            //Speed of Sound
	Type() int                 //Particle Type
	Mass() float64
}

//Defines General Particle Parameters - Implements Particle
type SPHParticle struct {
	Mass    float64
	Dens    Density
	Radius  float64
	KrnlRad float64
	EosK    float64
	EosG    float64
	Sos     float64
	Enum    int
}

//-----------Implements SPH PARTICLE ------------------------//

func (p SPHParticle) D0() Density {
	return p.Dens
}

func (p SPHParticle) Mass() float64 {
	return p.Mass
}

//MapEOS maps a particle density to pressure
func (p SPHParticle) MapEOS(d Density) Pressure {
	if d < p.Dens {
		d = p.Dens
	}
	R := d / p.Dens
	return Pressure(p.Stiffness() * (math.Pow(float64(R), p.EosG) - 1.0))
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
