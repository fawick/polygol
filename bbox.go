package polygol

type Bbox struct {
	ll Vector
	ur Vector
}

func (bbox Bbox) isInBbox(point Vector) bool {
	return bbox.ll.x.isLessThanOrEqualTo(point.x) &&
		point.x.isLessThanOrEqualTo(bbox.ur.x) &&
		bbox.ll.y.isLessThanOrEqualTo(point.y) &&
		point.y.isLessThanOrEqualTo(bbox.ur.y)

}

func (b1 Bbox) getBboxOverlap(b2 Bbox) *Bbox {
	// check if the bboxes overlap at all
	if b2.ur.x.isLessThan(b1.ll.x) ||
		b1.ur.x.isLessThan(b2.ll.x) ||
		b2.ur.y.isLessThan(b1.ll.y) ||
		b1.ur.y.isLessThan(b2.ll.y) {
		return nil
	}

	// find the middle two X values
	lowerX := b1.ll.x
	if b1.ll.x.isLessThan(b2.ll.x) {
		lowerX = b2.ll.x
	}
	upperX := b2.ur.x
	if b1.ur.x.isLessThan(b2.ur.x) {
		upperX = b1.ur.x
	}

	// find the middle two Y values
	lowerY := b1.ll.y
	if b1.ll.y.isLessThan(b2.ll.y) {
		lowerY = b2.ll.y
	}
	upperY := b2.ur.y
	if b1.ur.y.isLessThan(b2.ur.y) {
		upperY = b1.ur.y
	}

	return &Bbox{
		ll: Vector{x: lowerX, y: lowerY},
		ur: Vector{x: upperX, y: upperY},
	}
}
