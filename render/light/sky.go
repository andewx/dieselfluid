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
	"github.com/andewx/dieselfluid/sampler"
)

/*
Sky environment lighting model.

This model generates a sun position based on lat/long and solar times and then
we simulate atmospheric scattering processes via Rayleigh/Mie Scattering.

Compute Intensive Consider putting this into a job interface for progress and halts
*/
const (
	PI                 = 3.141529
	AXIAL_TILT         = 23.5      //Degrees
	CLEAR_LUX          = 105000.0  //Sun Lumens
	MDSL               = 1000.0    //Molecular density at sea level
	IOR_AIR            = 1.000293  //IOR
	RLH_440            = 0.0000331 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 440NM (BLUE)
	RLH_550            = 0.0000135 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 550NM (GREEN)
	RLH_680            = 0.0000058 //RAYLEIGTH SCATTER COEFFICIENT SEA LEVEL 680NM (RED)
	MIE                = 0.00210   //MIE SCATTER COEFFICIENT
	RAYLEIGH_SAMPLES   = 16        //RAYLEIGH sampling
	LIGHT_PATH_SAMPLES = 30        //MIE SAMPLING
	ATTENUATION        = 0.05      //Attenuation of Light in WATTS per Density
	AU                 = 150000000 //150 million KM
	HM                 = 1200      //m

)

//Sky Environment
type Sky struct {
	Light Light
	Spd   Spectrum
	Sun   mgl.Polar //Orbtial Solar System Earth 2 Sun Polar Coordinate
	Earth *EarthCoords
	Day   float32
	Dir   mgl.Vec //Euclidian Sun Direction
}

//Allocates Default Data Structure and Solar Coords Structs
func NewSky(lat float32, long float32) *Sky {
	sky := Sky{}
	sky.Earth = NewEarth(45.0, 0)
	sky.Sun = mgl.NewPolar(-AU)
	sky.Light = Directional{mgl.Vec{-1, 0, 0}, mgl.Vec{-1, 0, 0}, Source{mgl.Vec{1, 1, 1}, 20.5, WATTS}}
	sky.Spd = InitSunlight(20)
	sky.SetDay(1.0)
	return &sky
}

//Updates Frame of Reference Solar Coordinates with regards to Decimal Day Local Time
func (sky *Sky) StepDay(day float32) error {
	var err error
	axialMag := float32(2 * AXIAL_TILT)
	sky.Day += day
	u := (sky.Day / 12.0) / (365.0 / 12.0) //Normalized cos units
	axialTilt := AXIAL_TILT - float32(math.Cos(float64(u)))*axialMag
	sky.Sun.AddAzimuth(-sky.Day/365.0*mgl.DEG2RAD + (-sky.Day*24)*mgl.DEG2RAD + sky.Earth.Longitude) //Rotate sun directional azimuth
	sky.Sun.AddPolar(-axialTilt*mgl.DEG2RAD + sky.Earth.Latitude)
	sky.Dir, err = mgl.Sphere2Vec(sky.Sun)
	sky.Dir = sky.Dir.Norm()
	sky.Light.SetDir(sky.Dir)
	return err
}

//Updates Frame of Reference Solar Coordinates with regards to Decimal Day Local Time
func (sky *Sky) SetDay(day float32) error {
	var err error
	axialMag := float32(2 * AXIAL_TILT)
	sky.Day = day
	u := (sky.Day / 12.0) / (365.0 / 12.0) //Normalized cos units
	axialTilt := AXIAL_TILT - float32(math.Cos(float64(u)))*axialMag
	sky.Sun.AddAzimuth(-day/365.0*mgl.DEG2RAD + (-day*24)*mgl.DEG2RAD + sky.Earth.Longitude) //Rotate sun directional azimuth
	sky.Sun.AddPolar(-axialTilt*mgl.DEG2RAD + sky.Earth.Latitude)
	sky.Dir, err = mgl.Sphere2Vec(sky.Sun)
	sky.Dir = sky.Dir.Norm()
	sky.Light.SetDir(sky.Dir)
	return err
}

