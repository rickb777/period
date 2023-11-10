package period

import (
	"github.com/govalues/decimal"
	"time"
)

var (
	seven       = decimal.MustNew(7, 0)
	twelve      = decimal.MustNew(12, 0)
	twentyFour  = decimal.MustNew(24, 0)
	sixty       = decimal.MustNew(60, 0)
	threeSixSix = decimal.MustNew(366, 0)
	daysPerYear = decimal.MustNew(3652425, 4) // by the Gregorian rule
)

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
//
// See also NormaliseDaysToYears().
func (period Period) Normalise(precise bool) Period {
	// first phase - ripple large numbers to the left
	period.minutes, period.seconds = moveWholePartsLeft(period.minutes, period.seconds, sixty)
	period.hours, period.minutes = moveWholePartsLeft(period.hours, period.minutes, sixty)
	if !precise {
		period.days, period.hours = moveWholePartsLeft(period.days, period.hours, twentyFour)
	}
	period.weeks, period.days = moveWholePartsLeft(period.weeks, period.days, seven)
	period.years, period.months = moveWholePartsLeft(period.years, period.months, twelve)
	return period
}

// NormaliseDaysToYears tries to propagate large numbers of days (and corresponding weeks)
// to the years field. Based on the Gregorian rule, there are assumed to be 365.2425 days per year.
//
//   - Multiples of 365.2425 days become years
//
// If the calculations would lead to arithmetic errors, the current values are kept unaltered.
//
// A common use pattern would be to chain this after Normalise, i.e.
//
//	p.Normalise(false).NormaliseDaysToYears()
func (period Period) NormaliseDaysToYears() Period {
	if period.neg {
		return period.Negate().NormaliseDaysToYears().Negate()
	}

	days := period.DaysIncWeeks()

	if days.Cmp(threeSixSix) < 0 {
		return period
	}

	ey, rem, err := days.QuoRem(daysPerYear)
	if err != nil {
		return period
	}

	period.years, err = period.years.Add(ey)
	if err != nil {
		return period
	}

	period.weeks, period.days = moveWholePartsLeft(decimal.Zero, rem.Trim(0), seven)
	return period
}

func moveWholePartsLeft(larger, smaller, nd decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	if smaller.IsZero() {
		return larger, smaller
	}

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
func (period Period) Simplify(precise bool) Period {
	period.years, period.months = moveToRight(period.years, period.months, twelve)
	period.weeks, period.days = moveToRight(period.weeks, period.days, seven)
	if !precise {
		period.days, period.hours = moveToRight(period.days, period.hours, twentyFour)
	}
	period.hours, period.minutes = moveToRight(period.hours, period.minutes, sixty)
	period.minutes, period.seconds = moveToRight(period.minutes, period.seconds, sixty)
	return period
}

func moveToRight(larger, smaller, nd decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	if larger.IsZero() || isSimple(larger, smaller) {
		return larger, smaller
	}

	// first check whether it's actually simpler to keep things normalised
	lg1, sm1 := moveWholePartsLeft(larger, smaller, nd)
	if isSimple(lg1, sm1) {
		return lg1, sm1 // it's hard to beat this
	}

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
func (period Period) NormaliseSign() Period {
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

func (period Period) flipSign() Period {
	period.neg = !period.neg
	return period.negateAllFields()
}

func (period Period) negateAllFields() Period {
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
func (period Period) DurationApprox() time.Duration {
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
func (period Period) Duration() (time.Duration, bool) {
	sign := time.Duration(period.Sign())
	tdE9 := time.Duration(totalDaysApproxE9(period)) * secondsPerDay
	stE9 := totalSeconds(period)
	return sign * (tdE9 + stE9), tdE9 == 0
}

func totalDaysApproxE9(period Period) int64 {
	dd := fieldDuration(period.days, oneE9)
	ww := fieldDuration(period.weeks, 7*oneE9)
	mm := fieldDuration(period.months, daysPerMonthE6*oneE3)
	yy := fieldDuration(period.years, daysPerYearE6*oneE3)
	return dd + ww + mm + yy
}

func totalSeconds(period Period) time.Duration {
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

	// gregorianYearExtraSeconds is the extra seconds needed to convert years to days, there
	// being 365.2425 days per year by the Gregorian rule.
	gregorianYearExtraSeconds = 20952 // 0.2425 * 86,400 seconds

	oneE3 int64 = 1000
	oneE9 int64 = 1_000_000_000 // used for fractions because 0 < fraction <= 999_999_999
)
