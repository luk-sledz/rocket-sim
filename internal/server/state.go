package server

import "rocket-sim/internal/rocket"

// StateMsg to struktura przechowująca stan symulacji.
type StateMsg struct {
	T        float64 `json:"t"`  // czas symulacji [s]
	X        float64 `json:"x"`  // pozycja [m]
	Y        float64 `json:"y"`  // wysokość [m]
	Vx       float64 `json:"vx"` // prędkość [m/s]
	Vy       float64 `json:"vy"`
	Angle    float64 `json:"angle"`    // nachylenie [rad], 0 = pionowo
	Gimbal   float64 `json:"gimbal"`   // kąt gimbala [rad]
	Throttle float64 `json:"throttle"` // 0..1
	Fuel     float64 `json:"fuel"`     // paliwo [kg]
	Aero     float64 `json:"aero"`     // |siła aero| [N]
	Caught   bool    `json:"caught"`
	Crashed  bool    `json:"crashed"`
}

// Snapshot wysyła stan symulacji do przeglądarki.
func Snapshot(r *rocket.Rocket, t float64, caught, crashed bool) StateMsg {
	return StateMsg{
		T:        t,
		X:        r.Pos.X,
		Y:        r.Pos.Y,
		Vx:       r.Vel.X,
		Vy:       r.Vel.Y,
		Angle:    r.Angle,
		Gimbal:   r.GimbalAngle,
		Throttle: r.Throttle,
		Fuel:     r.Fuel,
		Aero:     r.AeroForce().Len(),
		Caught:   caught,
		Crashed:  crashed,
	}
}
