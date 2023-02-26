package pcisph

import (
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model/sph"
	"github.com/andewx/dieselfluid/render"
)

type PciMethod struct {
	system       *sph.SPH
	particle_vao uint32
	particle_vbo uint32
}

func NewPCIMethod(sys *sph.SPH, vao uint32, vbo uint32) *PciMethod {
	return &PciMethod{sys, vao, vbo}
}

func (pci *PciMethod) SetGL(vao uint32, vbo uint32) {
	pci.particle_vao = vao
	pci.particle_vbo = vbo
}

func (pci *PciMethod) Run(message chan string, hasGL bool, mRender *render.RenderSystem) {
	done := false
	refDensity := pci.system.Field().D0()
	field := pci.system.Field().Particles
	_vel := make([]float32, field.N()*3)
	_pos := make([]float32, field.N()*3)
	positions := field.Positions()
	velocities := field.Velocities()
	pressures := field.Pressures()
	t := pci.system.Time()
	m := field.Mass()
	d0 := field.D0()
	beta := (t * t * m * m) * (2 / (d0 * d0))
	beta = 1 / beta
	for index := 0; index < field.N()*3; index++ {
		_pos[index] = positions[index]
		_vel[index] = velocities[index]
	}

	for !done {
		pci.system.DensityAll()
		pci.system.ViscousAll()
		max_error_ratio := float32(0.0)
		density_error := float32(0.0)
		num := field.N()
		_maxIterations := int(5)
		_max_density_error_ratio := float32(0.01)

		for iter := 0; iter < _maxIterations; iter++ {

			max_error_ratio = 0

			//Predict Velocity / Position
			for index := 0; index < num; index++ {
				x := index * 3
				if x < field.N()*3 {
					t_pos := []float32{_pos[x], _pos[x+1], _pos[x+2]}
					t_vel := []float32{_vel[x], _vel[x+1], _vel[x+2]}
					ext_force := field.Force(index)
					accel := vector.Scale(ext_force, 1/field.Mass())
					t_vel = vector.Add(t_vel, vector.Scale(accel, pci.system.Time()))
					t_pos = vector.Add(t_pos, vector.Scale(t_vel, pci.system.Time()))
					_pos[x] = t_pos[0]
					_pos[x+1] = t_pos[1]
					_pos[x+2] = t_pos[2]
					_vel[x] = t_vel[0]
					_vel[x+1] = t_vel[1]
					_vel[x+2] = t_vel[2]
				}
			}

			//Compute Pressure From density Error
			for index := 0; index < field.N(); index++ {
				x := index * 3
				if x < field.N()*3 {
					nPos := []float32{_pos[x], _pos[x+1], _pos[x+2]}

					//Predict density and error updating pi()
					calc_density := pci.system.Field().DensityF(nPos, _pos)
					density_error = (calc_density - refDensity)
					abs_density_error := density_error / refDensity
					delta := pci.system.Delta()
					pressures[index] += density_error * delta

					if abs_density_error > max_error_ratio {
						max_error_ratio = abs_density_error
					}
				}
			}
			pci.system.GradientPressureForce()

			if max_error_ratio <= _max_density_error_ratio {
				//fmt.Printf("stable condition\n")
				break
			}
		}

		pci.system.Update()

		select {
		case msg := <-message:
			if msg == "QUIT" {
				done = true
			}
		default:
			//proceed
		}

		select {
		case message <- "SAMPLER_UPDATE":
		default:
		}
	}
	return
}
