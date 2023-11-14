// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"errors"
	"github.com/govalues/decimal"
	"time"
)

// AddTo adds the period to a time, returning the result.
// A flag is also returned that is true when the conversion was precise, and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
//
// Similarly, when the period specifies whole years, months, weeks and days (i.e. without fractions),
// the result is precise.
//
// However, when years, months or days contains fractions, the result is only an approximation (it
// assumes that all days are 24 hours and every year is 365.2425 days, as per Gregorian calendar rules).
func (period Period) AddTo(t time.Time) (time.Time, bool) {
	wholeYears := period.years.Scale() == 0
	wholeMonths := period.months.Scale() == 0
	wholeWeeks := period.weeks.Scale() == 0
	wholeDays := period.days.Scale() == 0

	if wholeYears && wholeMonths && wholeWeeks && wholeDays {
		// in this case, time.AddDate provides an exact solution

		years, _, ok1 := period.years.Int64(0)
		months, _, ok2 := period.months.Int64(0)
		weeks, _, ok3 := period.weeks.Int64(0)
		days, _, ok4 := period.days.Int64(0)

		hms, ok5 := totalHrMinSec(period)

		if period.neg {
			years = -years
			months = -months
			weeks = -weeks
			days = -days
			hms = -hms
		}

		t1 := t.AddDate(int(years), int(months), int(7*weeks+days)).Add(hms)
		return t1, ok1 && ok2 && ok3 && ok4 && ok5
	}

	// fractional years or months or weeks or days
	d, precise := period.Duration()
	return t.Add(d), precise
}

//-------------------------------------------------------------------------------------------------

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
// Arithmetic overflow will result in an error.
func (period Period) Add(other Period) (Period, error) {
	var left, right Period

	if period.neg {
		left = period.flipSign()
	} else {
		left = period
	}

	if other.neg {
		right = other.flipSign()
	} else {
		right = other
	}

	years, e1 := left.years.Add(right.years)
	months, e2 := left.months.Add(right.months)
	weeks, e3 := left.weeks.Add(right.weeks)
	days, e4 := left.days.Add(right.days)
	hours, e5 := left.hours.Add(right.hours)
	minutes, e6 := left.minutes.Add(right.minutes)
	seconds, e7 := left.seconds.Add(right.seconds)

	result := Period{years: years, months: months, weeks: weeks, days: days, hours: hours, minutes: minutes, seconds: seconds}.Normalise(true).normaliseSign()
	return result, errors.Join(e1, e2, e3, e4, e5, e6, e7)
}

// Subtract subtracts one period from another.
// Arithmetic overflow will result in an error.
func (period Period) Subtract(other Period) (Period, error) {
	return period.Add(other.Negate())
}

//-------------------------------------------------------------------------------------------------

// Mul multiplies a period by a factor. Obviously, this can both enlarge and shrink it,
// and change the sign if the factor is negative. The result is not normalised.
func (period Period) Mul(factor decimal.Decimal) (Period, error) {
	var years, months, weeks, days, hours, minutes, seconds decimal.Decimal
	var e1, e2, e3, e4, e5, e6, e7 error

	if period.years.Coef() != 0 {
		years, e1 = period.years.Mul(factor)
		years = years.Trim(0)
	}
	if period.months.Coef() != 0 {
		months, e2 = period.months.Mul(factor)
		months = months.Trim(0)
	}
	if period.weeks.Coef() != 0 {
		weeks, e3 = period.weeks.Mul(factor)
		weeks = weeks.Trim(0)
	}
	if period.days.Coef() != 0 {
		days, e4 = period.days.Mul(factor)
		days = days.Trim(0)
	}
	if period.hours.Coef() != 0 {
		hours, e5 = period.hours.Mul(factor)
		hours = hours.Trim(0)
	}
	if period.minutes.Coef() != 0 {
		minutes, e6 = period.minutes.Mul(factor)
		minutes = minutes.Trim(0)
	}
	if period.seconds.Coef() != 0 {
		seconds, e7 = period.seconds.Mul(factor)
		seconds = seconds.Trim(0)
	}

	result := Period{
		years:   years,
		months:  months,
		weeks:   weeks,
		days:    days,
		hours:   hours,
		minutes: minutes,
		seconds: seconds,
		neg:     period.neg,
	}

	return result.normaliseSign(), errors.Join(e1, e2, e3, e4, e5, e6, e7)
}
