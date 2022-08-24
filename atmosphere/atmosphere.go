package atmosphere

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"strconv"

	"github.com/andewx/dieselfluid/math/common"
	"github.com/andewx/dieselfluid/math/polar"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/render/light"
	"github.com/andewx/dieselfluid/sampler"
)

/*
Atmosphere environment lighting model.

This model generates a sun position based on lat/long and solar times and then
we simulate atmospheric scattering processes via Rayleigh/Mie Scattering.
*/
const (
	PI                 = 3.141529
	AXIAL_TILT         = 23.5      //Degrees
	RAYLEIGH_SAMPLES   = 24        //RAYLEIGH sampling
	LIGHT_PATH_SAMPLES = 24        //PATH SAMPLES
	AU                 = 150000000 //SUN
	HM                 = 1200      //AEROSOL MIE SCATTER HEIGHT
	HR                 = 8500      //RAYLEIGH SCATTER HEIGHT

)

type Domain struct {
	min float32
	max float32
}

//Maps a value with an associated domain to this Domain
func (m *Domain) Map(value float32, domain Domain) float32 {
	r := domain.max - domain.min
	norm := (r - value) / r
	return m.min + norm*(m.max-m.min)
}

//Atmosphere Environment
type Atmosphere struct {
	Light light.Light
	Spd   light.Spectrum
	Sun   polar.Polar //Orbtial Solar System Earth 2 Sun Polar Coordinate
	Earth *EarthCoords
	Day   float32
	Dir   vector.Vec //Euclidian Sun Direction
}

//Allocates Default Data Structure and Solar Coords Structs
func NewAtmosphere(lat float32, long float32) *Atmosphere {
	sky := Atmosphere{}
	sky.Earth = NewEarth(65.0, 0)
	sky.Sun = polar.NewPolar(-AU)
	sky.Light = light.Directional{vector.Vec{0, 0, 0}, vector.Vec{0, 0, 0}, light.Source{vector.Vec{1, 1, 1}, 18.5, light.WATTS}}
	sky.Spd = light.InitSunlight(20)
	sky.SetDay(1.0)
	return &sky
}

//Updates Frame of Reference Solar Coordinates with regards to Decimal Day Local Time
func (sky *Atmosphere) StepDay(day float32) error {
	var err error
	axialMag := float32(2 * AXIAL_TILT)
	sky.Day += day
	u := (sky.Day / 12.0) / (365.0 / 12.0) //Normalized cos units
	axialTilt := AXIAL_TILT - float32(math.Cos(float64(u)))*axialMag
	sky.Sun = sky.Sun.AddAzimuth(-sky.Day/365.0*common.DEG2RAD + (-sky.Day*24)*common.DEG2RAD + sky.Earth.Longitude) //Rotate sun directional azimuth
	sky.Sun = sky.Sun.AddPolar(-axialTilt*common.DEG2RAD + sky.Earth.Latitude)
	sky.Dir, err = polar.Sphere2Vec(sky.Sun)
	sky.Dir = sky.Dir.Norm()
	sky.Light.SetDir(sky.Dir)
	return err
}

//Updates Frame of Reference Solar Coordinates with regards to Decimal Day Local Time
func (sky *Atmosphere) SetDay(day float32) error {
	var err error
	axialMag := float32(2 * AXIAL_TILT)
	sky.Day = day
	u := (sky.Day / 12.0) / (365.0 / 12.0) //Normalized cos units
	axialTilt := AXIAL_TILT - float32(math.Cos(float64(u)))*axialMag
	sky.Sun.AddAzimuth(-day/365.0*common.DEG2RAD + (-day*24)*common.DEG2RAD + sky.Earth.Longitude) //Rotate sun directional azimuth
	sky.Sun.AddPolar(-axialTilt*common.DEG2RAD + sky.Earth.Latitude)
	sky.Dir, err = polar.Sphere2Vec(sky.Sun)
	sky.Dir = sky.Dir.Norm()
	sky.Light.SetDir(sky.Dir)
	return err
}

//Creates texture from from computed atmosphere, non-clamping allows for HDR storage
func (sky *Atmosphere) CreateTexture(width int, height int, clamp bool, filename string) {
	wd, _ := os.Getwd()
	fmt.Printf("Working Dir: %s\n", wd)
	rgbs := sky.ComputeAtmosphere(width, height) //pre-normalized (non-hdr)
	ImageFromPixels(rgbs, width, height, clamp, 0xff, filename)
}

