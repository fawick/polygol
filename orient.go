package polygol

func orient(a, b, c Vector) int {
	ax := a.x
	ay := a.y
	cx := c.x
	cy := c.y
	area2 := ay.minus(cy).times(b.x.minus(cx)).minus(ax.minus(cx).times(b.y.minus(cy)))
	if precisionEnabled {
		l := cx.minus(ax)
		r := cy.minus(ay)
		if area2.times(area2).isLessThanOrEqualTo(
			l.times(l).plus(r.times(r)).times(precisionEpsilon)) {
			return 0
		}
	}
	return area2.Cmp(bigZero())
}
