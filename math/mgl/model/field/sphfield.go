package field

//32 bit field Field

import (
	"github.com/andewx/dieselfluid/kernel"
	V "github.com/andewx/dieselfluid/math/mgl" //Diesel Vector Library - Simple Vec
	"github.com/andewx/dieselfluid/model"
	"github.com/andewx/dieselfluid/sampler"
)

//Field interface provides lagrangian analytical component to SPH
type Field interface {
	Density(positions []V.Vec, density_field []float32)
	Gradient(positions []V.Vec, scalar []float32, vector_gradient []V.Vec)
	Div(positions []V.Vec, vector_field []V.Vec, div_field []float32)
	Laplacian(positions []V.Vec, densities []float32, vector_field []V.Vec, force_field []V.Vec)
	Curl(positions []V.Vec, vector_field []V.Vec, scalar_field []float32)
}

//SPHField implements Field Interface
type SPHField struct {
	Kern  kernel.Kernel
	Smplr sampler.Sampler
	Part  model.Particle
}

//Density -- Computes density field for SPH Fiel
func (p SPHField) Density(positions []V.Vec, density_field []float32) {
	N := len(positions)
	for i := 0; i < N; i++ {
		sampleList := p.Smplr.GetSamples(i)
		weight := p.Kern.W0()

		for j := 0; j < len(sampleList); j++ {
			pIndex := sampleList[j]
			if i != j {
				dist := V.Dist(positions[i], positions[pIndex]) //Change to dist
				weight += p.Kern.F(dist)
			}
		}
		density := float32(p.Part.Mass() * weight)
		density_field[i] = density
	}
}

//Gradient computes SPH Particle scalar gradient dependent on density field
func (p SPHField) Gradient(positions []V.Vec, scalar []float32, densities []float32, vector_gradient []V.Vec) {

	for i := 0; i < len(positions); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		F := float32(0.0)
		p.Kern.Adjust(dens / p.Part.D0())
		mass := p.Part.Mass()
		mq2 := -mass * mass
		accumGrad := V.Vec{}
		//Conduct inner loop
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := (V.Sub(positions[samples[j]], positions[i]))
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := p.Kern.Grad(float32(dist), dir)
				F = ((scalar[i] / (dens * dens)) + (scalar[samples[j]] / (jDensity * jDensity)))
				accumGrad = V.Add(accumGrad, V.Scl(grad, -F))
			}
			//End Inner Loop
		}
		V.Add(vector_gradient[i], V.Scl(accumGrad, mq2))
	} //End Particle Loop

}

//Div computes vector field divergence
func (p SPHField) Div(positions []V.Vec, vector_field []V.Vec, densities []float32, div_field []float32) {
	for i := 0; i < len(positions); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		div := float32(0.0)
		p.Kern.Adjust(dens / p.Part.D0())
		mass := p.Part.Mass()

		//For all particle neighbors -- Non Symmetric
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := (V.Sub(positions[samples[j]], positions[i]))
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := p.Kern.Grad(dist, dir)
				scaleVec := V.Scl(vector_field[samples[j]], mass/jDensity)
				div = div + V.Dot(scaleVec, grad)
			}
		} //End J
		div_field[i] = div
	} //End Particle Loop
}

//Laplacian computes vector field laplacian which is formally known as the divergence of the gradient of F
func (p SPHField) Laplacian(positions []V.Vec, densities []float32, vector_field []V.Vec, force_field []V.Vec) {
	for i := 0; i < len(positions); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		lap := V.Vec{}
		p.Kern.Adjust(dens / p.Part.D0())
		mass := p.Part.Mass()

		//Conduct inner loop
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dist := V.Dist(positions[i], positions[j])
				vecDif := V.Sub(vector_field[i], vector_field[j])
				scaleVecDif := V.Scl(vecDif, (mass/jDensity)*p.Kern.O2D(dist))
				lap = V.Add(lap, scaleVecDif)
			}
			//End Inner Loop
		}
		force_field[i] = V.Add(force_field[i], lap)
	} //End Particle Loop
}

//Curl computes non-symmetric curl
func (p SPHField) Curl(positions []V.Vec, densities []float32, vector_field []V.Vec, curl_field []V.Vec) {
	for i := 0; i < len(positions); i++ {
		//For Each Particle Calculate Kernel Based Summation
		samples := p.Smplr.GetSamples(i)
		dens := densities[i]
		curl_vec := V.Vec{}
		p.Kern.Adjust(dens / p.Part.D0())
		mass := p.Part.Mass()

		//For all particle neighbors
		for j := 0; j < len(samples); j++ {
			jIndex := samples[j]
			if jIndex != i {
				jDensity := densities[samples[j]]
				dir := (V.Sub(positions[samples[j]], positions[i]))
				dist := V.Mag(dir)
				dir = V.Norm(dir) //Normalize
				grad := p.Kern.Grad(dist, dir)
				scaleVec := V.Scl(vector_field[samples[j]], mass/jDensity)
				curl_vec = V.Add(curl_vec, V.Cross(scaleVec, grad))
			}
		} //End J
		curl_field[i] = curl_vec
	} //End Particle Loop
}
