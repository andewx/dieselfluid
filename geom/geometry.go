package geom

//Global Interfaces List
import (
	Vec "dslfluid.com/dsl/math/math64"
)

//Particle Collider Interface Defines Particle Velocity Collisions with underlying object
//Note that SPH Particle Boundarys are handled within the internal SPH Calculations, no collision
//Checks are needed.
type Collider interface {
	//Collision takes a (Position, Velocity, Delta Time, Particle Radius)
	//Outputs: Normal, Barycentric Coords, Collision Point, Collision Bool
	Collision(P math32.Vec, V Vec.Vec, dt float64, r float32) (Vec.Vec, Vec.Vec, Vec.Vec, bool)
}
