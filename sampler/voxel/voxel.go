//Voxel Array Based Storage Structure O1 Lookup / Edit Times Using Voxel Index Bucket Technique
//Voxel Array Runs an Update Thread To Continuously Edit Particle Position Indexes
package sampler

import (
	"fmt"
	V "github.com/andewx/dieselfluid/math/math64"
)

const MAX_DIM = 15
const LOAD_FACTOR = 4.0
const VOXEL_SAMPLES = 60
const THREAD_ERROR = 404
const THREAD_RUN_SAMPLER = 60
const THREAD_WAIT_SAMPLER = 61

//Cubical Structure
type VRanges struct {
	Min       float64
	Max       float64
	Divisions int
	Buckets   int
	DivLength float64
	Particles int
}

//Storage Container For Particle Reference Indexes
type VoxelArray struct {
	PositionsRef            []V.Vec
	VoxelDescriptor         VRanges
	ComponentStorageSizeMax int
	Voxel                   [][]int
	PVoxelIdx               ParticleVoxelArray
	Utilized                int
	Samples                 int
}

//Stores the X,Y,Z, POSITION For Particle References in the Grid
//Header Positions Will Return X,Y,Z,0 Quartet
type VoxelIndex struct {
	VoxIndex [2]int
}

//Neighbor Volume Look Structure - Nil Refs Out Of Bounds - Includes Center
type VoxelLookup struct {
	Indexes []VoxelIndex
}

//Particle Index Voxel Index Reference
type ParticleVoxelArray struct {
	Indexes []VoxelIndex
}

//Allocates 1.5 x Number Particles Voxel Based Array At A Resolution Not To Exceed MAX_DIM
//If A Particle Cannot Be Added To the Voxel Bucket It will attempt to insert in Neighbor Grids
//If No Positions Are Found In The Neighbor Grid The Particle Index Won't Be Added Or Updated
func Allocate(positions []V.Vec, res int, scale_grid float64) *VoxelArray {
	//Determine Bucket Header length
	buckets := res * res * res
	num_particles := len(positions)
	bucketLength := int(float64(num_particles/buckets)*float64(LOAD_FACTOR) + 1.0)

	//Set Main Voxel Headers
	VoxelStorage := VoxelArray{}
	VoxelStorage.PositionsRef = positions

	VDescrip := VRanges{}
	VDescrip.Min = -1.0 * scale_grid
	VDescrip.Max = 1.0 * scale_grid
	VDescrip.Divisions = res
	VDescrip.Buckets = bucketLength
	VDescrip.DivLength = (VDescrip.Max - VDescrip.Min) / float64(res)
	VoxelStorage.VoxelDescriptor = VDescrip
	VoxelStorage.PVoxelIdx = ParticleVoxelArray{make([]VoxelIndex, len(positions))}
	VoxelStorage.VoxelDescriptor.Particles = num_particles
	//Loop Construct Storage Structure
	VoxelStorage.Voxel = make([][]int, res*res*res)
	for i := 0; i < res*res*res; i++ {
		VoxelStorage.Voxel[i] = make([]int, VDescrip.Buckets)
	}

	//Nil All Voxel Refs
	for i := 0; i < res*res*res; i++ {
		for t := 0; t < VDescrip.Buckets; t++ {
			VoxelStorage.Voxel[i][t] = -1
		}
	}

	VoxelStorage.Samples = 25

	//Set all Reference indexes to -1 For Initialization
	for i := 0; i < VoxelStorage.VoxelDescriptor.Particles; i++ {
		VoxelStorage.PVoxelIdx.Indexes[i].VoxIndex[0] = -1
		VoxelStorage.PVoxelIdx.Indexes[i].VoxIndex[1] = -1
	}

	return &VoxelStorage
}

func (v *VoxelIndex) IsNil() bool {
	if v.VoxIndex[0] == -1 {
		return true
	}
	return false
}

