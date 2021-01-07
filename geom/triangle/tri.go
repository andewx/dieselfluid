package triangle

import (
	Vec "dslfluid.com/dsl/math/math64"
)

//Define Geometry Type Structures
type Triangle struct {
	Verts [3]*Vec.Vec
}

//Triangle
func InitTriangle(a Vec.Vec, b Vec.Vec, c Vec.Vec) Triangle {
	V := Triangle{}
	V.Verts[0] = &a
	V.Verts[1] = &b
	V.Verts[2] = &c

	return V
}

func (tri *Triangle) Normal() Vec.Vec {
	N := Vec.Cross(Vec.Sub(*tri.Verts[1], *tri.Verts[0]), Vec.Sub(*tri.Verts[2], *tri.Verts[0]))
	return Vec.Norm(N)
}

//Barycentric Collission Test returns float64 distance, bool collision detected
func (t *Triangle) Collision(P *Vec.Vec) (Vec.Vec, bool) {

	//Measures if  a point projected into the triangles plane gives a barycentric coord
	coord, isBarycentric := t.Barycentric(P)
	return coord, isBarycentric

}

//Barycentric Focused Collision Test (returns Normal, Coords, CollisionPoint, Collision Bool)
func (t *Triangle) BarycentricCollision(P Vec.Vec, V Vec.Vec, n Vec.Vec, dt float64, r float64) (Vec.Vec, Vec.Vec, Vec.Vec, bool) {

	if Vec.Mag(V) == 0 {
		return n, Vec.Vec{}, Vec.Vec{}, false
	}

	t0 := *t.Verts[0]
	//Take Care of Normal
	v0 := Vec.Sub(t0, P)

	nDotRay := Vec.Dot(n, V)

	//Point Plane Distance Projected onto the Velocity
	if nDotRay == 0 {
		nDotRay = 0.0001
	}

	d := Vec.Dot(v0, n)
	k := (d) / (nDotRay)            //Project distance to velocity Vector
	p0 := Vec.Add(P, Vec.Scl(V, k)) //Projection to the plane
	dist := Vec.Mag(Vec.Sub(P, p0))

	//Check the P2 is crossed current plane
	//- TODO SEE IF PROJECTED VELOCITY IS A BARYCENTRIC COLLISION
	//	p1 := Vec.Add(P, Vec.Scl(V, dt)) //Actual projected vector
	//	p10 := Vec.Sub(p1, p0)
	//dotp10 := Vec.Dot(n, p10)
	//	dist2 := Vec.Mag(p10)

	//Point Crossed the plane in a time step. We don't care about the actual collision point
	//This needs to be scaled with velocity or time step needs to be decreased (dotp10 > 0 && dv0 < 0) ||
	if dist <= r {
		coord, collision := t.Barycentric(&P)
		P = Vec.Add(P, Vec.Scl(V, -1.0*dt))
		return n, coord, P, collision
	} else {
		return n, Vec.Vec{}, Vec.Vec{}, false
	}

}

//Project XY, XZ, YZ - Plane must be
func (t *Triangle) Barycentric(p *Vec.Vec) (Vec.Vec, bool) {

	v0 := Vec.Sub(*t.Verts[1], *t.Verts[0])
	v1 := Vec.Sub(*t.Verts[2], *t.Verts[0])
	v2 := Vec.Sub(*p, *t.Verts[0])
	d00 := Vec.Dot(v0, v0)
	d01 := Vec.Dot(v0, v1)
	d11 := Vec.Dot(v1, v1)
	d20 := Vec.Dot(v2, v0)
	d21 := Vec.Dot(v2, v1)
	denom := d00*d11 - d01*d01
	u := (d11*d20 - d01*d21) / denom
	v := (d00*d21 - d01*d20) / denom
	w := 1.0 - v - u
	coord := Vec.Vec{u, v, w}
	collision := false
	if u <= 1.0 && v <= 1.0 && w <= 1.0 && (u+v+w) <= 1.0 && u >= 0 && v >= 0 && w >= 0 {
		collision = true
	}

	return coord, collision

}

//Planar Projection Transform of a triangle onto a Normal Vector
func (t *Triangle) Project(N Vec.Vec) Triangle {
	nTri := Triangle{}
	a := Vec.ProjPlane(*t.Verts[0], N)
	b := Vec.ProjPlane(*t.Verts[1], N)
	c := Vec.ProjPlane(*t.Verts[2], N)
	nTri.Verts[0] = &a
	nTri.Verts[1] = &b
	nTri.Verts[2] = &c
	return nTri
}
