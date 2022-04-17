package camera

import (
	"fmt"

	"github.com/andewx/dieselfluid/math/matrix"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/render/transform"
)

const (
	YAW = -90.0
	RAD = 0.0174532925
	PI  = 3.1415729
)

//Camera World Identity View Matrix Performs TMRT^ Normalize
//Operation for matrix construction
type Camera struct {
	Transform  transform.Transform
	ViewMatrix matrix.Mat
	Pos        vector.Vec
	Rot        vector.Vec
	Front      vector.Vec
	Exposure   float32
}

//Camera FPS Euler Decomposition Matrix
func NewCamera(pos vector.Vec) Camera {
	cam := Camera{}
	cam.Transform.Matrix = matrix.Mat4(1.0)
	cam.Transform.Translate(vector.Scale(pos, 1.0))
	cam.Pos = pos
	cam.ViewMatrix = matrix.Mat4(1.0)
	cam.Rot = vector.Vec{0, 0, 0}
	cam.Front = vector.Vec{0, 0, 1, 0}
	return cam
}

//Return Size 3 Vector from the 4 Stored vector
func (cam *Camera) FrontVec() vector.Vec {
	return vector.Vec{cam.Front[0], cam.Front[1], cam.Front[2]}
}

//Updates matrix 4 vec for front from the 3 vec
func (cam *Camera) CopyFront(a vector.Vec) {
	cam.Front[0] = a[0]
	cam.Front[1] = a[1]
	cam.Front[2] = a[2]
}

//Updates the camera view matrix with a translation
//from an approrpiately scaled vector
func (cam *Camera) Translate(vec vector.Vec) {
	cam.Transform.Translate(vec)
}

func (cam *Camera) Log() {
	fmt.Printf("Camera View Mat:\n[")
	for i := 0; i < 16; i++ {

		if i%4 == 0 && i != 0 {
			fmt.Printf("]\n[")
		}
		fmt.Printf(" %.2f ", cam.Transform.Matrix[i])
	}
	fmt.Printf("]\n")
}

//Calculates transpose
func (cam *Camera) Transpose() matrix.Mat {
	return cam.Transform.Matrix.Transpose()
}

//Returns the inverse camera view matrix for scene usage
func (cam *Camera) Update() matrix.Mat {
	return cam.Transform.Matrix.Inv()
}

//Rotates camera around a give unit axis
func (cam *Camera) Rotate(axis vector.Vec, angle float32) {
	cam.Transform.Rotate(axis, angle*RAD)
}

/* RotateFPS treats camera rotation as a rotation of its front facing direction
around a fixed up axis and arbitrary and calculated X (right) axis*/
func (cam *Camera) RotateFPS(rot vector.Vec) {

	//DX/DY
	cam.Rot[0] = rot[0] //YAW
	cam.Rot[1] = rot[1] // PITCH

	//Compute Independent Rotation
	RotTransform := transform.Transform{}
	RotTransform.Matrix = matrix.Mat4(1.0)
	RotTransform.EulerRotate(0.0, cam.Rot[1]*RAD, cam.Rot[0]*RAD)

	//Cross with the 4x4 Transform Matrix (3D Transforms might be easier if they returned a compacted 3 MAT)
	nDir := RotTransform.Matrix.CrossVec(cam.Front)
	front := vector.Norm(nDir) //vector.Add(cam.Front, nDir)
	right := vector.Norm(vector.Cross(vector.Vec{0, 1, 0}, cam.FrontVec()))
	up := vector.Norm(vector.Cross(cam.FrontVec(), right))
	cam.CopyFront(front)
	cam.Transform.Matrix.Set(0, right)
	cam.Transform.Matrix.Set(1, up)
	cam.Transform.Matrix.Set(2, front)

}
