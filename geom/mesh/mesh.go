package mesh

import (
	T "dslfluid.com/dsl/geom/triangle"
	Vec "dslfluid.com/dsl/math/math32"
	"fmt"
)

//Triangle Mesh Storage - Implements Collider Interface
type Mesh struct {
	Vertexes []Vec.Vec
	Normals  []Vec.Vec
}

//Init Mesh Creates a linear triangle list from List of Vertices
func InitMesh(vertices []Vec.Vec, origin Vec.Vec) Mesh {
	nMesh := Mesh{}
	nMesh.Vertexes = vertices
	nMesh.Normals = make([]Vec.Vec, len(vertices)/3)
	index := 0
	//Makes the normals from triangle vertices - We would like all normals to be inward
	//Default Point towards zero vec
	for i := 0; i < len(vertices); i += 3 {
		thisTriangle := T.InitTriangle(vertices[i], vertices[i+1], vertices[i+2])
		n := thisTriangle.Normal()
		v0 := Vec.Sub(vertices[i], origin)
		dv0 := Vec.Dot(n, v0)
		if dv0 > 0 {
			Vec.Scl(n, -1.0)
		}
		nMesh.Normals[i/3] = n
		index++
	}

	return nMesh
}

//Collision checks for collision with underlying mesh
//Returns Normal, Barycentric Coords, Collision Point, Collision Bool
func (g *Mesh) Collision(P Vec.Vec, V Vec.Vec, dt float64, r float32) (Vec.Vec, Vec.Vec, Vec.Vec, bool) {

	VERTS := len(g.Vertexes)

	for i := 0; i < VERTS; i += 3 {
		normal := g.Normals[i/3]
		triangle := T.InitTriangle(g.Vertexes[i], g.Vertexes[i+1], g.Vertexes[i+2])
		fN, coord, p0, c0 := triangle.BarycentricCollision(P, V, normal, dt, r)

		if c0 {
			return fN, coord, p0, true
		}

	}

	return Vec.Vec{}, Vec.Vec{}, Vec.Vec{}, false
}

func (g *Mesh) PrintNormals() {
	fmt.Printf("Printing Triangle Normals: Order is {FRONT, BACK, BOTTOM, TOP,LEFT,RIGHT}\n\n")
	VERTS := len(g.Vertexes)
	for i := 0; i < VERTS; i += 3 {

		fmt.Printf("N: [%f, %f, %f]\n", g.Normals[i/3][0], g.Normals[i/3][1], g.Normals[i/3][2])
	}
}

//Triangle Mesh Box with 12 Triangles // 36 Vertexes
func Box(w float32, h float32, d float32, o Vec.Vec) Mesh {
	const TRIANGLES = 12
	var Verts = make([]Vec.Vec, 12*3)

	x := o[0]
	y := o[1]
	z := o[2]

	p := w / 2
	q := h / 2
	s := d / 2

	//FRONT FACE -Z
	Verts[0] = Vec.Vec{x - p, y - q, z + s} //LFB
	Verts[1] = Vec.Vec{x - p, y + q, z + s} //LFT
	Verts[2] = Vec.Vec{x + p, y + q, z + s} //RFT

	Verts[3] = Vec.Vec{x + p, y + q, z + s} //RFT
	Verts[4] = Vec.Vec{x + p, y - q, z + s} //RFB
	Verts[5] = Vec.Vec{x - p, y - q, z + s} //LFB

	//BACK FACE -Z
	Verts[6] = Vec.Vec{x - p, y - q, z - s} //LBB
	Verts[7] = Vec.Vec{x - p, y + q, z - s} //LBT
	Verts[8] = Vec.Vec{x + p, y - q, z - s} //RBB

	Verts[9] = Vec.Vec{x - p, y + q, z - s}  //LBT
	Verts[10] = Vec.Vec{x + p, y + q, z - s} //RBT
	Verts[11] = Vec.Vec{x + p, y - q, z - s} //RBB

	//BOTTOM FACE -Y
	Verts[12] = Vec.Vec{x - p, y - q, z + s} //LFB
	Verts[13] = Vec.Vec{x - p, y - q, z - s} //LBB
	Verts[14] = Vec.Vec{x + p, y - q, z - s} //RBB

	Verts[15] = Vec.Vec{x - p, y - q, z + s} //LFB
	Verts[16] = Vec.Vec{x + p, y - q, z - s} //RBB
	Verts[17] = Vec.Vec{x + p, y - q, z + s} //RFB

	//Top FACE - Y
	Verts[18] = Vec.Vec{x - p, y + q, z + s} //LFT
	Verts[19] = Vec.Vec{x - p, y + q, z - s} //LBT
	Verts[20] = Vec.Vec{x + p, y + q, z - s} //RBT

	Verts[21] = Vec.Vec{x + p, y + q, z - s} //RBT
	Verts[22] = Vec.Vec{x + p, y + q, z + s} //RFT
	Verts[23] = Vec.Vec{x - p, y + q, z + s} //LFT

	//LEFT FACE - X
	Verts[24] = Vec.Vec{x - p, y - q, z + s} //LFB
	Verts[25] = Vec.Vec{x - p, y - q, z - s} //LBB
	Verts[26] = Vec.Vec{x - p, y + q, z + s} //LTF

	Verts[27] = Vec.Vec{x - p, y - q, z - s} //LBB
	Verts[28] = Vec.Vec{x - p, y + q, z - s} //LBT
	Verts[29] = Vec.Vec{x - p, y + q, z + s} //LFT

	//Right FACE - X
	Verts[30] = Vec.Vec{x + p, y + q, z + s} //LFT
	Verts[31] = Vec.Vec{x + p, y - q, z + s} //LFB
	Verts[32] = Vec.Vec{x + p, y - q, z - s} //LBB

	Verts[33] = Vec.Vec{x + p, y + q, z + s} //RFT
	Verts[34] = Vec.Vec{x + p, y + q, z - s} //RFB
	Verts[35] = Vec.Vec{x + p, y - q, z - s} //RBB

	boxMesh := InitMesh(Verts, o)
	return boxMesh

}
