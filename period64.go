// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
	"io"
	"strings"
	"time"
)

// Period64 holds a period of time as a set of integers, one for each field in the ISO-8601
// period, and additional information to track any fraction.
//
// By conventional, all the fields have the same sign. However, this is not restricted,
// so each field after the first non-zero field can be independently positive or negative.
// Sometimes this makes sense, e.g. "P1DT-1S" is one second less than one day.
//
// The precision is large: all fields are scaled decimals using int64 internally for calculations, although
// the method inputs and outputs are int for convenience. Fractions are supported on the least significant
// non-zero field only.
//
// Instances are immutable.
type Period64 struct {
	years, months, weeks, days, hours, minutes, seconds decimal.Decimal

	// neg indicates a negative period, which negates all fields (even if they are already negative)
	neg bool
}

// Zero is the zero period.
var Zero = Period64{}

//-------------------------------------------------------------------------------------------------

// NewYMWD creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 12 months will not become 1 year. Use the Normalise method if you
// need to.
func NewYMWD(years, months, weeks, days int) Period64 {
	return New(years, months, weeks, days, 0, 0, 0)
}

// NewHMS creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise method
// if you need to.
func NewHMS(hours, minutes, seconds int) Period64 {
	return New(0, 0, 0, 0, hours, minutes, seconds)
}

// New creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise() method
// if you need to.
func New(years, months, weeks, days, hours, minutes, seconds int) Period64 {
	return Period64{
		years:   decimal.MustNew(int64(years), 0),
		months:  decimal.MustNew(int64(months), 0),
		weeks:   decimal.MustNew(int64(weeks), 0),
		days:    decimal.MustNew(int64(days), 0),
		hours:   decimal.MustNew(int64(hours), 0),
		minutes: decimal.MustNew(int64(minutes), 0),
		seconds: decimal.MustNew(int64(seconds), 0)}.NormaliseSign()
}

// NewDecimal creates a period from seven decimal values. The fields are trimmed but no normalisation
// is applied, e.g. 120 seconds will not become 2 minutes. Use the Normalise() method
// if you need to.
//
// Periods only allow the least-significant non-zero field to contain a fraction. If any of the
// more-significant fields is supplied with a fraction, an error will be returned. This can be safely
// ignored for non-standard behaviour.
func NewDecimal(years, months, weeks, days, hours, minutes, seconds decimal.Decimal) (period Period64, err error) {
	ymwd := make([]string, 0, 4)
	hms := make([]string, 0, 4)
	if years.Scale() > 0 {
		ymwd = append(ymwd, fmt.Sprintf("%sY", years))
	}
	if months.Scale() > 0 {
		ymwd = append(ymwd, fmt.Sprintf("%sM", months))
	}
	if weeks.Scale() > 0 {
		ymwd = append(ymwd, fmt.Sprintf("%sW", weeks))
	}
	if days.Scale() > 0 {
		ymwd = append(ymwd, fmt.Sprintf("%sD", days))
	}
	if hours.Scale() > 0 {
		hms = append(hms, fmt.Sprintf("%sH", hours))
	}
	if minutes.Scale() > 0 {
		hms = append(hms, fmt.Sprintf("%sM", minutes))
	}
	if seconds.Scale() > 0 {
		hms = append(hms, fmt.Sprintf("%sS", seconds))
	}
	if len(ymwd)+len(hms) > 1 {
		sep := ""
		if len(hms) > 0 {
			sep = "T"
		}
		err = fmt.Errorf("only the least significant field can have a fraction; found fractions in %s%s%s",
			strings.Join(ymwd, ""), sep, strings.Join(hms, ""))
	}

	return Period64{
		years:   years.Trim(0),
		months:  months.Trim(0),
		weeks:   weeks.Trim(0),
		days:    days.Trim(0),
		hours:   hours.Trim(0),
		minutes: minutes.Trim(0),
		seconds: seconds.Trim(0),
	}.NormaliseSign(), err
}

// NewOf converts a time duration to a Period64. The result just a number of seconds, possibly including
// a fraction. It is not normalised; see Normalise().
func NewOf(duration time.Duration) Period64 {
	seconds := decimal.MustNew(int64(duration), 9).Trim(0)
	return Period64{seconds: seconds}.NormaliseSign()
}

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

	if period.hours.Coef() == 0 && period.minutes.Coef() == 0 && period.seconds.Coef() == 0 {
		return uwSum(buf)
	}

	buf.WriteByte('T')

	writeField(buf, period.hours, Hour)
	writeField(buf, period.minutes, Minute)
	writeField(buf, period.seconds, Second)

	return uwSum(buf)
}

