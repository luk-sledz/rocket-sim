package control

const (
	DefaultPosKp           = 0.25 // (m/s) na 1 m błędu X → desiredVx
	DefaultMaxLateralSpeed = 25.0 // m/s  max prędkość w X
	DefaultVelKp           = 0.08 // rad na 1 m/s błędu prędkości → Angle
	DefaultMaxTilt         = 0.15 // rad ≈ 8.6°
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

// wylicza Vx na podstawie odległości od wieży
func LateralDesiredVx(posX float64) float64 {
	return clamp(-DefaultPosKp*posX, -DefaultMaxLateralSpeed, DefaultMaxLateralSpeed)
}

// wylicza kąt na podstawie błędu prędkości
func LateralAngleFromVx(velX, desiredVx float64) float64 {
	u := -DefaultVelKp * (velX - desiredVx)
	return clamp(u, -DefaultMaxTilt, DefaultMaxTilt)
}
