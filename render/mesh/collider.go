//LISCENSE HERE
//render/collider.go

package mesh

import (
	"github.com/andewx/dieselfluid/math/mgl"
)

//Collider Interface Represents Polymorphic types for Mesh Attached
//Colliders which may have an underlying vertice buffer/ parametric curve type surface
//for SDF representation
type Collider interface {

	//Origin() returns the world space transformed mesh origin
	Origin() mgl.Vec

	//UpdateOrigin() Updates the collider world space transformed origin
	UpdateOrigin(mgl.Vec)

	//ImplicitCollide() Parametric Solver Returns boolean for collision inside
	//plane and the float32 distance t for the collision. Typically if the distance
	//is within some epsilon value you would count a collision
	ImplicitCollide(point *mgl.Vec, dir *mgl.Vec) (collide bool, dist float32)

	//ExplicitCollide() non-paramteric tests if the ray crosses some threshold by
	//marching the ray some distance. Returns only whether a collision happened
	//due to sign changing
	ExplicitCollide(point *mgl.Vec, dir *mgl.Vec, dist float32) (collide bool)

	//RenderObjectID Associates a Mesh with its render object
	RenderObjectID() int

	//IsStatic recognizes if the collider is static in world space
	IsStatic() bool
}