package field

import (
	"github.com/andewx/dieselfluid/geom"
	"github.com/andewx/dieselfluid/geom/grid"
	"github.com/andewx/dieselfluid/kernel"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model"
	"github.com/andewx/dieselfluid/sampler"
	"math"
)

type SPHField struct {
	kern          kernel.Kernel
	smplr         sampler.Sampler
	particles     []model.Particle
	fields        map[string]Field
	tensor_fields map[string]TensorField
	densities     DensityField
	pressures     PressureField
	velocities    VelocityField
	forces        ForceField
	divergence    Field
	vort          Field
	mass          float32
	d0            float32
}

func InitSPH(parts []model.Particle, ref sampler.Sampler, kern kernel.Kernel, basis int) SPHField {
	mySPH := SPHField{}
	mySPH.kern = kern
	mySPH.smplr = ref
	mySPH.particles = parts
	mySPH.densities = DensityField{mySPH.particles}
	mySPH.pressures = PressureField{mySPH.particles}
	mySPH.velocities = VelocityField{mySPH.particles}
	mySPH.forces = ForceField{mySPH.particles}
	mySPH.mass = float32(1.0)
	mySPH.d0 = float32(1000.0)
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

func (p SPHField) Particles() []model.Particle {
	return p.particles
}

func (p SPHField) Field() *SPHField {
	return &p
}

func (p SPHField) Mass() float32 {
	return p.mass
}

func (p SPHField) D0() float32 {
	return p.d0
}

func (p SPHField) Kernel() kernel.Kernel {
	return p.kern
}

func (p SPHField) GetKernelLength() float32 {
	return p.kern.H0()
}

func (p SPHField) GetFields() map[string]Field {
	return p.fields
}

func (p SPHField) BoundaryParticles(colliders []geom.Collider) {
	//Make the boundary Particles
	colliders = colliders
	if colliders != nil {
		for i := 0; i < len(colliders); i++ {
			colliderPositions := colliders[i].GenerateBoundaryParticles(1 / p.GetKernelLength())
			boundary_particles := make([]model.Particle, len(colliderPositions))
			for i := 0; i < len(boundary_particles); i++ {
				boundary_particles[i].SetPosition(vector.Cast(colliderPositions[i]))
			}
			//Append the implicit colliders
			p.particles = append(p.particles, boundary_particles[:]...)
		}
	}
}

func (p SPHField) AlignWithGrid(mGrid grid.Grid) {
	x := int(mGrid.DimXYZ[0])
	y := int(mGrid.DimXYZ[1])
	z := int(mGrid.DimXYZ[2])
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < z; k++ {
				nPos := mGrid.GridPosition(i, j, k)
				id := mGrid.Index(i, j, k)
				p.particles[id].SetPosition(nPos)
			}
		}
	}
}

func (p SPHField) GetTensorFields() map[string]TensorField {
	return p.tensor_fields
}

//Sampler Nearest Neighbors Update
func (p SPHField) NN() {
	p.smplr.UpdateSampler()
}

func (p SPHField) GetSampler() sampler.Sampler {
	return p.smplr
}

//nterpolates a scalar field given a position giving a continuous field
func (p SPHField) Interpolate(position []float32, field Field) float32 {
	sampleList := p.smplr.GetRegionalSamples(p.smplr.Hash(vector.CastFixed(position)), 1)
	sum := float32(0.0)
	for i := 0; i < len(sampleList); i++ {
		part := p.particles[sampleList[i]]
		dist := vector.Dist(position, part.Position())
		weight := p.mass / part.Density() * p.kern.F(dist)
		sum += weight * field.Value(sampleList[i])
	}
	return sum
}

//Interpolates a scalar field given a position giving a continuous field
func (p SPHField) InterpolateVectors(position []float32, field TensorField) []float32 {
	sampleList := p.smplr.GetRegionalSamples(p.smplr.Hash(vector.CastFixed(position)), 2)
	sum := vector.Vec{0, 0, 0}
	for i := 0; i < len(sampleList); i++ {
		part := p.particles[sampleList[i]]
		dist := vector.Dist(part.Position(), position)
		weight := p.mass / part.Density() * p.kern.F(dist)
		vector.Add(sum, vector.Scale(field.Value(sampleList[i]), weight))
	}
	return sum
}

