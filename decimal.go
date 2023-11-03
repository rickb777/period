package period

import (
	"bytes"
	bigdecimal "github.com/shopspring/decimal"
	"strconv"
)

// decimal is a decimal number limited in the range from 10^-31 to math.Int64Max * 10^31.
// Unlike decimal.Decimal (github.com/shopspring/decimal), it is easy to compare.
type decimal struct {
	value int64
	exp   int32
}

func newFromInt(v int64) decimal {
	return decimal{value: v}.Canonical()
}

func newDecimal(d bigdecimal.Decimal) decimal {
	return decimal{value: d.CoefficientInt64(), exp: d.Exponent()}.Canonical()
}

func (d decimal) Shift(i int32) decimal {
	d.exp += i
	return d
}

func (d decimal) Add(other decimal) (sum decimal) {
	if d.exp > other.exp {
		for i := d.exp - other.exp; i > 0; i-- {
			d.value *= 10
		}
		sum.exp = other.exp
	} else if d.exp < other.exp {
		for i := other.exp - d.exp; i > 0; i-- {
			other.value *= 10
		}
		sum.exp = d.exp
	}
	sum.value = d.value + other.value
	return sum.Canonical()
}

func (d decimal) Canonical() decimal {
	for {
		if d.value == 0 || d.value%10 != 0 {
			return d
		} else {
			d.value /= 10
			d.exp++
		}
	}
}

func (d decimal) Decimal() bigdecimal.Decimal {
	return bigdecimal.New(d.value, d.exp)
}

func (d decimal) String() string {
	s := strconv.FormatInt(d.value, 10)

	if d.exp > 0 {
		bs := bytes.NewBufferString(s)
		for i := d.exp; i > 0; i-- {
			bs.WriteByte('0')
		}
		s = bs.String()

	} else if d.exp < 0 {
		dp := len(s) + int(d.exp)
		bs := &bytes.Buffer{}
		if dp <= 0 {
			bs.WriteByte('0')
			bs.WriteByte(DecimalPoint)
			for ; dp < 0; dp++ {
				bs.WriteByte('0')
			}
		} else {
			bs.WriteString(s[:dp])
			bs.WriteByte(DecimalPoint)
		}
		bs.WriteString(s[dp:])
		s = bs.String()
	}

	return s
}

var zero = decimal{}

// DecimalPoint is used when rendering decimal numbers containing fractional parts.
// It can be either '.' or ',' and is initialised to '.'.
var DecimalPoint uint8 = '.'
