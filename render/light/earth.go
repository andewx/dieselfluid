package light

import (
	"math"

	"github.com/andewx/dieselfluid/math/mgl"
)

const (
	EARTH_RAD = 6371.0
	HR        = 12.1
	DENSITY0  = 1.225 //KG/M^3
)

type EarthCoords struct {
	Day              float32    //Decimal Earth Day
	Latitude         float32    //Decimal Lat
	Longitude        float32    //Decimal Long
	PolarCoord       mgl.Polar  //Polar Axis offset
	StandardMeridian float32    //Longitude Standard Meridian For Local Time Offsets
	DomainOffset     [2]float32 //Domain offsets for polar sampling of sky depths
	GreaterSphere    mgl.Polar  //Atmospheric Polar Parameters Polar Parameters
}

//Declare New Sun Environment with Standard Merdian time set for NYC
func NewEarth(lat float32, long float32, day float32) *EarthCoords {
	myEarth := EarthCoords{}
	myEarth.Latitude = lat
	myEarth.Longitude = long
	myEarth.Day = day
	myEarth.PolarCoord = mgl.Polar{EARTH_RAD, Day2Rotation(day), mgl.DEG2RAD * lat} //KM
	myEarth.getPolarSamplerDomain()
	myEarth.GreaterSphere = mgl.Polar{EARTH_RAD + HR, 0, 0}
	return &myEarth
}

func Day2Rotation(day float32) float32 {
	return day / 24.0 * 2 * PI
}

func (earth *EarthCoords) GetRadius() float32 {
	return earth.PolarCoord[0]
}

//Takes clamped [U,V] polar coordinates from [-1,1.0] and returns the ray depth
//Returns vector with magnitude to fixed point in sky in valid coordinates
func (earth *EarthCoords) GetSample(uv [2]float32) mgl.Vec {
	uv[0] = mgl.Clamp1f(uv[0], -1.0, 1.0)
	uv[1] = mgl.Clamp1f(uv[1], -1.0, 1.0)
	atmosphereCoords := mgl.Polar{EARTH_RAD + HR, earth.PolarCoord[0] + uv[0]*earth.DomainOffset[0], earth.PolarCoord[1] + uv[1]*earth.DomainOffset[1]}
	rE_Vec, _ := mgl.Sphere2Vec(earth.PolarCoord)
	rSK_Vec, _ := mgl.Sphere2Vec(atmosphereCoords)
	return mgl.Sub(rSK_Vec, rE_Vec)
}

func (earth *EarthCoords) GetSampleDepth(sample mgl.Vec) float32 {
	mVec, _ := mgl.Sphere2Vec(earth.PolarCoord)
	return mgl.SinDot(mVec, sample) * sample.Mag()
}

func (earth *EarthCoords) GetSampleDensity(sample mgl.Vec) float32 {
	return DENSITY0 * float32(math.Exp(float64(-(earth.GetSampleDepth(sample) / HR))))
}

//Gets when rotated earth tangent vector and tangent 2 atomospheric perion vectors are parallel
func (earth *EarthCoords) getPolarSamplerDomain() [2]float32 {

	rE := earth.PolarCoord.Copy() // EARTH
	rE_Vec := mgl.Vec{}
	rE0 := earth.PolarCoord.Copy() //EARTH PRIME
	rE0_Vec := mgl.Vec{}
	rSK := earth.PolarCoord.Copy() //ATMOS
	rSK[0] = EARTH_RAD + 12.1
	samplerDomain := [2]float32{0, 0}
	incr := float32(2 * PI / 720) //0.5 degree intervals
	total_az := float32(0.0)
	total_pl := float32(0.0)
	tan := mgl.Vec{}
	rE2P := mgl.Vec{}
	lastTan := float32(0.0)
	P0, _ := mgl.Sphere2Vec(rSK)
	P0.Norm()

	//Check Azimuth Bounds
	for {

		rE0.AddAzimuth(incr)
		rE0_Vec, _ = mgl.Sphere2Vec(rE0)
		rE_Vec, _ = mgl.Sphere2Vec(rE)
		tan = mgl.Sub(rE0_Vec, rE_Vec).Norm()
		rE2P = mgl.Sub(P0, rE_Vec).Norm()
		compareTan := mgl.Dot(tan, rE2P)

		if compareTan*lastTan < 0.0 {
			break
		} else {
			lastTan = compareTan
			total_az += incr
			rE.AddAzimuth(incr)
		}
		samplerDomain[0] = total_az
	}

	rE = earth.PolarCoord.Copy() // EARTH
	rE_Vec = mgl.Vec{}
	rE0 = earth.PolarCoord.Copy() //EARTH PRIME
	rE0_Vec = mgl.Vec{}

	//Check Polar Bounds
	for {
		rE0.AddPolar(incr)
		rE0_Vec, _ = mgl.Sphere2Vec(rE0)
		rE_Vec, _ = mgl.Sphere2Vec(rE)
		tan = mgl.Sub(rE0_Vec, rE_Vec).Norm()
		rE2P = mgl.Sub(P0, rE_Vec).Norm()
		compareTan := mgl.Dot(tan, rE2P)

		if compareTan*lastTan < 0.0 {
			break
		} else {
			lastTan = compareTan
			total_pl += incr
			rE.AddPolar(incr)
		}

		samplerDomain[1] = total_pl

	}
	earth.DomainOffset = samplerDomain
	return samplerDomain
}