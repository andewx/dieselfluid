package imaging

import (
	"github.com/andewx/dieselfluid/math/vector"
)

//HDR interface packs and unpacks pixel values into 64 bit format assumed RGBX space with an
//X function used for the extensible RGBX X parameter
type HDR interface {
	Pack(a vector.Vec) []int16
	Unpack(b []int16) []float32
	XPack(x float32) int16
	XUnpack(x int16) float32
	MaxInt() int32
	NormalizationPoint() int32
}
