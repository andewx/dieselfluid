package glr

import (
	"fmt"
	"strings"

	"github.com/andewx/dieselfluid/math/matrix"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func checkGlError(op string) {
	error := gl.GetError()
	if error == gl.NO_ERROR {
		return
	}
	fmt.Printf(op+"GL Error %d: ", error)
}

/*
compileShader(source, type) -  from go-gl/cube creates GL shader object for specified shader type
from source string. compiles shader source and returns error log on compilation
failure
*/
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

//Constructs a Matrix from Translation scale rotation quat
func MatrixTRS(t []float32, r []float32, s []float32) []float32 {
	M := matrix.Mat4(1.0)

	//Trans Matrix Affine
	T := matrix.Mat4(1.0)
	T[12] = t[0]
	T[13] = t[1]
	T[14] = t[2]

	S := matrix.Mat4(1.0)
	S[0] = s[0]
	S[5] = s[1]
	S[10] = s[2]

	M = T.MulM(S)

	return M
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
