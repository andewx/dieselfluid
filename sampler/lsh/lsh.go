//Voxel Array Based Storage Structure O1 Lookup / Edit Times Using Voxel Index Bucket Technique
//Voxel Array Runs an Update Thread To Continuously Edit Particle Position Indexes
package lsh

import (
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model"
	"math/rand"
	"time"
)

const LOAD_FACTOR = float32(1.5)
const THREAD_ERROR = 404
const THREAD_RUN_SAMPLER = 60
const THREAD_WAIT_SAMPLER = 61

//Constructs locality sensitive hash using random projection vectors in the euclidian plane
type HashSampler struct {
	Table       [][]int
	Buckets     int
	Size        int
	Indexes     []int
	HashVectors []vector.Vec
	HashBits    int
	Particles   []model.Particle
}

//Dynamically allocates the hash table
func Allocate(num_particles int, buckets int, hash_bits int, particles []model.Particle) HashSampler {
	sampler := HashSampler{}
	time_seed := time.Now()
	factor := int(float32(num_particles/buckets) * LOAD_FACTOR)
	sampler.Table = make([][]int, buckets)
	r := rand.New(rand.NewSource(int64(time_seed.Second())))
	sampler.HashVectors = make([]vector.Vec, hash_bits)
	//Allocation All tables
	for i := 0; i < buckets; i++ {
		sampler.Table[i] = make([]int, factor)
	}

	//Allocate all Hash Vectors
	for i := 0; i < hash_bits; i++ {
		sampler.HashVectors[i] = vector.Vec{r.Float32() - 0.5, r.Float32() - 0.5, r.Float32() - 0.5}
	}

	sampler.Buckets = buckets
	sampler.HashBits = hash_bits
	sampler.Size = factor
	sampler.Particles = particles
	return sampler
}

func sgn(val float32) int {
	if val <= 0 {
		return 0
	}
	return 1
}

func (s HashSampler) GetHashSize() int {
	return s.HashBits
}
func (s HashSampler) GetElements() int {
	return len(s.Table) * s.Size
}

func (s HashSampler) GetData() [][]int {
	return s.Table
}

func (s HashSampler) GetVectors() []vector.Vec {
	return s.HashVectors
}

func (s HashSampler) GetBuckets() int {
	return s.Buckets
}

func (s HashSampler) BucketSize() int {
	return s.Size
}

//Locality Sensity Hash - Random Cosine Projection into a 32 bit integer which is
//Also the maximum amount of hash bits the function will consider
func (s HashSampler) Hash(pos [3]float32) int {
	p := vector.Cast(pos)
	hash := int(0)
	for i := 0; i < s.HashBits; i++ {
		r := s.HashVectors[i]
		hash = hash << 1
		hash += sgn(vector.Dot(p, r))
	}
	return hash % s.Buckets
}

func (s HashSampler) Insert(hash int, particle int) {
	if s.Indexes[hash] < s.Size {
		s.Table[hash][s.Indexes[hash]] = particle
		s.Indexes[hash]++
	}
}

func (s HashSampler) Reset() {
	for i := 0; i < s.Buckets; i++ {
		s.Indexes[i] = 0
	}
}

func (s HashSampler) UpdateSampler() {
	s.Reset()
	for i := 0; i < len(s.Particles); i++ {
		s.Insert(s.Hash(vector.CastFixed(s.Particles[i].Position())), i)
	}
}

func (s HashSampler) GetSamples(particle int) []int {
	return s.Table[s.Hash(vector.CastFixed(s.Particles[particle].Position()))]
}

func (s HashSampler) GetRegionalSamples(hash int, width int) []int {
	w2 := int(width / 2)
	start := hash - w2
	if start < 0 {
		start = s.Buckets - (w2 - hash)
	}
	samples := make([]int, s.Size*width)
	for i := 0; i < width; i++ {
		index := (start + i) % s.Buckets
		samples = append(samples, s.Table[index][:]...)
	}
	return samples
}

func (s HashSampler) Run(status chan int) {
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
