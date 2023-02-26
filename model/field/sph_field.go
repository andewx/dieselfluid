package field

import (
	"github.com/andewx/dieselfluid/geom/grid"
	"github.com/andewx/dieselfluid/geom/mesh"
	"github.com/andewx/dieselfluid/kernel"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model"
	"github.com/andewx/dieselfluid/sampler/lsh"
)

const SAMPLES = 150

//Manage SPH Field and Particle Field interactions
type SPHField struct {
	kern          kernel.Kernel
	smplr         *lsh.HashSampler
	Particles     model.ParticleArray
	fields        map[string]Field
	tensor_fields map[string]TensorField
	densities     DensityField
	pressures     PressureField
	velocities    VelocityField
	forces        ForceField
	divergence    Field
	vort          Field
}

func InitSPH(parts model.ParticleArray, ref *lsh.HashSampler, kern kernel.Kernel, basis int) SPHField {
	mySPH := SPHField{}
	mySPH.kern = kern
	mySPH.smplr = ref
	mySPH.Particles = parts
	mySPH.densities = DensityField{mySPH.Particles}
	mySPH.pressures = PressureField{mySPH.Particles}
	mySPH.velocities = VelocityField{mySPH.Particles}
	mySPH.forces = ForceField{mySPH.Particles}
	mySPH.divergence = ScalarField{make([]float32, basis)}
	mySPH.vort = ScalarField{make([]float32, basis)}
	mySPH.fields = make(map[string]Field, 10)
	mySPH.tensor_fields = make(map[string]TensorField, 10)
	mySPH.fields["density"] = mySPH.densities
	mySPH.fields["pressure"] = mySPH.pressures
	mySPH.fields["divergence"] = mySPH.divergence
	mySPH.fields["vorticity"] = mySPH.vort
	mySPH.tensor_fields["velocity"] = mySPH.velocities
	mySPH.tensor_fields["force"] = mySPH.forces
	return mySPH
}

func (p *SPHField) Field() *SPHField {
	return p
}

func (p *SPHField) Mass() float32 {
	return p.Particles.Mass()
}

func (p *SPHField) D0() float32 {
	return p.Particles.D0()
}

func (p *SPHField) Kernel() kernel.Kernel {
	return p.kern
}

func (p *SPHField) GetKernelLength() float32 {
	return p.kern.H0()
}

func (p *SPHField) GetFields() map[string]Field {
	return p.fields
}

func (p *SPHField) BoundaryParticles(colliders []*mesh.Mesh) []float32 {
	//Make the boundary Particles
	var colliderPositions []float32
	if colliders != nil {
		for i := 0; i < len(colliders); i++ {
			colliderPositions = colliders[i].GenerateBoundaryParticles(2.0)
			p.Particles.AddBoundaryParticles(colliderPositions)
		}
	}
	return colliderPositions
}

func (p *SPHField) AlignWithGrid(mGrid grid.Grid) {
	x := int(mGrid.Div[0])
	y := int(mGrid.Div[1])
	z := int(mGrid.Div[2])
	pos := p.Particles.Positions()
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < z; k++ {
				nPos := mGrid.GridPosition(i, j, k)
				id := mGrid.Index(i, j, k)
				particle := p.Particles.Get(id)
				particle.Position = vector.CastFixed(nPos)
				idx := id * 3
				if idx+2 < len(pos) {
					pos[idx] = particle.Position[0]
					pos[idx+1] = particle.Position[1]
					pos[idx+2] = particle.Position[2]
				}
			}
		}
	}
}

func (p *SPHField) GetTensorFields() map[string]TensorField {
	return p.tensor_fields
}

//Sampler Nearest Neighbors Update
func (p *SPHField) NN() {
	p.smplr.UpdateSampler()
}

func (p *SPHField) GetSampler() *lsh.HashSampler {
	return p.smplr
}

//nterpolates a scalar field given a position giving a continuous field
func (p *SPHField) Interpolate(position []float32, field Field) float32 {
	sampleList := p.smplr.GetSamplesFromPosition(position)
	sum := float32(0.0)
	mass := p.Mass()
	for i := 0; i < len(sampleList); i++ {
		part := p.Particles.Get(sampleList[i])
		dist := vector.Dist(position, part.Position[:])
		weight := mass / part.Density * p.kern.F(dist)
		sum += weight * field.Value(sampleList[i])
	}
	return sum
}

func (p *SPHField) DensityF(pos vector.Vec, positions []float32) float32 {
	sampleList := p.smplr.GetSamplesFromPosition(pos)
	density := p.kern.W0()
	mass := p.Field().Mass()

	for j := 0; j < len(sampleList); j++ {
		pIndex := sampleList[j]
		if pIndex < p.Particles.Total() {

			particle_j := p.Particles.Get(pIndex)
			dist := vector.Dist(pos, particle_j.Position[:]) //Change to dist
			density += mass * p.kern.F(dist)
		}
	}
	return density
}

