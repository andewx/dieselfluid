package light

import "github.com/andewx/dieselfluid/math/mgl"
import "math"

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
	NYC_SM           = 75.0      //Standard Meridian NYC
	MDSL             = 1000.0    //Molecular density at sea level
	IOR_AIR          = 1.000293  //IOR
	RLH_440          = 0.0000331 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 440NM (BLUE)
	RLH_550          = 0.0000135 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 550NM (GREEN)
	RLH_680          = 0.0000058 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 680NM (RED)
	HR               = 8.2500000 //HEIGHT ATMOSHPEHRE IN KM APPROXIMATE
	MIE              = 0.00210   //MIE SCATTER COEFFICIENT
)

//Sky Environment
type SkyEnvironment struct {
	Lgt    Light
	Spd    SampledSpectrum
	Coords SolarCoords
	Dir    mgl.Vec
}

type ScatterCoefficients struct {
	SK_440 float32
	SK_550 float32
	SK_680 float32
}

type SolarCoords struct {
	Lgt    Light   //Light Properties
	Jul    float64 //Julian Date
	Tm     float64 //Solar Time
	Decl   float64 //Solar Declination
	Atten  float32 //Light Luminance Attenuation Factor
	Zen    float64 //Solar Zenith
	Az     float64 //Solar Azimuth
	Sm     float64 //Standard Meridian Longitutde
	Lat    float64 //latitude
	Long   float64 //Longitude
	Dir    mgl.Vec //Ray Direction
	RayDir mgl.Vec //Ray Direction
}

//Scatter Coefficients
func Alloc_ScatterCoeff() *ScatterCoefficients {
	return &ScatterCoefficients{RLH_440, RLH_550, RLH_680}
}

//Allocates Default Data Structure and Solar Coords Structs
func Alloc_SkyEnv() *SkyEnvironment {
	return nil
}

func (sky *SkyEnvironment) BuildSkyBox() {
	//Update Sky Position
	//Initialize Spectral Sun Data
	//Initialize CIE XYZ Structures
	//Create Sphere Sampler Pattern - Map Sphere Coverage Pattern to 6 x 128 x 128 Cube Faces
	//Conduct Sampler Ray Pattern, Store RGB Results in Cube Map Relay Cube to Renderer
}

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
func (sun *SolarCoords) RayleighCoeff(h float64, km float64) float64 {
	if h <= 0 {
		h = 0.001
	}
	hr := 8.0 //Scale Height KM
	return (8 * PI * PI * PI * (IOR_AIR - 1) * (IOR_AIR - 1) * math.Exp(-h/hr)) / (3 * MDSL * km * km * km * km)
}

//Luminance&Color mapping to sun angle
func (sun *SkyEnvironment) Color(hour float32) Source {
	h := HR2RAD(hour)
	mLux := sun.Lgt.Lx()
	fac := 0.3 * float32(math.Sin(float64(h)))
	nRGB := mgl.Scale(mLux.RGB, fac)
	lum := mLux.Flux * fac
	return Source{nRGB, lum, 0}
}
