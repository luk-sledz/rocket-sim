package control
import (
	"testing"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

// testuje cały regulator poziomu
func TestLateralSoftCatch(t *testing.T) {
	r := rocket.NewStarship(vec.Vec2{X: 30, Y: 200}, vec.Vec2{})
	dt := 0.01
	for i := 0; i < 8000; i++ { 
		r.Throttle = AltitudeThrottle(
			r.Pos.Y, r.Vel.Y, rocket.TowerHeight,
			r.Mass(), r.MaxThrust, r.Angle, dt,
		)
		desiredVx := LateralDesiredVx(r.Pos.X)
		r.Angle = LateralAngleFromVx(r.Vel.X, desiredVx)
		r.Step(dt)
		if r.Caught(5, 5, 2) {
			return 
		}
	}
	if !r.Caught(5, 5, 2) {
		t.Fatalf("X=%.2f Y=%.2f Vx=%.2f Vy=%.2f", r.Pos.X, r.Pos.Y, r.Vel.X, r.Vel.Y)
	}
}
// testuje wyliczanie Vx na podstawie odległości od wieży
func TestLateralDesiredVx(t *testing.T) {
	// znak: X>0 → VxCmd < 0
	if v := LateralDesiredVx(40); v >= 0 {
		t.Fatalf("X=40 → desiredVx powinno być < 0, jest %v", v)
	}
	// saturacja: duży |X| → ±MaxLateralSpeed
	if v := LateralDesiredVx(200); v != -DefaultMaxLateralSpeed {
		t.Fatalf("X=200 → powinno saturować do %v, jest %v", -DefaultMaxLateralSpeed, v)
	}
	if v := LateralDesiredVx(-200); v != DefaultMaxLateralSpeed {
		t.Fatalf("X=-200 → powinno saturować do %v, jest %v", DefaultMaxLateralSpeed, v)
	}
}

// testuje wyliczanie kąta na podstawie błędu prędkości
func TestLateralAngleFromVx(t *testing.T) {
	// gdy Vx > VxCmd → Angle < 0
	if a := LateralAngleFromVx(10, 0); a >= 0 {
		t.Fatalf("velX=10, desired=0 → Angle < 0, jest %v", a)
	}
	// Czy saturowanie działa
	if a := LateralAngleFromVx(100, 0); a != -DefaultMaxTilt {
		t.Fatalf("duży błąd → Angle=%v, oczekiwane %v", a, -DefaultMaxTilt)
	}
}