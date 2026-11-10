package control

const (
	DefaultPosKp           = 0.25 // (m/s) na 1 m błędu X → desiredVx
	DefaultMaxLateralSpeed = 25.0 // m/s  max prędkość w X
	DefaultVelKp           = 0.08 // rad na 1 m/s błędu prędkości → Angle
	DefaultMaxTilt         = 0.15 // rad ≈ 8.6°
	DefaultAttKp     = 1.5  // rad gimbala na 1 rad błędu Angle
	DefaultAttKd     = 0.4  // rad gimbala na 1 rad/s ω (tłumienie)
	DefaultMaxGimbalAngle = 0.20 // rad — ten sam limit co w rocket
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

// Regulator PD ylicza kąt gimbalu na podstawie błędu kąta i prędkości kątowej 
func AttitudeGimbal(angle, angleCmd, angularVel float64) float64 {
	u := DefaultAttKp*(angleCmd-angle) - DefaultAttKd*angularVel
	return clamp(u, -DefaultMaxGimbalAngle, DefaultMaxGimbalAngle)
}