func writeField(w usefulWriter, field decimal.Decimal, fieldDesignator designator) {
	if field.Coef() != 0 {
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

// sign returns 1 if the period is zero or positive, -1 if it is negative.
func (period Period64) signI() int {
	if period.neg {
		return -1
	} else {
		return 1
	}
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
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
func (period Period64) IsNegative() bool {
	return period.neg
}

// IsPositive returns true if the period is positive or zero.
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

// YearsInt gets the whole number of years in the period.
func (period Period64) YearsInt() int {
	i, _, _ := period.years.Int64(0)
	return int(i) * period.signI()
}

// Years gets the number of years in the period, including any fraction present.
func (period Period64) Years() decimal.Decimal {
	return period.applySign(period.years)
}

// MonthsInt gets the whole number of months in the period.
func (period Period64) MonthsInt() int {
	i, _, _ := period.months.Int64(0)
	return int(i) * period.signI()
}

// Months gets the number of months in the period, including any fraction present.
func (period Period64) Months() decimal.Decimal {
	return period.applySign(period.months)
}

// WeeksInt gets the whole number of weeks in the period.
func (period Period64) WeeksInt() int {
	i, _, _ := period.weeks.Int64(0)
	return int(i) * period.signI()
}

// Weeks gets the number of weeks in the period, including any fraction present.
func (period Period64) Weeks() decimal.Decimal {
	return period.applySign(period.weeks)
}

// DaysInt gets the whole number of days in the period.
func (period Period64) DaysInt() int {
	i, _, _ := period.days.Int64(0)
	return int(i) * period.signI()
}

// Days gets the number of days in the period, including any fraction present.
func (period Period64) Days() decimal.Decimal {
	return period.applySign(period.days)
}

// HoursInt gets the whole number of hours in the period.
func (period Period64) HoursInt() int {
	i, _, _ := period.hours.Int64(0)
	return int(i) * period.signI()
}

// Hours gets the number of hours in the period, including any fraction present.
func (period Period64) Hours() decimal.Decimal {
	return period.applySign(period.hours)
}

// MinutesInt gets the whole number of minutes in the period.
func (period Period64) MinutesInt() int {
	i, _, _ := period.minutes.Int64(0)
	return int(i) * period.signI()
}

// Minutes gets the number of minutes in the period, including any fraction present.
func (period Period64) Minutes() decimal.Decimal {
	return period.applySign(period.minutes)
}

// SecondsInt gets the whole number of seconds in the period.
func (period Period64) SecondsInt() int {
	i, _, _ := period.seconds.Int64(0)
	return int(i) * period.signI()
}

// Seconds gets the number of seconds in the period, including any fraction present.
func (period Period64) Seconds() decimal.Decimal {
	return period.applySign(period.seconds)
}

// Seconds gets the number of seconds in the period, including any fraction present.
func (period Period64) applySign(field decimal.Decimal) decimal.Decimal {
	if period.neg {
		return field.Neg()
	}
	return field
}

//-------------------------------------------------------------------------------------------------

// Normalise simplifies the fields by propagating large values towards the more significant fields.
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time to or from the days element. To give control over this, there are two modes:
// it operates in either precise or imprecise mode.
//
//   - Multiples of 60 seconds become minutes - both modes.
//   - Multiples of 60 minutes become hours - both modes.
//   - Multiples of 24 hours become days - imprecise mode only
//   - Multiples of 7 days become weeks - both modes.
//   - Multiples of 12 months become years - both modes.
//
// Note that leap seconds are disregarded: every minute is assumed to have 60 seconds.
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
func (period Period64) Normalise(precise bool) Period64 {
	// first phase - ripple large numbers to the left
	period.minutes, period.seconds = moveWholePartsLeft(period.minutes, period.seconds, 60)
	period.hours, period.minutes = moveWholePartsLeft(period.hours, period.minutes, 60)
	if !precise {
		period.days, period.hours = moveWholePartsLeft(period.days, period.hours, 24)
	}
	period.weeks, period.days = moveWholePartsLeft(period.weeks, period.days, 7)
	period.years, period.months = moveWholePartsLeft(period.years, period.months, 12)
	return period
}

func moveWholePartsLeft(larger, smaller decimal.Decimal, n int64) (decimal.Decimal, decimal.Decimal) {
	if smaller.IsZero() {
		return larger, smaller
	}

	nd := decimal.MustNew(n, 0)

	q, r, err := smaller.QuoRem(nd)
	if err != nil {
		return larger, smaller
	}

	if !r.IsZero() && r.Prec() <= r.Scale() {
		return larger, smaller // more complex so no change
	}

	l2, err := larger.Add(q)
	if err != nil {
		return larger, smaller
	}

	return l2, r
}

// Simplify simplifies the fields by propagating large values towards the less significant fields.
// This is akin to converting mixed fractions to improper fractions, across the group of fields.
// However, existing values are not altered if they are a simple way of expression their period already.
//
// For example, "P2Y1M" simplifies to "P25M" but "P2Y" remains "P2Y".
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time to or from the days element. To give control over this, there are two modes:
// it operates in either precise or imprecise mode.
//
//   - Years may become multiples of 12 months if the number of months is non-zero - both modes.
//   - Weeks may become multiples of 7 days if the number of days is non-zero - both modes.
//   - Days may become multiples of 24 hours if the number of hours is non-zero - imprecise mode only
//   - Hours may become multiples of 60 minutes if the number of minutes is non-zero - both modes.
//   - Minutes may become multiples of 60 seconds if the number of seconds is non-zero - both modes.
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
func (period Period64) Simplify(precise bool) Period64 {
	period.years, period.months = moveToRight(period.years, period.months, 12)
	period.weeks, period.days = moveToRight(period.weeks, period.days, 7)
	if !precise {
		period.days, period.hours = moveToRight(period.days, period.hours, 24)
	}
	period.hours, period.minutes = moveToRight(period.hours, period.minutes, 60)
	period.minutes, period.seconds = moveToRight(period.minutes, period.seconds, 60)
	return period
}

func moveToRight(larger, smaller decimal.Decimal, n int64) (decimal.Decimal, decimal.Decimal) {
	if larger.IsZero() || isSimple(larger, smaller) {
		return larger, smaller
	}

	// first check whether it's actually simpler to keep things normalised
	lg1, sm1 := moveWholePartsLeft(larger, smaller, n)
	if isSimple(lg1, sm1) {
		return lg1, sm1 // it's hard to beat this
	}

	//extraDigits := int64(larger.Sign()) * int64(larger.Coef()) * n
	//extra, err := decimal.New(extraDigits, larger.Scale())
	nd := decimal.MustNew(n, 0)
	extra, err := larger.Mul(nd)
	if err != nil {
		return larger, smaller
	}

	sm2, err := smaller.Add(extra)
	if err != nil {
		return larger, smaller
	}

	sm2 = sm2.Trim(0)

	originalDigits := larger.Prec() + smaller.Prec()
	if sm2.Prec() > originalDigits {
		return larger, smaller // because we would just add more digits
	}

	return decimal.Zero, sm2
}

func isSimple(larger, smaller decimal.Decimal) bool {
	return smaller.IsZero() && larger.Scale() == 0
}

// NormaliseSign swaps the signs of all fields so that the largest non-zero field is positive and the overall sign
// indicates the original sign. Otherwise it has no effect.
func (period Period64) NormaliseSign() Period64 {
	if period.years.Sign() > 0 {
		return period
	} else if period.years.Sign() < 0 {
		return period.flipSign()
	}

	if period.months.Sign() > 0 {
		return period
	} else if period.months.Sign() < 0 {
		return period.flipSign()
	}

	if period.weeks.Sign() > 0 {
		return period
	} else if period.weeks.Sign() < 0 {
		return period.flipSign()
	}

	if period.days.Sign() > 0 {
		return period
	} else if period.days.Sign() < 0 {
		return period.flipSign()
	}

	if period.hours.Sign() > 0 {
		return period
	} else if period.hours.Sign() < 0 {
		return period.flipSign()
	}

	if period.minutes.Sign() > 0 {
		return period
	} else if period.minutes.Sign() < 0 {
		return period.flipSign()
	}

	if period.seconds.Sign() > 0 {
		return period
	} else if period.seconds.Sign() < 0 {
		return period.flipSign()
	}

	return Zero
}

func (period Period64) flipSign() Period64 {
	period.neg = !period.neg
	return period.negateAllFields()
}

func (period Period64) negateAllFields() Period64 {
	period.years = period.years.Neg()
	period.months = period.months.Neg()
	period.weeks = period.weeks.Neg()
	period.days = period.days.Neg()
	period.hours = period.hours.Neg()
	period.minutes = period.minutes.Neg()
	period.seconds = period.seconds.Neg()
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

func fieldDuration(field decimal.Decimal, factor int64) (d int64) {
	if field.Coef() == 0 {
		return 0
	}

	for i := field.Scale(); i > 0; i-- {
		factor /= 10
	}

	if factor != 0 {
		d += int64(field.Sign()) * int64(field.Coef()) * factor
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