//Density -- Computes density field for SPH Field - Boundary particle contribute to infinite density
func (p *SPHField) Density(i int) {
	sampleList := p.smplr.GetSamples(i)
	density := float32(0)
	particle := p.Particles.Get(i)
	lenSample := len(sampleList)
	mass := p.Field().Mass()
	for j := 0; j < lenSample; j++ {
		pIndex := sampleList[j]
		if i != pIndex && pIndex < p.Particles.Total() {

			particle_j := p.Particles.Get(pIndex)
			dist := vector.Dist(particle.Position[:], particle_j.Position[:]) //Change to dist
			density += mass * p.kern.F(dist)
		}
	}
	particle.Density = density
	p.Particles.Set(i, particle)
}

//Computes gradient vector at particle i given a scalar field
func (p *SPHField) Gradient(i int, field Field) []float32 {

	samples := p.smplr.GetSamples(i)
	F := float32(0.0)
	mass := p.Particles.Mass()
	accumGrad := vector.Vec{0, 0, 0}
	particle := p.Particles.Get(i)
	dens := particle.Density

	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			particle_j := p.Particles.Get(jIndex)
			jDensity := particle_j.Density
			dir := vector.Sub(particle_j.Position[:], particle.Position[:])
			dist := vector.Mag(dir)
			dir = vector.Norm(dir)
			grad := p.kern.Grad(float32(dist), dir)
			F = (field.Value(i) / (dens * dens)) + field.Value(samples[j])/(jDensity*jDensity)
			accumGrad = vector.Add(accumGrad, vector.Scale(grad, F))
		}
	}

	return vector.Scale(accumGrad, dens*mass)

}

//Computes the Divergence of a tensor field
func (p *SPHField) Div(i int, field TensorField) float32 {

	particle := p.Particles.Get(i)
	samples := p.smplr.GetSamples(i)
	div := float32(0.0)
	mass := p.Mass()

	//For all particle neighbors -- Non Symmetric
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			particle_j := p.Particles.Get(jIndex)
			jDensity := particle_j.Density
			dir := vector.Sub(particle_j.Position[:], particle.Position[:])
			dist := vector.Mag(dir)
			dir = vector.Norm(dir) //Normalize
			grad := p.kern.Grad(dist, dir)
			scaleVec := vector.Scale(field.Value(samples[j]), mass/jDensity)
			div += vector.Dot(scaleVec, grad)
		}
	} //End J

	return div

}

//Computes a laplacian value at the particle i for the given scalar field
func (p *SPHField) Laplacian(i int, field Field) float32 {

	particle := p.Particles.Get(i)
	samples := p.smplr.GetSamples(i)
	m := p.Mass()
	sum := float32(0.0)
	//Conduct inner loop
	for j := 0; j < len(samples); j++ {

		jIndex := samples[j]
		particle_j := p.Particles.Get(jIndex)
		if jIndex != i {
			jDensity := particle_j.Density
			dist := vector.Dist(particle.Position[:], particle_j.Position[:])
			sum += m * ((field.Value(samples[j]) - field.Value(i)) / jDensity) * p.kern.O2D(dist)
		}
	}
	return sum
}

//Computes a laplacian value at the particle i for the given scalar field
func (p *SPHField) LaplacianForce(i int, field TensorField) []float32 {

	particle := p.Particles.Get(i)
	samples := p.smplr.GetSamples(i)
	m := p.Mass()
	force := vector.Vec{0, 0, 0}
	//Conduct inner loop
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			particle_j := p.Particles.Get(jIndex)
			jDensity := particle_j.Density
			v := vector.Scale(vector.Sub(particle_j.Velocity[:], particle.Velocity[:]), 1/jDensity)
			dist := vector.Dist(particle.Position[:], particle_j.Position[:])
			force = force.Add(v.Scale(p.kern.O2D(dist))).Scale(m)
		}
	}
	return force
}

//Curl computes non-symmetric curl
func (p *SPHField) Curl(i int, field TensorField) []float32 {

	particle := p.Particles.Get(i)
	samples := p.smplr.GetSamples(i)
	curl_vec := vector.Vec{0, 0, 0}
	mass := p.Mass()

	//For all particle neighbors
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			particle_j := p.Particles.Get(jIndex)
			jDensity := particle_j.Density
			dir := vector.Sub(particle_j.Position[:], particle.Position[:])
			dist := vector.Mag(dir)
			dir = vector.Norm(dir) //Normalize
			grad := p.kern.Grad(dist, dir)
			scaleVec := vector.Scale(field.Value(samples[j]), mass/jDensity)
			curl_vec = vector.Add(curl_vec, vector.Cross(scaleVec, grad))
		}
	} //End J
	return curl_vec
}