//Density -- Computes density field for SPH Field
func (p SPHField) Density(i int) {
	sampleList := p.smplr.GetSamples(i)
	weight := p.kern.W0()

	for j := 0; j < len(sampleList); j++ {
		pIndex := sampleList[j]
		if i != pIndex {
			dist := vector.Dist(p.particles[i].Position(), p.particles[pIndex].Position()) //Change to dist
			weight += p.kern.F(dist)
		}
	}
	p.particles[i].SetDensity(float32(p.mass * weight))
}

//Computes gradient vector at particle i given a scalar field
func (p SPHField) Gradient(i int, field Field) []float32 {

	samples := p.smplr.GetSamples(i)
	dens := p.particles[i].Density()
	F := float32(0.0)
	mass := p.mass
	accumGrad := vector.Vec{}

	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			jDensity := p.particles[jIndex].Density()
			dir := vector.Sub(p.particles[samples[j]].Position(), p.particles[i].Position())
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
func (p SPHField) Div(i int, field TensorField) float32 {

	samples := p.smplr.GetSamples(i)
	div := float32(0.0)
	mass := p.mass

	//For all particle neighbors -- Non Symmetric
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			jDensity := p.particles[samples[j]].Density()
			dir := vector.Sub(p.particles[samples[j]].Position(), p.particles[i].Position())
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
func (p SPHField) Laplacian(i int, field Field) float32 {

	samples := p.smplr.GetSamples(i)
	m := p.mass
	sum := float32(0.0)
	//Conduct inner loop
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			jDensity := p.particles[samples[j]].Density()
			dist := vector.Dist(p.particles[i].Position(), p.particles[samples[j]].Position())
			sum += m * ((field.Value(samples[j]) - field.Value(i)) / jDensity) * p.kern.O2D(dist)
		}
	}
	return sum
}

//Computes a laplacian value at the particle i for the given scalar field
func (p SPHField) LaplacianForce(i int, field TensorField) []float32 {

	samples := p.smplr.GetSamples(i)
	m := p.mass
	force := vector.Vec{0, 0, 0}
	//Conduct inner loop
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			jDensity := p.particles[samples[j]].Density()
			v := vector.Scale(vector.Sub(p.particles[samples[j]].Velocity(), p.particles[i].Velocity()), 1/jDensity)
			dist := vector.Dist(p.particles[i].Position(), p.particles[samples[j]].Position())
			force = force.Add(v.Scale(p.kern.O2D(dist))).Scale(m)
		}
	}
	return force
}

//Curl computes non-symmetric curl
func (p SPHField) Curl(i int, field TensorField) []float32 {

	samples := p.smplr.GetSamples(i)
	curl_vec := vector.Vec{}
	mass := p.mass

	//For all particle neighbors
	for j := 0; j < len(samples); j++ {
		jIndex := samples[j]
		if jIndex != i {
			jDensity := p.particles[samples[j]].Density()
			dir := vector.Sub(p.particles[samples[j]].Position(), p.particles[i].Position())
			dist := vector.Mag(dir)
			dir = vector.Norm(dir) //Normalize
			grad := p.kern.Grad(dist, dir)
			scaleVec := vector.Scale(field.Value(samples[j]), mass/jDensity)
			curl_vec = vector.Add(curl_vec, vector.Cross(scaleVec, grad))
		}
	} //End J
	return curl_vec
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

func (p SPHField) EOSGamma(x float32, c0 float32, d0 float32, gamma float32, p0 float32) float32 {
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
func (p SPHField) TaitEOS(x float32, p0 float32) float32 {
	g := float32(7.16)
	w := float32(2.15)
	d0 := p.d0
	y := (w/g)*float32(math.Pow(float64(x/d0), float64(g))-1) + p0
	if y <= p0 {
		return p0
	}
	return y
}
