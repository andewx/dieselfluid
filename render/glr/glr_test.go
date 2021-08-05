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
	glRender := InitRenderer()
	runtime.LockOSThread()
	glRender.GLHandle = InitGLFW()
	defer glfw.Terminate()
	InitOpenGL()
	glRender.Setup("Minimal2.gltf", 1024, 740)
	glRender.CompileShaders()
	for !glRender.GLHandle.ShouldClose() {
		glRender.Draw()
	}

}
