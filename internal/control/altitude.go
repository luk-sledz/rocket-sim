package control

const (
	DefaultKp            = 0.02  // ile „dokręcić” Throttle na 1 m błędu
	DefaultKd            = 0.08  // ile „dokręcić” Throttle na 1 m/s prędkości
	g                    = 9.81 // przyspieszenie ziemskie
)

//thrust musi byc w zakresie 0-1
func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// AltitudeThrottle liczy zadanie ciągu dla jednej pętli wysokości.
// posY, velY — stan; target — np. 90 m; dt — na razie nieużywane.
func AltitudeThrottle(posY, velY, target, mass, maxThrust, dt float64) float64 {
	_ = dt
	hover := mass * g / maxThrust
	err := target - posY
	u := hover + DefaultKp*err - DefaultKd*velY
	return clamp01(u)
}