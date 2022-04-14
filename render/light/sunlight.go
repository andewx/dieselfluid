package light

import "math"

//Sunlight Distribution Variables
const (
	SPD_B = 9.6
	SPD_C = 3.1
	SPD_K = 1.68
	SPD_G = 0.45
)

//------------INLINE FUNCTIONS FOR RETURNING SUN COEFFICIENT SPECTRUM---------
func SunlightSpectrum(nm float64) float64 {
	x := nm / 100
	return SPD_B * SPD_B * SPD_C * SPD_K * (math.Pow(SPD_G*x, SPD_C-1)) / math.Pow(SPD_B+math.Pow(SPD_G*x, SPD_C-1), SPD_K+1)
}

//SPD Burr Sunlight Distribution Model Attenuated by k_atten factor
func SunlightSpectrumAttenuate(nm float64, k_atten float64) float64 {
	x := nm / 100
	return SPD_B * SPD_B * SPD_C * SPD_K * k_atten * (math.Pow(SPD_G*x, SPD_C-1)) / math.Pow(SPD_B+math.Pow(SPD_G*x, SPD_C-1), SPD_K*k_atten+1)
}

//------------------Coefficient Spectrum--------------------------------
func InitSunlight(steps int) Spectrum {
	cSpectrum := CoefficientSpectrum{}
	//Check error
	if steps <= 0 {
		return nil
	}
	cSpectrum.Samples = make([]float32, steps)
	sp_width := SP_RED - SP_VIOLET
	sp_step := sp_width / float64(steps)

	for i := 0; i < steps; i++ {
		nm := SP_VIOLET + (float64(i) * sp_step)
		cSpectrum.Samples[i] = float32(SunlightSpectrumAttenuate(nm, 1.086))
	}

	return &cSpectrum
}
