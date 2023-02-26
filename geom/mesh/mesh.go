package mesh

import (
	"fmt"

	T "github.com/andewx/dieselfluid/geom/triangle"
	"github.com/andewx/dieselfluid/math/vector"
)

//Triangle Mesh Storage - Implements Collider Interface
type Mesh struct {
	Vertexes []vector.Vec
	Normals  []vector.Vec
}

//Init Mesh Creates a linear triangle list from List of Vertices
func InitMesh(vertices []vector.Vec, origin vector.Vec) Mesh {
	nMesh := Mesh{}
	nMesh.Vertexes = vertices
	nMesh.Normals = make([]vector.Vec, len(vertices)/3)
	index := 0
	//Makes the normals from triangle vertices - We would like all normals to be inward
	//Default Point towards zero vec
	for i := 0; i < len(vertices)-3; i += 3 {
		thisTriangle := T.InitTriangle(vertices[i], vertices[i+1], vertices[i+2])
		n := thisTriangle.Normal()
		v0 := vector.Sub(vertices[i], origin)
		dv0 := vector.Dot(n, v0)
		if dv0 > 0 {
			vector.Scale(n, -1.0)
		}
		nMesh.Normals[i/3] = n
		index++
	}

	return nMesh
}

//Collision checks for collision with underlying mesh
//Returns Normal, Barycentric Coords, Collision Point, Collision Bool
func (g *Mesh) Collision(P vector.Vec, V vector.Vec, dt float64, r float32) (vector.Vec, vector.Vec, vector.Vec, bool) {

	VERTS := len(g.Vertexes)

	for i := 0; i < VERTS; i += 3 {
		normal := g.Normals[i/3]
		triangle := T.InitTriangle(g.Vertexes[i], g.Vertexes[i+1], g.Vertexes[i+2])
		fN, coord, p0, c0 := triangle.BarycentricCollision(P, V, normal, dt, r)

		if c0 {
			return fN, coord, p0, true
		}

	}

	return vector.Vec{}, vector.Vec{}, vector.Vec{}, false
}

//Generates particles in world space given the world space triangles
func (g *Mesh) GenerateBoundaryParticles(density float32) []float32 {

	num_particles := len(g.Vertexes)

	particle_list := make([]float32, num_particles*3)

	for index := 0; index < len(g.Vertexes); index++ {
		p0 := g.Vertexes[index]
		x := index * 3
		if x < len(particle_list)-3 {
			particle_list[x] = p0[0]
			particle_list[x+1] = p0[1]
			particle_list[x+2] = p0[2]
		}
	}
	return particle_list
}

func (g *Mesh) PrintNormals() {
	fmt.Printf("Printing Triangle Normals: Order is {FRONT, BACK, BOTTOM, TOP,LEFT,RIGHT}\n\n")
	VERTS := len(g.Vertexes)
	for i := 0; i < VERTS; i += 3 {

		fmt.Printf("N: [%f, %f, %f]\n", g.Normals[i/3][0], g.Normals[i/3][1], g.Normals[i/3][2])
	}
}

//Triangle Mesh Box with 12 Triangles // 36 Vertexes
func Box(w float32, h float32, d float32, o vector.Vec) Mesh {
	const TRIANGLES = 12
	var Verts = make([]vector.Vec, 12*3)

	x := o[0]
	y := o[1]
	z := o[2]

	p := w / 2
	q := h / 2
	s := d / 2

	//FRONT FACE -Z
	Verts[0] = vector.Vec{x - p, y - q, z + s} //LFB
	Verts[1] = vector.Vec{x - p, y + q, z + s} //LFT
	Verts[2] = vector.Vec{x + p, y + q, z + s} //RFT

	Verts[3] = vector.Vec{x + p, y + q, z + s} //RFT
	Verts[4] = vector.Vec{x + p, y - q, z + s} //RFB
	Verts[5] = vector.Vec{x - p, y - q, z + s} //LFB

	//BACK FACE -Z
	Verts[6] = vector.Vec{x - p, y - q, z - s} //LBB
	Verts[7] = vector.Vec{x - p, y + q, z - s} //LBT
	Verts[8] = vector.Vec{x + p, y - q, z - s} //RBB

	Verts[9] = vector.Vec{x - p, y + q, z - s}  //LBT
	Verts[10] = vector.Vec{x + p, y + q, z - s} //RBT
	Verts[11] = vector.Vec{x + p, y - q, z - s} //RBB

	//BOTTOM FACE -Y
	Verts[12] = vector.Vec{x - p, y - q, z + s} //LFB
	Verts[13] = vector.Vec{x - p, y - q, z - s} //LBB
	Verts[14] = vector.Vec{x + p, y - q, z - s} //RBB

	Verts[15] = vector.Vec{x - p, y - q, z + s} //LFB
	Verts[16] = vector.Vec{x + p, y - q, z - s} //RBB
	Verts[17] = vector.Vec{x + p, y - q, z + s} //RFB

	//Top FACE - Y
	Verts[18] = vector.Vec{x - p, y + q, z + s} //LFT
	Verts[19] = vector.Vec{x - p, y + q, z - s} //LBT
	Verts[20] = vector.Vec{x + p, y + q, z - s} //RBT

	Verts[21] = vector.Vec{x + p, y + q, z - s} //RBT
	Verts[22] = vector.Vec{x + p, y + q, z + s} //RFT
	Verts[23] = vector.Vec{x - p, y + q, z + s} //LFT

	//LEFT FACE - X
	Verts[24] = vector.Vec{x - p, y - q, z + s} //LFB
	Verts[25] = vector.Vec{x - p, y - q, z - s} //LBB
	Verts[26] = vector.Vec{x - p, y + q, z + s} //LTF

	Verts[27] = vector.Vec{x - p, y - q, z - s} //LBB
	Verts[28] = vector.Vec{x - p, y + q, z - s} //LBT
	Verts[29] = vector.Vec{x - p, y + q, z + s} //LFT

	//Right FACE - X
	Verts[30] = vector.Vec{x + p, y + q, z + s} //LFT
	Verts[31] = vector.Vec{x + p, y - q, z + s} //LFB
	Verts[32] = vector.Vec{x + p, y - q, z - s} //LBB

	Verts[33] = vector.Vec{x + p, y + q, z + s} //RFT
	Verts[34] = vector.Vec{x + p, y + q, z - s} //RFB
	Verts[35] = vector.Vec{x + p, y - q, z - s} //RBB

	boxMesh := InitMesh(Verts, o)
	return boxMesh

}
