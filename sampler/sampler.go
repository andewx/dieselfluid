package sampler

type Sampler interface {
	UpdateSampler()
	Run()
	GetSamples(i int)
}
