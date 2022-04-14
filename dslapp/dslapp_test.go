package dslapp

import (
	"testing"
)

func TestMain(t *testing.T) {

	electronApplication, err := New()

	if err != nil {
		t.Errorf("Failed to launch application %s", err)
		return

	}

	err = electronApplication.Init() /*Application Blocks*/

	if err != nil {
		t.Errorf("Failed to launch window %s", err)
		return
	}

}
