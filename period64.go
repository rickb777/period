// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	bigdecimal "github.com/shopspring/decimal"
	"io"
	"strings"
	"time"
)

// Period64 holds a period of time as a set of integers, one for each field in the ISO-8601
// period, and additional information to track any fraction.
// Conventionally, all the fields normally have the same sign. However, this is not restricted,
// so each field after the first non-zero field can be independently positive or negative.
// The precision is almost unlimited (int64 is used for all fields for calculations).
// However, the fraction part can have up to 9 digits, allowing at best nanosecond resolution.
//
// Instances are immutable.
type Period64 struct {
	years, months, weeks, days, hours, minutes, seconds decimal

	// neg indicates a negative period, which negates all fields (even if they are already negative)
	neg bool
}

// Zero is the zero length period.
var Zero = Period64{}

//-------------------------------------------------------------------------------------------------

// Period converts the period to ISO-8601 string form.
// If there is a decimal fraction, it will be rendered using a decimal point separator
// (not a comma).
func (period Period64) Period() Period {
	return Period(period.String())
}

// String converts the period to ISO-8601 string form.
// If there is a decimal fraction, it will be rendered using a decimal point separator.
// (not a comma).
func (period Period64) String() string {
	if period == Zero {
		return "P0D"
	}

	buf := &strings.Builder{}
	period.WriteTo(buf)
	return buf.String()
}

// WriteTo converts the period to ISO-8601 form.
func (period Period64) WriteTo(w io.Writer) (int64, error) {
	buf := adapt(w)

	if period.neg {
		buf.WriteByte('-')
	}

	buf.WriteByte('P')

	writeField(buf, period.years, Year)
	writeField(buf, period.months, Month)
	writeField(buf, period.weeks, Week)
	writeField(buf, period.days, Day)

	if period.hours.value == 0 && period.minutes.value == 0 && period.seconds.value == 0 {
		return uwSum(buf)
	}

	buf.WriteByte('T')

	writeField(buf, period.hours, Hour)
	writeField(buf, period.minutes, Minute)
	writeField(buf, period.seconds, Second)

	return uwSum(buf)
}

func writeField(w usefulWriter, field decimal, fieldDesignator designator) {
	if field.value != 0 {
		w.WriteString(field.String())
		w.WriteByte(fieldDesignator.Byte())
	}
}

//-------------------------------------------------------------------------------------------------

