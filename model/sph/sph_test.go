package sph

import "testing"
import "github.com/andewx/dieselfluid/math/vector"
import "runtime"
import "github.com/go-gl/glfw/v3.3/glfw"

const N = 16

func TestGPUCompile(t *testing.T) {
	runtime.LockOSThread()
	sph := InitSPH(1.0, vector.Vec{}, nil, N)
	w := sph.FluidShader(false)

	if w != nil {
		for !w.ShouldClose() {
			w.SwapBuffers()
			glfw.PollEvents()
		}
	} else {
		t.Errorf("Requested window not passed early return")
	}

}
