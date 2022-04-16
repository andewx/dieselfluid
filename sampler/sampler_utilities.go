package sampler

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/andewx/dieselfluid/math/mgl"
)

/*-----------------------Sampler Dataset Utility Functions--------------------*/

//Linearly Interpolates between two values
func Lerp(c float32, v0 float32, v1 float32) float32 {
	r := v1 - v0
	if c < 0.0 || c > 1.0 {
		c = 0.0
	}
	return v0 + r*c
}

// @no-mutate / @:sampler / @:interpolation / @:functional
//Maps linear interpolation to a exponential function where x is in the [0,1]
// in [0,1 domain]. w is an exponential scale parameter for e^(-wq)
func Ease(x float32, w float32) float32 {
	x = mgl.Clamp1f(x, 0, 1)
	return float32(math.Exp(float64(w*x - w)))
}

//For now just looks in resource folder
func ImportSampler(resource string) (*SamplerJSON, error) {
	myImporter := SamplerJSON{}
	content, err := ioutil.ReadFile("../" + resource)
	if err != nil {
		fmt.Printf("Unable to load GLTF File\n")
		log.Fatal(err)
	}
	jsonErr := myImporter.UnmarshalJSON(content)

	if jsonErr != nil {
		fmt.Println(jsonErr)
		return nil, jsonErr
	}

	return &myImporter, nil
}

//Returns an average value under a sampling domain given value slices and their
//domain areas we assume the data provided is sorted. Out r domains sample
//endpoint in the domain. We can pass in the sample rs as slices
func SampleAverage1D(domain []float32, values []float32, n int, start_domain float32, end_domain float32, bindIndex *int) (float32, error) {

	sum := float32(0.0)
	r := end_domain - start_domain

	//Parameter error
	if n < 0 || n > len(domain) || n > len(values) || *bindIndex > n {
		return 0.0, fmt.Errorf("sampler.SampleAverage1D() Paramter error see function definition")
	}

	//Bounds single value or out of r
	if n == 1 || (domain[0] >= end_domain) {
		return values[0], fmt.Errorf("sampler.SampleAverage1D()-Domain Window Error")
	}

	//Domain completely out of r
	if domain[n-1] <= start_domain {
		return values[n-1], fmt.Errorf("sampler.SampleAverage1D()-Domain window error")
	}

	//Contributions of constant segments before/after sample rs
	if start_domain < domain[0] {
		sum += values[0] * (domain[0] - start_domain)
	}
	if end_domain > domain[n-1] {
		sum += values[n-1] * (end_domain - domain[n-1])
	}

	found := false
	i := *bindIndex
	//Search starting contribution
	for i < n && !found {
		if start_domain < domain[i] {
			found = true
		}
		i++
	}

	//Compute interpolated segments for specified domain
	for i < n-1 && domain[i] >= start_domain && domain[i] <= end_domain {
		segStart := max(start_domain, domain[i])
		segEnd := min(end_domain, domain[i+1])
		v0 := values[i]
		v1 := values[i+1]
		half := Lerp(0.5, v0, v1)
		sum += half * (segEnd - segStart)
		i++
	}
	*bindIndex = i
	return sum / r, nil
}

func max(a float32, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func min(a float32, b float32) float32 {
	if a > b {
		return b
	}
	return a
}
