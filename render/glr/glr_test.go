package glr

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"testing"
)

//Test Render Loop---
func Test_Render(t *testing.T) {
	fmt.Printf("Starting...\n")
	runtime.LockOSThread()
	defer glfw.Terminate()
	glRender := InitRenderer()
	glRender.GLHandle = InitGLFW()
	InitOpenGL()

	if err := glRender.Setup("PbrSphere.gltf", 1024, 740); err != nil {
		t.Errorf("Exiting...\n")
		return
	}

	glRender.CompileShaders()

	for !glRender.GLHandle.ShouldClose() {
		glRender.Draw()
	}

}
