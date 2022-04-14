package light

import (
	"math"

	"github.com/andewx/dieselfluid/math/mgl"
)

/*
Sky environment lighting model.

This model generates a sun position based on lat/long and solar times and then
we simulate atmospheric scattering processes via Rayleigh/Mie Scattering.

The resulting model can be sampled generating 3D Sampler Textures producing a
realistic sky environment mode
*/
const (
	PI               = 3.141529
	SEASONAL_PERIOD  = 181.8
	AXIAL_TILT       = 23.5      //Degrees
	CLEAR_LUX        = 105000.0  //Sun Lumens
	ANGULAR_VELOCITY = 15.0      //per hour
	MDSL             = 1000.0    //Molecular density at sea level
	IOR_AIR          = 1.000293  //IOR
	RLH_440          = 0.0000331 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 440NM (BLUE)
	RLH_550          = 0.0000135 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 550NM (GREEN)
	RLH_680          = 0.0000058 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 680NM (RED)
	HR               = 8.2500000 //HEIGHT ATMOSHPEHRE IN KM APPROXIMATE
	MIE              = 0.00210   //MIE SCATTER COEFFICIENT
)

//Sky Environment
type Sky struct {
	Lgt     Light
	Spd     Spectrum
	Coords0 mgl.Polar //Orbtial Solar System Earth 2 Sun Polar Coordinate
	Earth   *EarthCoords
	Day     float32
	Dir     mgl.Vec //Euclidian Sun Direction
}

type ScatterCoefficients struct {
	SK_440 float32
	SK_550 float32
	SK_680 float32
}

//Scatter Coefficients
func NewScatterCoefficients() *ScatterCoefficients {
	return &ScatterCoefficients{RLH_440, RLH_550, RLH_680}
}

//Allocates Default Data Structure and Solar Coords Structs
func NewSky() *Sky {
	sky := Sky{}
	sky.Earth = NewEarth(0, 0, 0)
	sky.Dir = mgl.Vec{}
	sky.Coords0 = mgl.Polar{}
	sky.Lgt = Directional{mgl.Vec{0, 1, 0}, Source{mgl.Vec{1, 1, 1}, 150, WATTS}}
	sky.Spd = InitSunlight(20)
	sky.Dir = mgl.Vec{}
	return &sky
}

//Updates Frame of Reference Solar Coordinates with regards to Decimal Day Local Time
func (sky *Sky) UpdateDay(day float32) error {
	var err error
	sky.Day = day
	polar := mgl.Polar{1, 0, 0}
	//Orbital Coordinates + Earth Axial Tilt
	solarRotation := mgl.Polar{1, -day / 365.0 * 2 * PI, -23.44 * mgl.DEG2RAD}
	//Rotate Azimuthal Day
	polar.Add(solarRotation).AddAzimuth(Day2Rotation(day))
	//Rotate Lat / Long Coords - Minus Standard Meridian
	polar.Add(sky.Earth.PolarCoord).AddAzimuthDegrees(-sky.Earth.StandardMeridian)
	sky.Coords0 = polar
	sky.Dir, err = mgl.Sphere2Vec(sky.Coords0)
	return err
}

func (sky *Sky) BuildSkyBox() {
	//Update Sky Position
	//Initialize Spectral Sun Data
	//Initialize CIE XYZ Structures
	//Create Sphere Sampler Pattern - Map Sphere Coverage Pattern to 6 x 128 x 128 Cube Faces
	//Conduct Sampler Ray Pattern, Store RGB Results in Cube Map Relay Cube to Renderer
}

//Updates scaterring coefficients based on the parameter height calcualtes
//approx air density as exponential parameter normalized in respect to p0 sea level density
func (strct *ScatterCoefficients) UpdateHeight(h float32) {
	hg := float64(h)
	k := float32(math.Exp(-hg / HR))
	strct.SK_440 = RLH_440 * k
	strct.SK_550 = RLH_550 * k
	strct.SK_680 = RLH_680 * k
}

func RayleigthPhase(u float32) float32 {
	return (3 / (16 * PI)) * (1 + (u * u))
}

func MiePhase(u float32) float32 {
	g := float32(0.76)
	num := (1 - (g * g)) * (1 + (u * u))
	denom := float32(math.Pow(float64((2+g*g)*(1+g*g-2*g*u)), 1.5))
	return (3 / (8 * PI)) * (num / denom)
}

//Returns rayleigh scatter coefficients for depth and wavelength of light
//This is a utility that must be integrated across an SPD
func (sun *EarthCoords) RayleighCoeff(h float64, km float64) float64 {
	if h <= 0 {
		h = 0.001
	}
	hr := 8.0 //Scale Height KM
	return (8 * PI * PI * PI * (IOR_AIR - 1) * (IOR_AIR - 1) * math.Exp(-h/hr)) / (3 * MDSL * km * km * km * km)
}
