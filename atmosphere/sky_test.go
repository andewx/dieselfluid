package atmosphere

import (
	"strconv"
	"testing"
)

func TestPNGCreate(t *testing.T) {
	mSky := NewAtmosphere(45.0, 0.0)
	base := "sky_"
	mSky.UpdatePosition(365 * 44 / 48)
	for i := 44; i < 48; i++ {
		filename := base + strconv.FormatInt(int64(i), 10) + ".png"
		mSky.UpdatePosition(365 / 48)
		mSky.CreateTexture(512, 512, true, filename)
	}
}

/*
func TestEnvBox(t *testing.T) {
	mSky := NewAtmosphere(45.0, 0.0)
	mSky.UpdatePosition(1.5)
	mSky.CreateEnvBox(512, 512, true)
}
*/
