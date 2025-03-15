package polygol

var (
	precisionEnabled bool
	precisionEpsilon = newBigNumber(float64(7.)/3 - float64(4.)/3 - float64(1.))
)

func setPrecision(eps float64) {
	precisionEpsilon = newBigNumber(eps)
	precisionEnabled = true
}

func resetPrecision() {
	precisionEnabled = false
}

func compare(a, b BigNumber) int {
	if precisionEnabled {
		if b.minus(a).abs().isLessThanOrEqualTo(precisionEpsilon) {
			return 0
		}
	}
	return a.Cmp(b)
}