//Voxel Array Update
func (v *VoxelArray) UpdateSampler() {
	v.Utilized = 0
	for i := 0; i < v.VoxelDescriptor.Particles; i++ {
		v.Update(i)
	}
	if v.Utilization() < 0.90 {
		fmt.Printf("\nError Voxel Particles Exclusion: %f ", v.Utilization())
	}
}

func (v *VoxelArray) Update(pindex int) {
	//Find POSITION

	refIndex := v.PVoxelIdx.Indexes[pindex]
	r := refIndex.VoxIndex[0]
	s := refIndex.VoxIndex[1]

	vIndex := v.VoxelHash(v.PositionsRef[pindex])
	x := vIndex.VoxIndex[0]

	//No Change In Particle Bucket
	if r == x && !refIndex.IsNil() {
		v.Utilized++
		return
	}

	//Place particle in bucket if available
	for i := 0; i < len(v.Voxel[x]); i++ {
		if v.Voxel[x][i] == -1 {
			v.Voxel[x][i] = pindex
			v.PVoxelIdx.Indexes[pindex] = VoxelIndex{[2]int{x, i}}
			if !refIndex.IsNil() {
				v.Voxel[r][s] = -1
			}
			v.Utilized++
			return
		}
	}

	//Resize the Bucket Array If Not Found- Keep track of this for resize rates
	nLength := len(v.Voxel[x]) * 2
	nBuffer := make([]int, nLength)
	if nLength > v.VoxelDescriptor.Buckets*4 {
		fmt.Printf("Voxel[%d] usage beyond limit %d\n", x, nLength)
	}
	//Copy All Elements
	locatedIndex := 0
	for i := 0; i < nLength; i++ {
		if i < int(nLength/2) {
			nBuffer[i] = v.Voxel[x][i]
			locatedIndex = i + 1
		} else {
			nBuffer[i] = -1
		}
	}
	//Place Particle & Reset Buffer Pointer
	nBuffer[locatedIndex] = pindex
	if !refIndex.IsNil() {
		v.Voxel[r][s] = -1
	}
	v.Voxel[x] = nBuffer
	v.Utilized++

	return

}

func (v *VoxelArray) Utilization() float64 {
	return float64(v.Utilized) / float64(v.VoxelDescriptor.Particles)
}

//Modulus Hashes A Position Into a Voxel Index Bucket - Lazy Hash Method
func (v *VoxelArray) VoxelHash(pos V.Vec) VoxelIndex {
	x := pos[0]
	y := pos[1]
	z := pos[2]

	x = x - v.VoxelDescriptor.Min
	y = y - v.VoxelDescriptor.Min
	z = z - v.VoxelDescriptor.Min
	//Hash Position Indexes
	x0 := int(x/v.VoxelDescriptor.DivLength) % v.VoxelDescriptor.Divisions
	y0 := int(y/v.VoxelDescriptor.DivLength) % v.VoxelDescriptor.Divisions
	z0 := int(z/v.VoxelDescriptor.DivLength) % v.VoxelDescriptor.Divisions

	//All Indexes are Positive
	if x0 < 0 {
		x0 = x0 + v.VoxelDescriptor.Divisions
	}
	if y0 < 0 {
		y0 = y0 + v.VoxelDescriptor.Divisions
	}
	if z0 < 0 {
		z0 = z0 + v.VoxelDescriptor.Divisions
	}

	VIndex := VoxelIndex{}
	VIndex.VoxIndex[0] = (z0*v.VoxelDescriptor.Divisions+y0)*v.VoxelDescriptor.Divisions + x0
	VIndex.VoxIndex[1] = 0

	return VIndex

}

func (v *VoxelArray) GetSamples(idx int) []int {
	return v.GetSampleVoxels(v.PositionsRef[idx], idx)
}

