package geom

import "github.com/andewx/dieselfluid/math/vector"

//Octal Tree Structure/Object acts as a classifier for points and point sets
//Whereby point representation in space is assigned a containing volume code by
//the encoded pseudo-octet. We don't actually compress here in this libary to the
//bit wise encoding but with depths usually small depths this shouldn't be an issue
type OctalTree struct {
	Bounds   [3]float32 //Width, Height, Depth, (X,Y,Z)
	Origin   [3]float32 //Initial Original of Octal
	MaxDepth int        //Octal Tree Depth
	Map      map[string][]int
}

func InitOctalTree(w float32, h float32, d float32) *OctalTree {
	return &OctalTree{[3]float32{w, h, d}, [3]float32{0, 0, 0}, 6, make(map[string][]int, 10)}
}

//Octet Lazy Encoding - No Bit Manipulation for Space Concerns
//Use uint8 encoding for space partitioning sequence and depth is determined
//By the sequence length
func (p OctalTree) GetCentroid(encoding []uint8) [3]float32 {
	current_centroid := p.Origin
	current_bounds := p.Bounds
	for i := 0; i < len(encoding); i += 3 {
		xb := (float32(encoding[i]) - 0.5)
		yb := (float32(encoding[i+1]) - 0.5)
		zb := (float32(encoding[i+2]) - 0.5)
		current_centroid[0] = current_centroid[0] + (xb * current_bounds[0])
		current_centroid[1] = current_centroid[1] + (yb * current_bounds[1])
		current_centroid[2] = current_centroid[2] + (zb * current_bounds[2])
		current_bounds[0] = current_bounds[0] / 2
		current_bounds[1] = current_bounds[1] / 2
		current_bounds[2] = current_bounds[2] / 2
	}
	return current_centroid
}

func (p OctalTree) GetParent(encoding []uint8) string {
	if len(encoding) >= 3 {
		return string(encoding[0 : len(encoding)-3])
	} else {
		return string(encoding[:])
	}
}

//TODO Add Sibling Nodes
func (p OctalTree) GetNeighbors(encoding []uint8) []int {
	encode_str := string(encoding[:])
	enc_par := p.GetParent(encoding)
	myNeighbors := append(p.Map[encode_str][:], p.Map[enc_par][:]...)
	return myNeighbors
}

//inserts
func (p OctalTree) InsertPoint(point vector.Vec, unique_id int, desired_depth int) {

	encoding := p.EncodePoint(point, desired_depth)
	str := string(encoding[:])
	p.Map[str] = append(p.Map[str][:], unique_id)
}

func (p OctalTree) RemovePoint(encoding string, unique_id int) {
	index := 0
	found := false
	for !found && index < len(p.Map[encoding]) {
		if p.Map[encoding][index] == unique_id {
			found = true
		} else {
			index++
		}
	}
	if index < len(p.Map[encoding])-1 {
		p.Map[encoding] = append(p.Map[encoding][0:index], p.Map[encoding][index+1:]...)
	} else {
		p.Map[encoding] = append(p.Map[encoding][0:index])
	}
}

/*Encodes a single point into a octal encoding with max depth*/
func (p OctalTree) EncodePoint(point vector.Vec, desired_depth int) []uint8 {
	encoding := make([]uint8, p.MaxDepth*3)
	current_centroid := p.Origin
	current_bounds := p.Bounds
	for i := 0; i < p.MaxDepth*3; i += 3 {
		if point[0] > current_centroid[0] {
			encoding[i] = uint8(1)
			current_centroid[0] = current_centroid[0] + (0.5 * current_bounds[0])
		} else {
			current_centroid[0] = current_centroid[0] + (-0.5 * current_bounds[0])
		}
		if point[1] > current_centroid[1] {
			encoding[i+1] = uint8(1)
			current_centroid[1] = current_centroid[1] + (0.5 * current_bounds[1])
		} else {
			current_centroid[1] = current_centroid[1] + (-0.5 * current_bounds[1])
		}
		if point[2] > current_centroid[2] {
			encoding[i+2] = uint8(1)
			current_centroid[2] = current_centroid[2] + (0.5 * current_bounds[2])
		} else {
			current_centroid[2] = current_centroid[2] + (-0.5 * current_bounds[2])
		}

		current_bounds[0] = current_bounds[0] / 2
		current_bounds[1] = current_bounds[1] / 2
		current_bounds[2] = current_bounds[2] / 2
	}
	return encoding
}

/* Generalizes the encoding of a sequential point mesh grouping into a single
octal encoding*/
func (p OctalTree) EncodePointGroup(group []vector.Vec, depth_encode int) []uint8 {
	encodes := make([][]uint8, len(group))
	max_group_encode := make([]uint8, p.MaxDepth*3)
	for i := 0; i < len(group); i++ {
		encodes[i] = p.EncodePoint(group[i], depth_encode)
	}
	depth := 0
	similar := true
	for similar && depth < p.MaxDepth*3 {
		first := encodes[0][depth]
		second := encodes[0][depth+1]
		third := encodes[0][depth+2]
		for j := 0; j < len(group); j++ {
			if first != encodes[j][depth] {
				similar = false
			}
			if second != encodes[j][depth+1] {
				similar = false
			}
			if third != encodes[j][depth+2] {
				similar = false
			}
		}
		//We always encode at least the first group of non-similarity according to the first points grouping
		//If we want to get the point group encodings actual centroid we decode to depth - 1
		max_group_encode[depth] = first
		max_group_encode[depth+1] = second
		max_group_encode[depth+2] = third
		if similar {
			depth += 3
		}
	}

	r_encode := make([]uint8, depth)
	for i := 0; i < depth; i++ {
		r_encode[i] = max_group_encode[i]
	}

	return r_encode
}

//Returns tree depth similarity between two encodings
func (p OctalTree) DepthSimilarity(a []uint8, b []uint8) int {
	depth := 0
	similar := true

	min := len(b)
	if len(a) < min {
		min = len(a)
	}

	for similar && depth < min {
		first := a[depth]
		second := a[depth+1]
		third := a[depth+2]

		if first != b[depth] {
			similar = false
		}
		if second != b[depth+1] {
			similar = false
		}
		if third != b[depth+2] {
			similar = false
		}

		if similar {
			depth += 3
		}

	}
	return depth / 3
}
