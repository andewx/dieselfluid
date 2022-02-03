package light

import "dslfluid.com/dsl/math/mgl"
import "dslfluid.com/dsl/render/transform"
import "math"

/*
Approximation model for sun elevation/azimuth given latitude and localized decimal
time approximation
*/
const (
	SEASONAL_PERIOD  = 181.8
	AXIAL_TILT       = 23.5 //Degrees
	CLEAR_LUX        = 105000.0
	ANGULAR_VELOCITY = 15.0 //per hour
)

//World/Sun light lux with sky position transform (110,000 Afternoon Bright Day Lux)
type Environment struct {
	Light       Directional
	Day         float64
	Time        float32
	Declination float32
}

func DEG2RAD(x float64) float64 {
	return math.Pi / 180 * x
}

//Calculates Sun Position based on Latitudinal position at mid-day for that season
//Involves a few transforms since we are dealing with solar/earth coordinates
//rather than local fram coordinates only
func NewEnvironmentLight(lat float32, day float32) *Environment {
	myLight := Directional{mgl.Vec{1.0, 1.0, 0.95}, Lux{mgl.Vec{1,1,1},float32(CLEAR_LUX)}
	//Interval in time of hours from local solar noon to sunset
	decl := Declination(day)
	w0 := math.Acos(-math.Tan(DEG2RAD(float64(lat)))*math.Tan(decl)) / DEG2RAD(15.0)

	myLightEnv := Environment{myLight, day, 0.0, w0, decl}
	return &myLightEnv
}

//Declination angle in radians of sun position declination is equal to 0 at the
//spring and fall equinoxes. day is the number of days since the beginning of
//the year
func Declination(day float32) {
	return DEG2RAD(-23.5 * math.Cos(float64(360/(365*(day+10)))))
}
