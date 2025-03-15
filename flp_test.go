package polygol

import "testing"

func TestFlpCompare(t *testing.T) {
	setPrecision(NumberEPSILON)

	var a, b BigNumber
	// exactly equal
	a = newBigNumber(1)
	b = newBigNumber(1)
	expect(t, compare(a, b) == 0)

	// flp equal
	a = newBigNumber(1)
	b = newBigNumber(1).plus(newBigNumber(NumberEPSILON))
	expect(t, compare(a, b) == 0)

	// barely less than
	a = newBigNumber(1)
	b = newBigNumber(1).plus(newBigNumber(NumberEPSILON).times(newBigNumber(2)))
	expect(t, compare(a, b) == -1)

	// less than
	a = newBigNumber(1)
	b = newBigNumber(2)
	expect(t, compare(a, b) == -1)

	// barely more than
	a = newBigNumber(1).plus(newBigNumber(NumberEPSILON).times(newBigNumber(2)))
	b = newBigNumber(1)
	expect(t, compare(a, b) == 1)

	// more than
	a = newBigNumber(2)
	b = newBigNumber(1)
	expect(t, compare(a, b) == 1)

	// both flp equal to 0
	a = bigZero()
	b = newBigNumber(NumberEPSILON).minus(newBigNumber(NumberEPSILON).times(newBigNumber(NumberEPSILON)))
	expect(t, compare(a, b) == 0)

	// really close to 0
	a = newBigNumber(NumberEPSILON)
	b = newBigNumber(NumberEPSILON).plus(newBigNumber(NumberEPSILON).times(newBigNumber(NumberEPSILON)).times(newBigNumber(2)))
	expect(t, compare(a, b) == 0)
	resetPrecision()
}