//Gets Samples from position
func (v *VoxelArray) GetSampleVoxels(pos V.Vec, idx int) []int {
	mNeighbors := v.VolumeLookup(pos)
	sampleIndexes := make([]int, v.Samples)
	index := 0
	nLength := len(mNeighbors.Indexes)
	pindex := v.VoxelHash(pos)
	//Get Initial Samples From This Particle Volume
	for i := 0; i < len(v.Voxel[pindex.VoxIndex[0]]) && index < v.Samples; i++ {
		sIndex := v.Voxel[pindex.VoxIndex[0]][i]
		if sIndex != -1 && sIndex != idx {
			sampleIndexes[index] = v.Voxel[pindex.VoxIndex[0]][i]
			index++
		}
	}

	//Get Neighbor Samples
	for i := 0; i < nLength && index < v.Samples; i++ {
		r := mNeighbors.Indexes[i].VoxIndex[0]
		for q := 0; q < len(v.Voxel[r]) && index < v.Samples; q++ {
			sIndex := v.Voxel[r][q]
			if sIndex != -1 {
				sampleIndexes[index] = v.Voxel[r][q]
				index++
			}
		}
	}

	return sampleIndexes
}

func (v *VoxelArray) GetAllNeighbors(idx int) []int {
	mNeighbors := v.VolumeLookup(v.PositionsRef[idx])
	listCount := 0

	//Count all Neighbors for ALLOCATION
	for i := 0; i < len(mNeighbors.Indexes); i++ {
		q := mNeighbors.Indexes[i].VoxIndex[0]
		for j := 0; j < len(v.Voxel[q]); j++ {
			particleIndex := v.Voxel[q][j]
			if particleIndex != -1 && particleIndex != idx {
				listCount++
			}
		}
	}

	neigh := make([]int, listCount)
	index := 0

	for i := 0; i < len(mNeighbors.Indexes); i++ {
		q := mNeighbors.Indexes[i].VoxIndex[0]
		for j := 0; j < len(v.Voxel[q]); j++ {
			particleIndex := v.Voxel[q][j]

			if index >= listCount {
				return neigh
			}

			if particleIndex != -1 && particleIndex != idx {
				neigh[index] = v.Voxel[q][j]
				index++
				if index >= listCount {
					return neigh
				}
			}
		}
	}

	return neigh

}

//Constructs Neighbor Voxels with position hashes
func (v *VoxelArray) VolumeLookup(pos V.Vec) VoxelLookup {

	vxNeighbors := VoxelLookup{make([]VoxelIndex, 27)}
	i := 0
	//Set Neighbor Volume Indexes
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			for z := -1; z < 2; z++ {
				nPos := V.Scl(pos, 1.0)
				addVec := V.Vec{float64(x) * v.VoxelDescriptor.DivLength, float64(y) * v.VoxelDescriptor.DivLength, float64(z) * v.VoxelDescriptor.DivLength}
				nPos = V.Add(nPos, addVec)
				vxlHash := v.VoxelHash(nPos)
				vxNeighbors.Indexes[i] = vxlHash
				i++
			}
		}
	}
	return vxNeighbors

}

//Thread Run
//Runs sampler interface run. Performs the core sampler maintence taskings
func (s VoxelArray) Run(status chan int) {
	done := false
	synced := true

	for !done {
		if synced {
			s.UpdateSampler()
			synced = false
			status <- THREAD_WAIT_SAMPLER //Thread has been synced
			nStatus := <-status           //Wait on the next status update
			if nStatus == THREAD_RUN_SAMPLER {
				synced = true
			}
		}
	}
}

//Estimated Storage Requirement - Changes when resizes are called
func (v *VoxelArray) PrintStorageRequirements() {
	particleBytes := v.VoxelDescriptor.Particles * 16
	cells := v.VoxelDescriptor.Divisions * v.VoxelDescriptor.Divisions * v.VoxelDescriptor.Divisions
	voxelIndexBytes := cells * 4 * v.VoxelDescriptor.Buckets
	totalBytes := particleBytes + voxelIndexBytes
	kilobytes := totalBytes / 1024

	fmt.Printf("Voxel Grid Storage: %dkB\nBucket Size: %d\nVoxel Cells: %d\n", kilobytes, v.VoxelDescriptor.Divisions, cells)

}
