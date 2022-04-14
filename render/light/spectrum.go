//dieselfluid/render/light/spectrum.go
//Physically based lighting support structures / classes / and helper methods
// TODO
// Refactor - Defer to Spectrum type interface
// Implement - XYZ Conversion
// Implement - RGB Conversion
// Implement - RGB Spectrum Type
// Move - Sunlight Specific Spectrum Constructors decouple Sunlight production

package light

import "math"
import "fmt"
import "github.com/andewx/dieselfluid/sampler"

const (
	SP_VIOLET = 380.0
	SP_RED    = 625.0
	WATT      = 0
	LUX       = 1
)

//----------------------------------------------------------------------------

type Spectrum interface {
	InitSpectrum(samples int, total_power float32) Spectrum
	Sample(n int) float32
	Set(index int, value float32)
	Wavelength(n int) float32
	Units() int
	IsWatts() bool
	IsBlack() bool
	IsNaN() bool
	Add(Spectrum) Spectrum
	Mul(Spectrum) Spectrum
	Div(Spectrum) Spectrum
	Sub(Spectrum) Spectrum
	Neg() Spectrum
	Sqrt() Spectrum
	Lerp(Spectrum, float32) Spectrum
	Pow(k float32) Spectrum
	Clamp(low float32, high float32) Spectrum
	WriteJSON(filename string, name string, id int)
}

//Coefficient Spectrum
type CoefficientSpectrum struct {
	Samples []float32
	Type    int
	Wv_a    float32 //Wavelength start violet
	Wv_b    float32 //Wavelngth end red
}

//Sampled Spectrum holds an SPD and references to CIE XYZ sampled Spectrums
//cieX , cieY, cieZ should be references only.
type SampledSpectrum struct {
	SPD                 Spectrum
	Samples             int
	CoefficientSpectrum CoefficientSpectrum
	cieX                *CoefficientSpectrum
	cieY                *CoefficientSpectrum
	cieZ                *CoefficientSpectrum
	yint                float32
}

//------------------Coefficient Spectrum--------------------------------------
func (ref *CoefficientSpectrum) InitSpectrum(steps int, total_power float32) Spectrum {

	if steps <= 0 {
		return nil
	}

	ref.Wv_a = SP_VIOLET
	ref.Wv_b = SP_RED
	ref.Samples = make([]float32, steps)
	pow := total_power / float32(steps)

	for i := 0; i < steps; i++ {
		ref.Samples[i] = pow
	}

	return ref
}

//Sets sampler coefficient
func (ref *CoefficientSpectrum) Set(index int, value float32) {
	ref.Samples[index] = value
}

//Gets the Sampler Photon from sampler slice
func (ref *CoefficientSpectrum) Sample(n int) float32 {
	if n < 0 || n > len(ref.Samples) {
		return 0.0
	}
	return ref.Samples[n]
}

//Gets the Sampler Photon from sampler slice
func (ref *CoefficientSpectrum) Wavelength(n int) float32 {
	samples := len(ref.Samples)
	return ref.Wv_a + ((ref.Wv_b-ref.Wv_a)/float32(samples))*float32(n)
}

func (ref *CoefficientSpectrum) Units() int {
	return ref.Type
}

func (ref *CoefficientSpectrum) IsWatts() bool {
	if ref.Type == WATTS {
		return true
	}
	return false
}

func (ref *CoefficientSpectrum) IsBlack() bool {
	for i := 0; i < len(ref.Samples); i++ {
		if ref.Samples[i] != 0.0 {
			return false
		}
	}
	return true
}
func (ref *CoefficientSpectrum) IsNaN() bool {
	for i := 0; i < len(ref.Samples); i++ {
		if math.IsNaN(float64(ref.Samples[i])) {
			return true
		}
	}
	return false
}

//Adds to samples together no error checking out of bounds responsible for caller
func (ref *CoefficientSpectrum) Add(m Spectrum) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = ref.Samples[i] + m.Sample(i)
	}
	return &newSpectrum
}

//Multiplies samples together no error checking out of bounds responsible for caller
func (ref *CoefficientSpectrum) Mul(m Spectrum) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = ref.Samples[i] * m.Sample(i)
	}
	return &newSpectrum
}

//Divides samples together no error checking out of bounds responsible for caller
func (ref *CoefficientSpectrum) Div(m Spectrum) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = ref.Samples[i] / m.Sample(i)
	}
	return &newSpectrum
}

//Subs samples together no error checking out of bounds responsible for caller
func (ref *CoefficientSpectrum) Sub(m Spectrum) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = ref.Samples[i] - m.Sample(i)
	}
	return &newSpectrum
}

func (ref *CoefficientSpectrum) Neg() Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = -ref.Samples[i]
	}
	return &newSpectrum
}

func (ref *CoefficientSpectrum) Sqrt() Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = float32(math.Sqrt(float64(ref.Samples[i])))
	}
	return &newSpectrum
}

func (ref *CoefficientSpectrum) Lerp(m Spectrum, t float32) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = (1.0-t)*ref.Samples[i] + (t * m.Sample(i))
	}
	return &newSpectrum
}

func (ref *CoefficientSpectrum) Pow(k float32) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {
		newSpectrum.Samples[i] = float32(math.Pow(float64(ref.Samples[i]), float64(k)))
	}
	return &newSpectrum
}

