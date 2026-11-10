package main

import (
	"fmt"
	"rocket-sim/internal/control"
	"rocket-sim/internal/rocket"
	"rocket-sim/internal/vec"
)

func main() {
	// WSTĘPNE WARUNKI:
	pos := vec.Vec2{X: 100, Y: 1000}
	vel := vec.Vec2{Y: -100}
	r := rocket.NewStarship(pos, vel)
	dt := 0.01
	t := 0.0
	maxT := 120.0 // limit czasu [s]
	for t < maxT {
		r.Throttle = control.AltitudeThrottle(
			r.Pos.Y, r.Vel.Y, rocket.TowerHeight,
			r.Mass(), r.MaxThrust, r.Angle, dt,
		)
		desiredVx := control.LateralDesiredVx(r.Pos.X)
		r.AngleCmd = control.LateralAngleFromVx(r.Vel.X, desiredVx)
		r.GimbalCmd = control.AttitudeGimbal(r.Angle, r.AngleCmd, r.AngularVel)
		r.Step(dt)
		t += dt
		// log co 0.5 s
		step := int(t/dt + 0.5)
		if step%50 == 0 {
			fmt.Printf("t=%5.1fs  X=%7.1f  Y=%7.1f  Vx=%6.2f  Vy=%6.2f  thr=%.2f  fuel=%.0f  angle=%.2f  gimbal=%.2f\n",
				t, r.Pos.X, r.Pos.Y, r.Vel.X, r.Vel.Y, r.Throttle, r.Fuel, r.Angle, r.GimbalAngle)
		}
		if r.Caught(rocket.DefaultCatchTolX, rocket.DefaultCatchTolY, rocket.DefaultCatchMaxSpeed) {
			fmt.Printf("CAUGHT @ t=%.1fs  X=%.2f  Y=%.2f  Vx=%.2f  Vy=%.2f  angle=%.2f\n", t, r.Pos.X, r.Pos.Y, r.Vel.X, r.Vel.Y, r.Angle)
			return
		}
		if !r.Alive() {
			fmt.Println("CRASH (Y=0)")
			return
		}
	}
	fmt.Printf("TIMEOUT  X=%.1f  Y=%.1f  Vx=%.2f  Vy=%.2f  angle=%.2f\n", r.Pos.X, r.Pos.Y, r.Vel.X, r.Vel.Y, r.Angle)
}