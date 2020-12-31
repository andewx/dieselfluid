package sampler

import "github.com/andewx/dieselfluid/sampler/voxel"

//Sampler type interface defines sampler function interface
type Sampler interface {
	UpdateSampler()
	Run()
	GetSamples(i int) []int
}

type VoxelSampler struct {
	Voxel voxel.VoxelArray
}
