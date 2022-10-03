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
		if err != nil {
			t.Errorf("Sphere2Vec Error %s", err.Error())
		}
	} else {
		fmt.Printf("---Pass: Sphere2Vec---\n")
	}
}

func TestInitialization(t *testing.T) {

	bVec := vector.Vec{1.0, 0.0, 0.0}
	cVec := vector.Vec{0.0, 1.0, 0.0}

	bVec2Polar, _ := Vec2Sphere(bVec)
	cVec2Polar, _ := Vec2Sphere(cVec)

	bPolar2Vec, _ := Sphere2Vec(bVec2Polar)
	cPolar2Vec, _ := Sphere2Vec(cVec2Polar)

	if !bVec.Equal(bPolar2Vec) {
		t.Errorf(bVec.ToString() + "!= " + bPolar2Vec.ToString())
	}

	if !cVec.Equal(cPolar2Vec) {
		t.Errorf(cVec.ToString() + "!= " + cPolar2Vec.ToString())
	}
}
