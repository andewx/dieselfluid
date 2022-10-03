package common

import "testing"
import "os"

func TestPaths(t *testing.T) {
	wd, _ := os.Getwd()
	if ProjectRelativePath("common") != wd {
		t.Errorf("Unable to reconstruct directory")
	}
}
