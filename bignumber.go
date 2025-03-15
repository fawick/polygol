package polygol

import (
	"github.com/cockroachdb/apd/v3"
)

var context = apd.BaseContext.WithPrecision(50)

func init() {
	context.Rounding = apd.RoundHalfUp
}

type BigNumber struct {
	f *apd.Decimal
}

func (b BigNumber) String() string {
	return b.f.Text('f')
}

func newBigNumber(f float64) BigNumber {
	b := BigNumber{f: new(apd.Decimal)}
	b.f.SetFloat64(f)
	return b
}

func bigZero() BigNumber {
	b := BigNumber{f: new(apd.Decimal)}
	// b.f.SetPrec(1024)
	return b
}

func (b BigNumber) plus(other BigNumber) BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Add(r.f, b.f, other.f)
	return r
}

func (b BigNumber) minus(other BigNumber) BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Sub(r.f, b.f, other.f)
	return r
}

func (b BigNumber) times(other BigNumber) BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Mul(r.f, b.f, other.f)
	return r
}

func (b BigNumber) div(other BigNumber) BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Quo(r.f, b.f, other.f)
	return r

}
func (b BigNumber) abs() BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Abs(r.f, b.f)
	return r
}

func (b BigNumber) Cmp(other BigNumber) int {
	return b.f.Cmp(other.f)
}

func (b BigNumber) equalTo(other BigNumber) bool {
	return b.f.Cmp(other.f) == 0
}

func (b BigNumber) isLessThan(other BigNumber) bool {
	return b.f.Cmp(other.f) == -1
}

func (b BigNumber) isLessThanOrEqualTo(other BigNumber) bool {
	return b.f.Cmp(other.f) != 1
}

func (b BigNumber) isGreaterThanOrEqualTo(other BigNumber) bool {
	return b.f.Cmp(other.f) != -1
}

func (b BigNumber) isGreaterThan(other BigNumber) bool {
	return b.f.Cmp(other.f) == 1
}

func (b BigNumber) notEqualTo(other BigNumber) bool {
	return b.f.Cmp(other.f) != 0
}

func (b BigNumber) closeTo(other BigNumber) bool {
	return b.minus(other).abs().isLessThanOrEqualTo(newBigNumber(2e-16))
}

func (b BigNumber) number() float64 {
	f, _ := b.f.Float64()
	return f
}

func (b BigNumber) sqrt() BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	context.Sqrt(r.f, b.f)
	return r
}

func (b BigNumber) negated() BigNumber {
	r := BigNumber{f: new(apd.Decimal)}
	r.f.Neg(b.f)
	return r
}

func (b BigNumber) isZero() bool {
	return b.f.Cmp(new(apd.Decimal)) == 0
}

func bigInf(setToNegative bool) BigNumber {
	if setToNegative {
		return newBigNumber(-1e99)
	}
	return newBigNumber(1e99)
}
