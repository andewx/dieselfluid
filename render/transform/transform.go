package transform

import "dslfluid.com/dsl/math/mgl"
import "math"

//Holds Transform Matrix
type Transform struct {
	Matrix mgl.Mat
}

//Augments Transform Matrix with Translation
func (T *Transform) Translate(vec mgl.Vec) mgl.Mat {
	TransMat := mgl.Mat4(1.0)
	TransMat[12] = vec[0]
	TransMat[13] = vec[1]
	TransMat[14] = vec[2]
	T.Matrix = mgl.MulM(T.Matrix, TransMat)
	return T.Matrix
}

//Rotate applies matrix rotation of the angle theta which should be in radians,
//The rotation axis is given by the Vector axis (l, m, n)
func (T *Transform) Rotate(axis mgl.Vec, angle float32) mgl.Mat {
	rot := mgl.Mat4(1.0)
	axis = mgl.Norm(axis)
	l := axis[0]
	m := axis[1]
	n := axis[2]
	mCos := float32(1.0 - math.Cos(float64(angle)))
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	rot[0] = l*l*mCos + cos
	rot[1] = m*l*mCos - (n * sin)
	rot[2] = n*l*mCos + (m * sin)
	rot[4] = l*m*mCos + (n * sin)
	rot[5] = m*m*mCos + cos
	rot[6] = n*m*mCos - (l * sin)
	rot[8] = l*n*mCos - m*sin
	rot[9] = m*n*mCos + (l * sin)
	rot[10] = n*n*mCos + cos
	T.Matrix = mgl.MulM(T.Matrix, rot)
	return T.Matrix
}

//EulerRotate compactifies the intended yaw / pitch / roll values into a rotation matrix form
func (T *Transform) EulerRotate(yaw float32, pitch float32, roll float32) {
	y := float64(yaw)
	p := float64(pitch)
	r := float64(roll)

	cosZ := float32(math.Cos(y))
	sinZ := float32(math.Sin(y))
	cosY := float32(math.Cos(p))
	sinY := float32(math.Sin(p))
	cosX := float32(math.Cos(r))
	sinX := float32(math.Sin(r))

	euler := mgl.Mat{
		cosZ * cosY, cosZ*sinY*sinX - sinZ*cosX, cosZ*sinY*cosX + sinZ*sinX, 0,
		sinZ * cosY, sinZ*sinY*sinX + cosZ*cosX, sinZ*sinY*cosX - cosZ*sinX, 0,
		-sinY, cosY * sinX, cosY * cosX, 0, 0, 0, 0, 1,
	}

	T.Matrix = euler

}