//Creates 6-textures from from computed atmosphere, non-clamping allows for HDR storage
//width and height are the per texture width height values of the image size
func (sky *Atmosphere) CreateEnvBox(width int, height int, clamp bool) {
	wd, _ := os.Getwd()
	fmt.Printf("Working Dir: %s\n", wd)
	if height != width || height <= 0 || width <= 0 {
		fmt.Printf("Computed region must be to a square texture. Parameters don't pass safeguard\n")
		return
	}
	if height%4 != 0 || width%4 != 0 {
		fmt.Printf("Computed region must be modulo 4\n")
		return
	}
	//Generate each side face
	region_width := int(width / 2)
	region_height := int(height / 2)
	x_corner := 0
	y_corner := 0

	//Faces loop
	for i := 0; i < 4; i++ {
		rgbs := sky.ComputeRegion(width, height, x_corner, y_corner, region_width, region_height)
		ImageFromPixels(rgbs, region_width, region_height, clamp, 0xff, "ENVBOX_"+strconv.FormatInt(int64(i), 10)+".png")

		//Move computed x,y region
		if x_corner < region_width {
			x_corner += region_width
		} else {
			if y_corner < region_height {
				y_corner += region_height
			}
		}
	}

	//Generate top and bottom dregion
	rgbs := sky.ComputeRegion(width, height, region_width/2, region_height/2, region_width, region_height)
	ImageFromPixels(rgbs, region_width, region_height, clamp, 0xff, "ENVBOX_4.png")
	ImageFromPixels(rgbs, region_width, region_height, clamp, 0x44, "ENVBOX_5.png")

}

//Utility function creates an image from a set of pixels
func ImageFromPixels(pixels []vector.Vec, width int, height int, clamp bool, alpha uint8, filename string) {
	corner := image.Point{0, 0}
	bottom := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{corner, bottom})
	index := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := pixels[index]
			if clamp {
				pixel.Clamp(0, 1)
			}
			r := uint8(pixel[0] * 255)
			g := uint8(pixel[1] * 255)
			b := uint8(pixel[2] * 255)
			img.Set(x, y, color.RGBA{r, g, b, alpha})
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
func (sky *Atmosphere) ComputeAtmosphere(uSampleDomain int, vSampleDomain int) []vector.Vec {
	sizeT := uSampleDomain * vSampleDomain
	tex := make([]vector.Vec, sizeT+1)
	index := 0

	for x := 0; x < uSampleDomain; x++ {
		u := 2.0*(float64(x)+0.5)/(float64(uSampleDomain-1)) - 1.0
		for y := 0; y < vSampleDomain; y++ {
			v := 2.0*(float64(y)+0.5)/(float64(vSampleDomain-1)) - 1.0
			z2 := u*u + v*v
			phi := math.Atan2(v, u)
			theta := math.Acos(1 - z2)
			sample := vector.Vec{float32(math.Sin(theta) * math.Cos(phi)), float32(math.Sin(theta) * math.Sin(phi)), float32(math.Cos(theta))}
			if index < len(tex) {
				tex[index] = sky.VolumetricScatterRay(sample, vector.Vec{0, 1, 0})
			}
			index++
		}
	}
	return tex
}

//Maps texel coordinates to spherical coordinate sampler values (-1,1) and stores
func (sky *Atmosphere) ComputeRegion(uSampleDomain int, vSampleDomain int, x_corner int, y_corner int, width int, height int) []vector.Vec {
	sizeT := width * height
	tex := make([]vector.Vec, sizeT+1)
	index := 0
	for x := x_corner; x < x_corner+width; x++ {
		u := (2.0*(float64(x)+0.5)/(float64(uSampleDomain-1)) - 1.0)
		for y := y_corner; y < y_corner+height; y++ {
			v := (2.0*(float64(y)+0.5)/(float64(vSampleDomain-1)) - 1.0)
			z2 := u*u + v*v
			phi := math.Atan2(v, u)
			theta := math.Acos(1 - z2)
			sample := vector.Vec{float32(math.Sin(theta) * math.Cos(phi)), float32(math.Sin(theta) * math.Sin(phi)), float32(math.Cos(theta))}
			if index < len(tex) {
				tex[index] = sky.VolumetricScatterRay(sample, vector.Vec{0, 1, 0})
			}
			index++
		}
	}
	return tex
}

