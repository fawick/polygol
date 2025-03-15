package polygol

import (
	"testing"
)

func equalBbox(bb1, bb2 Bbox) bool {
	return bb1.ll.x.equalTo(bb2.ll.x) &&
		bb1.ll.y.equalTo(bb2.ll.y) &&
		bb1.ur.x.equalTo(bb2.ur.x) &&
		bb1.ur.y.equalTo(bb2.ur.y)
}

const NumberEPSILON = float64(7.)/3 - float64(4.)/3 - float64(1.)

// const epsilon = 2e-12

func TestIsInBbox(t *testing.T) {
	var b Bbox

	// outside
	b = Bbox{ll: Vector{x: newBigNumber(1), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}
	expect(t, !b.isInBbox(Vector{x: newBigNumber(0), y: newBigNumber(3)}))
	expect(t, !b.isInBbox(Vector{x: newBigNumber(3), y: newBigNumber(30)}))
	expect(t, !b.isInBbox(Vector{x: newBigNumber(3), y: newBigNumber(-30)}))
	expect(t, !b.isInBbox(Vector{x: newBigNumber(9), y: newBigNumber(3)}))

	// inside
	b = Bbox{ll: Vector{x: newBigNumber(1), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}
	expect(t, b.isInBbox(Vector{x: newBigNumber(1), y: newBigNumber(2)}))
	expect(t, b.isInBbox(Vector{x: newBigNumber(5), y: newBigNumber(6)}))
	expect(t, b.isInBbox(Vector{x: newBigNumber(1), y: newBigNumber(6)}))
	expect(t, b.isInBbox(Vector{x: newBigNumber(5), y: newBigNumber(2)}))
	expect(t, b.isInBbox(Vector{x: newBigNumber(3), y: newBigNumber(4)}))

	// barely inside & outside
	b = Bbox{ll: Vector{x: newBigNumber(1), y: newBigNumber(0.8)}, ur: Vector{x: newBigNumber(1.2), y: newBigNumber(6)}}
	expect(t, b.isInBbox(Vector{x: newBigNumber(1.2).minus(newBigNumber(NumberEPSILON)), y: newBigNumber(6)}))
	expect(t, !b.isInBbox(Vector{x: newBigNumber(1.2).plus(newBigNumber(NumberEPSILON)), y: newBigNumber(6)}))
	expect(t, b.isInBbox(Vector{x: newBigNumber(1), y: newBigNumber(0.8).plus(newBigNumber(NumberEPSILON))}))
	expect(t, !b.isInBbox(Vector{x: newBigNumber(1), y: newBigNumber(0.8).minus(newBigNumber(NumberEPSILON))}))
}

func TestBboxOverlap(t *testing.T) {

	var b1, b2 Bbox
	b1 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}

	// disjoint - none
	// above
	b2 = Bbox{ll: Vector{x: newBigNumber(7), y: newBigNumber(7)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(8)}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// left
	b2 = Bbox{ll: Vector{x: newBigNumber(1), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(3), y: newBigNumber(8)}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// down
	b2 = Bbox{ll: Vector{x: newBigNumber(2), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(3), y: newBigNumber(3)}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// right
	b2 = Bbox{ll: Vector{x: newBigNumber(12), y: newBigNumber(1)}, ur: Vector{x: newBigNumber(14), y: newBigNumber(9)}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// touching - one Vector
	// upper right corner of 1
	b2 = Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(8)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}))
	// upper left corner of 1
	b2 = Bbox{ll: Vector{x: newBigNumber(3), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(8)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(6)}}))
	// lower left corner of 1
	b2 = Bbox{ll: Vector{x: newBigNumber(0), y: newBigNumber(0)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(4)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(4)}}))
	// lower right corner of 1
	b2 = Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(0)}, ur: Vector{x: newBigNumber(12), y: newBigNumber(4)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(4)}}))

	// overlapping - two Vectors

	// full overlap

	// matching bboxes
	expect(t, equalBbox(*b1.getBboxOverlap(b1), b1))

	// one side & two corners matching
	b2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}))

	// one corner matching, part of two sides
	b2 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(5)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(5)}}))

	// part of a side matching, no corners
	b2 = Bbox{ll: Vector{x: newBigNumber(4.5), y: newBigNumber(4.5)}, ur: Vector{x: newBigNumber(5.5), y: newBigNumber(6)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4.5), y: newBigNumber(4.5)}, ur: Vector{x: newBigNumber(5.5), y: newBigNumber(6)}}))

	// completely enclosed - no side or corner matching
	b2 = Bbox{ll: Vector{x: newBigNumber(4.5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(5.5), y: newBigNumber(5.5)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), b2))

	// partial overlap

	// full side overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(3), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}))

	// partial side overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(4.5)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(5.5)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(4.5)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(5.5)}}))

	// corner overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(7)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}))

	// line bboxes

	// vertical line & normal

	// no overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(7), y: newBigNumber(3)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(6)}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// Vector overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(0)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(4)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(4)}}))

	// line overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(0)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(9)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}))

	// horizontal line & normal

	// no overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(3), y: newBigNumber(7)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(7)}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// Vector overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(1), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(6)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(6)}}))

	// line overlap
	b2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}))

	// two vertical lines
	var v1, v2 Bbox
	v1 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(6)}}

	// no overlap
	v2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(7)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(8)}}
	expect(t, v1.getBboxOverlap(v2) == nil)

	// Vector overlap
	v2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(3)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(4)}}
	expect(t, equalBbox(*v1.getBboxOverlap(v2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(4)}}))

	// line overlap
	v2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(3)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(5)}}
	expect(t, equalBbox(*v1.getBboxOverlap(v2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(5)}}))

	// two horizontal lines
	var h1, h2 Bbox
	h1 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(6)}}

	// no overlap
	h2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(5)}}
	expect(t, h1.getBboxOverlap(h2) == nil)

	// Vector overlap
	h2 = Bbox{ll: Vector{x: newBigNumber(7), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(6)}}
	expect(t, equalBbox(*h1.getBboxOverlap(h2), Bbox{ll: Vector{x: newBigNumber(7), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(6)}}))

	// line overlap
	h2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(6)}}
	expect(t, equalBbox(*h1.getBboxOverlap(h2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(7), y: newBigNumber(6)}}))

	// horizontal & vertical lines

	// no overlap
	h1 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(6)}}
	v1 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(7)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(9)}}
	expect(t, h1.getBboxOverlap(v1) == nil)

	// Vector overlap
	h1 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(6)}}
	v1 = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(9)}}
	expect(t, equalBbox(*h1.getBboxOverlap(v1), Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(6)}}))

	// produced line box

	// horizontal
	b2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(8)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}))

	// vertical
	b2 = Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(8), y: newBigNumber(8)}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), Bbox{ll: Vector{x: newBigNumber(6), y: newBigNumber(4)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(6)}}))

	// Vector bboxes
	var p Bbox

	// Vector & normal

	// no overlap
	p = Bbox{ll: Vector{x: newBigNumber(2), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(2), y: newBigNumber(2)}}
	expect(t, b1.getBboxOverlap(p) == nil)

	// Vector overlap
	p = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(5)}}
	expect(t, equalBbox(*b1.getBboxOverlap(p), p))

	// Vector & line
	var l Bbox

	// no overlap
	p = Bbox{ll: Vector{x: newBigNumber(2), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(2), y: newBigNumber(2)}}
	l = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(8)}}
	expect(t, l.getBboxOverlap(p) == nil)

	// Vector overlap
	p = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(5)}}
	l = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(6), y: newBigNumber(5)}}
	expect(t, equalBbox(*l.getBboxOverlap(p), p))

	// Vector & Vector
	var p1, p2 Bbox

	// no overlap
	p1 = Bbox{ll: Vector{x: newBigNumber(2), y: newBigNumber(2)}, ur: Vector{x: newBigNumber(2), y: newBigNumber(2)}}
	p2 = Bbox{ll: Vector{x: newBigNumber(4), y: newBigNumber(6)}, ur: Vector{x: newBigNumber(4), y: newBigNumber(6)}}
	expect(t, p1.getBboxOverlap(p2) == nil)

	// Vector overlap
	p = Bbox{ll: Vector{x: newBigNumber(5), y: newBigNumber(5)}, ur: Vector{x: newBigNumber(5), y: newBigNumber(5)}}
	expect(t, equalBbox(*p.getBboxOverlap(p), p))
}
