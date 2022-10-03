package field

import "github.com/andewx/dieselfluid/math/vector"
import "github.com/andewx/dieselfluid/model"

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
	Ref []model.Particle
}

func (p DensityField) Value(i int) float32 {
	return p.Ref[i].Density()
}

func (p DensityField) Set(x float32, i int) {
	p.Ref[i].SetDensity(x)
}

//PRESSURE

type PressureField struct {
	Ref []model.Particle
}

func (p PressureField) Value(i int) float32 {
	return p.Ref[i].Pressure()
}

func (p PressureField) Set(x float32, i int) {
	p.Ref[i].SetPressure(x)
}

//FORCE

type ForceField struct {
	Ref []model.Particle
}

func (p ForceField) Value(i int) []float32 {
	return p.Ref[i].Force()
}

func (p ForceField) Set(x []float32, i int) {
	p.Ref[i].SetForce(x)
}

type VelocityField struct {
	Ref []model.Particle
}

func (p VelocityField) Value(i int) []float32 {
	return p.Ref[i].Velocity()
}

func (p VelocityField) Set(x []float32, i int) {
	p.Ref[i].SetVelocity(x)
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
	return vector.Cast(p.Values[i])
}

func (p Vector3Field) Set(x []float32, i int) {
	p.Values[i] = vector.CastFixed(x)
}
