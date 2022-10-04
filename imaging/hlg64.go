package imaging

/*
Sigmoid-Log HDR Compression scheme for HDRR 64 bit (16-bit color channel representation)
with addtional 16 bit encoding for an extensible additional photon parameter. This is less HDR
than SDR with luminance pixel packing however, the relative luminance brightness could be used
to adjust the Sigmoid Gamma and Bias so we get localized tone mapping
*/
import (
	"math"

	"github.com/andewx/dieselfluid/math/vector"
)

//Describes Sigmoid-Dynamic Log packing and unpacking into a 16-bit NRGBA64 space
//source pixels define SDR (0-1) with a gamma profile and (1+) with the logaraithmic packing
//Use this package to import and export back into 32 bit float pixel lumen values
type SigmoidLog64 struct {
	Gamma      float64
	Bias       float64
	MaxLogBase float64
	MxLower    int
	MxUpper    int
	ColorBits  int
}

//Explicit Pixel structure
type Pixel16 [4]uint16

//Unpacked RAW Image Data in 16 bit float format
type RGBAImage64 struct {
	Filename string
	Meta     string
	Width    int32
	Height   int32
	top      int32
	left     int32
	pixels   [][4]uint16
}

//Get image pixel at some value. Unbounded area returns zero pixel
func (p *RGBAImage64) PixelAt(x int, y int) Pixel16 {
	if x*y >= p.Width*p.Height || x*y < 0 || x > p.Width || y > p.Height {
		return Pixel16{}
	}
	index := y*p.Width + x
	return Pixel16{RGBAImage64.pixels[index][0], RGBAImage64.pixels[index][1], RGBAImage64.pixels[index][2], RGBAImage64.pixels[index][3]}
}

func NewRGBAImage64(f string, w int32, h int32) {
	mRet := RGBAImage64{f, "hlg64", w, h, 0, 0, make([][4]uint16, x*y)}
	return mRet
}

//Reference Hybrid log gamma curve follows 12 stops of Dynamic Range
func NewSigmoidLog64(colorBits int, lumBits int) *SigmoidLog64 {
	if colorBits+lumBits != 16 {
		return nil
	}
	mxL := int(math.Pow(2, float64(colorBits)))
	mxU := int(math.Pow(2, float64(16-lumBits)))
	return &SigmoidLog64(3.3, 0.0, 0.0, mxL, mxU, colorBits)
}

//Sigmoid Tone Mapping applied to visible pixels with a linear super-luminance output
func (p *SigmoidLog64) OOTF(x float64) float64 {
	y := 0
	if x < 0.0 {
		x = 0.0
	}

	if x <= 1 {
		op := -(x*math.E - (math.E+p.Bias)/2) * p.Gamma
		y = 1 / (1 + math.Exp(op))
	} else {
		y = x
	}
	return y
}

//Combine signal channels from compressed /tone mapped signal into bit channel
func (p *SigmoidLog64) OETF(e float64) uint16 {
	mapped := uint16(0)
	pow := uint16(0)
	if e >= 0 && e <= 1.0 {

		max_mid := uint16(p.MaxLower - 1)
		mapped := uint16(e * max_mid)
	} else {
		pow_map := uint16(math.Logb(p.MaxLogBase, e))
		if pow_map > p.MxUpper-1 {
			pow = p.MxUpper - 1
		} else {
			pow = pow_map
		}
	}
	//Shift and Combine into return int
	nb := pow_map << p.Pxbits
	nb = nb & mapped
	return nb
}

//Maps the stored function back into its original form
func (p *SigmoidLog64) EOTF(element uint16) float64 {
	retVal := 0.0
	lower_int := element
	mask := uint16(0xFFFF)>>16 - p.ColorBits
	lower_int = element & mask
	mask = !mask
	upper_int = element & mask
	if upper_int > 0 {
		retVal = math.Pow(p.MaxLogBase, float64(upper_int))
	} else {
		retVal = float64(element) / float64(p.MxLower)
	}
}

//Packs a signal vector into the 16 bit space
func (p *SigmoidLog64) Pack(a vector.Vec) Pixel16 {
	px := Pixel16{}
	for i := 0; i < 3; i++ {
		y := p.OOTF(float64(a[i]))
		n := p.OETF(y)
		px[i] = n
	}
	px[3] = uint16(a[3])

	return px
}

//Unpacks 64bit encoding bytes into continuous signal
func (p *SigmoidLog64) Unpack(a Pixel16) vector.Vec {
	retVal := vector.Vec{0, 0, 0, 0}
	for i := 0; i < 3; i++ {
		retVal[i] = float32(p.EOTF(a[i]))
	}

	retVal[3] = float32(a[3])
	return retVal
}
