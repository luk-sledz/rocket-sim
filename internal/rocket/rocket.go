package rocket

import (
	"math"
	"rocket-sim/internal/vec"
)

// Domyślne parametry (Starship przy lądowaniu, ~3 Raptory).
const (
	G = 9.81
	DefaultDryMass   = 1.2e5     // kg (~120 t)
	DefaultFuel      = 3.0e4     // kg — zapas paliwana hamowanie, nie pełny tank
	DefaultMaxThrust = 3 * 2.3e6 // N — ~3 × Raptor
	DefaultIsp       = 330.0     // s — ile ciągu z 1 kg paliwa (Raptor, atmosfera)
	TowerHeight = 90.0 // wysokość ramion wieży chwytającej [m] 
	DefaultCatchTolY     = 0.5 // zakres tolerancji w Y podczas catch
	DefaultCatchTolX     = 0.5 // zakres tolerancji w X podczas catch
	DefaultCatchMaxSpeed = 1.0 // m/s — |Vel.Y|
	DefaultAngleTau = 0.2 // s — czas regulacji kąta
	DefaultMaxAngularRate = 1.0 // rad/s — maksymalna prędkość kątowa
)

// Rocket to stan rakiety w symulacji 2D (jednostki SI).
type Rocket struct {
	Pos       vec.Vec2 // pozycja [m]; Y = wysokość nad ziemią
	Vel       vec.Vec2 // prędkość [m/s]
	DryMass   float64  // masa bez paliwa [kg]
	Fuel      float64  // masa paliwa [kg]
	MaxThrust float64  // maksymalny ciąg [N], w górę (+Y)
	Isp       float64  // impuls właściwy [s]
	Throttle  float64  // zadanie ciągu w [0, 1]
	Angle     float64  // kąt nachylenia rakiety [rad]
	AngleCmd  float64  // kąt nachylenia rakiety żądany [rad]
	AngularVel float64 // prędkość kątowa rakiety [rad/s]
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// New tworzy rakietę. + sprawdza poprawność parametrów.
func New(pos, vel vec.Vec2, dryMass, fuel, maxThrust, isp float64) *Rocket {
	if dryMass <= 0 {
		panic("rocket: dryMass musi być > 0")
	}
	if fuel < 0 {
		panic("rocket: fuel nie może być ujemne")
	}
	if maxThrust < 0 {
		panic("rocket: maxThrust nie może być ujemne")
	}
	if isp <= 0 {
		panic("rocket: isp musi być > 0")
	}
	return &Rocket{
		Pos:       pos,
		Vel:       vel,
		DryMass:   dryMass,
		Fuel:      fuel,
		MaxThrust: maxThrust,
		Isp:       isp,
		Throttle:  0, // start z wyłączonym silnikiem
		Angle:     0, // start rakieta pionowo
	}
}

// NewStarship tworzy rakietę z parametrami w skali Starshipa.
func NewStarship(pos, vel vec.Vec2) *Rocket {
	return New(pos, vel, DefaultDryMass, DefaultFuel, DefaultMaxThrust, DefaultIsp)
}

// Mass zwraca aktualną masę [kg]: sucha + paliwo.
func (r *Rocket) Mass() float64 {
	return r.DryMass + r.Fuel
}

// Alive zwraca true, gdy rakieta nad ziemią.
func (r *Rocket) Alive() bool {
	return r.Pos.Y > 0
}

// Step symuluje kolejny krok symulacji.
func (r *Rocket) Step(dt float64) {
	if !r.Alive() || dt <= 0 {
		return
	}

	// 1) postawa: Angle dogania AngleCmd 
	omegaCmd := (r.AngleCmd - r.Angle) / DefaultAngleTau
	r.AngularVel = clamp(omegaCmd, -DefaultMaxAngularRate, DefaultMaxAngularRate)
	r.Angle += r.AngularVel * dt
	
	th := r.Throttle
	if th < 0 {
		th = 0
	}
	if th > 1 {
		th = 1
	}
	thrust := th * r.MaxThrust
	if r.Fuel <= 0 {
		thrust = 0
	}
	
	m := r.Mass()
	fx := thrust * math.Sin(r.Angle)
	fy := thrust * math.Cos(r.Angle)
	force := vec.Vec2{X: fx, Y: fy - m*G} // rozkład sił w zależności od kąta nachylenia
	a := force.Scale(1 / m)                   // a = F/m
	r.Vel = r.Vel.Add(a.Scale(dt))
	r.Pos = r.Pos.Add(r.Vel.Scale(dt))
	
	// zużycie paliwa: ṁ = F / (Isp * g0)
	if thrust > 0 {
		r.Fuel -= thrust / (r.Isp * G) * dt
		if r.Fuel < 0 {
			r.Fuel = 0
		}
	}
	
	// 3) ziemia
	if r.Pos.Y < 0 {
		r.Pos.Y = 0
		r.Vel = vec.Vec2{}
	}
}
// Caught zwraca true, gdy rakieta jest w oknie chwytu ramion wieży.
// tolY — tolerancja wysokości [m], maxSpeed — max |Vel.Y| [m/s].
func (r *Rocket) Caught(tolX, tolY, maxSpeed float64) bool {
	if !r.Alive() {
		return false
	}
	nearX := math.Abs(r.Pos.X) <= tolX
	nearY := math.Abs(r.Pos.Y-TowerHeight) <= tolY
	slowX := math.Abs(r.Vel.X) <= maxSpeed
	slowY := math.Abs(r.Vel.Y) <= maxSpeed
	return nearX && nearY && slowX && slowY
}