package rocket

import (
	"math"
	"rocket-sim/internal/vec"
)

// Domyślne parametry (Starship przy lądowaniu, ~3 Raptory).
const (
	G                     = 9.81
	DefaultDryMass        = 1.2e5     // kg (~120 t)
	DefaultFuel           = 3.0e4     // kg — zapas paliwana hamowanie, nie pełny tank
	DefaultMaxThrust      = 3 * 2.3e6 // N — ~3 × Raptor
	DefaultIsp            = 330.0     // s — ile ciągu z 1 kg paliwa (Raptor, atmosfera)
	TowerHeight           = 90.0      // wysokość ramion wieży chwytającej [m]
	DefaultCatchTolY      = 0.5       // zakres tolerancji w Y podczas catch
	DefaultCatchTolX      = 0.5       // zakres tolerancji w X podczas catch
	DefaultCatchMaxSpeed  = 1.0       // m/s — |Vel.Y|
	DefaultInertia        = 2.5e7     // kg·m² ≈ m·L²/12; m~120 t, L~50 m(przybliżenie)
	DefaultLeverArm       = 15.0      // m — długość od środka masy do dyszy
	DefaultMaxGimbalAngle = 0.20      // rad ≈ 11.5° — limit kąta gimbala
	DefaultMaxGimbalRate  = 1.0       // rad/s prędkość maksymalna gimbala
	AirDensity      = 1.225 // kg/m³ — ρ na poziomie morza, stałe
    DefaultDragCd   = 1.0   // opór gdy powietrze napływa na dół rakiety
	DefaultDragCd90   = 7.08  // Cd przy α = π/2; ≈ Cd0·(450/63.6) przeskalowane dla cylindra (Opór gdy powietrze napływa na bok rakiety)

    DefaultDragArea = 63.6  // m² — przekrój ≈ π·(4.5 m)²
)

