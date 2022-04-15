package geom

//Global Interfaces List
import (
	Vec "github.com/andewx/dieselfluid/math/mgl"
)

//Particle Collider Interface Defines Particle Velocity Collisions with underlying object
//Note that SPH Particle Boundarys are handled within the internal SPH Calculations, no collision
//Checks are needed.
type Collider interface {
	//Collision takes a (Position, Velocity, Delta Time, Particle Radius)
	//Outputs: Normal, Barycentric Coords, Collision Point, Collision Bool
	Collision(P Vec.Vec, V Vec.Vec, dt float64, r float32) (Vec.Vec, Vec.Vec, Vec.Vec, bool)
}

//GridPoint Returns single point from 3D Index Grid Reference
type GridPoint interface {
	GridPosition(i int, j int, k int) Vec.Vec
}

//Grid Face returns the X, Y, Z point origin positions for a grid centroid
type GridFace interface {
	GridFaces(i int, j int, k int) [3]Vec.Vec
}
