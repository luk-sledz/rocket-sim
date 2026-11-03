package main

import (
	"fmt"
	"rocket-sim/internal/control"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

func main() {
	// start wyżej, lekko w dół
	pos := vec.Vec2{Y: 1000}
	vel := vec.Vec2{Y: -100}
	r := rocket.NewStarship(pos, vel)
	dt := 0.01
	t := 0.0
	maxT := 120.0 // limit czasu [s]
	for t < maxT {
		r.Throttle = control.AltitudeThrottle(
			r.Pos.Y, r.Vel.Y, rocket.TowerHeight,
			r.Mass(), r.MaxThrust, dt,
		)
		r.Step(dt)
		t += dt
		// log co 0.5 s
		step := int(t/dt + 0.5)
		if step%50 == 0 {
			fmt.Printf("t=%5.1fs  Y=%7.1f  Vy=%6.2f  thr=%.2f  fuel=%.0f\n",
				t, r.Pos.Y, r.Vel.Y, r.Throttle, r.Fuel)
		}
		if r.Caught(rocket.DefaultCatchTolY, rocket.DefaultCatchMaxSpeed) {
			fmt.Printf("CAUGHT @ t=%.1fs  Y=%.2f  Vy=%.2f\n", t, r.Pos.Y, r.Vel.Y)
			return
		}
		if !r.Alive() {
			fmt.Println("CRASH (Y=0)")
			return
		}
	}
	fmt.Printf("TIMEOUT  Y=%.1f  Vy=%.2f\n", r.Pos.Y, r.Vel.Y)
}