//Given a sampling vector and a viewing direction calculate RGB stimulus return
//Based on the Attenuation/Mie Phase Scatter/RayleighScatter Terms
func (sky *Atmosphere) VolumetricScatterRay(sample vector.Vec, view vector.Vec) vector.Vec {

	//Declare volumetric scatter ray vars
	intersects := polar.RaySphereIntersect(sample, sky.Earth.GetPosition(), sky.Earth.GreaterSphere)
	viewRay := vector.Scale(sample, polar.Priority(intersects))
	viewRayMag := viewRay.Mag()

	rgb := vector.Vec{0, 0, 0} //Pixel Output

	betaR := vector.Vec{0.0000088, .0000135, 0.0000331}
	betaM := vector.Vec{0.000021, 0.000021, 0.000021}
	sumR := vector.Vec{0, 0, 0} //rayleigh
	sumM := vector.Vec{0, 0, 0} //mie

	//Compute rayleight/mie coefficients from sample vector and sun direction
	u := vector.Dot(sample, sky.Dir)
	mu := float64(u)
	phaseR := float32(3.0 / (16.0 * PI) * (1.0 + mu*mu))
	g := 0.79
	phaseM := float32(3.0 / (8.0 * PI) * ((1.0 - g*g) * (1.0 + mu*mu)))
	phaseM = phaseM / float32((2+g*g)*math.Pow((1+g*g-2*g*mu), 1.5))
	var opticalDepthR, opticalDepthM float32
	var vmag0, vmag1, vds float32

	//Rayleigh Scatter Computation
	sampleStep := float32(1.0 / RAYLEIGH_SAMPLES)
	for i := 1; i <= RAYLEIGH_SAMPLES && len(intersects) > 0; i++ {

		//Generate Sample Rays along sample view ray path- assume sample ray is normalized
		w := float32(4.0)
		sampleScale := sampler.Ease(float32(i)*sampleStep, w)

		viewSample := vector.Scale(viewRay, sampleScale)
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
		viewSampleOrigin := vector.Add(viewSample, sky.Earth.GetPosition())
		lightIntersects := polar.RaySphereIntersect(vector.Scale(sky.Dir, -1.0), viewSampleOrigin, sky.Earth.GreaterSphere)
		viRay := viewSampleOrigin.Sub(sky.Dir)
		if len(lightIntersects) == 0 {
			fmt.Printf("No Ray Sphere Intersection")
			return vector.Vec{0, 0, 0}
		}

		//Light Path Transmittance + Attenutation
		lightRay := vector.Scale(viRay, polar.Priority(lightIntersects))
		lightRayMag := lightRay.Mag()
		lightPathSampleStep := float32(1.0/LIGHT_PATH_SAMPLES) * lightRayMag

		//Light Transmittance
		ds := float32(0.0)   //differential magnitude
		mag0 := float32(0.0) //for calculating differential
		mag1 := float32(0.0) //for calculating differntial

		var opticalDepthLightM, opticalDepthLightR float32

		//Compute light path to sample position
		for j := 0; j < LIGHT_PATH_SAMPLES; j++ {
			pathScale := sampler.Ease(lightPathSampleStep*float32(j), w)
			lightPath := vector.Scale(lightRay, pathScale)
			mag1 = pathScale * lightRayMag
			ds = mag1 - mag0
			mag0 = mag1
			lightPathSamplePosition := vector.Add(viewSample, lightPath)
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
		attenuation := vector.Vec{float32(math.Exp(float64(-tau[0]))), float32(math.Exp(float64(-tau[1]))), float32(math.Exp(float64(-tau[2])))}
		sumR = sumR.Add(attenuation.Scale(hr))
		sumM = sumM.Add(attenuation.Scale(hm))
	}
	rayliegh := sumR.Mul(betaR).Scale(phaseR)
	mie := sumM.Mul(betaM).Scale(phaseM)
	rgb = rayliegh.Add(mie).Scale(sky.Light.Lx().Flux).Mul(sky.Light.Lx().RGB)

	return rgb
}
