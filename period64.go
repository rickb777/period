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

// NewOf converts a time duration to a Period64.
// The result just a number of seconds, possibly including a fraction. It is not normalised; see Normalise().
func NewOf(duration time.Duration) Period64 {
	seconds := decimal.MustNew(int64(duration), 9).Trim(0)
	return Period64{seconds: seconds}.NormaliseSign()
}

//-------------------------------------------------------------------------------------------------

// Between converts the span between two times to a period. Based on the Gregorian conversion
// algorithms of `time.Time`, the resultant period is precise.
//
// If t2 is before t1, the result is a negative period.
//
// The result just a number of seconds, possibly including a fraction. It is not normalised; see Normalise().
//
// Remember that the resultant period does not retain any knowledge of the calendar, so any subsequent
// computations applied to the period can only be precise if they concern either the date (year, month,
// day) part, or the clock (hour, minute, second) part, but not both.
func Between(t1, t2 time.Time) Period64 {
	if t2.Before(t1) {
		return NewOf(t2.Sub(t1))
	}

	return NewOf(t1.Sub(t2)).Negate()
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
	buf := &strings.Builder{}
	_, _ = period.WriteTo(buf)
	return buf.String()
}

// WriteTo converts the period to ISO-8601 form.
func (period Period64) WriteTo(w io.Writer) (int64, error) {
	aw := adapt(w)

	if period == Zero {
		_, _ = aw.WriteString(string(CanonicalZero))
		return uwSum(aw)
	}

	if period.neg {
		_ = aw.WriteByte('-')
	}

	_ = aw.WriteByte('P')

	writeField(aw, period.years, Year)
	writeField(aw, period.months, Month)
	writeField(aw, period.weeks, Week)
	writeField(aw, period.days, Day)

	if period.hours.Coef() != 0 || period.minutes.Coef() != 0 || period.seconds.Coef() != 0 {
		_ = aw.WriteByte('T')

		writeField(aw, period.hours, Hour)
		writeField(aw, period.minutes, Minute)
		writeField(aw, period.seconds, Second)
	}

	return uwSum(aw)
}

func writeField(w usefulWriter, field decimal.Decimal, fieldDesignator designator) {
	if field.Coef() != 0 {
		_, _ = w.WriteString(field.String())
		_ = w.WriteByte(fieldDesignator.Byte())
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

// DaysIncWeeks gets the number of days in the period, including all the weeks and including any
// fraction present. The result is d + (w * 7), given d days and w weeks.
func (period Period64) DaysIncWeeks() decimal.Decimal {
	wdays, _ := period.weeks.Mul(seven)
	days, _ := wdays.Add(period.days)
	return period.applySign(days)
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