func (sky *Sky) CreateTexture(width int, height int, filename string) {
	wd, _ := os.Getwd()
	fmt.Printf("Working Dir: %s\n", wd)
	rgbs := sky.ComputeAtmosphere(width, height)
	corner := image.Point{0, 0}
	bottom := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{corner, bottom})
	index := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := uint8(rgbs[index][0] * 255)
			g := uint8(rgbs[index][1] * 255)
			b := uint8(rgbs[index][2] * 255)
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
			index++
		}
	}

	//Encode as PNG
	f, _ := os.Create(filename)
	err := png.Encode(f, img)

	if err != nil {
		fmt.Printf("Error writing image to %s\n", filename)
	}

}

//Maps texel coordinates to spherical coordinate sampler values (-1,1) and stores
//resultant map in single texture.
func (sky *Sky) ComputeAtmosphere(uSampleDomain int, vSampleDomain int) []mgl.Vec {
	sizeT := uSampleDomain * vSampleDomain
	tex := make([]mgl.Vec, sizeT+1)
	index := 0

	for x := 0; x < uSampleDomain; x++ {
		u := 2.0*(float64(x)+0.5)/(float64(uSampleDomain-1)) - 1.0
		for y := 0; y < vSampleDomain; y++ {
			v := 2.0*(float64(y)+0.5)/(float64(vSampleDomain-1)) - 1.0
			z2 := u*u + v*v
			phi := math.Atan2(v, u)
			theta := math.Acos(1 - z2)
			sample := mgl.Vec{float32(math.Sin(theta) * math.Cos(phi)), float32(math.Cos(theta)), float32(math.Sin(theta) * math.Sin(phi))}
			if index < len(tex) {
				tex[index] = sky.VolumetricScatterRay(sample, mgl.Vec{0, 1, 0})
			}
			index++
		}
	}
	return tex
}

