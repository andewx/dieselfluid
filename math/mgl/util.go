package mgl

func Clamp1f(a float32, min float32, max float32) float32 {
	if a >= min && a <= max {
		return a
	}

	if a < min {
		return min
	}

	if a > max {
		return max
	}

	return a
}
