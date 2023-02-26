//Voxel Array Based Storage Structure O1 Lookup / Edit Times Using Voxel Index Bucket Technique
//Voxel Array Runs an Update Thread To Continuously Edit Particle Position Indexes
package lsh

import (
	"math/rand"
	"time"

	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/model"
)

const LOAD_FACTOR = float32(1.5)
const THREAD_ERROR = 404
const THREAD_RUN_SAMPLER = 60
const THREAD_WAIT_SAMPLER = 61
const SAMPLES = 100

//Constructs locality sensitive hash using random projection vectors in the euclidian plane using 8 bit hash
type HashSampler struct {
	Table       [][]int
	Buckets     int
	Size        int
	HashVectors []vector.Vec
	HashBits    int
	particles   *model.ParticleArray
}

//Dynamically allocates the hash table
func Allocate(num_particles int, buckets int, hash_bits int, particles *model.ParticleArray) *HashSampler {
	sampler := HashSampler{}
	time_seed := time.Now()
	factor := int(float32(num_particles/buckets) * LOAD_FACTOR)
	sampler.Table = make([][]int, buckets)
	r := rand.New(rand.NewSource(int64(time_seed.Second())))
	sampler.HashVectors = make([]vector.Vec, hash_bits)

	//Allocate all Hash Vectors
	for i := 0; i < hash_bits; i++ {
		sampler.HashVectors[i] = vector.Vec{r.Float32() - 0.5, r.Float32() - 0.5, r.Float32() - 0.5}
	}

	sampler.Buckets = buckets
	sampler.HashBits = hash_bits
	sampler.Size = factor
	sampler.particles = particles

	return &sampler
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

//Constructs finite sampler
func (s HashSampler) GetData1D() []int {
	nArray := make([]int, s.Buckets*s.Size)
	for i := 0; i < s.Buckets; i++ {
		for j := 0; j < s.Size; j++ {
			if s.Table[i] != nil && j < len(s.Table[i]) {
				nArray[i*s.Size+j] = s.Table[i][j]
			}
		}
	}
	return nArray
}

func (s HashSampler) GetVectors() []float32 {
	nArray := make([]float32, len(s.HashVectors)*3)
	for i := 0; i < len(s.HashVectors); i++ {
		for j := 0; j < 3; j++ {
			nArray[i*3+j] = s.HashVectors[i][j]
		}
	}
	return nArray
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
	return hash % (s.Buckets)
}

func (s HashSampler) Insert(hash int, particle int) {
	if s.Table[hash] == nil {
		s.Table[hash] = make([]int, 0)
	}
	s.Table[hash] = append(s.Table[hash], particle)
}

func (s HashSampler) Reset() {
	for i := 0; i < s.Buckets; i++ {
		s.Table[i] = nil
	}
}

func (s HashSampler) UpdateSampler() {
	s.Reset()

	for i := 0; i < s.particles.Total(); i++ {
		particle := s.particles.Get(i)
		s.Insert(s.Hash(particle.Position), i)
	}
}

//FINDS # SAMPLES
func (s HashSampler) GetSamples(x int) []int {
	particle := s.particles.Get(x)
	samples := make([]int, SAMPLES)
	num := int(0)
	index := s.Hash(particle.Position)

	for num < SAMPLES {
		if index > len(s.Table) {
			index = 0
		}
		table := s.Table[index]
		if table == nil {
			index++
			index = index % s.Buckets
		} else {
			for j := 0; j < len(table) && num < SAMPLES; j++ {
				samples[num] = table[j]
				num++
			}
		}
	}
	return samples
}

func (s HashSampler) GetSamplesFromPosition(pos []float32) []int {
	samples := make([]int, SAMPLES)
	num := int(0)
	index := s.Hash(vector.CastFixed(pos))

	for num < SAMPLES {
		if index > len(s.Table) {
			index = 0
		}
		table := s.Table[index]
		if table == nil {
			index++
			index = index % s.Buckets
		} else {
			for j := 0; j < len(table) && num < SAMPLES; j++ {
				samples[num] = table[j]
				num++
			}
		}
	}
	return samples
}

func (s HashSampler) Run(status chan string) {
	done := false
	for !done {
		st := <-status
		if st == "SAMPLER_UPDATE" {
			s.UpdateSampler()
		}
	}
}
