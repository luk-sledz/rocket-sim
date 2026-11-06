package control

const (
	DefaultAngleKp = 0.02 // rad na 1 m błędu X
	DefaultMaxTilt = 0.15  // rad ≈ 8.6° — mały tilt, żeby cos≈1
	DefaultAngleKd = 0.08  // rad na 1 m/s błędu X
)

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// regulator kąta kp + kd
func LateralAngle(posX, velX float64) float64 {
	u := -DefaultAngleKp*posX - DefaultAngleKd*velX
	return clamp(u, -DefaultMaxTilt, DefaultMaxTilt)
}