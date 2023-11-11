// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
	"strings"
	"time"
)

// Period holds a period of time as a set of decimal numbers, one for each field in the ISO-8601
// period.
//
// By conventional, all the fields should have the same sign. However, this is not restricted,
// so each field after the first non-zero field can be independently positive or negative.
// Sometimes this makes sense, e.g. "P1YT-1S" is one second less than one year.
//
// The precision is large: all fields are scaled decimals using int64 internally for calculations.
// The value of each field can have up to 19 digits (the range of int64), of which up to 19 digits
// can be a decimal fraction. So the range is much wider than that of time.Duration.
//
// For convenience, the method inputs and outputs are int.
//
// Fractions are supported on the least significant non-zero field only. It is an error for
// more-significant fields to have fractional values too.
//
// Instances are immutable.
type Period struct {
	years, months, weeks, days, hours, minutes, seconds decimal.Decimal

	// neg indicates a negative period, which negates all fields (even if they are already negative)
	neg bool
}

// Zero is the zero period.
var Zero = Period{}

//-------------------------------------------------------------------------------------------------

// NewYMWD creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 12 months will not become 1 year. Use the Normalise method if you
// need to.
func NewYMWD(years, months, weeks, days int) Period {
	return New(years, months, weeks, days, 0, 0, 0)
}

// NewHMS creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise method
// if you need to.
func NewHMS(hours, minutes, seconds int) Period {
	return New(0, 0, 0, 0, hours, minutes, seconds)
}

// New creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise() method
// if you need to.
func New(years, months, weeks, days, hours, minutes, seconds int) Period {
	return Period{
		years:   decimal.MustNew(int64(years), 0),
		months:  decimal.MustNew(int64(months), 0),
		weeks:   decimal.MustNew(int64(weeks), 0),
		days:    decimal.MustNew(int64(days), 0),
		hours:   decimal.MustNew(int64(hours), 0),
		minutes: decimal.MustNew(int64(minutes), 0),
		seconds: decimal.MustNew(int64(seconds), 0)}.normaliseSign()
}

// NewDecimal creates a period from seven decimal values. The fields are trimmed but no normalisation
// is applied, e.g. 120 seconds will not become 2 minutes. Use the Normalise() method
// if you need to.
//
// Periods only allow the least-significant non-zero field to contain a fraction. If any of the
// more-significant fields is supplied with a fraction, an error will be returned. This can be safely
// ignored for non-standard behaviour.
func NewDecimal(years, months, weeks, days, hours, minutes, seconds decimal.Decimal) (period Period, err error) {
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

	return Period{
		years:   years.Trim(0),
		months:  months.Trim(0),
		weeks:   weeks.Trim(0),
		days:    days.Trim(0),
		hours:   hours.Trim(0),
		minutes: minutes.Trim(0),
		seconds: seconds.Trim(0),
	}.normaliseSign(), err
}

// NewOf converts a time duration to a Period.
// The result just a number of seconds, possibly including a fraction. It is not normalised; see Normalise().
func NewOf(duration time.Duration) Period {
	seconds := decimal.MustNew(int64(duration), 9).Trim(0)
	return Period{seconds: seconds}.normaliseSign()
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
func Between(t1, t2 time.Time) Period {
	if t2.Before(t1) {
		return NewOf(t2.Sub(t1))
	}

	return NewOf(t1.Sub(t2)).Negate()
}

//-------------------------------------------------------------------------------------------------

// IsZero returns true if applied to a period of zero length.
func (period Period) IsZero() bool {
	period.neg = false
	return period == Zero
}

// sign returns 1 if the period is zero or positive, -1 if it is negative.
func (period Period) signI() int {
	if period.neg {
		return -1
	} else {
		return 1
	}
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
func (period Period) Sign() int {
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
func (period Period) IsNegative() bool {
	return period.neg
}

// IsPositive returns true if the period is positive or zero.
func (period Period) IsPositive() bool {
	return !period.neg
}

// Abs converts a negative period to a positive period.
func (period Period) Abs() Period {
	period.neg = false
	return period
}

// Negate changes the sign of the period. Zero is not altered.
func (period Period) Negate() Period {
	if period.IsZero() {
		return Zero
	}
	period.neg = !period.neg
	return period
}

//-------------------------------------------------------------------------------------------------

// YearsInt gets the whole number of years in the period.
func (period Period) YearsInt() int {
	i, _, _ := period.years.Int64(0)
	return int(i) * period.signI()
}

// Years gets the number of years in the period, including any fraction present.
func (period Period) Years() decimal.Decimal {
	return period.applySign(period.years)
}

// MonthsInt gets the whole number of months in the period.
func (period Period) MonthsInt() int {
	i, _, _ := period.months.Int64(0)
	return int(i) * period.signI()
}

// Months gets the number of months in the period, including any fraction present.
func (period Period) Months() decimal.Decimal {
	return period.applySign(period.months)
}

// WeeksInt gets the whole number of weeks in the period.
func (period Period) WeeksInt() int {
	i, _, _ := period.weeks.Int64(0)
	return int(i) * period.signI()
}

// Weeks gets the number of weeks in the period, including any fraction present.
func (period Period) Weeks() decimal.Decimal {
	return period.applySign(period.weeks)
}

// DaysInt gets the whole number of days in the period.
func (period Period) DaysInt() int {
	i, _, _ := period.days.Int64(0)
	return int(i) * period.signI()
}

// Days gets the number of days in the period, including any fraction present.
func (period Period) Days() decimal.Decimal {
	return period.applySign(period.days)
}

// DaysIncWeeksInt gets the number of days in the period, including all the weeks but truncating
// any fractions present. The result is d + (w * 7), given d days and w weeks.
func (period Period) DaysIncWeeksInt() int {
	i, _, _ := period.DaysIncWeeks().Int64(0)
	return int(i) * period.signI()
}

// DaysIncWeeks gets the number of days in the period, including all the weeks and including any
// fraction present. The result is d + (w * 7), given d days and w weeks.
func (period Period) DaysIncWeeks() decimal.Decimal {
	wdays, _ := period.weeks.Mul(seven)
	days, _ := wdays.Add(period.days)
	return period.applySign(days)
}

// HoursInt gets the whole number of hours in the period.
func (period Period) HoursInt() int {
	i, _, _ := period.hours.Int64(0)
	return int(i) * period.signI()
}

// Hours gets the number of hours in the period, including any fraction present.
func (period Period) Hours() decimal.Decimal {
	return period.applySign(period.hours)
}

// MinutesInt gets the whole number of minutes in the period.
func (period Period) MinutesInt() int {
	i, _, _ := period.minutes.Int64(0)
	return int(i) * period.signI()
}

// Minutes gets the number of minutes in the period, including any fraction present.
func (period Period) Minutes() decimal.Decimal {
	return period.applySign(period.minutes)
}

// SecondsInt gets the whole number of seconds in the period.
func (period Period) SecondsInt() int {
	i, _, _ := period.seconds.Int64(0)
	return int(i) * period.signI()
}

// Seconds gets the number of seconds in the period, including any fraction present.
func (period Period) Seconds() decimal.Decimal {
	return period.applySign(period.seconds)
}

// Seconds gets the number of seconds in the period, including any fraction present.
func (period Period) applySign(field decimal.Decimal) decimal.Decimal {
	if period.neg {
		return field.Neg()
	}
	return field
}
