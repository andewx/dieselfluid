package light

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math"
	"os"

	"github.com/andewx/dieselfluid/math/mgl"
)

/*
Sky environment lighting model.

This model generates a sun position based on lat/long and solar times and then
we simulate atmospheric scattering processes via Rayleigh/Mie Scattering.

The resulting model can be sampled generating 3D Sampler Textures producing a
realistic sky environment mode

Compute Intensive Consider putting this into a job interface for progress and halts
*/
const (
	PI               = 3.141529
	AXIAL_TILT       = 23.5      //Degrees
	CLEAR_LUX        = 105000.0  //Sun Lumens
	MDSL             = 1000.0    //Molecular density at sea level
	IOR_AIR          = 1.000293  //IOR
	RLH_440          = 0.0000331 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 440NM (BLUE)
	RLH_550          = 0.0000135 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 550NM (GREEN)
	RLH_680          = 0.0000058 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 680NM (RED)
	MIE              = 0.00210   //MIE SCATTER COEFFICIENT
	RAYLEIGH_SAMPLES = 25        //RAYLEIGH sampling
	ATTEN_SAMPLES    = 10        //MIE SAMPLING
	ATTENUATION      = 0.00250   //Attenuation of Light in WATTS per Density

)

//Sky Environment
type Sky struct {
	Light   Light
	Spd     Spectrum
	Coords0 mgl.Polar //Orbtial Solar System Earth 2 Sun Polar Coordinate
	Earth   *EarthCoords
	Day     float32
	Dir     mgl.Vec //Euclidian Sun Direction
}

//Allocates Default Data Structure and Solar Coords Structs
func NewSky() *Sky {
	sky := Sky{}
	sky.Earth = NewEarth(0, 0, 0)
	sky.Dir = mgl.Vec{}
	sky.Coords0 = mgl.Polar{}
	sky.Light = Directional{mgl.Vec{0, 1, 0}, Source{mgl.Vec{1, 1, 1}, 150, WATTS}}
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

func (sky *Sky) CreateTexture() {
	wd, _ := os.Getwd()
	fmt.Printf("td:Add Texture Folder Select Prompt\n")
	fmt.Printf("Working Dir: %s\n", wd)
	sp := "/Users/briananderson/go/src/github.com/dieselfluid/SKY.png"
	rgbs := sky.ComputeAtmosphere(45, 45)
	corner := image.Point{0, 0}
	bottom := image.Point{45, 45}
	img := image.NewRGBA(image.Rectangle{corner, bottom})
	index := 0
	for x := 0; x < 45; x++ {
		for y := 0; y < 45; y++ {
			r := uint8(rgbs[index][0] * 255)
			g := uint8(rgbs[index][1] * 255)
			b := uint8(rgbs[index][2] * 255)
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
			index++ //(x*y)+y ??????
		}
	}

	//Encode as PNG
	f, _ := os.Create(sp)
	png.Encode(f, img)

}

//Maps texel coordinates to spherical coordinate sampler values (-1,1) and stores
//resultant map in single texture.
func (sky *Sky) ComputeAtmosphere(uSampleDomain int, vSampleDomain int) []mgl.Vec {
	sizeT := uSampleDomain * vSampleDomain
	tex := make([]mgl.Vec, sizeT)
	index := 0
	for x := -1.0; x < 1.0; x += 2.0 / float64(uSampleDomain) {
		for y := -1.0; y < 1.0; y += 2.0 / float64(vSampleDomain) {
			var uv [2]float32
			uv[0] = float32(x)
			uv[1] = float32(y)
			sample := sky.Earth.GetSample(uv)
			tex[index] = sky.VolumetricScatterRay(sample, mgl.Vec{0, 1, 0})
			index++
		}
	}
	return tex
}

//Given a sampling vector and a viewing direction calculate RGB stimulus return
//Based on the Attenuation/Mie Phase Scatter/RayleighScatter Terms
func (sky *Sky) VolumetricScatterRay(sample mgl.Vec, view mgl.Vec) mgl.Vec {

	sampleStep := float32(1.0 / RAYLEIGH_SAMPLES)
	viewSample := mgl.Vec{}
	rayleighDensity := float32(0.0)
	rgb := mgl.Vec{}

	//Construct initial sampler ray
	rE_Vec, _ := mgl.Sphere2Vec(sky.Earth.PolarCoord)

	//Calculate Ray Attenuation (Mie Scatter-Attenuation and Rayleigh Scatter)
	//Density are related to scatter density
	for i := 1; i <= RAYLEIGH_SAMPLES; i++ {
		viewSample = mgl.Scale(sample, float32(i)*sampleStep)
		sampleDensity := sky.Earth.GetSampleDensity(viewSample)
		rayleighDensity += sampleDensity //Rayleight
		viewSampleOrigin := mgl.Add(viewSample, rE_Vec)
		viewSampleSphereIntersection, flag := mgl.RaySphereIntersection(sky.Dir, viewSampleOrigin,
			mgl.Vec{0, 0, 0}, sky.Earth.GreaterSphere.Radius())

		if !flag {
			fmt.Printf("No Ray Sphere Intersection")
			return mgl.Vec{}
		}

		lightRaySegment := mgl.Sub(viewSampleOrigin, viewSampleSphereIntersection)
		//Compute Optical Depth For Light Attenuation
		attenuationSample := float32(sampleDensity)
		attenuationSampleStep := float32(1.0 / ATTEN_SAMPLES)

		//Compute Light Ray Attenuation
		for j := 1; j <= ATTEN_SAMPLES; j++ {
			attenuationSampleVec := mgl.Scale(lightRaySegment, attenuationSampleStep*float32(j))
			attenutationOrigin := mgl.Add(viewSample, attenuationSampleVec)
			opticalDepth := sky.Earth.GetSampleDensity(attenutationOrigin)
			attenuationSample += opticalDepth
		}

		u := mgl.Dot(view, sky.Dir)

		//Rayleigh Scatter this watts coefficient and accumalte RGB Scattering
		watts := attenuationSample * ATTENUATION * sky.Light.Lx().Flux
		watts *= MiePhase(u) * sampleDensity
		rgb.Add(sky.Light.Lx().RGB.Scale(attenuationSample * ATTENUATION * MiePhase(u) * sampleDensity))

		//Compute the Rayleigh Contribution
		rayleighPhase := RayleighPhase(u)
		raylieghRGB := mgl.Vec{sky.Light.Lx().RGB[0] * sampleDensity * rayleighPhase * RLH_440,
			sky.Light.Lx().RGB[1] * sampleDensity * rayleighPhase * RLH_550,
			sky.Light.Lx().RGB[2] * sampleDensity * rayleighPhase * RLH_680}

		//Accumulate Rayleigh
		rgb.Add(raylieghRGB)

	}

	return rgb
}

func RayleighPhase(u float32) float32 {
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
