package light

import (
	"fmt"
	"testing"
)

func TestPNGCreate(t *testing.T) {
	mSky := NewSky()
	mSky.UpdateDay(200.5)
	mSky.CreateTexture()
	fmt.Printf("YOU ROCK. Bring it home boys\n")
}
