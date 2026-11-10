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
	r.AngleCmd = math.Pi / 2
	for i := 0; i < 100; i++ {
		r.Step(0.01)
	}
	if r.Vel.X <= 0 {
		t.Errorf("przy Angle=π/2 Vel.X powinno być > 0, jest %v", r.Vel.X)
	}
}

// testuje, czy kąt rakiety dogoni żądany kąt
func TestStepAngleTracksCmd(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Angle = 0
	r.AngleCmd = 0.1 // małe żądanie — poniżej limitu rate
	for i := 0; i < 300; i++ { // 3 s
		r.Step(0.01)
	}
	err := math.Abs(r.Angle - r.AngleCmd)
	if err > 0.02 {
		t.Fatalf("Angle nie dogonił AngleCmd: Angle=%v AngleCmd=%v err=%v",
			r.Angle, r.AngleCmd, err)
	}
	if math.Abs(r.AngularVel) > 0.05 {
		t.Fatalf("przy celu AngularVel powinno być małe, jest %v", r.AngularVel)
	}
}

// testuje, czy prędkość kątowa rakiety jest ograniczona
func TestStepAngularRateLimit(t *testing.T) {
	r := NewStarship(vec.Vec2{Y: 100}, vec.Vec2{})
	r.Angle = 0
	r.AngleCmd = 1.0 // duży skok → omegaCmd = 1/0.2 = 5 > MaxRate
	r.Step(0.01)
	if math.Abs(r.AngularVel) > DefaultMaxAngularRate+1e-9 {
		t.Fatalf("|AngularVel|=%v przekroczyło MaxRate=%v",
			r.AngularVel, DefaultMaxAngularRate)
	}
	// przy takim skoku powinno siedzieć na limicie
	if math.Abs(r.AngularVel-DefaultMaxAngularRate) > 1e-9 {
		t.Fatalf("oczekiwane AngularVel=%v, jest %v",
			DefaultMaxAngularRate, r.AngularVel)
	}
}