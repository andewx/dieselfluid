package triangle

import (
	"github.com/andewx/dieselfluid/math/vector"
)

//Define Geometry Type Structures
type Triangle struct {
	Verts [3]*vector.Vec
}

//Triangle
func InitTriangle(a vector.Vec, b vector.Vec, c vector.Vec) Triangle {
	V := Triangle{}
	V.Verts[0] = &a
	V.Verts[1] = &b
	V.Verts[2] = &c

	return V
}

func (tri *Triangle) Normal() vector.Vec {
	N := vector.Cross(vector.Sub(*tri.Verts[1], *tri.Verts[0]), vector.Sub(*tri.Verts[2], *tri.Verts[0]))
	return vector.Norm(N)
}

//Barycentric Collission Test returns float64 distance, bool collision detected
func (t *Triangle) Collision(P vector.Vec) (vector.Vec, bool) {

	//Measures if  a point projected into the triangles plane gives a barycentric coord
	coord, isBarycentric := t.Barycentric(P)
	return coord, isBarycentric

}

//Barycentric Focused Collision Test (returns Normal, Coords, CollisionPoint, Collision Bool)
func (t *Triangle) BarycentricCollision(P vector.Vec, V vector.Vec, n vector.Vec, dt float64, r float32) (vector.Vec, vector.Vec, vector.Vec, bool) {

	if vector.Mag(V) == 0 {
		return n, vector.Vec{}, vector.Vec{}, false
	}

	t0 := *t.Verts[0]
	//Take Care of Normal
	v0 := vector.Sub(t0, P)

	nDotRay := vector.Dot(n, V)

	//Point Plane Distance Projected onto the Velocity
	if nDotRay == 0 {
		nDotRay = 0.0001
	}

	d := vector.Dot(v0, n)
	k := (d) / (nDotRay)                    //Project distance to velocity Vector
	p0 := vector.Add(P, vector.Scale(V, k)) //Projection to the plane
	dist := vector.Mag(vector.Sub(P, p0))

	//Check the P2 is crossed current plane
	//- TODO SEE IF PROJECTED VELOCITY IS A BARYCENTRIC COLLISION
	//	p1 := vector.Add(P, vector.Scale(V, dt)) //Actual projected vector
	//	p10 := vector.Sub(p1, p0)
	//dotp10 := vector.Dot(n, p10)
	//	dist2 := vector.Mag(p10)

	//Point Crossed the plane in a time step. We don't care about the actual collision point
	//This needs to be scaled with velocity or time step needs to be decreased (dotp10 > 0 && dv0 < 0) ||
	if dist <= r {
		coord, collision := t.Barycentric(P)
		P = vector.Add(P, vector.Scale(V, -1.0*float32(dt)))
		return n, coord, P, collision
	} else {
		return n, vector.Vec{}, vector.Vec{}, false
	}

}

//Project XY, XZ, YZ - Plane must be
func (t *Triangle) Barycentric(p vector.Vec) (vector.Vec, bool) {

	v0 := vector.Sub(*t.Verts[1], *t.Verts[0])
	v1 := vector.Sub(*t.Verts[2], *t.Verts[0])
	v2 := vector.Sub(p, *t.Verts[0])
	d00 := vector.Dot(v0, v0)
	d01 := vector.Dot(v0, v1)
	d11 := vector.Dot(v1, v1)
	d20 := vector.Dot(v2, v0)
	d21 := vector.Dot(v2, v1)
	denom := d00*d11 - d01*d01
	u := (d11*d20 - d01*d21) / denom
	v := (d00*d21 - d01*d20) / denom
	w := 1.0 - v - u
	coord := vector.Vec{u, v, w}
	collision := false
	if u <= 1.0 && v <= 1.0 && w <= 1.0 && (u+v+w) <= 1.0 && u >= 0 && v >= 0 && w >= 0 {
		collision = true
	}

	return coord, collision

}

//Planar Projection Transform of a triangle onto a Normal Vector
func (t *Triangle) Project(N vector.Vec) Triangle {
	nTri := Triangle{}
	a := vector.ProjPlane(*t.Verts[0], N)
	b := vector.ProjPlane(*t.Verts[1], N)
	c := vector.ProjPlane(*t.Verts[2], N)
	nTri.Verts[0] = &a
	nTri.Verts[1] = &b
	nTri.Verts[2] = &c
	return nTri
}
