package render

import (
	"runtime"
	"testing"
)

func TestGraphics(t *testing.T) {

	runtime.LockOSThread()
	Sys, _ := Init("MaterialSphere.gltf")

	if err := Sys.Init(1024, 720, "andewx/diselfluid/render"); err != nil {
		t.Error(err)
	}

	if err := Sys.CompileLink(); err != nil {
		t.Error(err)
	}

	if err := Sys.Meshes(); err != nil {
		t.Error(err)
	}

	if err := Sys.Run(); err != nil {
		t.Error(err)
	}
}
