package polygol

import (
	"fmt"
)

type segment struct {
	id              int
	leftSE          *sweepEvent
	rightSE         *sweepEvent
	rings           []*ringIn
	windings        []int
	ringOut         *ringOut
	consumedBy      *segment
	inResult        bool
	forceInResult   bool
	doneInResult    bool
	prev            *segment
	prevSegInResult *segment
	after           *state
	before          *state
	op              *operation
}

func (o *operation) newSegment(leftSE, rightSE *sweepEvent, rings []*ringIn, windings []int) *segment {
	o.segmentID++

	s := &segment{}
	s.id = o.segmentID
	s.leftSE = leftSE

	leftSE.segment = s
	leftSE.otherSE = rightSE

	s.rightSE = rightSE

	rightSE.segment = s
	rightSE.otherSE = leftSE

	s.rings = rings
	s.windings = windings

	s.op = o

	return s
}

func segmentCompare(a, b interface{}) int {

	aSeg := a.(*segment)
	bSeg := b.(*segment)

	alx := aSeg.leftSE.point.x
	blx := bSeg.leftSE.point.x
	arx := aSeg.rightSE.point.x
	brx := bSeg.rightSE.point.x

	// check if they're even in the same vertical plane
	if brx.isLessThan(alx) {
		return 1
	}
	if arx.isLessThan(blx) {
		return -1
	}

	aly := aSeg.leftSE.point.y
	bly := bSeg.leftSE.point.y
	ary := aSeg.rightSE.point.y
	bry := bSeg.rightSE.point.y

	// is left endpoint of segment B the right-more?
	if alx.isLessThan(blx) {

		// are the two segments in the same horizontal plane?
		if bly.isLessThan(aly) && bly.isLessThan(ary) {
			return 1
		}
		if bly.isGreaterThan(aly) && bly.isGreaterThan(ary) {
			return -1
		}

		// is the B left endpoint colinear to segment A?
		aCmpBLeft := aSeg.comparePoint(bSeg.leftSE.point)
		if aCmpBLeft < 0 {
			return 1
		}
		if aCmpBLeft > 0 {
			return -1
		}

		// is the A right endpoint colinear to segment B ?
		bCmpARight := bSeg.comparePoint(aSeg.rightSE.point)
		if bCmpARight != 0 {
			return bCmpARight
		}

		// colinear segments, consider the one with left-more
		// left endpoint to be first (arbitrary?)
		return -1
	}

	// is left endpoint of segment A the right-more?
	if alx.isGreaterThan(blx) {

		if aly.isLessThan(bly) && aly.isLessThan(bry) {
			return -1
		}
		if aly.isGreaterThan(bly) && aly.isGreaterThan(bry) {
			return 1
		}

		// is the A left endpoint colinear to segment B?
		bCmpALeft := bSeg.comparePoint(aSeg.leftSE.point)
		if bCmpALeft != 0 {
			return bCmpALeft
		}

		// is the B right endpoint colinear to segment A?
		aCmpBRight := aSeg.comparePoint(bSeg.rightSE.point)
		if aCmpBRight < 0 {
			return 1
		}
		if aCmpBRight > 0 {
			return -1
		}

		// colinear segments, consider the one with left-more
		// left endpoint to be first (arbitrary?)
		return 1
	}

	// if we get here, the two left endpoints are in the same
	// vertical plane, ie alx === blx

	// consider the lower left-endpoint to come first
	if aly.isLessThan(bly) {
		return -1
	}
	if aly.isGreaterThan(bly) {
		return 1
	}

	// left endpoints are identical
	// check for colinearity by using the left-more right endpoint

	// is the A right endpoint more left-more?
	if arx.isLessThan(brx) {
		bCmpARight := bSeg.comparePoint(aSeg.rightSE.point)
		if bCmpARight != 0 {
			return bCmpARight
		}
	}

	// is the B right endpoint more left-more?
	if arx.isGreaterThan(brx) {
		aCmpBRight := aSeg.comparePoint(bSeg.rightSE.point)
		if aCmpBRight < 0 {
			return 1
		}
		if aCmpBRight > 0 {
			return -1
		}
	}

	if !arx.equalTo(brx) {
		// are these two [almost] vertical segments with opposite orientation?
		// if so, the one with the lower right endpoint comes first
		ay := ary.minus(aly)
		ax := arx.minus(alx)
		by := bry.minus(bly)
		bx := brx.minus(blx)
		if ay.isGreaterThan(ax) && by.isLessThan(bx) {
			return 1
		}
		if ay.isLessThan(ax) && by.isGreaterThan(bx) {
			return -1
		}
	}

	// we have colinear segments with matching orientation
	// consider the one with more left-more right endpoint to be first
	if arx.isGreaterThan(brx) {
		return 1
	}
	if arx.isLessThan(brx) {
		return -1
	}

	// if we get here, two two right endpoints are in the same
	// vertical plane, ie arx === brx

	// consider the lower right-endpoint to come first
	if ary.isLessThan(bry) {
		return -1
	}
	if ary.isGreaterThan(bry) {
		return 1
	}

	// right endpoints identical as well, so the segments are identical
	// fall back on creation order as consistent tie-breaker
	if aSeg.id < bSeg.id {
		return -1
	}
	if aSeg.id > bSeg.id {
		return 1
	}

	// identical segment, ie a === b
	return 0
}