// Rocket to stan rakiety w symulacji 2D (jednostki SI).
type Rocket struct {
	Pos         vec.Vec2 // pozycja [m]; Y = wysokość nad ziemią
	Vel         vec.Vec2 // prędkość [m/s]
	DryMass     float64  // masa bez paliwa [kg]
	Fuel        float64  // masa paliwa [kg]
	MaxThrust   float64  // maksymalny ciąg [N], w górę (+Y)
	Isp         float64  // impuls właściwy [s]
	Throttle    float64  // zadanie ciągu w [0, 1]
	Angle       float64  // kąt nachylenia rakiety [rad]
	AngleCmd    float64  // kąt nachylenia rakiety żądany [rad]
	AngularVel  float64  // prędkość kątowa rakiety [rad/s]
	GimbalAngle float64  // kąt gimbala względem kadłuba [rad]
	Inertia     float64  // I
	LeverArm    float64  // L — dźwignia od środka masy do dyszy
	GimbalCmd   float64  // żądany kąt gimbala [rad]
	DragCd   float64 // współczynnik oporu Cd
	DragCd90 float64 // Cd przy napływie z boku (α = π/2)
	DragArea float64 // powierzchnia odniesienia A [m²]
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
		Inertia:   DefaultInertia,
		LeverArm:  DefaultLeverArm,
		DragCd:    DefaultDragCd,
		DragCd90:  DefaultDragCd90,
		DragArea:  DefaultDragArea,
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

// Zwraca kąt z którego napływa powietrze na rakietę 
func (r *Rocket) AttackAngle() float64 {
	speed := r.Vel.Len()
	if speed == 0 {
		return 0
	}
	// oś kadłuba: Angle=0 → dziób w +Y
	bx := math.Sin(r.Angle)
	by := math.Cos(r.Angle)
	// kierunek ruchu
	ux := r.Vel.X / speed 
	uy := r.Vel.Y / speed 
	dot := bx*ux + by*uy // iloczyn skalarny wektorów b i u (a·b = cosα) 
	dot = clamp(dot, -1, 1) // ochrona przed acos NaN z błędów float
	return math.Acos(dot) // rad — kąt między dziobem kadłuba a kierunkiem Vel
}

// DragCoeff zwraca Cd(α): Cd0 + (Cd90−Cd0)·sin²(α).
// α=0 lub π → Cd0; α=π/2 → Cd90.
func (r *Rocket) DragCoeff() float64 {
	s := math.Sin(r.AttackAngle())
	return r.DragCd + (r.DragCd90-r.DragCd)*s*s
}

// LiftCoeff zwraca Cl(α) = ½·(Cd90−Cd0)·sin(2α),
// α znakowane: atan2(û_vel×û_body, û_body·û_vel) ∈ (−π, π].
func (r *Rocket) LiftCoeff() float64 {
	speed := r.Vel.Len()
	if speed == 0 {
		return 0
	}
	bx := math.Sin(r.Angle)
	by := math.Cos(r.Angle)
	ux := r.Vel.X / speed
	uy := r.Vel.Y / speed
	dot := bx*ux + by*uy
	cross := ux*by - uy*bx
	alpha := math.Atan2(cross, dot)
	return 0.5 * (r.DragCd90 - r.DragCd) * math.Sin(2*alpha)
}

// Zwraca siłę dragu [N], Fd = −q·Cd·S·û_vel; 

func (r *Rocket) DragForce() vec.Vec2 {
	speed := r.Vel.Len()
	if speed == 0 {
		return vec.Vec2{}
	}
	q := 0.5 * AirDensity * speed * speed
	ux := r.Vel.X / speed
	uy := r.Vel.Y / speed
	fd := -q * r.DragCoeff() * r.DragArea
	return vec.Vec2{X: fd * ux, Y: fd * uy}
}

//Zwraca siłę liftu [N], Fl = q·Cl·S·û_perp; û_perp = (−uy, ux).
func (r *Rocket) LiftForce() vec.Vec2 {
	speed := r.Vel.Len()
	if speed == 0 {
		return vec.Vec2{}
	}
	q := 0.5 * AirDensity * speed * speed
	ux := r.Vel.X / speed
	uy := r.Vel.Y / speed
	fl := q * r.LiftCoeff() * r.DragArea
	return vec.Vec2{X: fl * (-uy), Y: fl * ux}
}

// AeroForce zwraca siłę aero [N] w osiach świata.
// Uzasadnienie modelu i źródła: docs/physics/aero.md
func (r *Rocket) AeroForce() vec.Vec2 {
	return r.DragForce().Add(r.LiftForce())
}

// Step symuluje kolejny krok symulacji.
func (r *Rocket) Step(dt float64) {
	if !r.Alive() || dt <= 0 {
		return
	}

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

	// 1a) GimbalAngle dogania GimbalCmd z limitem rate
	gimbalCmd := clamp(r.GimbalCmd, -DefaultMaxGimbalAngle, DefaultMaxGimbalAngle) //ograniczenie żądanego kąta gimbala do max i min
	maxDelta := DefaultMaxGimbalRate * dt // maksymalny przyrost kąta gimbala w jednym kroku
	delta := gimbalCmd - r.GimbalAngle // błąd kąta gimbala
	if delta > maxDelta {
		delta = maxDelta
	} else if delta < -maxDelta {
		delta = -maxDelta
	}
	r.GimbalAngle += delta
	r.GimbalAngle = clamp(r.GimbalAngle, -DefaultMaxGimbalAngle, DefaultMaxGimbalAngle)
	gimbalAngle := r.GimbalAngle

	// 1b) TVC: moment od wychylonego ciągu → obrót (α = τ/I)
	tau := thrust * r.LeverArm * math.Sin(gimbalAngle) // Moment siły od wychylonego ciągu rośnie z kątem gimbala
	alpha := tau / r.Inertia                           // rad/s² przyspieszenie kątowe w zależności od momentu siły a=τ/I
	r.AngularVel += alpha * dt
	r.Angle += r.AngularVel * dt

	// 2) ruch liniowy: siła od ciągu → przyspieszenie → prędkość → pozycja
	m := r.Mass()
	fx := thrust * math.Sin(r.Angle+gimbalAngle)
	fy := thrust * math.Cos(r.Angle+gimbalAngle)
	force := vec.Vec2{X: fx, Y: fy - m*G} // rozkład sił w zależności od kąta nachylenia
	force = force.Add(r.AeroForce())      // + aero (środek masy, bez momentu) wektorowo
	a := force.Scale(1 / m)               // a = F/m
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