// IsZero returns true if applied to a period of zero length.
func (period Period64) IsZero() bool {
	period.neg = false
	return period == Zero
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
// If the period has not had its sign normalised, the result is
// undefined (see NormaliseSign or Normalise).
func (period Period64) Sign() int {
	switch {
	case period.neg:
		return -1
	case period != Zero:
		return 1
	default:
		return 0
	}
}

// IsNegative returns true if the period is negative.
// If the period has not had its sign normalised, the result is
// undefined (see NormaliseSign or Normalise).
func (period Period64) IsNegative() bool {
	return period.neg
}

// IsPositive returns true if the period is positive or zero.
// If the period has not had its sign normalised, the result is
// undefined (see NormaliseSign or Normalise).
func (period Period64) IsPositive() bool {
	return !period.neg
}

// Abs converts a negative period to a positive period.
func (period Period64) Abs() Period64 {
	period.neg = false
	return period
}

// Negate changes the sign of the period. Zero is not altered.
func (period Period64) Negate() Period64 {
	if period.IsZero() {
		return Zero
	}
	period.neg = !period.neg
	return period
}

//-------------------------------------------------------------------------------------------------

// Normalise attempts to simplify the fields. It operates in either precise or imprecise mode.
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time to or from the days element, or to transfer days to or from the months
// element. To give control over this, there are two modes.
//
// In precise mode:
// Multiples of 60 seconds become minutes.
// Multiples of 60 minutes become hours.
// Multiples of 7 days become weeks.
// Multiples of 12 months become years.
//
// Additionally, in imprecise mode:
// Multiples of 24 hours become days.
//
// Note that leap seconds are disregarded: every minute is assumed to have 60 seconds.
func (period Period64) Normalise(precise bool) Period64 {
	period.minutes, period.seconds = normaliseN(period.minutes, period.seconds, 60)
	period.hours, period.minutes = normaliseN(period.hours, period.minutes, 60)
	if !precise {
		period.days, period.hours = normaliseN(period.days, period.hours, 24)
	}
	period.weeks, period.days = normaliseN(period.weeks, period.days, 7)
	period.years, period.months = normaliseN(period.years, period.months, 12)
	return period.NormaliseSign()
}

func normaliseN(larger, smaller decimal, n int64) (decimal, decimal) {
	if smaller.value == 0 {
		return larger, smaller
	}

	sd := smaller.Decimal()
	nd := bigdecimal.NewFromInt(n)

	q := sd.Div(nd)

	ip := q.IntPart()
	if ip == 0 {
		return larger, smaller // no change
	}

	r := bigdecimal.NewFromInt(ip).Mul(nd)

	l2 := larger.Add(decimal{value: ip})

	s2 := sd.Sub(r)
	if s2.IntPart() == 0 && !s2.IsZero() {
		return larger, smaller // no change
	}

	return l2, newDecimal(s2)
}

// NormaliseSign swaps the signs of all fields so that the largest non-zero field is positive and the overall sign
// indicates the original sign. Otherwise it has no effect.
func (period Period64) NormaliseSign() Period64 {
	if period.years.value > 0 {
		return period
	} else if period.years.value < 0 {
		return period.flipSign()
	}

	if period.months.value > 0 {
		return period
	} else if period.months.value < 0 {
		return period.flipSign()
	}

	if period.weeks.value > 0 {
		return period
	} else if period.weeks.value < 0 {
		return period.flipSign()
	}

	if period.days.value > 0 {
		return period
	} else if period.days.value < 0 {
		return period.flipSign()
	}

	if period.hours.value > 0 {
		return period
	} else if period.hours.value < 0 {
		return period.flipSign()
	}

	if period.minutes.value > 0 {
		return period
	} else if period.minutes.value < 0 {
		return period.flipSign()
	}

	if period.seconds.value > 0 {
		return period
	} else if period.seconds.value < 0 {
		return period.flipSign()
	}

	return Zero
}

func (period Period64) flipSign() Period64 {
	period.neg = !period.neg
	return period.negateAllFields()
}

func (period Period64) negateAllFields() Period64 {
	period.years.value = -period.years.value
	period.months.value = -period.months.value
	period.weeks.value = -period.weeks.value
	period.days.value = -period.days.value
	period.hours.value = -period.hours.value
	period.minutes.value = -period.minutes.value
	period.seconds.value = -period.seconds.value
	return period
}

//-------------------------------------------------------------------------------------------------

// DurationApprox converts a period to the equivalent duration in nanoseconds.
// When the period specifies hours, minutes and seconds only, the result is precise.
// however, when the period specifies years, months, weeks and days, it is impossible to
// be precise because the result may depend on knowing date and timezone information. So
// the duration is estimated on the basis of a year being 365.2425 days (as per Gregorian
// calendar rules) and a month being 1/12 of a that; days are all assumed to be 24 hours long.
func (period Period64) DurationApprox() time.Duration {
	d, _ := period.Duration()
	return d
}

// Duration converts a period to the equivalent duration in nanoseconds.
// A flag is also returned that is true when the conversion was precise, and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
// However, when the period specifies years, months, weeks and days, it is impossible to
// be precise because the result may depend on knowing date and timezone information. So
// the duration is estimated on the basis of a year being 365.2425 days (as per Gregorian
// calendar rules) and a month being 1/12 of a that; days are all assumed to be 24 hours long.
func (period Period64) Duration() (time.Duration, bool) {
	sign := time.Duration(period.Sign())
	tdE9 := time.Duration(totalDaysApproxE9(period)) * secondsPerDay
	stE9 := totalSeconds(period)
	return sign * (tdE9 + stE9), tdE9 == 0
}

func totalDaysApproxE9(period Period64) int64 {
	dd := fieldDuration(period.days, oneE9)
	ww := fieldDuration(period.weeks, 7*oneE9)
	mm := fieldDuration(period.months, daysPerMonthE6*oneE3)
	yy := fieldDuration(period.years, daysPerYearE6*oneE3)
	return dd + ww + mm + yy
}

func totalSeconds(period Period64) time.Duration {
	hh := fieldDuration(period.hours, int64(time.Hour))
	mm := fieldDuration(period.minutes, int64(time.Minute))
	ss := fieldDuration(period.seconds, int64(time.Second))
	return time.Duration(hh + mm + ss)
}

func fieldDuration(field decimal, factor int64) (d int64) {
	if field.value == 0 {
		return 0
	}

	for i := field.exp; i < 0; i++ {
		factor /= 10
	}

	if factor != 0 {
		d += field.value * factor
	}

	return d
}

const (
	secondsPerDay = 24 * 60 * 60 // assuming 24-hour day

	daysPerYearE6  = 365242500          // 365.2425 days by the Gregorian rule
	daysPerMonthE6 = daysPerYearE6 / 12 // 30.436875 days per month

	oneE3 int64 = 1000
	oneE9 int64 = 1_000_000_000 // used for fractions because 0 < fraction <= 999_999_999
)
