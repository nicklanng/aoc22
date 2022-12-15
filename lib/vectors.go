package lib

import (
	"fmt"
	"math"
)

type Vector2i struct{ X, Y int }

func (v Vector2i) ToIndex(width int) int             { return v.Y*width + v.X }
func (v Vector2i) Add(v2 Vector2i) Vector2i          { return Vector2i{X: v.X + v2.X, Y: v.Y + v2.Y} }
func (v Vector2i) Sub(v2 Vector2i) Vector2i          { return Vector2i{X: v.X - v2.X, Y: v.Y - v2.Y} }
func (v Vector2i) Magnitude() float64                { return math.Sqrt(float64(v.X*v.X + v.Y*v.Y)) }
func (v Vector2i) Distance(v2 Vector2i) float64      { return v.Sub(v2).Magnitude() }
func (v Vector2i) ManhattanDistance(v2 Vector2i) int { return IntAbs(v2.X-v.X) + IntAbs(v2.Y-v.Y) }
func (v Vector2i) UnitNormalize() Vector2i {
	if v.X != 0 {
		v.X = v.X / IntAbs(v.X)
	}
	if v.Y != 0 {
		v.Y = v.Y / IntAbs(v.Y)
	}
	return v
}
func (v *Vector2i) ParseString(s string) error {
	_, err := fmt.Sscanf(s, "%d,%d", &v.X, &v.Y)
	return err
}