func (o *operation) newSegmentFromRing(pt1, pt2 *point, ring *ringIn) (*segment, error) {
	var leftPt, rightPt *point
	var winding int

	cmpPts := sweepEventComparePoints(pt1, pt2)
	if cmpPts < 0 {
		leftPt = pt1
		rightPt = pt2
		winding = 1
	} else if cmpPts > 0 {
		leftPt = pt2
		rightPt = pt1
		winding = -1
	} else {
		return nil, fmt.Errorf("tried to create degenerate segment at [%f,%f].", pt1.x.number(), pt1.y.number())
	}

	leftSE := newSweepEvent(leftPt, true)
	rightSE := newSweepEvent(rightPt, false)

	return o.newSegment(leftSE, rightSE, []*ringIn{ring}, []int{winding}), nil
}

func (s *segment) replaceRightSE(newRightSE *sweepEvent) {
	s.rightSE = newRightSE
	s.rightSE.segment = s
	s.rightSE.otherSE = s.leftSE
	s.leftSE.otherSE = s.rightSE
}

func (s *segment) bbox() Bbox {

	y1 := s.leftSE.point.y
	y2 := s.rightSE.point.y

	lly := y2
	if y1.isLessThan(y2) {
		lly = y1
	}

	ury := y2
	if y1.isGreaterThan(y2) {
		ury = y1
	}

	return Bbox{
		ll: Vector{x: s.leftSE.point.x, y: lly},
		ur: Vector{x: s.rightSE.point.x, y: ury},
	}
}

func (s *segment) vector() Vector {
	return Vector{
		x: s.rightSE.point.x.minus(s.leftSE.point.x),
		y: s.rightSE.point.y.minus(s.leftSE.point.y),
	}
}

func (s *segment) isAnEndpoint(point *point) bool {
	if s == nil {
		return false
	}
	if point == nil {
		return false
	}
	return (point.x.equalTo(s.leftSE.point.x) && point.y.equalTo(s.leftSE.point.y)) ||
		(point.x.equalTo(s.rightSE.point.x) && point.y.equalTo(s.rightSE.point.y))
	// if s.leftSE != nil {
	// 	// if almostEqual(point.x, s.leftSE.point.x) && almostEqual(point.y, s.leftSE.point.y) {
	// 	if point.x == s.leftSE.point.x && point.y == s.leftSE.point.y {
	// 		return true
	// 	}
	// }
	// if s.rightSE != nil {
	// 	// if almostEqual(point.x, s.rightSE.point.x) && almostEqual(point.y, s.rightSE.point.y) {
	// 	if point.x == s.rightSE.point.x && point.y == s.rightSE.point.y {
	// 		return true
	// 	}
	// }
	// return false
}

