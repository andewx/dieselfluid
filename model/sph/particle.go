package sph

import "math"

//Particle Const Enum Defs
const PARTICLE_SPH = 0
const PARTICLE_BOUND = 1
const PARTICLE_GHOST = 2

//Particle Interface
interface Particle{
  func (p Particle)D0()float64 //Target Density
  func (p Particle)MapEOS(d Density)Pressure //Maps Pressure
  func (p Particle)KernelVolume()float64 //Kernel Volume
  func (p Particle)KernelRad()float64 //Kernel Radius
  func (p Particle)Rad() float64 //Particle Radius
  func (p Particle)Vol() float64 //Particle Volume
  func (p Particle)Stiffness() float64 //Stifness Constant
  func (p Particle)Gamma() float64 //Gamma Constant
  func (p Particle)Sound() float64 //Speed of Sound
  func (p Particle)Type() int //Particle Type
}


//Defines General Particle Parameters - Implements Particle
type SPHParticle struct{
  D0 Density
  Rad float64
  KrnlRad float64
  EosK  float64
  EosG  float64
  Sound float64
  Type int
}

//-----------Implements SPH PARTICLE ------------------------//

func (p SPHParticle)D0()float64{
  return p.D0
}

//MapEOS maps a particle density to pressure
func (p SPHParticle)MapEOS(d Density)Pressure{
  if d < p.D0{
    d = p.D0
  }
  R := d/p.D0
  return p.Stiffness() * (math.Pow(R,p.EosG) - 1.0)
}

//Kernel Volume Returns kernel Volume
func (p SPHParticle)KernelVolume()float64{
  d := 2 * p.KrnlRad
  return d * d * d
}

//Kernel Rad returns Particle Kernel Radius
func (p SPHParticle)KernelRad()float64{
  return p.KrnlRad
}

//Rad returns the particle radius
func (p SPHParticle)Rad()float64{
  return p.Rad
}

//Vol returns particle volume
func (p SPHParticle)Vol() float64{
  d := 2 * p.Rad
  return d * d * d
}

//Stiffness return Eos equation K Parameter
func (p SPHParticle)Stiffness() float64{
  return (p.D0 * p.EosK)/p.EosG
}
//Gamma returns particle Eos equation gamma
func (p SPHParticle)Gamma() float64{
  return p.EosG
}

//Sound returns speed of sound
func (p SPHParticle)Sound() float64{
  return p.Sound
}

func (p SPHParticle)Type() int{
  return p.Type
}

//SPHParticle returns new SPHParticle with water parameters given particle radius
//under the Weak EOS Equation K = 1 G = 1
func Build_SPHParticle(rad float64)SPHParticle{
  return SPHParticle{1000, rad, 1.5 * rad, 1, 1, 4800,0}
}

func Build_SPHParticle0(rad float64, stiff float64, gamma float64, sound float64)SPHParticle{
  return SPHParticle(1000, rad, 1.5 * rad, stiff, gamma, sound,0)
}
