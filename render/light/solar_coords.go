package light

import "math"
import "github.com/andewx/dieselfluid/math/mgl"

//Declare New Sun Environment with Standard Merdian time set for NYC
func InitSolarCoordinates(lat float64, day float64) *SolarCoords {
	mySun := SolarCoords{}
	mySun.Lat = lat
	mySun.Jul = day
	mySun.SolarDeclination()
	mySun.SolarTime(12.0)
	mySun.Azimuth()
	mySun.Zenith()
	return &mySun
}

//Solar Declination calculation
func (sun *SolarCoords) SolarDeclination() float64 {
	sun.Decl = 0.4093 * math.Sin(2*PI*(sun.Jul-81/368))
	return sun.Decl
}

//Solar time in PI Radians Time factors longitudes with standard meridians
func (sun *SolarCoords) SolarTime(standard_time float64) float64 {
	p1 := float64(4 * PI * (sun.Jul - 80) / 373)
	p2 := float64(2 * PI * (sun.Jul - 8) / 355)
	sun.Tm = standard_time + 0.17*math.Sin(p1) - 0.129*math.Sin(p2) + 12*(sun.Sm-sun.Long)/PI
	return sun.Tm
}

//Calculates Sun Zenith
func (sun *SolarCoords) Zenith() float64 {
	cosD := math.Cos(sun.Decl)
	sinD := math.Sin(sun.Decl)
	cosL := math.Cos(sun.Lat)
	sinL := math.Sin(sun.Lat)
	cos12 := math.Cos(PI * sun.Tm / 12)
	sin12 := math.Sin(PI * sun.Tm / 12)
	sun.Az = math.Atan((-cosD * sin12) / (cosL*sinD - sinL*cosD*cos12))
	return sun.Az
}

//Calculates Sun Azimuth
func (sun *SolarCoords) Azimuth() float64 {
	cosD := math.Cos(sun.Decl)
	sinD := math.Sin(sun.Decl)
	cosL := math.Cos(sun.Lat)
	sinL := math.Sin(sun.Lat)
	cos12 := math.Cos((PI * sun.Tm / 12))
	sun.Zen = PI/2 - math.Asin(sinL*sinD-cosL*cosD*cos12)
	return sun.Zen
}

//Calculates Suns Direction Vector f(0,p) = cos^2(0)
func (sun *SolarCoords) Vector() mgl.Vec {
	cos2Th := math.Cos(sun.Zen) * math.Cos(sun.Zen)
	sin2Th := math.Sin(sun.Zen) * math.Sin(sun.Zen)
	cos2Az := math.Cos(sun.Az) * math.Cos(sun.Az)
	sun.Dir = mgl.Vec{float32(cos2Th * cos2Az), float32(sin2Th),
		float32(cos2Az * sin2Th)}
	sun.RayDir = sun.Dir.Scale(-1.0)
	return mgl.Norm(sun.Dir)
}

//Updates Sun Calculations per local time
func (sun *SolarCoords) UpdatePosition(local_time float64) mgl.Vec {
	sun.SolarTime(local_time)
	sun.SolarDeclination()
	sun.Zenith()
	sun.Azimuth()
	return sun.Vector()
}

//Calcuates atmospheric depth based azimuth (rads) in KMb - WAIITITITITITITT
//WE WILL USE A RAY CAST METHOD TO A GREATER SPHERE
func (sun *SolarCoords) GetDepth(azimuth float32) float32 {
	pi2 := PI / 2
	return 8.0 + float32(math.Tan(pi2+float64(azimuth)))
}

//Hour to Radians - 0 - 2PI with hours 0-24
func HR2RAD(hour float32) float32 {
	return 2 * PI * hour / 24
}
