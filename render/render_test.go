package render

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/andewx/dieselfluid/common"
)

func TestGraphics(t *testing.T) {

	runtime.LockOSThread()
	Render, err := Init(common.ProjectRelativePath("data/meshes/materialbowl/ceramic_bowl.gltf"))

	if err != nil {
		fmt.Printf("Unable to initiate scene graph\n%v\n", err)
		t.Errorf("Failed")
		return
	}

	if err := Render.Init(1024, 720, common.ProjectRelativePath("render"), false); err != nil {
		t.Error(err)
	}

	if err := Render.CompileLink(); err != nil {
		t.Error(err)
	}

	if err := Render.Meshes(); err != nil {
		t.Error(err)
	}
	message := make(chan string)

	if err := Render.Run(message, nil); err != nil {
		t.Error(err)
	}
}
