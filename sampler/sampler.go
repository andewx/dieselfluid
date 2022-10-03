package sampler

import "github.com/andewx/dieselfluid/math/vector"

//Sampler represents abstracted sampler class for domains
type Sampler interface {
	UpdateSampler()
	Run(status chan int)
	Hash([3]float32) int
	GetSamples(i int) []int
	GetRegionalSamples(hash int, width int) []int
	GetData() [][]int
	GetElements() int
	GetVectors() []vector.Vec
	GetHashSize() int
	GetBuckets() int
	BucketSize() int
}
