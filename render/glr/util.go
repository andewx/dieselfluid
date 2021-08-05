package glr

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"strings"
)

// initGlfw initializes glfw and returns a Window to use.
func InitGLFW() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(1024, 740, "DSLFLUID", nil, nil)
	checkError(err)
	window.MakeContextCurrent()
	window.SetKeyCallback(ProcessInput)
	window.SetMouseButtonCallback(ProcessMouse)
	window.SetCursorPosCallback(ProcessCursor)

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() {

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func checkError(err error) bool {
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("Error loading asset\n")
		return true
	}
	return false
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength = int32(1000)
		fmt.Printf("Log length %d\n", logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Printf("%s", log)
		return 0, fmt.Errorf("GLSL Shader failed to compile\n: %v", log)
	}
	free()
	return shader, nil
}

func SizeGL(typeID string) uint32 {
	if typeID == "SCALAR" {
		return 1
	}
	if typeID == "VEC3" {
		return 3
	}
	if typeID == "VEC2" {
		return 2
	}
	if typeID == "VEC4" {
		return 4
	}
	return 1
}
