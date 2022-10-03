package geom

import "testing"

//import "fmt"
import "github.com/andewx/dieselfluid/math/vector"

func TestOctal(t *testing.T) {
	oct := InitOctalTree(100.0, 100.0, 100.0)
	pointA := vector.Vec{2.5, 1.56, 2.61}
	pointB := vector.Vec{3.61, 5.10, 2.43}
	/*
		pointC := vector.Vec{2.55, 3.55, 4.62}
		group := make([]vector.Vec, 3)
		group[0] = pointA
		group[1] = pointB
		group[2] = pointC
	*/
	encA := oct.EncodePoint(pointA, 4)
	encB := oct.EncodePoint(pointB, 4)
	//	enc_group := oct.EncodePointGroup(group)

	abSim := oct.DepthSimilarity(encA, encB)
	//	enc_group_depth := len(enc_group) / 3
	//	fmt.Printf("Group Depth %d\n", enc_group_depth)
	//	fmt.Printf("Point A, Point B Octal Depth Similarity:%d\nPointA Encoding[", abSim)
	/*
		for i := 0; i < len(encA); i++ {
			d := uint32(encA[i])
			if i%3 == 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%d", d)
		}
	*/
	//	fmt.Printf("]\nPoint B Encoding[")
	/*
		for i := 0; i < len(encB); i++ {
			d := uint32(encB[i])
			if i%3 == 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%d", d)
		}
		fmt.Printf("]\n")
	*/
	if abSim != 5 {
		t.Errorf("Points A & B Depth Similarity not 5\n")
	}
}

func BenchmarkDepth(b *testing.B) {
	oct := InitOctalTree(100.0, 100.0, 100.0)
	pointA := vector.Vec{2.5, 1.56, 2.61}
	pointB := vector.Vec{3.61, 5.10, 2.43}
	encA := oct.EncodePoint(pointA, 4)
	encB := oct.EncodePoint(pointB, 4)

	for i := 0; i < b.N; i++ {
		oct.DepthSimilarity(encA, encB)
	}

}

func BenchmarkEncodeTriangle(b *testing.B) {
	oct := InitOctalTree(100.0, 100.0, 100.0)
	pointA := vector.Vec{2.5, 1.56, 2.61}
	pointB := vector.Vec{3.61, 5.10, 2.43}
	pointC := vector.Vec{2.55, 3.55, 4.62}
	group := make([]vector.Vec, 3)
	group[0] = pointA
	group[1] = pointB
	group[2] = pointC

	for i := 0; i < b.N; i++ {
		oct.EncodePointGroup(group, 4)
	}
}
