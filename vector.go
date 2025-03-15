package polygol

import "fmt"

type Vector struct {
	x BigNumber
	y BigNumber
}

func (v Vector) String() string {
	return fmt.Sprintf("(%s, %s)", v.x, v.y)
}

func newVectorLit(x, y float64) Vector {
	return Vector{x: newBigNumber(x), y: newBigNumber(y)}
}

func crossProduct(a, b Vector) BigNumber {
	return a.x.times(b.y).minus(a.y.times(b.x))
}

func dotProduct(a, b Vector) BigNumber {
	return a.x.times(b.x).plus(a.y.times(b.y))
}

func length(v Vector) BigNumber {
	return dotProduct(v, v).sqrt()
}

func sineOfAngle(pShared, pBase, pAngle Vector) BigNumber {
	vBase := Vector{x: pBase.x.minus(pShared.x), y: pBase.y.minus(pShared.y)}
	vAngle := Vector{x: pAngle.x.minus(pShared.x), y: pAngle.y.minus(pShared.y)}
	return crossProduct(vAngle, vBase).div(length(vAngle)).div(length(vBase))
}

func cosineOfAngle(pShared, pBase, pAngle Vector) BigNumber {
	vBase := Vector{x: pBase.x.minus(pShared.x), y: pBase.y.minus(pShared.y)}
	vAngle := Vector{x: pAngle.x.minus(pShared.x), y: pAngle.y.minus(pShared.y)}
	return dotProduct(vAngle, vBase).div(length(vAngle)).div(length(vBase))
}

func perpendicular(v Vector) Vector {
	return Vector{x: v.y.negated(), y: v.x}
}

func verticalIntersection(pt, v Vector, x BigNumber) *Vector {
	if v.x.isZero() {
		return nil
	}
	return &Vector{x: x, y: pt.y.plus((v.y.div(v.x)).times(x.minus(pt.x)))}
}

func horizontalIntersection(pt, v Vector, y BigNumber) *Vector {
	if v.y.isZero() {
		return nil
	}
	return &Vector{x: pt.x.plus((v.x.div(v.y)).times(y.minus(pt.y))), y: y}
}

func intersection(pt1, v1, pt2, v2 Vector) *Vector {
	// take some shortcuts for vertical and horizontal lines
	// this also ensures we don't calculate an intersection and then discover
	// it's actually outside the bounding box of the line
	if v1.x.isZero() {
		return verticalIntersection(pt2, v2, pt1.x)
	}
	if v2.x.isZero() {
		return verticalIntersection(pt1, v1, pt2.x)
	}
	if v1.y.isZero() {
		return horizontalIntersection(pt2, v2, pt1.y)
	}
	if v2.y.isZero() {
		return horizontalIntersection(pt1, v1, pt2.y)
	}

	// General case for non-overlapping segments.
	// This algorithm is based on Schneider and Eberly.
	// http://www.cimec.org.ar/~ncalvo/Schneider_Eberly.pdf - pg 244

	kross := crossProduct(v1, v2)
	if kross.isZero() {
		return nil
	}

	ve := Vector{x: pt2.x.minus(pt1.x), y: pt2.y.minus(pt1.y)}
	d1 := crossProduct(ve, v1).div(kross)
	d2 := crossProduct(ve, v2).div(kross)

	// take the average of the two calculations to minimize rounding error
	x1 := pt1.x.plus(d2.times(v1.x))
	x2 := pt2.x.plus(d1.times(v2.x))
	y1 := pt1.y.plus(d2.times(v1.y))
	y2 := pt2.y.plus(d1.times(v2.y))
	x := x1.plus(x2).div(newBigNumber(2))
	y := y1.plus(y2).div(newBigNumber(2))
	return &Vector{x: x, y: y}
}