func (s *segment) comparePoint(point *point) int {
	return orient(s.leftSE.point.Vector, point.vector(), s.rightSE.point.Vector)
}

func (s *segment) getIntersection(other *segment) *point {

	if s == nil || other == nil {
		return nil
	}

	// If bboxes don't overlap, there can't be any intersections
	segBbox := s.bbox()
	otherBbox := other.bbox()

	bboxOverlap := segBbox.getBboxOverlap(otherBbox)
	if bboxOverlap == nil {
		return nil
	}

	// We first check to see if the endpoints can be considered intersections.
	// This will 'snap' intersections to endpoints if possible, and will
	// handle cases of colinearity.

	tlp := s.leftSE.point
	trp := s.rightSE.point
	olp := other.leftSE.point
	orp := other.rightSE.point

	// does each endpoint touch the other segment?
	// note that we restrict the 'touching' definition to only allow segments
	// to touch endpoints that lie forward from where we are in the sweep line pass
	touchesOtherLSE := segBbox.isInBbox(olp.Vector) && s.comparePoint(olp) == 0
	touchesThisLSE := otherBbox.isInBbox(tlp.Vector) && other.comparePoint(tlp) == 0
	touchesOtherRSE := segBbox.isInBbox(orp.Vector) && s.comparePoint(orp) == 0
	touchesThisRSE := otherBbox.isInBbox(trp.Vector) && other.comparePoint(trp) == 0

	// do left endpoints match?
	if touchesThisLSE && touchesOtherLSE {
		// these two cases are for colinear segments with matching left
		// endpoints, and one segment being longer than the other
		if touchesThisRSE && !touchesOtherRSE {
			return trp
		}
		if !touchesThisRSE && touchesOtherRSE {
			return orp
		}
		// either the two segments match exactly (two trival intersections)
		// or just on their left endpoint (one trivial intersection
		return nil
	}

	// does this left endpoint matches (other doesn't)
	if touchesThisLSE {
		// check for segments that just intersect on opposing endpoints
		if touchesOtherRSE {
			if tlp.x.equalTo(orp.x) && tlp.y.equalTo(orp.y) {
				return nil
			}
		}
		// t-intersection on left endpoint
		return tlp
	}

	// does other left endpoint matches (this doesn't)
	if touchesOtherLSE {
		// check for segments that just intersect on opposing endpoints
		if touchesThisRSE {
			if trp.x.equalTo(olp.x) && trp.y.equalTo(olp.y) {
				return nil
			}
		}
		// t-intersection on left endpoint
		return olp
	}

	// trivial intersection on right endpoints
	if touchesThisRSE && touchesOtherRSE {
		return nil
	}

	// t-intersections on just one right endpoint
	if touchesThisRSE {
		return trp
	}
	if touchesOtherRSE {
		return orp
	}

	// None of our endpoints intersect. Look for a general intersection between
	// infinite lines laid over the segments

	pt := intersection(
		tlp.vector(),
		s.vector(),
		olp.vector(),
		other.vector(),
	)

	// are the segments parallel? Note that if they were colinear with overlap,
	// they would have an endpoint intersection and that case was already handled above
	if pt == nil {
		return nil
	}
	ptInter := &point{Vector: *pt}

	// is the intersection found between the lines not on the segments?
	if !bboxOverlap.isInBbox(ptInter.Vector) {
		return nil
	}

	// round the the computed point if needed
	return s.op.rounder.round(ptInter.x, ptInter.y)
}

