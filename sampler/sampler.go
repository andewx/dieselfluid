package sampler

import "github.com/andewx/dieselfluid/sampler/voxel"

//Sampler represents abstracted sampler classes for 1D,2D,3D samplers with given
//PDF sampling functions.
type Sampler interface {
	UpdateSampler()
	Run()
	GetSamples(i int) []int
}

type VoxelSampler struct {
	Voxel voxel.VoxelArray
}
