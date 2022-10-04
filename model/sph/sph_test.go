package sph

import "testing"
import "github.com/andewx/dieselfluid/math/vector"

const N = 16

func TestGPUCompile(t *testing.T) {

	Init(1.0, vector.Vec{}, nil, N, true)

}
