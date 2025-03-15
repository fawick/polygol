package polygol

import "fmt"

type point struct {
	Vector
	events []*sweepEvent
}

func newPoint(x, y float64) *point {
	return &point{
		Vector: Vector{
			x: newBigNumber(x),
			y: newBigNumber(y),
		},
	}
}

func newPointBN(x, y BigNumber) *point {
	return &point{
		Vector: Vector{x: x, y: y},
	}
}

// TODO rename to vector
func (p point) vector() Vector {
	return p.Vector
}

func (p point) equal(other point) bool {
	return p.x.equalTo(other.x) && p.y.equalTo(other.y)
}

func (p point) String() string {
	return fmt.Sprintf("(%s, %s)", p.x, p.y)
}
