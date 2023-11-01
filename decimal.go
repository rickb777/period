package period

import bigdecimal "github.com/shopspring/decimal"

// decimal is a decimal number limited in the range from 10^-31 to math.Int64Max * 10^31.
// Unlike decimal.Decimal, it is easy to compare.
type decimal struct {
	value int64
	exp   int32
}

func newDecimal(d bigdecimal.Decimal) decimal {
	return decimal{value: d.CoefficientInt64(), exp: d.Exponent()}
}

func (d decimal) Shift(i int32) decimal {
	d.exp += i
	return d
}

func (d decimal) Add(other decimal) decimal {
	sum := d.Decimal().Add(other.Decimal())
	return newDecimal(sum)
}

func (d decimal) Decimal() bigdecimal.Decimal {
	return bigdecimal.New(d.value, d.exp)
}

var zero = decimal{}