//Given a sampling vector and a viewing direction calculate RGB stimulus return
//Based on the Attenuation/Mie Phase Scatter/RayleighScatter Terms
func (sky *Sky) VolumetricScatterRay(sample mgl.Vec, view mgl.Vec) mgl.Vec {

	//Declare volumetric scatter ray vars
	intersects := mgl.RaySphereIntersect(sample, sky.Earth.GetPosition(), sky.Earth.GreaterSphere)
	rgb := mgl.Vec{0, 0, 0} //Pixel Output

	betaR := mgl.Vec{0.0000088, .0000135, 0.0000331}
	betaM := mgl.Vec{0.000021, 0.000021, 0.000021}
	sumR := mgl.Vec{0, 0, 0}
	sumM := mgl.Vec{0, 0, 0}
	u := mgl.Dot(sample, sky.Dir) / (sample.Mag() * sky.Dir.Mag())
	mu := float64(u)
	phaseR := float32(3.0 / (16.0 * PI) * (1.0 + mu*mu))
	g := 0.76
	phaseM := float32(3.0 / (8.0 * PI) * ((1.0 - g*g) * (1.0 + mu*mu)))
	var opticalDepthR, opticalDepthM float32
	var vmag0, vmag1, vds float32

	//Scatter Along View Rays
	sampleStep := float32(1.0 / RAYLEIGH_SAMPLES)
	for i := 1; i <= RAYLEIGH_SAMPLES && len(intersects) > 0; i++ {

		//Generate Sample Rays along view path - assume sample ray is normalized
		w := float32(5.0)
		sampleScale := sampler.Ease(float32(i)*sampleStep, w)
		viewRay := mgl.Scale(sample, mgl.Priority(intersects))
		viewRayMag := viewRay.Mag()
		viewSample := mgl.Scale(viewRay, sampleScale)
		depth := sky.Earth.GetSampleDepth(viewSample)

		//Compute the view ray ds parameters
		vmag1 = viewRayMag * sampleScale
		vds = vmag1 - vmag0
		vmag0 = vmag1

		//Get optical depth
		hr := float32(math.Exp(float64(-depth/HR))) * vds
		hm := float32(math.Exp(float64(-depth/HM))) * vds
		opticalDepthR += hr
		opticalDepthM += hm

		//Constructs Light Path Rays from Viewer Sample Positions and Calculates
		viewSampleOrigin := mgl.Add(viewSample, sky.Earth.GetPosition())
		lightIntersects := mgl.RaySphereIntersect(mgl.Scale(sky.Dir, -1.0), viewSampleOrigin, sky.Earth.GreaterSphere)

		if len(lightIntersects) == 0 {
			fmt.Printf("No Ray Sphere Intersection")
			return mgl.Vec{0, 0, 0}
		}

		//Light Path Transmittance + Attenutation
		lightRay := mgl.Scale(mgl.Scale(sky.Dir, -1.0), mgl.Priority(lightIntersects))
		lightRayMag := lightRay.Mag()
		lightPathSampleStep := float32(1.0 / LIGHT_PATH_SAMPLES)

		//Light Transmittance
		ds := float32(0.0)   //differential magnitude
		mag0 := float32(0.0) //for calculating differential
		mag1 := float32(0.0) //for calculating differntial

		var opticalDepthLightM, opticalDepthLightR float32

		//Compute light path to sample position
		for j := 1; j <= LIGHT_PATH_SAMPLES; j++ {
			pathScale := sampler.Ease(lightPathSampleStep*float32(j), w)
			lightPath := mgl.Scale(lightRay, pathScale)
			mag1 = pathScale * lightRayMag
			ds = mag1 - mag0
			mag0 = mag1
			lightPathSamplePosition := mgl.Add(viewSample, lightPath)
			lightPathDepth := sky.Earth.GetSampleDepth(lightPathSamplePosition)

			if lightPathDepth < 0 {
				break
			}

			//Accumlate Light Path Transmittance
			opticalDepthLightR += float32(math.Exp(float64(-lightPathDepth/HR))) * ds
			opticalDepthLightM += float32(math.Exp(float64(-lightPathDepth/HM))) * ds
		}

		//Compute Contributions and Accumulate
		tau := betaR.Scale(opticalDepthR + opticalDepthLightR).Add(betaM.Scale(1.1).Scale(opticalDepthM + opticalDepthLightM))
		attenuation := mgl.Vec{float32(math.Exp(float64(-tau[0]))), float32(math.Exp(float64(-tau[1]))), float32(math.Exp(float64(-tau[2])))}
		sumR = sumR.Add(attenuation.Scale(hr))
		sumM = sumM.Add(attenuation.Scale(hm))
	}
	rayliegh := sumR.Mul(betaR).Scale(phaseR)
	mie := sumM.Mul(betaM).Scale(phaseM)
	rgb = rayliegh.Add(mie).Scale(sky.Light.Lx().Flux).Mul(sky.Light.Lx().RGB)
	rgb.Clamp(0, 1) //Dev Note: Add in a mechanism to produce HDR with excess luminance
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

func RayleighCoefficient(wv float32, u float32, h float32) float32 {
	n := float32(1.00029)
	N := 2.50 * float32(math.Pow(10, 25))
	if h <= 0 {
		h = 0.001
	}
	wv4 := wv * wv * wv * wv
	res := 1 / wv4 * 1 / N
	exp := float32(math.Exp(float64(-h / HR)))
	coeff := ((PI * PI) * (n*n - 1) * (n*n - 1) / 2) * (exp * res) * (1 + u*u)
	return coeff
}

//Returns rayleigh scatter coefficients for depth and wavelength of light
//This is a utility that must be integrated across an SPD
func (sun *EarthCoords) RayleighCoeff(h float64, km float64) float64 {
	if h <= 0 {
		h = 0.001
	}
	hr := 12.0 //Scale Height KM
	return (8 * PI * PI * PI * (IOR_AIR - 1) * (IOR_AIR - 1) * math.Exp(-h/hr)) / (3 * MDSL * km * km * km * km)
}
