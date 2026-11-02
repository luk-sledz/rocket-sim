// Package vec zawiera podstawowe operacje na wektorach 2D.
//
// Wszystkie wartości są w jednostkach SI - znaczenie wektora zależy od kontekstu
// (pozycja [m], prędkość [m/s], przyspieszenie [m/s^2], siła [N]).
//
// Konwencja osi: X - poziomo, Y - w górę (wysokość nad ziemią).
package vec

import "math"

// Vec2 to wektor dwuwymiarowy.
//
// Typ jest używany przez wartość (nie przez wskaźnik) - metody nie modyfikują
// odbiorcy, tylko zwracają nowy wektor.
type Vec2 struct {
	X, Y float64
}

// Add zwraca sumę wektorów v + o.
func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{X: v.X + o.X, Y: v.Y + o.Y}
}

// Sub zwraca różnicę wektorów v - o.
func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{X: v.X - o.X, Y: v.Y - o.Y}
}

// Scale zwraca wektor pomnożony przez skalar k.
func (v Vec2) Scale(k float64) Vec2 {
	return Vec2{X: v.X * k, Y: v.Y * k}
}

// Len zwraca długość (normę euklidesową) wektora.
func (v Vec2) Len() float64 {
	return math.Hypot(v.X, v.Y)
}