func (s *segment) split(point *point) []*sweepEvent {

	newEvents := []*sweepEvent{}
	alreadyLinked := point.events != nil

	newLeftSE := newSweepEvent(point, true)
	newRightSE := newSweepEvent(point, false)
	oldRightSE := s.rightSE

	s.replaceRightSE(newRightSE)
	newEvents = append(newEvents, newRightSE)
	newEvents = append(newEvents, newLeftSE)

	newRings := make([]*ringIn, len(s.rings))
	copy(newRings, s.rings)

	newWindings := make([]int, len(s.windings))
	copy(newWindings, s.windings)

	newSeg := s.op.newSegment(newLeftSE, oldRightSE, newRings, newWindings)

	// when splitting a nearly vertical downward-facing segment,
	// sometimes one of the resulting new segments is vertical, in which
	// case its left and right events may need to be swapped
	if sweepEventComparePoints(newSeg.leftSE.point, newSeg.rightSE.point) > 0 {
		newSeg.swapEvents()
	}
	if sweepEventComparePoints(s.leftSE.point, s.rightSE.point) > 0 {
		s.swapEvents()
	}

	// in the point we just used to create new sweep events with was already
	// linked to other events, we need to check if either of the affected
	// segments should be consumed
	if alreadyLinked {
		newLeftSE.checkForConsuming()
		newRightSE.checkForConsuming()
	}

	return newEvents
}

func (s *segment) swapEvents() {
	s.rightSE, s.leftSE = s.leftSE, s.rightSE
	s.leftSE.isLeft = true
	s.rightSE.isLeft = false
	for i := 0; i < len(s.windings); i++ {
		s.windings[i] *= -1
	}
}

func (s *segment) consume(otherSeg *segment) {
	consumer := s
	consumee := otherSeg

	for consumer.consumedBy != nil {
		consumer = consumer.consumedBy
	}
	for consumee.consumedBy != nil {
		consumee = consumee.consumedBy
	}

	cmp := segmentCompare(consumer, consumee)
	if cmp == 0 {
		return // already consumed
	}

	// the winner of the consumption is the earlier segment
	// according to sweep line ordering
	if cmp > 0 {
		consumer, consumee = consumee, consumer
	}

	// make sure a segment doesn't consume its prev
	if consumer.prev == consumee {
		consumer, consumee = consumee, consumer
	}

	for i := 0; i < len(consumee.rings); i++ {
		ring := consumee.rings[i]
		winding := consumee.windings[i]
		index := ring.indexOf(consumer.rings)
		if index == -1 {
			consumer.rings = append(consumer.rings, ring)
			consumer.windings = append(consumer.windings, winding)
		} else {
			consumer.windings[index] += winding
		}
	}
	consumee.rings = nil
	consumee.windings = nil
	consumee.consumedBy = consumer

	// mark sweep events consumed as to maintain ordering in sweep event queue
	consumee.leftSE.consumedBy = consumer.leftSE
	consumee.rightSE.consumedBy = consumer.rightSE
}

func (s *segment) prevInResult() *segment {
	if s.prevSegInResult != nil {
		return s.prevSegInResult
	}
	if s.prev == nil {
		s.prevSegInResult = nil
	} else if s.prev.isInResult() {
		s.prevSegInResult = s.prev
	} else {
		s.prevSegInResult = s.prev.prevInResult()
	}
	return s.prevSegInResult
}

type state struct {
	rings      []*ringIn
	windings   []int
	multiPolys []*multiPolyIn
}

func (s *segment) beforeState() *state {
	if s.before != nil {
		return s.before
	}
	if s.prev == nil {
		s.before = &state{
			rings:      []*ringIn{},
			windings:   []int{},
			multiPolys: []*multiPolyIn{},
		}
	} else {
		seg := s.prev.consumedBy
		if s.prev.consumedBy == nil {
			seg = s.prev
		}
		s.before = seg.afterState()
	}
	return s.before
}

