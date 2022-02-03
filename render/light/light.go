package light

import "dslfluid.com/dsl/math/mgl"

//----------------------------------------------------------------------------//
//  DslFluid lighting model uses SI lm/m^2 units for radiance calculations    //
//  Which is a lux lighting model for point/attenuated/directional lights     //
//  Area lights are approximated via Lighting Rigs whose light sources can be //
//  useful for approximation of area lights and light probes                  //
//----------------------------------------------------------------------------//

//----------------------------------------------------------------------------//
// Note that most lights are intended to be used and calc'd inside the render //
// shading system. We can only pass the basic uniforms that describe the      //
// lighting model. So most of the function parametres are useful for offline  //
// render computations
//----------------------------------------------------------------------------//

const (
	ATTENUATED_LIGHT  = 1
	DIRECTIONAL_LIGHT = 2
	AREA_LIGHT        = 3
)

type Light interface {
	Luminosity(p mgl.Vec) float32
	Position() mgl.Vec
	Rgb() mgl.Vec
	Type() int
}

type Lux struct {
	RGB mgl.Vec
	Lm2 float32
}

//Attenuated lights cast a total luminance across a unit sphere which amounts
//inverse squared attenuation factor given the total luminal output in lux.
// Attenuated lights are synonymous with point lights.
type Attenuated struct {
	Pos mgl.Vec
	Lx  Lux
}

//Directional lights are infinite lights with apparent-non local divergence for
//ray paths. Typically position is not a factor for directional lights but we
//may model occlusions etc with a position
type Directional struct {
	Pos  mgl.Vec
	Lx   Lux
	Norm mgl.Vec
}

//Area lights are attenuated lights with a restrictive projected solid
//angle. For a point to recieve light it must be encapsulate by the projecting
//solid angle. Lux is taken to a total output across this solid angle max
type Area struct {
	Pos     mgl.Vec
	Norm    mgl.Vec
	Cuttoff float32
	Lx      Lux
}

//-------------Attenuated Point light------------------------//

func (p *Attenuated) Luminosity(p1 mgl.Vec) float32 {
	dist := mgl.Mag(p1.Sub(p.Pos))
	return p.Lx.Lm2 / (dist * dist)
}

func (p *Attenuated) Rgb() mgl.Vec {
	return p.Lx.RGB
}

func (p *Attenuated) Position() mgl.Vec {
	return p.Pos
}

func (p *Attenuated) Type() int {
	return ATTENUATED_LIGHT
}

//-----------Directional Light----------------//

func (p *Directional) Luminosity(p1 mgl.Vec) float32 {
	return p.Lx.Lm2
}

func (p *Directional) Rgb() mgl.Vec {
	return p.Lx.RGB
}

func (p *Directional) Position() mgl.Vec {
	return p.Pos
}

func (p *Directional) Type() int {
	return DIRECTIONAL_LIGHT
}

//-----------Area Light----------------//

func (p *Area) Luminosity(p1 mgl.Vec) float32 {
	dot := mgl.Dot(p.Pos.Sub(p1), p.Norm)
	if dot < p.Cuttoff {
		dist := mgl.Mag(p.Pos.Sub(p1))
		return p.Lx.Lm2 / (dist * dist)
	} else {
		return 0.0
	}
}

func (p *Area) Rgb() mgl.Vec {
	return p.Lx.RGB
}

func (p *Area) Position() mgl.Vec {
	return p.Pos
}

func (p *Area) Type() int {
	return AREA_LIGHT
}
