package period

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Period64 holds a period of time as a set of integers, one for each field in the ISO-8601
// period and additional information to track the fraction.
// The precision is almost unlimited (int64 is used for all fields for calculations). Fractions
// can hold up to 9 decimal places, therefore the finest grain is one nanosecond.
type Period64 struct {
	// always positive values
	years, months, weeks, days, hours, minutes, seconds, fraction int64

	// true if the period is negative
	neg bool

	// the fraction applies to this field
	lastField designator
}

//-------------------------------------------------------------------------------------------------

func (p64 Period64) String() string {
	neg := p64.neg
	p64.neg = false

	if p64 == (Period64{}) {
		return "P0D"
	}

	buf := &strings.Builder{}
	if neg {
		buf.WriteByte('-')
	}

	buf.WriteByte('P')

	if writeField64(buf, p64.years, p64.fraction, Year, p64.lastField) {
		return buf.String()
	}
	if writeField64(buf, p64.months, p64.fraction, Month, p64.lastField) {
		return buf.String()
	}
	if writeField64(buf, p64.weeks, p64.fraction, Week, p64.lastField) {
		return buf.String()
	}
	if writeField64(buf, p64.days, p64.fraction, Day, p64.lastField) {
		return buf.String()
	}

	buf.WriteByte('T')

	if writeField64(buf, p64.hours, p64.fraction, Hour, p64.lastField) {
		return buf.String()
	}
	if writeField64(buf, p64.minutes, p64.fraction, Minute, p64.lastField) {
		return buf.String()
	}
	writeField64(buf, p64.seconds, p64.fraction, Second, p64.lastField)

	return buf.String()
}

func writeField64(w usefulWriter, field, fraction int64, fieldDesignator, lastField designator) bool {
	if field != 0 || (fraction != 0 && fieldDesignator == lastField) {
		fmt.Fprintf(w, "%d", field)
		if fieldDesignator == lastField {
			writeFraction(w, fraction)
		}
		w.WriteByte(fieldDesignator.Byte())
	}
	return fieldDesignator == lastField
}

func writeFraction(w usefulWriter, fraction int64) {
	if fraction > 0 {
		w.WriteByte('.')
		switch {
		case fraction < 10:
			w.WriteString("00000000")
		case fraction < 100:
			w.WriteString("0000000")
		case fraction < 1000:
			w.WriteString("000000")
		case fraction < 10000:
			w.WriteString("00000")
		case fraction < 100000:
			w.WriteString("0000")
		case fraction < 1000000:
			w.WriteString("000")
		case fraction < 10000000:
			w.WriteString("00")
		case fraction < 100000000:
			w.WriteString("0")
		}
		s := strconv.FormatInt(fraction, 10)
		i := len(s) - 1
		for ; i > 0; i-- {
			if s[i] != '0' {
				break
			}
		}
		s = s[:i+1]
		w.WriteString(s)
	}
}

type usefulWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

//-------------------------------------------------------------------------------------------------

// Abs converts a negative period to a positive one.
func (p64 Period64) Abs() Period64 {
	p64.neg = false
	return p64
}

// Negate changes the sign of the period.
func (p64 Period64) Negate() Period64 {
	p64.neg = !p64.neg
	return p64
}

// IsZero returns true if applied to a zero-length period.
func (p64 Period64) IsZero() bool {
	p64.neg = false
	return p64 == Period64{}
}

//func (p64 *Period64) toPeriod() (Period, error) {
//	var f []string
//	if p64.years > math.MaxInt16 {
//		f = append(f, "years")
//	}
//	if p64.months > math.MaxInt16 {
//		f = append(f, "months")
//	}
//	if p64.days > math.MaxInt16 {
//		f = append(f, "days")
//	}
//	if p64.hours > math.MaxInt16 {
//		f = append(f, "hours")
//	}
//	if p64.minutes > math.MaxInt16 {
//		f = append(f, "minutes")
//	}
//	if p64.seconds > math.MaxInt16 {
//		f = append(f, "seconds")
//	}
//
//	if len(f) > 0 {
//		if p64.input == "" {
//			p64.input = p64.String()
//		}
//		return Period{}, fmt.Errorf("%s: integer overflow occurred in %s", p64.input, strings.Join(f, ","))
//	}
//
//	if p64.neg {
//		return Period{
//			int16(-p64.years), int16(-p64.months), int16(-p64.days),
//			int16(-p64.hours), int16(-p64.minutes), int16(-p64.seconds),
//		}, nil
//	}
//
//	return Period{
//		int16(p64.years), int16(p64.months), int16(p64.days),
//		int16(p64.hours), int16(p64.minutes), int16(p64.seconds),
//	}, nil
//}

//func (p64 *Period64) normalise64(precise bool) *Period64 {
//	return p64.rippleUp(precise).moveFractionToRight()
//}

//func (p64 *Period64) rippleUp(precise bool) *Period64 {
//	// remember that the fields are all fixed-point 1E1
//
//	p64.minutes += (p64.seconds / 600) * 10
//	p64.seconds = p64.seconds % 600
//
//	p64.hours += (p64.minutes / 600) * 10
//	p64.minutes = p64.minutes % 600
//
//	// 32670-(32670/60)-(32670/3600) = 32760 - 546 - 9.1 = 32204.9
//	if !precise || p64.hours > 32204 {
//		p64.days += (p64.hours / 240) * 10
//		p64.hours = p64.hours % 240
//	}
//
//	if !precise || p64.days > 32760 {
//		dE6 := p64.days * oneE5
//		p64.months += (dE6 / daysPerMonthE6) * 10
//		p64.days = (dE6 % daysPerMonthE6) / oneE5
//	}
//
//	p64.years += (p64.months / 120) * 10
//	p64.months = p64.months % 120
//
//	return p64
//}

// moveFractionToRight attempts to remove fractions in higher-order fields by moving their value to the
// next-lower-order field. For example, fractional years become months.
//func (p64 *Period64) moveFractionToRight() *Period64 {
//	// remember that the fields are all fixed-point 1E1
//
//	y10 := p64.years % 10
//	if y10 != 0 && (p64.months != 0 || p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
//		p64.months += y10 * 12
//		p64.years = (p64.years / 10) * 10
//	}
//
//	m10 := p64.months % 10
//	if m10 != 0 && (p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
//		p64.days += (m10 * daysPerMonthE6) / oneE6
//		p64.months = (p64.months / 10) * 10
//	}
//
//	d10 := p64.days % 10
//	if d10 != 0 && (p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
//		p64.hours += d10 * 24
//		p64.days = (p64.days / 10) * 10
//	}
//
//	hh10 := p64.hours % 10
//	if hh10 != 0 && (p64.minutes != 0 || p64.seconds != 0) {
//		p64.minutes += hh10 * 60
//		p64.hours = (p64.hours / 10) * 10
//	}
//
//	mm10 := p64.minutes % 10
//	if mm10 != 0 && p64.seconds != 0 {
//		p64.seconds += mm10 * 60
//		p64.minutes = (p64.minutes / 10) * 10
//	}
//
//	return p64
//}