func (s *segment) afterState() *state {
	if s.after != nil {
		return s.after
	}

	beforeState := s.beforeState()

	ringsBefore := make([]*ringIn, len(beforeState.rings))
	copy(ringsBefore, beforeState.rings)
	windingsBefore := make([]int, len(beforeState.windings))
	copy(windingsBefore, beforeState.windings)

	s.after = &state{
		rings:      ringsBefore,
		windings:   windingsBefore,
		multiPolys: []*multiPolyIn{},
	}

	// calculate ringsAfter, windingsAfter
	for i := 0; i < len(s.rings); i++ {
		ring := s.rings[i]
		winding := s.windings[i]
		index := ring.indexOf(s.after.rings)
		if index == -1 {
			s.after.rings = append(s.after.rings, ring)
			s.after.windings = append(s.after.windings, winding)
		} else {
			s.after.windings[index] += winding
		}
	}

	// calculate polysAfter
	polysAfter := []*polyIn{}
	polysExclude := []*polyIn{}
	for i := 0; i < len(s.after.rings); i++ {
		if s.after.windings[i] == 0 { // non-zero rule
			continue
		}
		ring := s.after.rings[i]
		poly := ring.poly
		index := poly.indexOf(polysExclude)
		if index != -1 {
			continue
		}
		if ring.isExterior { // exterior ring
			polysAfter = append(polysAfter, poly)
		} else { // interior ring
			if poly.indexOf(polysExclude) == -1 {
				polysExclude = append(polysExclude, poly)
			}
			index := ring.poly.indexOf(polysAfter)
			if index != -1 {
				polysAfter = append(polysAfter[:index], polysAfter[index+1:]...) // splice(index,1)
			}
		}
	}

	// calculate multiPolysAfter
	for i := 0; i < len(polysAfter); i++ {
		mp := polysAfter[i].multiPoly
		if mp.indexOf(s.after.multiPolys) == -1 {
			s.after.multiPolys = append(s.after.multiPolys, mp)
		}
	}

	return s.after
}

func (s *segment) isInResult() bool {
	// if we've been consumed, we're not in the result
	if s == nil {
		return false
	}
	if s.consumedBy != nil {
		return false
	}
	if s.forceInResult {
		return s.inResult
	}
	if s.doneInResult {
		return s.inResult
	}

	mpsBefore := s.beforeState().multiPolys
	mpsAfter := s.afterState().multiPolys

	switch s.op.opType {
	case "union":
		// UNION - included iff:
		//  * On one side of us there is 0 poly interiors AND
		//  * On the other side there is 1 or more.
		noBefores := len(mpsBefore) == 0
		noAfters := len(mpsAfter) == 0
		s.inResult = noBefores != noAfters
	case "intersection":
		// INTERSECTION - included iff:
		//  * on one side of us all multipolys are rep. with poly interiors AND
		//  * on the other side of us, not all multipolys are repsented
		//    with poly interiors
		var least, most int
		if len(mpsBefore) < len(mpsAfter) {
			least = len(mpsBefore)
			most = len(mpsAfter)
		} else {
			least = len(mpsAfter)
			most = len(mpsBefore)
		}
		s.inResult = most == s.op.numMultiPolys && least < most
	case "xor":
		// XOR - included iff:
		//  * the difference between the number of multipolys represented
		//    with poly interiors on our two sides is an odd number
		diff := abs(len(mpsBefore) - len(mpsAfter))
		s.inResult = diff%2 == 1
	case "difference":
		// DIFFERENCE included iff:
		//  * on exactly one side, we have just the subject
		isJustSubject := func(mps []*multiPolyIn) bool {
			return len(mps) == 1 && mps[0].isSubject
		}
		s.inResult = isJustSubject(mpsBefore) != isJustSubject(mpsAfter)
	default:
		fmt.Printf("Unrecognized operation type found %s", s.op.opType)
	}
	s.doneInResult = true
	return s.inResult
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