func (ref *CoefficientSpectrum) Clamp(low float32, high float32) Spectrum {
	newSpectrum := CoefficientSpectrum{}
	newSpectrum.InitSpectrum(len(ref.Samples), 0.0)
	for i := 0; i < len(ref.Samples); i++ {

		value := ref.Samples[i]
		if value < low {
			value = low
		}

		if value > high {
			value = high
		}
		newSpectrum.Samples[i] = value
	}
	return &newSpectrum
}

/*-----------------------Sampled Spectrum------------------------------------*/
//Returns Sampled spectrum with default constructions - returns nil for invalid
func InitSampledSpectrum(v float32, n float32) *SampledSpectrum {
	if n <= 0 {
		return nil
	}
	smp := SampledSpectrum{}
	smp.CoefficientSpectrum = CoefficientSpectrum{}
	smp.Samples = int(n)
	smp.SPD = smp.CoefficientSpectrum.InitSpectrum(int(n), v*n)
	return &smp
}

//We assume that the sample arrays are sorted - returns error if we can detect
// an unordered sample sapce or the data sets arent correlated
func (ref *SampledSpectrum) FromSampled(wv []float32, v []float32) error {
	if len(wv) != len(v) {
		return fmt.Errorf("Spectrum sample space not equivlant")
	}

	bindIndex := 0

	//!!!!!!!!          REVISE 				!!!!!!!!
	//Calculate sample space averages
	for i := 0; i < ref.Samples; i++ {
		g0 := ref.SPD.Wavelength(i)
		g1 := ref.SPD.Wavelength(i + 1)
		avg, err := sampler.SampleAverage1D(wv, v, 2, g0, g1, &bindIndex)
		if err != nil {
			return err
		}
		ref.SPD.Set(i, avg)
	}
	return nil
}

//Imports SPD into currently allocated Sampled Spectrum Type
func (ref *SampledSpectrum) FromFile(file string) error {
	mySamplerX, err := sampler.ImportSampler(file)
	if err != nil {
		return err
	}
	ref.FromSampled(mySamplerX.Samples_1D.Domain, mySamplerX.Samples_1D.Values)
	return nil
}

//Allocates CIEX infrastructre
func (ref *SampledSpectrum) Add_CIEX(spd_ciex *CoefficientSpectrum) {
	ref.cieX = spd_ciex
}

//Allocates CIEY infrastructre
func (ref *SampledSpectrum) Add_CIEY(spd_ciey *CoefficientSpectrum) {
	ref.cieY = spd_ciey
}

//Allocates CIEZ infrastructure
func (ref *SampledSpectrum) Add_CIEZ(spd_ciez *CoefficientSpectrum) {
	ref.cieZ = spd_ciez
}

//Returns SPD in form of device independent XYZ color coordinates
func (ref *SampledSpectrum) ToXYZ() [3]float32 {
	var xyz [3]float32

	for i := 0; i < ref.Samples; i++ {
		xyz[0] += ref.cieX.Samples[i] * ref.SPD.Sample(i)
		xyz[1] += ref.cieY.Samples[i] * ref.SPD.Sample(i)
		xyz[2] += ref.cieZ.Samples[i] * ref.SPD.Sample(i)
	}

	xyz[0] = xyz[0] / ref.yint
	xyz[1] = xyz[1] / ref.yint
	xyz[2] = xyz[2] / ref.yint
	return xyz
}

func (ref *SampledSpectrum) YY() float32 {
	y := float32(0.0)
	for i := 0; i < ref.Samples; i++ {
		y += ref.cieY.Samples[i] * ref.SPD.Sample(i)
	}
	return y / ref.yint
}

func (ref *SampledSpectrum) ToRGB() [3]float32 {
	xyz := ref.ToXYZ()
	return XYZToRGB(xyz)
}

func XYZToRGB(xyz [3]float32) [3]float32 {
	var rgb [3]float32
	rgb[0] = 3.240479*xyz[0] - 1.537150*xyz[1] - 0.498535*xyz[2]
	rgb[1] = -0.969256*xyz[0] + 1.875991*xyz[1] + 0.041556*xyz[2]
	rgb[2] = 0.055648*xyz[0] - 0.204043*xyz[1] + 1.057311*xyz[2]
	return rgb
}

func RGBToXYZ(rgb [3]float32) [3]float32 {
	var xyz [3]float32
	xyz[0] = 0.412453*rgb[0] + 0.357580*rgb[1] + 0.180423*rgb[2]
	xyz[1] = 0.212671*rgb[0] + 0.715160*rgb[1] + 0.072169*rgb[2]
	xyz[2] = 0.019334*rgb[0] + 0.119193*rgb[1] + 0.950227*rgb[2]
	return xyz
}

//JSON Writing
func (ref *SampledSpectrum) WriteJSON(filename string, name string, id int) {
	ref.CoefficientSpectrum.WriteJSON(filename, name, id)
}

func (ref *CoefficientSpectrum) WriteJSON(filename string, name string, id int) {
	smpJSON := sampler.SamplerJSON{}
	smpJSON.Meta.Name = name
	smpJSON.Meta.SamplerID = id
	nSamples := len(ref.Samples)
	domain := make([]float32, nSamples)
	values := ref.Samples
	domainStep := (ref.Wv_a - ref.Wv_b) / float32(nSamples)
	wva := ref.Wv_a
	for i := 0; i < nSamples; i++ {
		domain[i] = wva + (float32(i) * domainStep)
	}
	smpJSON.Samples_1D.Domain = domain
	smpJSON.Samples_1D.Values = values

	smpJSON.ExportJSON(filename)

}
