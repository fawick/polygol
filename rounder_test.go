package polygol

import (
	"testing"
)

func TestRounderRound(t *testing.T) {
	// t.Parallel()

	var pt1, pt2, pt3 *point

	rounder := newPtRounder()

	// no overlap
	t.Run("no-overlap", func(t *testing.T) {
		resetPrecision()
		pt1 = newPoint(3, 4)
		pt2 = newPoint(4, 5)
		pt3 = newPoint(5, 5)
		expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
		expect(t, rounder.round(pt2.x, pt2.y).equal(*pt2))
		expect(t, rounder.round(pt3.x, pt3.y).equal(*pt3))
	})

	// exact overlap
	t.Run("exact-overlap", func(t *testing.T) {
		resetPrecision()
		pt1 = newPoint(3, 4)
		pt2 = newPoint(4, 5)
		pt3 = newPoint(3, 4)
		expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
		expect(t, rounder.round(pt2.x, pt2.y).equal(*pt2))
		expect(t, rounder.round(pt3.x, pt3.y).equal(*pt3))
	})

	// rounding one coordinate
	t.Run("rounding-one-coordinate", func(t *testing.T) {
		setPrecision(NumberEPSILON)
		pt1 = newPoint(3, 4)
		pt2 = &point{
			Vector: Vector{
				x: newBigNumber(3).plus(newBigNumber(NumberEPSILON)),
				y: newBigNumber(4),
			},
		}
		pt3 = &point{
			Vector: Vector{
				x: newBigNumber(3),
				y: newBigNumber(4).plus(newBigNumber(NumberEPSILON)),
			},
		}
		expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
		expect(t, rounder.round(pt2.x, pt2.y).equal(*pt1))
		expect(t, rounder.round(pt3.x, pt3.y).equal(*pt1))
		resetPrecision()
	})

	// rounding both coordinates
	t.Run("rounding-both-coordinates", func(t *testing.T) {
		setPrecision(NumberEPSILON)
		pt1 = newPoint(3, 4)
		pt2 = &point{
			Vector: Vector{
				x: newBigNumber(3).plus(newBigNumber(NumberEPSILON)),
				y: newBigNumber(4).plus(newBigNumber(NumberEPSILON)),
			},
		}
		expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
		expect(t, rounder.round(pt2.x, pt2.y).equal(*pt1))
		resetPrecision()
	})

	// preseed with 0
	t.Run("preseed-with-zero", func(t *testing.T) {
		setPrecision(NumberEPSILON)
		pt1 = &point{
			Vector: Vector{
				x: newBigNumber(NumberEPSILON).div(newBigNumber(2)),
				y: newBigNumber(-NumberEPSILON).div(newBigNumber(2)),
			},
		}
		expect(t, !pt1.x.isZero())
		expect(t, !pt1.y.isZero())
		expect(t, rounder.round(pt1.x, pt1.y).equal(*newPoint(0, 0)))
		resetPrecision()
	})
}
