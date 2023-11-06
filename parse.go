// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
)

// MustParse is as per Parse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParse[S Period | string](isoPeriod S) Period64 {
	p, err := Parse(isoPeriod)
	if err != nil {
		panic(err)
	}
	return p
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func Parse[S Period | string](isoPeriod S) (Period64, error) {
	p := Period64{}
	err := p.Parse(string(isoPeriod))
	return p, err
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func (period *Period64) Parse(isoPeriod string) error {
	if isoPeriod == "" {
		return fmt.Errorf(`cannot parse a blank string as a period`)
	}

	p := Zero

	remaining := isoPeriod
	if remaining[0] == '-' {
		p.neg = true
		remaining = remaining[1:]
	} else if remaining[0] == '+' {
		remaining = remaining[1:]
	}

	switch remaining {
	case "P0Y", "P0M", "P0W", "P0D", "PT0H", "PT0M", "PT0S":
		*period = Zero
		return nil // zero case
	case "":
		return fmt.Errorf(`cannot parse a blank string as a period`)
	}

	if remaining[0] != 'P' {
		return fmt.Errorf("%s: expected 'P' period mark at the start", isoPeriod)
	}
	remaining = remaining[1:]

	var haveFraction bool
	var number decimal.Decimal
	var years, months, weeks, days, hours, minutes, seconds itemState
	var des, previous designator
	var err error
	nComponents := 0

	years, months, weeks, days = Armed, Armed, Armed, Armed

	isHMS := false
	for len(remaining) > 0 {
		if remaining[0] == 'T' {
			if isHMS {
				return fmt.Errorf("%s: 'T' designator cannot occur more than once", isoPeriod)
			}
			isHMS = true

			years, months, weeks, days = Unready, Unready, Unready, Unready
			hours, minutes, seconds = Armed, Armed, Armed

			remaining = remaining[1:]

		} else {
			number, des, remaining, err = parseNextField(remaining, isoPeriod, isHMS)
			if err != nil {
				return err
			}

			if haveFraction && number.Coef() != 0 {
				return fmt.Errorf("%s: '%c' & '%c' only the last field can have a fraction", isoPeriod, previous.Byte(), des.Byte())
			}

			switch des {
			case Year:
				years, err = years.testAndSet(number, Year, &p.years, isoPeriod)
			case Month:
				months, err = months.testAndSet(number, Month, &p.months, isoPeriod)
			case Week:
				weeks, err = weeks.testAndSet(number, Week, &p.weeks, isoPeriod)
			case Day:
				days, err = days.testAndSet(number, Day, &p.days, isoPeriod)
			case Hour:
				hours, err = hours.testAndSet(number, Hour, &p.hours, isoPeriod)
			case Minute:
				minutes, err = minutes.testAndSet(number, Minute, &p.minutes, isoPeriod)
			case Second:
				seconds, err = seconds.testAndSet(number, Second, &p.seconds, isoPeriod)
			default:
				panic(fmt.Errorf("unreachable %s: '%c'", isoPeriod, des.Byte()))
			}
			nComponents++

			if err != nil {
				return err
			}

			if number.Scale() > 0 {
				haveFraction = true
				previous = des
			}
		}
	}

	if nComponents == 0 {
		return fmt.Errorf("%s: expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", isoPeriod)
	}

	*period = p.NormaliseSign()
	return nil
}

//-------------------------------------------------------------------------------------------------

type itemState int

const (
	Unready itemState = iota
	Armed
	Set
)

func (i itemState) testAndSet(number decimal.Decimal, des designator, result *decimal.Decimal, original string) (itemState, error) {
	switch i {
	case Unready:
		return i, fmt.Errorf("%s: '%c' designator cannot occur here", original, des.Byte())
	case Set:
		return i, fmt.Errorf("%s: '%c' designator cannot occur more than once", original, des.Byte())
	}

	*result = number
	return Set, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(str, original string, isHMS bool) (decimal.Decimal, designator, string, error) {
	number, i := scanDigits(str)
	switch i {
	case noNumberFound:
		return decimal.Zero, 0, "", fmt.Errorf("%s: expected a number but found '%c'", original, str[0])
	case stringIsAllNumeric:
		return decimal.Zero, 0, "", fmt.Errorf("%s: missing designator at the end", original)
	}

	dec, err := decimal.Parse(number)
	if err != nil {
		panic(fmt.Errorf("unreachable: %s: %w", original, err))
	}

	des, err := asDesignator(str[i], isHMS)
	if err != nil {
		return decimal.Zero, 0, "", fmt.Errorf("%s: %w", original, err)
	}

	return dec, des, str[i+1:], err
}

// scanDigits finds the index of the first non-digit character after some digits.
func scanDigits(s string) (string, int) {
	rs := []rune(s)
	number := make([]rune, 0, len(rs))

	for i, c := range rs {
		if i == 0 && c == '-' {
			number = append(number, c)
		} else if c == '.' || c == ',' {
			number = append(number, '.') // next step needs decimal point not comma
		} else if '0' <= c && c <= '9' {
			number = append(number, c)
		} else if len(number) > 0 {
			return string(number), i // index of the next non-digit character
		} else {
			return "", noNumberFound
		}
	}
	return "", stringIsAllNumeric
}

const (
	noNumberFound      = -1
	stringIsAllNumeric = -2
)
