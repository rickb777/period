// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import "errors"

// Subtract subtracts one period from another.
// Arithmetic overflow will result in an error.
func (period Period64) Subtract(other Period64) (Period64, error) {
	return period.Add(other.Negate())
}

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
// Arithmetic overflow will result in an error.
func (period Period64) Add(other Period64) (Period64, error) {
	var left, right Period64

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

	result := Period64{years: years, months: months, weeks: weeks, days: days, hours: hours, minutes: minutes, seconds: seconds}.Normalise(true).NormaliseSign()
	return result, errors.Join(e1, e2, e3, e4, e5, e6, e7)
}

//-------------------------------------------------------------------------------------------------

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
//func (period Period) AddTo(t time.Time) (time.Time, bool) {
//	wholeYears := (period.years % 10) == 0
//	wholeMonths := (period.months % 10) == 0
//	wholeWeeks := (period.weeks % 10) == 0
//	wholeDays := (period.days % 10) == 0
//
//	if wholeYears && wholeMonths && wholeWeeks && wholeDays {
//		// in this case, time.AddDate provides an exact solution
//		stE3 := totalSeconds(period)
//		t1 := t.AddDate(int(period.years/10), int(period.months/10), 7*int(period.weeks/10)+int(period.days/10))
//		return t1.Add(stE3 * time.Millisecond), true
//	}
//
//	d, precise := period.Duration()
//	return t.Add(d), precise
//}

//-------------------------------------------------------------------------------------------------

// Scale a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if the factor is negative. The result is normalised, but integer overflows
// are silently ignored.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
//func (period Period) Scale(factor float32) Period {
//	result, _ := period.ScaleWithOverflowCheck(factor)
//	return result
//}

// ScaleWithOverflowCheck scales a period by a multiplication factor. Obviously, this can both
// enlarge and shrink it, and change the sign if negative. The result is normalised. An error
// is returned if integer overflow happened.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with one
// decimal place; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
//func (period Period) ScaleWithOverflowCheck(factor float32) (Period, error) {
//	ap, neg := period.absNeg()
//
//	if -0.5 < factor && factor < 0.5 {
//		d, pr1 := ap.Duration()
//		mul := float64(d) * float64(factor)
//		p2, pr2 := NewOf(time.Duration(mul))
//		return p2.Normalise(pr1 && pr2), nil
//	}
//
//	y := int64(float32(ap.years) * factor)
//	m := int64(float32(ap.months) * factor)
//	w := int64(float32(ap.weeks) * factor)
//	d := int64(float32(ap.days) * factor)
//	hh := int64(float32(ap.hours) * factor)
//	mm := int64(float32(ap.minutes) * factor)
//	ss := int64(float32(ap.seconds) * factor)
//
//	p64 := &period64{years: y, months: m, weeks: w, days: d, hours: hh, minutes: mm, seconds: ss, neg: neg, denormal: true}
//	n64 := p64.normalise64(true)
//	return n64.toPeriod(), n64.checkOverflow()
//}

// RationalScale scales a period by a rational multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised. An error is returned if integer overflow
// happened.
//
// If the divisor is zero, a panic will arise.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
//func (period Period) RationalScale(multiplier, divisor int) (Period, error) {
//	return period.rationalScale64(int64(multiplier), int64(divisor))
//}

// moveFractionToRight attempts to remove fractions in higher-order fields by moving their value to the
// next-lower-order field. For example, fractional years become months.
//func (period *Period64) moveFractionToRight() *Period64 {
//	// remember that the fields are all fixed-point 1E1
//
//	if period.lastField == Year && (period.fraction != 0) {
//		f := int64(period.fraction) * 12
//		period.lastField = Month
//		period.months = int32(f / 1_000_000_000)
//		period.fraction = int32(f % 1_000_000_000)
//	}
//
//	//m10 := period.months % 10
//	//if m10 != 0 && (period.weeks != 0 || period.days != 0 || period.hours != 0 || period.minutes != 0 || period.seconds != 0) {
//	//	period.weeks += (m10 * weeksPerMonthE6) / oneE6
//	//	period.months = (period.months / 10) * 10
//	//}
//
//	if period.lastField == Week && (period.fraction != 0) {
//		f := int64(period.fraction) * 7
//		period.lastField = Day
//		period.days = int32(f / 1_000_000_000)
//		period.fraction = int32(f % 1_000_000_000)
//	}
//
//	//d10 := period.days % 10
//	//if d10 != 0 && (period.hours != 0 || period.minutes != 0 || period.seconds != 0) {
//	//	period.hours += d10 * 24
//	//	period.days = (period.days / 10) * 10
//	//}
//
//	if period.lastField == Hour && (period.fraction != 0) {
//		f := int64(period.fraction) * 60
//		period.lastField = Minute
//		period.minutes = int32(f / 1_000_000_000)
//		period.fraction = int32(f % 1_000_000_000)
//	}
//
//	if period.lastField == Minute && (period.fraction != 0) {
//		f := int64(period.fraction) * 60
//		period.lastField = Second
//		period.minutes = int32(f / 1_000_000_000)
//		period.fraction = int32(f % 1_000_000_000)
//	}
//
//	return period
//}
