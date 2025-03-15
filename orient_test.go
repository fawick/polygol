package polygol

import "testing"

func TestCompareVectorAngles(t *testing.T) {
	t.Run("colinear", func(t *testing.T) {
		pt1 := newVectorLit(1, 1)
		pt2 := newVectorLit(2, 2)
		pt3 := newVectorLit(3, 3)
		expect(t, orient(pt1, pt2, pt3) == 0)
		expect(t, orient(pt2, pt1, pt3) == 0)
		expect(t, orient(pt2, pt3, pt1) == 0)
		expect(t, orient(pt3, pt2, pt1) == 0)
	})
	t.Run("offset", func(t *testing.T) {
		pt1 := newVectorLit(0, 0)
		pt2 := newVectorLit(1, 1)
		pt3 := newVectorLit(1, 0)
		expect(t, orient(pt1, pt2, pt3) == 1)
		expect(t, orient(pt2, pt1, pt3) == -1)
		expect(t, orient(pt2, pt3, pt1) == 1)
		expect(t, orient(pt3, pt2, pt1) == -1)
	})
}
