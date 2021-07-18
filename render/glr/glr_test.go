package glr

import (
	"fmt"
	"testing"
)

//Test Render Loop---
func Test_Render(t *testing.T) {
	fmt.Printf("Starting...\n")
	glRender := InitRenderer()
	fmt.Printf("Initiating Render...\nLoading GLTF")
	glRender.Setup("Minimal.gltf", 640, 400)
	fmt.Printf("Initiate GL\n")
	glRender.Init()

}
