package polygol

import (
	"testing"
)

func TestSweepEventCompare(t *testing.T) {
	// t.Parallel()

	var se1, se2, se3 *sweepEvent
	var seg1, seg2 *segment
	var err error

	op := newOperation("")

	// favor earlier x in point
	se1 = newSweepEvent(newPoint(-5, 4), false)
	se2 = newSweepEvent(newPoint(5, 1), false)
	expect(t, sweepEventCompare(se1, se2) == -1)
	expect(t, sweepEventCompare(se2, se1) == 1)

	// then favor earlier y in point
	se1 = newSweepEvent(newPoint(5, -4), false)
	se2 = newSweepEvent(newPoint(5, 4), false)
	expect(t, sweepEventCompare(se1, se2) == -1)
	expect(t, sweepEventCompare(se2, se1) == 1)

	// then favor right events over left
	seg1, err = op.newSegmentFromRing(newPoint(5, 4), newPoint(3, 2), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(5, 4), newPoint(6, 5), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.rightSE, seg2.leftSE) == -1)
	expect(t, sweepEventCompare(seg2.leftSE, seg1.rightSE) == 1)

	// then favor non-vertical segments for left events
	seg1, err = op.newSegmentFromRing(newPoint(3, 2), newPoint(3, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 2), newPoint(5, 4), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.rightSE) == -1)
	expect(t, sweepEventCompare(seg2.rightSE, seg1.leftSE) == 1)

	// then favor vertical segments for right events
	seg1, err = op.newSegmentFromRing(newPoint(3, 4), newPoint(3, 2), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(3, 4), newPoint(1, 2), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.rightSE) == -1)
	expect(t, sweepEventCompare(seg2.rightSE, seg1.leftSE) == 1)

	// then favor lower segment
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 6), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.rightSE) == -1)
	expect(t, sweepEventCompare(seg2.rightSE, seg1.leftSE) == 1)

	// Sometimes from one segment's perspective it appears colinear
	// to another segment, but from that other segment's perspective
	// they aren't colinear. This happens because a longer segment
	// is able to better determine what is and is not colinear.
	// and favor barely lower segment
	seg1, err = op.newSegmentFromRing(newPoint(-75.725, 45.357), newPoint(-75.72484615384616, 45.35723076923077), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-75.725, 45.357), newPoint(-75.723, 45.36), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.leftSE) == 1)
	expect(t, sweepEventCompare(seg2.leftSE, seg1.leftSE) == -1)

	// then favor lower ring id
	seg1, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(4, 4), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(0, 0), newPoint(5, 5), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.leftSE) == -1)
	expect(t, sweepEventCompare(seg2.leftSE, seg1.leftSE) == 1)

	// identical equal
	se1 = newSweepEvent(newPoint(0, 0), false)
	se3 = newSweepEvent(newPoint(3, 3), false)
	op.newSegment(se1, se3, nil, nil)
	op.newSegment(se1, se3, nil, nil)
	expect(t, sweepEventCompare(se1, se1) == 0)

	// totally equal but not identical events are consistent
	se1 = newSweepEvent(newPoint(0, 0), false)
	se2 = newSweepEvent(newPoint(0, 0), false)
	se3 = newSweepEvent(newPoint(3, 3), false)
	op.newSegment(se1, se3, nil, nil)
	op.newSegment(se2, se3, nil, nil)
	result := sweepEventCompare(se1, se2)
	expect(t, sweepEventCompare(se1, se2) == result)
	expect(t, sweepEventCompare(se2, se1) == -result)

	// events are linked as side effect
	se1 = newSweepEvent(newPoint(0, 0), false)
	se2 = newSweepEvent(newPoint(0, 0), false)
	op.newSegment(se1, newSweepEvent(newPoint(2, 2), false), nil, nil)
	op.newSegment(se2, newSweepEvent(newPoint(3, 4), false), nil, nil)
	expect(t, se1.point.equal(*se2.point))
	sweepEventCompare(se1, se2)
	expect(t, se1.point.equal(*se2.point))

	// consistency edge case
	// harvested from https://github.com/mfogel/polygon-clipping/issues/62
	seg1, err = op.newSegmentFromRing(newPoint(-71.0390933353125, 41.504475), newPoint(-71.0389879, 41.5037842), nil)
	terr(t, err)
	seg2, err = op.newSegmentFromRing(newPoint(-71.0390933353125, 41.504475), newPoint(-71.03906280974431, 41.5042756), nil)
	terr(t, err)
	expect(t, sweepEventCompare(seg1.leftSE, seg2.leftSE) == -1)
	expect(t, sweepEventCompare(seg2.leftSE, seg1.leftSE) == 1)
}

