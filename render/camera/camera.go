package camera

import "dslfluid.com/dsl/gltf"
import "dslfluid.com/dsl/math/math32"
import "math"

type Camera struct {
	Node    *gltf.Node
	Fov     float32
	Aspect  float32
	Near    float32
	Far     float32
	ProjMat *math32.Mat4
}

//Initialize camera based on window width and height
func InitCamera(width float32, height float32) Camera {
	myCamera := Camera{nil, 45, width / height, 1, 1000, nil}
	top := float32(math.Tan(float64(myCamera.Fov)/2.0)) * myCamera.Near
	bottom := -top
	right := top * myCamera.Aspect
	left := bottom
	projMat := math32.ProjectionMatrix(left, right, top, bottom, myCamera.Near, myCamera.Far)
	myCamera.ProjMat = &projMat
	return myCamera
}

func (cam *Camera) ResizeCamera(width float32, height float32) {
	cam.Aspect = width / height
}

func (cam *Camera) SetNode(node *gltf.Node) {
	cam.Node = node
}
