package sampler

//Sampler type interface defines sampler function interface
type Sampler interface {
	UpdateSampler()
	Run()
	GetSamples(i int) []int
}