func TestNewSweepEvent(t *testing.T) {
	// t.Parallel()

	var (
		se1, se2 *sweepEvent
		p1       *point
	)

	// events created from same point are already linked
	p1 = newPoint(0, 0)
	se1 = newSweepEvent(p1, false)
	se2 = newSweepEvent(p1, false)
	expect(t, se1.point.equal(*p1))
	expect(t, equalSweepEvents(se1.point.events, se2.point.events))
}

func TestSweepEventLink(t *testing.T) {
	// t.Parallel()

	var se1, se2, se3, se4 *sweepEvent
	var p1, p2 *point
	var err error

	op := newOperation("")

	// no linked events
	se1 = newSweepEvent(newPoint(0, 0), false)
	expect(t, equalSweepEvents(se1.point.events, []*sweepEvent{se1}))

	// link events already linked with others
	p1 = newPoint(1, 2)
	p2 = newPoint(1, 2)
	se1 = newSweepEvent(p1, false)
	se2 = newSweepEvent(p1, false)
	se3 = newSweepEvent(p2, false)
	se4 = newSweepEvent(p2, false)
	op.newSegment(se1, newSweepEvent(newPoint(5, 5), false), nil, nil)
	op.newSegment(se2, newSweepEvent(newPoint(6, 6), false), nil, nil)
	op.newSegment(se3, newSweepEvent(newPoint(7, 7), false), nil, nil)
	op.newSegment(se4, newSweepEvent(newPoint(8, 8), false), nil, nil)
	err = se1.link(se3)
	terr(t, err)
	expect(t, len(se1.point.events) == 4)
	expect(t, se1.point.equal(*se2.point))
	expect(t, se1.point.equal(*se3.point))
	expect(t, se1.point.equal(*se4.point))

	// same event twice
	p1 = newPoint(0, 0)
	se1 = newSweepEvent(p1, false)
	se2 = newSweepEvent(p1, false)
	expect(t, se2.link(se1) != nil)
	expect(t, se1.link(se2) != nil)

	// unavailable linked events do not show up
	p1 = newPoint(0, 0)
	se1 = newSweepEvent(p1, false)
	seAlreadyProcessed := newSweepEvent(p1, false)
	seAlreadyProcessed.segment = &segment{forceInResult: true, inResult: true, ringOut: &ringOut{}}
	seNotInResult := newSweepEvent(p1, false)
	seNotInResult.segment = &segment{forceInResult: true, inResult: false}
	expect(t, equalSweepEvents(se1.getAvailableLinkedEvents(), []*sweepEvent{}))

	// available linked events show up
	p1 = newPoint(0, 0)
	se1 = newSweepEvent(p1, false)
	seOkay := newSweepEvent(p1, false)
	seOkay.segment = &segment{forceInResult: true, inResult: true}
	expect(t, equalSweepEvents(se1.getAvailableLinkedEvents(), []*sweepEvent{seOkay}))

	// link goes both ways
	p1 = newPoint(0, 0)
	seOkay1 := newSweepEvent(p1, false)
	seOkay2 := newSweepEvent(p1, false)
	seOkay1.segment = &segment{forceInResult: true, inResult: true}
	seOkay2.segment = &segment{forceInResult: true, inResult: true}
	expect(t, equalSweepEvents(seOkay1.getAvailableLinkedEvents(), []*sweepEvent{seOkay2}))
	expect(t, equalSweepEvents(seOkay2.getAvailableLinkedEvents(), []*sweepEvent{seOkay1}))
}

