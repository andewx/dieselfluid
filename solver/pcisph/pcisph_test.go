package pcisph

import "testing"

func TestOpenCompute(t *testing.T) {
	gpu, err := New_GPUPredictorCorrector(8, nil)
	if err != nil {
		t.Errorf("Failed gpu PCISPH Implementation %v", err)
	}
	gpu.Run()

}
