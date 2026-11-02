package vec

import (
	"math"
	"testing"
)

func TestVec2Add(t *testing.T) {
	tests := []struct {
		name string
		a, b Vec2
		want Vec2
	}{
		{"zera", Vec2{0, 0}, Vec2{0, 0}, Vec2{0, 0}},
		{"dodatnie", Vec2{1, 2}, Vec2{3, 4}, Vec2{4, 6}},
		{"ujemne", Vec2{1, 2}, Vec2{-1, -2}, Vec2{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Add(tt.b); got != tt.want {
				t.Errorf("%v.Add(%v) = %v, chcemy %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestVec2Sub(t *testing.T) {
	got := Vec2{5, 7}.Sub(Vec2{2, 3})
	want := Vec2{3, 4}
	if got != want {
		t.Errorf("Sub = %v, chcemy %v", got, want)
	}
}

func TestVec2Scale(t *testing.T) {
	tests := []struct {
		name string
		v    Vec2
		k    float64
		want Vec2
	}{
		{"podwojenie", Vec2{1, -2}, 2, Vec2{2, -4}},
		{"zero", Vec2{1, -2}, 0, Vec2{0, 0}},
		{"odwrocenie", Vec2{1, -2}, -1, Vec2{-1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Scale(tt.k); got != tt.want {
				t.Errorf("%v.Scale(%v) = %v, chcemy %v", tt.v, tt.k, got, tt.want)
			}
		})
	}
}

func TestVec2Len(t *testing.T) {
	tests := []struct {
		name string
		v    Vec2
		want float64
	}{
		{"zero", Vec2{0, 0}, 0},
		{"trojkat 3-4-5", Vec2{3, 4}, 5},
		{"ujemne skladowe", Vec2{-3, -4}, 5},
		{"os Y", Vec2{0, -2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Len(); math.Abs(got-tt.want) > 1e-12 {
				t.Errorf("%v.Len() = %v, chcemy %v", tt.v, got, tt.want)
			}
		})
	}
}
