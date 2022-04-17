package polar

import (
	"fmt"
	"testing"

	"github.com/andewx/dieselfluid/math/vector"
)

func TestPolar(t *testing.T) {
	a := NewPolar(1)
	v, err := Sphere2Vec(a)
	fmt.Printf("---Testing: Sphere2Vec---\n")
	if v.Equal(vector.Vec{0, 0, 0}) || err != nil {
		t.Errorf("v is zero vector")
		if err != nil {
			t.Errorf("Sphere2Vec Error %s", err.Error())
		}
	} else {
		fmt.Printf("---Pass: Sphere2Vec---\n")
	}
}
