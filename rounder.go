package polygol

import (
	splaytree "github.com/engelsjk/splay-tree"
)

type ptRounder struct {
	xRounder *coordRounder
	yRounder *coordRounder
}

func newPtRounder() *ptRounder {
	ptr := new(ptRounder)
	ptr.reset()
	return ptr
}

func (pr *ptRounder) reset() {
	pr.xRounder = newCoordRounder()
	pr.yRounder = newCoordRounder()
}

func (pr *ptRounder) roundFloat(x, y float64) *point {
	return pr.round(newBigNumber(x), newBigNumber(y))
}

func (pr *ptRounder) round(x, y BigNumber) *point {
	x = pr.xRounder.round(x)
	y = pr.yRounder.round(y)
	return newPointBN(x, y)
}

type coordRounder struct {
	tree *splaytree.SplayTree
}

func newCoordRounder() *coordRounder {
	less := func(a, b interface{}) int {
		af := a.(BigNumber)
		bf := b.(BigNumber)
		return compare(af, bf)
	}
	cr := &coordRounder{
		tree: splaytree.New(less),
	}
	cr.round(bigZero())
	return cr
}

func (cr *coordRounder) round(coord BigNumber) BigNumber {
	node := cr.tree.Add(coord)
	item := node.Item().(BigNumber)

	prevNode := cr.tree.Prev(node)
	if prevNode != nil {
		prevItem := prevNode.Item().(BigNumber)
		if item.equalTo(prevItem) {
			cr.tree.Remove(coord)
			return prevItem
		}
	}

	nextNode := cr.tree.Next(node)
	if nextNode != nil {
		nextItem := nextNode.Item().(BigNumber)
		if item.equalTo(nextItem) {
			cr.tree.Remove(coord)
			return nextItem
		}
	}

	return item
}