func TestSweepEventgetLeftMostComparator(t *testing.T) {
	// t.Parallel()

	var prevEvent, event *sweepEvent
	var comparator func(a, b *sweepEvent) int
	var se1, se2, se3, se4, se5 *sweepEvent

	op := newOperation("")

	// after a segment straight to the right
	prevEvent = newSweepEvent(newPoint(0, 0), false)
	event = newSweepEvent(newPoint(1, 0), false)
	comparator = event.getLeftMostComparator(prevEvent)

	se1 = newSweepEvent(newPoint(1, 0), false)
	op.newSegment(se1, newSweepEvent(newPoint(0, 1), false), nil, nil)

	se2 = newSweepEvent(newPoint(1, 0), false)
	op.newSegment(se2, newSweepEvent(newPoint(1, 1), false), nil, nil)

	se3 = newSweepEvent(newPoint(1, 0), false)
	op.newSegment(se3, newSweepEvent(newPoint(2, 0), false), nil, nil)

	se4 = newSweepEvent(newPoint(1, 0), false)
	op.newSegment(se4, newSweepEvent(newPoint(1, -1), false), nil, nil)

	se5 = newSweepEvent(newPoint(1, 0), false)
	op.newSegment(se5, newSweepEvent(newPoint(0, -1), false), nil, nil)

	expect(t, comparator(se1, se2) == -1)
	expect(t, comparator(se2, se3) == -1)
	expect(t, comparator(se3, se4) == -1)
	expect(t, comparator(se4, se5) == -1)

	expect(t, comparator(se2, se1) == 1)
	expect(t, comparator(se3, se2) == 1)
	expect(t, comparator(se4, se3) == 1)
	expect(t, comparator(se5, se4) == 1)

	expect(t, comparator(se1, se3) == -1)
	expect(t, comparator(se1, se4) == -1)
	expect(t, comparator(se1, se5) == -1)

	expect(t, comparator(se1, se1) == 0)

	// after a down and to the left
	prevEvent = newSweepEvent(newPoint(1, 1), false)
	event = newSweepEvent(newPoint(0, 0), false)
	comparator = event.getLeftMostComparator(prevEvent)

	se1 = newSweepEvent(newPoint(0, 0), false)
	op.newSegment(se1, newSweepEvent(newPoint(0, 1), false), nil, nil)

	se2 = newSweepEvent(newPoint(0, 0), false)
	op.newSegment(se2, newSweepEvent(newPoint(1, 0), false), nil, nil)

	se3 = newSweepEvent(newPoint(0, 0), false)
	op.newSegment(se3, newSweepEvent(newPoint(0, -1), false), nil, nil)

	se4 = newSweepEvent(newPoint(0, 0), false)
	op.newSegment(se4, newSweepEvent(newPoint(-1, 0), false), nil, nil)

	expect(t, comparator(se1, se2) == 1)
	expect(t, comparator(se1, se3) == 1)
	expect(t, comparator(se1, se4) == 1)

	expect(t, comparator(se2, se1) == -1)
	expect(t, comparator(se2, se3) == -1)
	expect(t, comparator(se2, se4) == -1)

	expect(t, comparator(se3, se1) == -1)
	expect(t, comparator(se3, se2) == 1)
	expect(t, comparator(se3, se4) == -1)

	expect(t, comparator(se4, se1) == -1)
	expect(t, comparator(se4, se2) == 1)
	expect(t, comparator(se4, se3) == 1)
}
