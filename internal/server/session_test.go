package server

import (
	"testing"

	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

func TestTelemetryMovesDuringCatching(t *testing.T) {
	s := NewSession()
	// Ustaw rakietę już w oknie chwytu, wolną.
	s.r = rocket.NewStarship(vec.Vec2{X: 0, Y: rocket.TowerHeight}, vec.Vec2{})
	s.running = true

	s0 := s.State()
	dt := 0.005
	for i := 0; i < 100; i++ { // 0.5 s
		s.Tick(dt)
	}
	s1 := s.State()

	if s1.Caught {
		t.Fatal("po 0.5 s nie powinno być jeszcze caught (hold=3s)")
	}
	if !s1.Catching {
		t.Fatal("powinno być catching")
	}
	if s1.T <= s0.T {
		t.Fatalf("czas stoi: t0=%v t1=%v", s0.T, s1.T)
	}
	// paliwo albo throttle powinny żyć (hover zużywa paliwo)
	if s1.Fuel >= s0.Fuel && s1.Throttle == 0 {
		t.Fatalf("brak oznak działania regulatora: fuel %v→%v thr=%v", s0.Fuel, s1.Fuel, s1.Throttle)
	}
}
