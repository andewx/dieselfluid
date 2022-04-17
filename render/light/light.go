package light

import (
	"fmt"

	"github.com/andewx/dieselfluid/math/vector"
)

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
	WATTS             = 0 // Default unit
	LUMENS            = 1
)

type Light interface {
	Lx() *Source
	Position() vector.Vec
	Type() int
	SetRGB(vector.Vec)
	GetRGB() vector.Vec
	GetFlux() float32
	SetFlux(float32)
	SetPos(vector.Vec)
	GetDir() vector.Vec
	SetDir(vector.Vec)
}

//Source - Luminous Flux and Light Color in RGB Color Space. Systems dependent on
//physical light interaction such as absorption, scattering should use light.SPD
//Units are either based in Watts or Lumens per meter area
type Source struct {
	RGB  vector.Vec
	Flux float32
	Unit int
}

//Attenuated lights cast a total luminance across a unit sphere which amounts
//inverse squared attenuation factor given the total luminal output in lux.
// Attenuated lights are synonymous with point lights.
type Attenuated struct {
	Pos vector.Vec
	Lum Source
}

//Directional lights are infinite lights with apparent-non local divergence for
//ray paths. Typically position is not a factor for directional lights but we
//may model occlusions etc with a position
type Directional struct {
	Pos vector.Vec
	Dir vector.Vec
	Lum Source
}

//Area lights are attenuated lights with a restrictive projected solid
//angle. For a point to recieve light it must be encapsulate by the projecting
//solid angle. Lux is taken to a total output across this solid angle max
type Area struct {
	Pos     vector.Vec
	Norm    vector.Vec
	Cuttoff float32
	Lum     Source
}

//-------------Attenuated Point light------------------------//

func (p Attenuated) Lx() *Source {
	return &p.Lum
}

func (p Attenuated) Position() vector.Vec {
	return p.Pos
}

func (p Attenuated) Type() int {
	return ATTENUATED_LIGHT
}

func (p Attenuated) SetRGB(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Lum.RGB = a
	}
}
func (p Attenuated) GetRGB() vector.Vec {
	return p.Lum.RGB
}
func (p Attenuated) GetFlux() float32 {
	return p.Lum.Flux
}
func (p Attenuated) SetFlux(a float32) {
	p.Lum.Flux = a
}
func (p Attenuated) SetPos(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Pos = a
	}
}
func (p Attenuated) GetDir() vector.Vec {
	return nil
}

func (p Attenuated) SetDir(a vector.Vec) {

}

//-----------Directional Light----------------//

func (p Directional) Lx() *Source {
	return &p.Lum
}

func (p Directional) Position() vector.Vec {
	return p.Pos
}

func (p Directional) Type() int {
	return DIRECTIONAL_LIGHT
}

func (p Directional) SetRGB(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Lum.RGB = a
	}
}
func (p Directional) GetRGB() vector.Vec {
	return p.Lum.RGB
}
func (p Directional) GetFlux() float32 {
	return p.Lum.Flux
}
func (p Directional) SetFlux(a float32) {
	p.Lum.Flux = a
}
func (p Directional) SetPos(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Pos = a
	}
}
func (p Directional) GetDir() vector.Vec {
	return p.Dir
}

func (p Directional) SetDir(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Dir = a
	}
}

//-----------Area Light----------------//

func (p Area) Lx() *Source {
	return &p.Lum
}

func (p Area) Position() vector.Vec {
	return p.Pos
}

func (p Area) Type() int {
	return AREA_LIGHT
}

func (p Area) SetRGB(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Lum.RGB = a
	}
}
func (p Area) GetRGB() vector.Vec {
	return p.Lum.RGB
}
func (p Area) GetFlux() float32 {
	return p.Lum.Flux
}
func (p Area) SetFlux(a float32) {
	p.Lum.Flux = a
}
func (p Area) SetPos(a vector.Vec) {
	if a != nil && len(a) == 3 {
		p.Pos = a
	}
}
func (p Area) GetDir() vector.Vec {
	return nil
}

func (p Area) SetDir(a vector.Vec) {
	fmt.Printf("Area light has no direction")
}
