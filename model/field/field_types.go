package field

import (
	"github.com/andewx/dieselfluid/model"
)

type Field interface {
	Value(i int) float32
	Set(x float32, i int)
}

type TensorField interface {
	Value(i int) []float32
	Set(value []float32, i int)
}

//DENSITY

type DensityField struct {
	Ref model.ParticleField
}

func (p DensityField) Value(i int) float32 {
	return p.Ref.Get(i).Density
}

func (p DensityField) Set(x float32, i int) {
	particle := p.Ref.Get(i)
	particle.Density = x
	p.Ref.Set(i, particle)
}

//PRESSURE

type PressureField struct {
	Ref model.ParticleField
}

func (p PressureField) Value(i int) float32 {
	return p.Ref.Get(i).Pressure(model.FLUID_DENSITY, 0.0)
}

func (p PressureField) Set(x float32, i int) {
	particle := p.Ref.Get(i)
	particle.Hpa = x
	p.Ref.Set(i, particle)
}

//FORCE

type ForceField struct {
	Ref model.ParticleField
}

func (p ForceField) Value(i int) []float32 {
	particle := p.Ref.Get(i)
	return particle.Force[:]
}

func (p ForceField) Set(x []float32, i int) {
	particle := p.Ref.Get(i)
	particle.Force[0] = x[0]
	particle.Force[1] = x[1]
	particle.Force[2] = x[2]
	p.Ref.Set(i, particle)
}

type VelocityField struct {
	Ref model.ParticleField
}

func (p VelocityField) Value(i int) []float32 {
	particle := p.Ref.Get(i)
	return particle.Velocity[:]
}

func (p VelocityField) Set(x []float32, i int) {
	particle := p.Ref.Get(i)
	particle.Velocity[0] = x[0]
	particle.Velocity[1] = x[1]
	particle.Velocity[2] = x[2]
	p.Ref.Set(i, particle)
}

//Scalar Fields

type ScalarField struct {
	Values []float32
}

func (p ScalarField) Value(i int) float32 {
	return p.Values[i]
}

func (p ScalarField) Set(x float32, i int) {
	p.Values[i] = x
}

//Vector 3 Field implements Tensor Field Interface
type Vector3Field struct {
	Values [][3]float32
}

func (p Vector3Field) Value(i int) []float32 {
	return p.Values[i][:]
}

func (p Vector3Field) Set(x []float32, i int) {
	p.Values[i][0] = x[0]
	p.Values[i][1] = x[1]
	p.Values[i][2] = x[2]
}
