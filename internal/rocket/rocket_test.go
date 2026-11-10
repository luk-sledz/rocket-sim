package rocket

import (
	"testing"
	"rocket-sim/internal/vec"
	"math"
)
// TestStepFreeFall testuje spadek swobodny rakiety
func TestStepFreeFall(t *testing.T) {
	r := New(vec.Vec2{Y: 100}, vec.Vec2{}, DefaultDryMass, DefaultFuel, DefaultMaxThrust, DefaultIsp)
	for i := 0; i < 100; i++ {
		r.Step(0.01) // łącznie 1 s
	}
	// po 1 s: Vel.Y ≈ -9.81, Pos.Y ≈ 100 - 0.5*9.81
	if r.Vel.Y >= 0 {
		t.Errorf("przy spadku Vel.Y powinno być < 0, jest %v", r.Vel.Y)
	}
}

// TestStepThrustUp testuje przyspieszanie rakiety w górę
func TestStepThrustUp(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Throttle = 1
	for i := 0; i < 100; i++ {
		r.Step(0.01) // 1 s
	}
	if r.Vel.Y <= 0 {
		t.Errorf("przy pełnym ciągu Vel.Y powinno być > 0, jest %v", r.Vel.Y)
	}
}

// TestCaught testuje czy rakieta jest w oknie chwytu ramion wieży
func TestCaught(t *testing.T) {
	// OK: przy ramionach, prawie stoi
	r := NewStarship(vec.Vec2{Y: TowerHeight}, vec.Vec2{})
	if !r.Caught(2, 2, 1) {
		t.Fatal("Y=90, Vel=0 → powinno być Caught")
	}
	// za nisko
	r.Pos.Y = 50
	r.Vel = vec.Vec2{}
	if r.Caught(2, 2, 1) {
		t.Fatal("Y=50 → nie powinno być Caught")
	}
	// dobra wysokość, za duża prędkość
	r.Pos.Y = TowerHeight
	r.Vel = vec.Vec2{Y: -5} // |−5| > maxSpeed=1
	if r.Caught(2, 2, 1) {
		t.Fatal("za szybki przelot → nie powinno być Caught")
	}
}

func TestStepThrustAngled(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Throttle = 1
	r.Angle = math.Pi / 2 //ciąg w +X
	r.GimbalAngle = 0
	r.AngularVel = 0
	for i := 0; i < 100; i++ {
		r.Step(0.01)
	}
	if r.Vel.X <= 0 {
		t.Errorf("przy Angle=π/2 Vel.X powinno być > 0, jest %v", r.Vel.X)
	}
}

//testuje znak momentu siły od wychylonego ciągu
func TestStepGimbalTorqueSign(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Throttle = 1
	r.Angle = 0
	r.AngularVel = 0
	r.GimbalCmd = 0.1
	r.GimbalAngle = 0.1
	for i := 0; i < 50; i++ {
		r.Step(0.01)
	}
	if r.AngularVel <= 0 || r.Angle <= 0 {
		t.Fatalf("+gimbal powinno dać +ω i +Angle; ω=%v Angle=%v",
			r.AngularVel, r.Angle)
	}
}

//testuje czy rakieta nie obraca się gdy gimbal jest 0
func TestStepZeroGimbalNoTorque(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Throttle = 1
	r.Angle = 0.05
	r.AngularVel = 0
	angle0 := r.Angle
	r.GimbalCmd = 0
	r.GimbalAngle = 0
	r.Step(0.01)
	if math.Abs(r.Angle-angle0) > 1e-12 || math.Abs(r.AngularVel) > 1e-12 {
		t.Fatalf("gimbal=0, ω=0 → brak obrotu; Angle=%v ω=%v", r.Angle, r.AngularVel)
	}
}

// testuje czy przesunięcie gimbala jest ograniczone 
func TestStepGimbalRateLimit(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.GimbalAngle = 0
	r.GimbalCmd = 0.20 // skok na max
	dt := 0.01
	r.Step(dt)
	maxDelta := DefaultMaxGimbalRate * dt // 0.01
	if math.Abs(r.GimbalAngle-maxDelta) > 1e-12 { // 
		t.Fatalf("oczekiwane GimbalAngle=%v, jest %v", maxDelta, r.GimbalAngle)
	}
}