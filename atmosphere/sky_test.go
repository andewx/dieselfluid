package atmosphere

import (
	"strconv"
	"testing"
)

func TestPNGCreate(t *testing.T) {
	mSky := NewAtmosphere(45.0, 0.0)
	mSky.SetDay(200.2)
	base := "sky_"
	for i := 0; i < 10; i++ {
		filename := base + strconv.FormatInt(int64(i), 10) + ".png"
		mSky.StepDay(1 + float32(i)/10)
		mSky.CreateTexture(256, 256, true, filename)
	}
}

func TestEnvBox(t *testing.T) {
	mSky := NewAtmosphere(45.0, 0.0)
	mSky.SetDay(200.2)
	mSky.CreateEnvBox(256, 256, true)
}
