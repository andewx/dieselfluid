package sampler

//Sampler represents abstracted sampler class for domains
type Sampler interface {
	UpdateSampler()
	Run(status chan int)
	Hash([3]float32) int
	GetSamples(i int) []int
	GetRegionalSamples(hash int, width int) []int
	GetData() [][]int
	GetElements() int
	GetVectors() []float32
	GetHashSize() int
	GetBuckets() int
	BucketSize() int
	GetData1D() []int
}
