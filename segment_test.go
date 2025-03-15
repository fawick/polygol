package polygol

import (
	"testing"
)

func equalRingIn(r1 []*ringIn, r2 []*ringIn) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}
	return true
}

func equalWinding(w1, w2 []int) bool {
	if len(w1) != len(w2) {
		return false
	}
	for i := range w1 {
		if w1[i] != w2[i] {
			return false
		}
	}
	return true
}

func equalVector(v1, v2 Vector) bool {
	return v1.x.equalTo(v2.x) && v1.y.equalTo(v2.y)
}

func TestSegmentNew(t *testing.T) {
	// t.Parallel()

	var leftSE, rightSE *sweepEvent

	op := newOperation("")

	// general
	leftSE = newSweepEvent(newPoint(0, 0), true)
	rightSE = newSweepEvent(newPoint(1, 1), false)
	rings := []*ringIn{}
	windings := []int{}
	seg := op.newSegment(leftSE, rightSE, rings, windings)

	expect(t, equalRingIn(seg.rings, rings))
	expect(t, equalWinding(seg.windings, windings))
	expect(t, seg.leftSE == leftSE)
	expect(t, seg.leftSE.otherSE == rightSE)
	expect(t, seg.rightSE == rightSE)
	expect(t, seg.rightSE.otherSE == leftSE)
	expect(t, seg.ringOut == nil)
	expect(t, seg.prev == nil)
	expect(t, seg.consumedBy == nil)

	// segment Id increments
	leftSE = newSweepEvent(newPoint(0, 0), true)
	rightSE = newSweepEvent(newPoint(1, 1), false)
	seg1 := op.newSegment(leftSE, rightSE, []*ringIn{}, nil)
	seg2 := op.newSegment(leftSE, rightSE, []*ringIn{}, nil)
	expect(t, seg2.id-seg1.id == 1)
}

func TestSegmentNewFromRing(t *testing.T) {
	// t.Parallel()

	var p1, p2 *point
	var seg *segment
	var err error

	op := newOperation("")

	// correct point on left and right 1
	p1, p2 = newPoint(0, 0), newPoint(0, 1)
	seg, err = op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)

	expect(t, seg.leftSE.point.equal(*p1))
	expect(t, seg.rightSE.point.equal(*p2))

	// correct point on left and right 1
	p1, p2 = newPoint(0, 0), newPoint(-1, 0)
	seg, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	expect(t, seg.leftSE.point.equal(*p2))
	expect(t, seg.rightSE.point.equal(*p1))

	// attempt create segment with same points
	p1, p2 = newPoint(0, 0), newPoint(0, 0)
	_, err = op.newSegmentFromRing(p1, p2, nil)
	expect(t, err != nil)
}

