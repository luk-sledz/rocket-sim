package control
import (
	"testing"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)
func TestAltitudeSoftCatch(t *testing.T) {
	r := rocket.NewStarship(vec.Vec2{Y: 200}, vec.Vec2{})
	dt := 0.01
	for i := 0; i < 4000; i++ { // 40 s
		r.Throttle = AltitudeThrottle(r.Pos.Y, r.Vel.Y, rocket.TowerHeight, r.Mass(), r.MaxThrust, r.Angle, dt)
		r.Step(dt)
	}
	if !r.Caught(5, 5, 2) {
		t.Fatalf("Y=%.1f Vel.Y=%.2f", r.Pos.Y, r.Vel.Y)
	}
}