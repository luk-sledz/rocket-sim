package control
import (
	"testing"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

func TestLateralSoftCatch(t *testing.T) {
	r := rocket.NewStarship(vec.Vec2{X: 30, Y: 200}, vec.Vec2{})
	dt := 0.01
	for i := 0; i < 8000; i++ { 
		r.Throttle = AltitudeThrottle(
			r.Pos.Y, r.Vel.Y, rocket.TowerHeight,
			r.Mass(), r.MaxThrust, dt,
		)
		r.Angle = LateralAngle(r.Pos.X, r.Vel.X)
		r.Step(dt)
		if r.Caught(5, 5, 2) {
			return 
		}
	}
	if !r.Caught(5, 5, 2) {
		t.Fatalf("X=%.2f Y=%.2f Vx=%.2f Vy=%.2f", r.Pos.X, r.Pos.Y, r.Vel.X, r.Vel.Y)
	}
}