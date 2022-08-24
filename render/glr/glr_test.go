package glr

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestRender(t *testing.T) {

	runtime.LockOSThread()
	fmt.Printf("%s\n", glfw.GetVersionString())

	ogl := Renderer()

	if err := ogl.Setup(640, 480, "Pkg glr - TestRender"); err != nil {
		t.Error(err)
	}
	for !ogl.Window.ShouldClose() {
		ogl.Window.SwapBuffers()
		glfw.PollEvents()
	}
}
