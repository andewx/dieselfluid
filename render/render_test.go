package render

import (
	"runtime"
	"testing"

	"github.com/andewx/dieselfluid/common"
)

func TestGraphics(t *testing.T) {

	runtime.LockOSThread()
	Sys, _ := Init(common.ProjectRelativePath("data/meshes/materialsphere/MaterialSphere.gltf"))

	if err := Sys.Init(1024, 720, common.ProjectRelativePath("render"), false); err != nil {
		t.Error(err)
	}

	if err := Sys.CompileLink(); err != nil {
		t.Error(err)
	}

	if err := Sys.Meshes(); err != nil {
		t.Error(err)
	}
	message := make(chan string)

	if err := Sys.Run(message); err != nil {
		t.Error(err)
	}
}
