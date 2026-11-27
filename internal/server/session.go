package server

import (
	"sync"
	"time"

	"rocket-sim/internal/control"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

// Session to jedna symulacja: rakieta, czas i flagi, chronione mutexem.
type Session struct {
	mu      sync.Mutex
	r       *rocket.Rocket
	t       float64
	running bool
	// catching: w oknie chwytu (ramiona się zamykają); fizyka i regulatory dalej działają.
	// hold = czas ciągle w oknie; po CatchHoldTime → caught i stop.
	catching bool
	hold     float64
	caught   bool
	crashed  bool
}

// newRocket tworzy rakietę z tymi samymi warunkami początkowymi co CLI.
func newRocket() *rocket.Rocket {
	pos := vec.Vec2{X: -100, Y: 1000}
	vel := vec.Vec2{Y: -100}
	return rocket.NewStarship(pos, vel)
}

// NewSession tworzy zapauzowaną sesję.
func NewSession() *Session {
	return &Session{r: newRocket()}
}

// Start wznawia symulację, chyba że już się skończyła (caught/crashed).
func (s *Session) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.caught && !s.crashed {
		s.running = true
	}
}

// Pause zatrzymuje symulację bez zmiany stanu rakiety.
func (s *Session) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// Reset wraca do warunków początkowych i pauzy.
func (s *Session) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.r = newRocket()
	s.t = 0
	s.running = false
	s.catching = false
	s.hold = 0
	s.caught = false
	s.crashed = false
}

// State zwraca kopię stanu do wysłania do przeglądarki.
func (s *Session) State() StateMsg {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Snapshot(s.r, s.t, s.catching, s.caught, s.crashed)
}

// Tick wykonuje jeden krok symulacji o długości dt [s] i aktualizuje stan.
func (s *Session) Tick(dt float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.caught {
		return
	}

	r := s.r
	// Regulatory i Step działają aż do caught — także podczas zamykania ramion.
	r.Throttle = control.AltitudeThrottle(
		r.Pos.Y, r.Vel.Y, rocket.TowerHeight,
		r.Mass(), r.MaxThrust, r.Angle, dt,
	)
	desiredVx := control.LateralDesiredVx(r.Pos.X)
	r.AngleCmd = control.LateralAngleFromVx(r.Vel.X, desiredVx)
	r.GimbalCmd = control.AttitudeGimbal(r.Angle, r.AngleCmd, r.AngularVel)
	r.Step(dt)
	s.t += dt

	if !r.Alive() {
		s.crashed = true
		s.running = false
		s.catching = false
		s.hold = 0
		return
	}

	// 3 s ciągle w oknie → caught (stop). Wypadnięcie z okna zeruje licznik.
	if r.Caught(rocket.DefaultCatchTolX, rocket.DefaultCatchTolY, rocket.DefaultCatchMaxSpeed) {
		s.catching = true
		s.hold += dt
		if s.hold >= rocket.CatchHoldTime {
			s.caught = true
			s.running = false
		}
	} else {
		s.catching = false
		s.hold = 0
	}
}

// Run wywołuje Tick co dt
func (s *Session) Run(dt float64) {
	ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
	defer ticker.Stop()
	for range ticker.C {
		s.Tick(dt)
	}
}