func TestSegmentSplit(t *testing.T) {
	// t.Parallel()

	var seg *segment
	var pt *point
	var evt, otherEvt *sweepEvent
	var evts []*sweepEvent
	var err error

	op := newOperation("")

	// on interior point
	seg, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(10, 10), nil)
	terr(t, err)

	pt = newPoint(5, 5)
	evts = seg.split(pt)
	expect(t, evts[0].segment == seg)
	expect(t, evts[0].point.equal(*pt))
	expect(t, !evts[0].isLeft)
	expect(t, evts[0].otherSE.otherSE == evts[0])
	expect(t, evts[1].segment.leftSE.segment == evts[1].segment)
	expect(t, evts[1].segment != seg)
	expect(t, evts[1].point.equal(*pt))
	expect(t, evts[1].isLeft)
	expect(t, evts[1].otherSE.otherSE == evts[1])
	expect(t, evts[1].segment.rightSE.segment == evts[1].segment)

	// on close-to-but-not-exactly interior point
	seg, err = op.newSegmentFromRing(newPoint(0, 10), newPoint(10, 0), nil)
	terr(t, err)

	pt = &point{
		Vector: Vector{
			x: newBigNumber(5).plus(newBigNumber(NumberEPSILON)),
			y: newBigNumber(5),
		},
	}
	evts = seg.split(pt)
	expect(t, evts[0].segment == seg)
	expect(t, evts[0].point.equal(*pt))
	expect(t, !evts[0].isLeft)
	expect(t, evts[1].segment != seg)
	expect(t, evts[1].point.equal(*pt))
	expect(t, evts[1].isLeft)
	expect(t, evts[1].segment.rightSE.segment == evts[1].segment)

	// on three interior points
	seg, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(10, 0), nil)
	terr(t, err)

	sPt1, sPt2, sPt3 := newPoint(2, 2), newPoint(4, 4), newPoint(6, 6)
	orgLeftEvt, orgRightEvt := seg.leftSE, seg.rightSE
	newEvts3 := seg.split(sPt3)
	newEvts2 := seg.split(sPt2)
	newEvts1 := seg.split(sPt1)
	newEvts := append(newEvts1, newEvts2...)
	newEvts = append(newEvts, newEvts3...)

	expect(t, len(newEvts) == 6)

	expect(t, seg.leftSE == orgLeftEvt)
	for _, e := range newEvts {
		if e.point == sPt1 && !e.isLeft {
			evt = e
			break
		}
	}
	expect(t, seg.rightSE == evt)

	for _, e := range newEvts {
		if e.point == sPt1 && e.isLeft {
			evt = e
			break
		}
	}
	for _, e := range newEvts {
		if e.point == sPt2 && !e.isLeft {
			otherEvt = e
			break
		}
	}
	expect(t, evt.segment == otherEvt.segment)

	for _, e := range newEvts {
		if e.point == sPt2 && e.isLeft {
			evt = e
			break
		}
	}
	for _, e := range newEvts {
		if e.point.equal(*sPt3) && !e.isLeft {
			otherEvt = e
			break
		}
	}
	expect(t, evt.segment == otherEvt.segment)

	for _, e := range newEvts {
		if e.point.equal(*sPt3) && e.isLeft {
			evt = e
			break
		}
	}
	expect(t, evt.segment == orgRightEvt.segment)
}
func TestSegmentSimplePropertiesBboxVector(t *testing.T) {
	// t.Parallel()

	var seg *segment
	var err error

	op := newOperation("")

	// general
	seg, err = op.newSegmentFromRing(newPoint(1, 2), newPoint(3, 4), nil)
	terr(t, err)

	expect(t, equalBbox(seg.bbox(), Bbox{ll: newVectorLit(1, 2), ur: newVectorLit(3, 4)}))
	expect(t, equalVector(seg.vector(), newVectorLit(2, 2)))

	// horizontal
	seg, err = op.newSegmentFromRing(newPoint(1, 4), newPoint(3, 4), nil)
	terr(t, err)

	expect(t, equalBbox(seg.bbox(), Bbox{ll: newVectorLit(1, 4), ur: newVectorLit(3, 4)}))
	expect(t, equalVector(seg.vector(), newVectorLit(2, 0)))

	// vertical
	seg, err = op.newSegmentFromRing(newPoint(3, 2), newPoint(3, 4), nil)
	terr(t, err)

	expect(t, equalBbox(seg.bbox(), Bbox{ll: newVectorLit(3, 2), ur: newVectorLit(3, 4)}))
	expect(t, equalVector(seg.vector(), newVectorLit(0, 2)))
}
func TestSegmentConsume(t *testing.T) {
	// t.Parallel()

	var p1, p2 *point
	var seg1, seg2 *segment
	var err error

	op := newOperation("")

	// not automatically consumed
	p1, p2 = newPoint(0, 0), newPoint(1, 0)
	seg1, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	seg2, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	expect(t, seg1.consumedBy == nil)
	expect(t, seg2.consumedBy == nil)

	// basic case
	p1, p2 = newPoint(0, 0), newPoint(1, 0)
	seg1, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	seg2, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	seg1.consume(seg2)
	expect(t, seg2.consumedBy == seg1)
	expect(t, seg1.consumedBy == nil)

	// earlier in sweep line  sorting consumes later
	p1, p2 = newPoint(0, 0), newPoint(1, 0)
	seg1, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	seg2, err = op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	seg2.consume(seg1)
	expect(t, seg2.consumedBy == seg1)
	expect(t, seg1.consumedBy == nil)

	// consuming cascades
	p1, p2 = newPoint(0, 0), newPoint(0, 0)
	p3, p4 := newPoint(1, 0), newPoint(1, 0)
	seg1, err = op.newSegmentFromRing(p1, p3, nil)
	terr(t, err)

	seg2, err = op.newSegmentFromRing(p1, p3, nil)
	terr(t, err)

	seg3, err := op.newSegmentFromRing(p2, p4, nil)
	terr(t, err)

	seg4, err := op.newSegmentFromRing(p2, p4, nil)
	terr(t, err)

	seg5, err := op.newSegmentFromRing(p2, p4, nil)
	terr(t, err)

	seg1.consume(seg2)
	seg4.consume(seg2)
	seg3.consume(seg2)
	seg3.consume(seg5)
	expect(t, seg1.consumedBy == nil)
	expect(t, seg2.consumedBy == seg1)
	expect(t, seg3.consumedBy == seg1)
	expect(t, seg4.consumedBy == seg1)
	expect(t, seg5.consumedBy == seg1)
}
func TestSegmentIsAnEndpoint(t *testing.T) {
	// t.Parallel()

	op := newOperation("")

	p1, p2 := newPoint(0, -1), newPoint(1, 0)
	seg, err := op.newSegmentFromRing(p1, p2, nil)
	terr(t, err)

	// yup
	expect(t, seg.isAnEndpoint(p1))
	expect(t, seg.isAnEndpoint(p2))

	// nope
	expect(t, !seg.isAnEndpoint(newPoint(-34, 46)))
	expect(t, !seg.isAnEndpoint(newPoint(0, 0)))
}
func TestSegmentComparisonWithPoint(t *testing.T) {
	// t.Parallel()

	var seg *segment
	var err error
	var pt *point

	op := newOperation("")

	t.Run("general", func(t *testing.T) {
		seg1, err := op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
		terr(t, err)

		seg2, err := op.newSegmentFromRing(newPoint(0, 1), newPoint(0, 0), nil)
		terr(t, err)

		expect(t, seg1.comparePoint(newPoint(0, 1)) == 1)
		expect(t, seg1.comparePoint(newPoint(1, 2)) == 1)
		expect(t, seg1.comparePoint(newPoint(0, 0)) == 0)
		expect(t, seg1.comparePoint(newPoint(5, -1)) == -1)

		expect(t, seg2.comparePoint(newPoint(0, 1)) == 0)
		expect(t, seg2.comparePoint(newPoint(1, 2)) == -1)
		expect(t, seg2.comparePoint(newPoint(0, 0)) == 0)
		expect(t, seg2.comparePoint(newPoint(5, -1)) == -1)
	})

	// barely above
	t.Run("barely-above", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
		terr(t, err)
		pt = &point{
			Vector: Vector{
				x: newBigNumber(2),
				y: newBigNumber(1).minus(newBigNumber(NumberEPSILON)),
			},
		}
		expect(t, seg.comparePoint(pt) == -1)
	})

	// barely above
	t.Run("barely-below", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
		terr(t, err)
		pt = &point{
			Vector: Vector{
				x: newBigNumber(2),
				y: newBigNumber(1).plus(newBigNumber(NumberEPSILON).times(newBigNumber(3)).div(newBigNumber(2))),
			},
		}
		expect(t, seg.comparePoint(pt) == 1)
	})

	// vertical before
	t.Run("vertical-before", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(1, 3), nil)
		terr(t, err)
		pt = newPoint(0, 0)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// vertical after
	t.Run("vertical-after", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(1, 3), nil)
		terr(t, err)
		pt = newPoint(2, 0)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// vertical on
	t.Run("vertical-on", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(1, 3), nil)
		terr(t, err)
		pt = newPoint(1, 0)
		expect(t, seg.comparePoint(pt) == 0)
	})

	// horizontal below
	t.Run("horizontal-below", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(0, 0)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// horizontal above
	t.Run("horizontal-above", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(0, 2)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// horizontal on
	t.Run("horizontal-on", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(0, 1)
		expect(t, seg.comparePoint(pt) == 0)
	})

	// in vertical plane below
	t.Run("in-vertical-plane-below", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 3), nil)
		terr(t, err)
		pt = newPoint(2, 0)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// in vertical plane above
	t.Run("in-vertical-plane-above", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 3), nil)
		terr(t, err)
		pt = newPoint(2, 4)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// in horizontal plane upward sloping before
	t.Run("in-horizontal-plane-upward-sloping-before", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 3), nil)
		terr(t, err)
		pt = newPoint(0, 2)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// in horizontal plane upward sloping after
	t.Run("in-horizontal-plane-upward-sloping-after", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 3), nil)
		terr(t, err)
		pt = newPoint(4, 2)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// in horizontal plane downward sloping below
	t.Run("in-horizontal-plane-downward-sloping-below", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 3), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(0, 2)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// in horizontal plane downward sloping after
	t.Run("in-horizontal-plane-downward-sloping-after", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 3), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(4, 2)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// upward more vertical before
	t.Run("upward-more-vertical-before", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 6), nil)
		terr(t, err)
		pt = newPoint(0, 2)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// upward more vertical after
	t.Run("upward-more-vertical-after", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 6), nil)
		terr(t, err)
		pt = newPoint(4, 2)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// downward more vertical before
	t.Run("downward-more-vertical-before", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 6), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(0, 2)
		expect(t, seg.comparePoint(pt) == -1)
	})

	// downward more vertical after
	t.Run("downward-more-vertical-after", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(1, 6), newPoint(3, 1), nil)
		terr(t, err)
		pt = newPoint(4, 2)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// downward-sloping segment with almost touching point - from issue 37
	t.Run("downward-sloping-segment-with-almost-touching-point", func(t *testing.T) {
		seg, err = op.newSegmentFromRing(newPoint(0.523985, 51.281651), newPoint(0.5241, 51.281651000100005), nil)
		terr(t, err)
		pt = newPoint(0.5239850000000027, 51.281651000000004)
		expect(t, seg.comparePoint(pt) == 1)
	})

	// avoid splitting loops on near vertical segments - from issue 60-2
	t.Run("avoid-splitting-loops-on-near-vertical-segments", func(t *testing.T) {
		setPrecision(NumberEPSILON)
		seg, err = op.newSegmentFromRing(newPoint(-45.3269382, -1.4059341), newPoint(-45.326737413921656, -1.40635), nil)
		terr(t, err)
		pt = newPoint(-45.326833968900424, -1.40615)
		expect(t, seg.comparePoint(pt) == 0)
		resetPrecision()
	})
}
func TestSegmentgetIntersections2(t *testing.T) {
	// t.Parallel()

	var seg1, seg2, s3 *segment
	var err error
	var inter *point

	op := newOperation("")

	// colinear full overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// colinear partial overlap upward slope
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(2, 2), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 3), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear partial overlap downward slope
	seg1, err = op.newSegmentFromRing(newPoint(0, 2), newPoint(2, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, 3), newPoint(1, 1), nil)
	terr(t, err)
	inter = newPoint(0, 2)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear partial overlap horizontal
	seg1, err = op.newSegmentFromRing(newPoint(0, 1), newPoint(2, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(3, 1), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear partial overlap vertical
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(0, 3), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 2), newPoint(0, 4), nil)
	terr(t, err)
	inter = newPoint(0, 2)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear endpoint overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(2, 2), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// colinear no overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 3), newPoint(4, 4), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// parallel no overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 3), newPoint(1, 4), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// intersect general
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(2, 2), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 2), newPoint(2, 0), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// T-intersect with an endpoint
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(2, 2), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(5, 4), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// intersect with vertical
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 5), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 0), newPoint(3, 44), nil)
	terr(t, err)
	inter = newPoint(3, 3)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// intersect with horizontal
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 5), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 3), newPoint(23, 3), nil)
	terr(t, err)
	inter = newPoint(3, 3)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// horizontal and vertical T-intersection
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 0), newPoint(3, 5), nil)
	terr(t, err)
	inter = newPoint(3, 0)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// horizontal and vertical general intersection
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, -5), newPoint(3, 5), nil)
	terr(t, err)
	inter = newPoint(3, 0)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// no intersection not even close
	seg1, err = op.newSegmentFromRing(newPoint(1000, 10002), newPoint(2000, 20002), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-234, -123), newPoint(-12, -23), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// no intersection kinda close
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 10), newPoint(10, 0), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// no intersection with vertical touching bbox
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(2, -5), newPoint(2, 0), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// shared point 1 (endpoint)
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 1), newPoint(0, 0), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// shared point 2 (endpoint)
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 1), newPoint(1, 1), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// T-crossing left endpoint
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0.5, 0.5), newPoint(1, 0), nil)
	terr(t, err)
	inter = newPoint(0.5, 0.5)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// T-crossing right endpoint
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 1), newPoint(0.5, 0.5), nil)
	terr(t, err)
	inter = newPoint(0.5, 0.5)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// full overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(10, 10), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(5, 5), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// shared point + overlap
	seg1, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(10, 10), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(5, 5), nil)
	terr(t, err)
	inter = newPoint(5, 5)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// mutual overlap
	seg1, err = op.newSegmentFromRing(newPoint(3, 3), newPoint(10, 10), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 5), nil)
	terr(t, err)
	inter = newPoint(3, 3)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// full overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// full overlap, orientation
	seg1, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(0, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// colinear, shared point
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(2, 2), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// colinear, shared other point
	seg1, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(0, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(2, 2), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// colinear, one encloses other
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(2, 2), nil)
	terr(t, err)
	inter = newPoint(1, 1)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear, one encloses other 2
	seg1, err = op.newSegmentFromRing(newPoint(4, 0), newPoint(0, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 1), newPoint(1, 3), nil)
	terr(t, err)
	inter = newPoint(1, 3)
	expect(t, seg1.getIntersection(seg2).equal(*inter))
	expect(t, seg2.getIntersection(seg1).equal(*inter))

	// colinear, no overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(2, 2), newPoint(4, 4), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// parallel
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -1), newPoint(1, 0), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// parallel, orientation
	seg1, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(0, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -1), newPoint(1, 0), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// parallel, position
	seg1, err = op.newSegmentFromRing(newPoint(0, -1), newPoint(1, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// endpoint intersections should be consistent - issue 60
	// If segment A T-intersects segment B, then the non-intersecting endpoint
	// of segment A should be irrelevant to the intersection of the two segs
	// From https://github.com/mfogel/polygon-clipping/issues/60
	setPrecision(NumberEPSILON)
	x, y := -91.41360941065206, 29.53135
	seg1, err = op.newSegmentFromRing(newPoint(x, y), newPoint(-91.4134943, 29.5310677), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(x, y), newPoint(-91.413, 29.5315), nil)
	terr(t, err)
	s3, err = op.newSegmentFromRing(newPoint(-91.4137213, 29.5316244), newPoint(-91.41352785864918, 29.53115), nil)
	terr(t, err)
	pt := newPoint(x, y)
	expect(t, seg1.getIntersection(s3).equal(*pt))
	expect(t, seg2.getIntersection(s3).equal(*pt))
	expect(t, s3.getIntersection(seg1).equal(*pt))
	expect(t, s3.getIntersection(seg2).equal(*pt))
	resetPrecision()

	// endpoint intersection takes priority - issue 60-5
	endX, endY := 55.31, -0.23544126113
	seg1, err = op.newSegmentFromRing(newPoint(18.60315316392773, 10.491431056669754), newPoint(endX, endY), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-32.42, 55.26), newPoint(endX, endY), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// endpoint intersection between very short and very vertical segment
	seg1, err = op.newSegmentFromRing(newPoint(-10.000000000000004, 0), newPoint(-9.999999999999995, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-10.000000000000004, 0), newPoint(-9.999999999999995, 1000), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

	// avoid intersection - issue 79
	seg1, err = op.newSegmentFromRing(newPoint(145.854148864746, -41.99816840491791), newPoint(145.85421323776, -41.9981723915721), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(145.854148864746, -41.998168404918), newPoint(145.8543, -41.9982), nil)
	terr(t, err)
	expect(t, seg1.getIntersection(seg2) == nil)
	expect(t, seg2.getIntersection(seg1) == nil)

}
func TestSegmentCompareSegments(t *testing.T) {
	// t.Parallel()

	var seg1, seg2, seg3 *segment
	var err error

	op := newOperation("")

	// non intersecting

	// not in same vertical space
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(1, 1), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(4, 3), newPoint(6, 7), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// in same vertical space, earlier is below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, -4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 1), newPoint(6, 7), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// in same vertical space, later is below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, -4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-5, -5), newPoint(6, -7), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// with left points in same vertical line
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -1), newPoint(5, -5), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// with earlier right point directly under later left point
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -1), newPoint(5, -5), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// with eariler right point directly over earlier left point
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-5, 5), newPoint(0, 3), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// intersecting not on endpoint

	// earlier c omes up from before & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, -5), newPoint(1, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// earlier comes up from directly over & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -2), newPoint(3, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// earlier comes up from after & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, -2), newPoint(3, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// later comes down from before & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, 5), newPoint(1, -2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// later comes up from directly over & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 2), newPoint(3, -2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// later comes up from after & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 2), newPoint(3, -2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// with a vertical
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, -2), newPoint(1, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// intersect but not share on an endpoint

	// with a vertical
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(2, -2), newPoint(6, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// intersect on left from above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-2, 2), newPoint(2, -2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// intersect on left from below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-2, -2), newPoint(2, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// intersect on left from vertical
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -2), newPoint(0, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// share right endpoint

	// earlier comes up from before & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, -5), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// earlier comes up from directly over & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, -2), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// earlier comes up from after & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, -2), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// later comes down from before & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, 5), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// later comes up from directly over & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 2), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// later comes up from after & above
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 0), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(1, 2), newPoint(4, 0), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// share left endpoint but not colinear

	// earlier comes up from before & below
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// one vertical, other not
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(0, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// one segment thinks theyre colinear, but the other says no
	setPrecision(NumberEPSILON)
	seg1, err = op.newSegmentFromRing(newPoint(-60.6876, -40.83428174062278), newPoint(-60.6841701, -40.83491), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-60.6876, -40.83428174062278), newPoint(-60.6874, -40.83431837489067), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)
	resetPrecision()
	// colinear

	// partial mutual overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, -1), newPoint(2, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// complete overlap
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, -1), newPoint(5, 5), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// right endpoints match
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-1, -1), newPoint(4, 4), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)

	// left endpoints match - should be length
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(3, 3), nil)
	terr(t, err)
	seg3, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 5), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == 1)
	expect(t, segmentCompare(seg2, seg1) == -1)
	expect(t, segmentCompare(seg2, seg3) == -1)
	expect(t, segmentCompare(seg3, seg2) == 1)
	expect(t, segmentCompare(seg1, seg3) == -1)
	expect(t, segmentCompare(seg3, seg1) == 1)

	// exactly equal segments should be sorted by ring id
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// exactly equal segments (but not identical) are consistent
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	result := segmentCompare(seg1, seg2)
	expect(t, segmentCompare(seg1, seg2) == result)
	expect(t, segmentCompare(seg2, seg1) == -result)

	// segment consistency - from #60
	seg1, err = op.newSegmentFromRing(newPoint(-131.57153657554915, 55.01963125), newPoint(-131.571478, 55.0187174), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-131.57153657554915, 55.01963125), newPoint(-131.57152375603846, 55.01943125), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)

	// ensure transitive - part of issue 60
	seg1, err = op.newSegmentFromRing(newPoint(-10.000000000000018, -9.17), newPoint(-10.000000000000004, -8.79), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-10.000000000000016, 1.44), newPoint(-9, 1.5), nil)
	terr(t, err)
	seg3, err = op.newSegmentFromRing(newPoint(-10.00000000000001, 1.75), newPoint(-9, 1.5), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg3) == -1)
	expect(t, segmentCompare(seg1, seg3) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)
	expect(t, segmentCompare(seg3, seg2) == 1)
	expect(t, segmentCompare(seg3, seg1) == 1)

	// ensure transitive 2 - also part of issue 60
	seg1, err = op.newSegmentFromRing(newPoint(-10.000000000000002, 1.8181818181818183), newPoint(-9.999999999999996, -3), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-10.000000000000002, 1.8181818181818183), newPoint(0, 0), nil)
	terr(t, err)
	seg3, err = op.newSegmentFromRing(newPoint(-10.000000000000002, 1.8181818181818183), newPoint(-10.000000000000002, 2), nil)
	terr(t, err)
	expect(t, segmentCompare(seg1, seg2) == -1)
	expect(t, segmentCompare(seg2, seg3) == -1)
	expect(t, segmentCompare(seg1, seg3) == -1)
	expect(t, segmentCompare(seg2, seg1) == 1)
	expect(t, segmentCompare(seg3, seg2) == 1)
	expect(t, segmentCompare(seg3, seg1) == 1)
}
