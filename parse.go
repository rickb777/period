package period

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (p32 *Period32) Parse(isoPeriod string) error {
	if isoPeriod == "" {
		return fmt.Errorf(`cannot parse a blank string as a period`)
	}

	*p32 = Period32{}

	if isoPeriod == "P0" {
		return nil // special case
	}

	remaining := isoPeriod
	if remaining[0] == '-' {
		p32.neg = true
		remaining = remaining[1:]
	} else if remaining[0] == '+' {
		remaining = remaining[1:]
	}

	if remaining == "" {
		return fmt.Errorf(`cannot parse a blank string as a period`)
	} else if remaining[0] != 'P' {
		return fmt.Errorf("%s: expected 'P' period mark at the start", isoPeriod)
	}
	remaining = remaining[1:]

	var integer, fraction, prevFraction int32
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
			integer, fraction, des, remaining, err = parseNextField(remaining, isoPeriod, isHMS)
			if err != nil {
				return err
			}

			if prevFraction != 0 && (integer != 0 || fraction != 0) {
				return fmt.Errorf("%s: '%c' & '%c' only the last field can have a fraction", isoPeriod, previous.Byte(), des.Byte())
			}

			switch des {
			case Year:
				years, err = years.testAndSet(integer, fraction, Year, p32, &p32.years, isoPeriod)
			case Month:
				months, err = months.testAndSet(integer, fraction, Month, p32, &p32.months, isoPeriod)
			case Week:
				weeks, err = weeks.testAndSet(integer, fraction, Week, p32, &p32.weeks, isoPeriod)
			case Day:
				days, err = days.testAndSet(integer, fraction, Day, p32, &p32.days, isoPeriod)
			case Hour:
				hours, err = hours.testAndSet(integer, fraction, Hour, p32, &p32.hours, isoPeriod)
			case Minute:
				minutes, err = minutes.testAndSet(integer, fraction, Minute, p32, &p32.minutes, isoPeriod)
			case Second:
				seconds, err = seconds.testAndSet(integer, fraction, Second, p32, &p32.seconds, isoPeriod)
			default:
				return fmt.Errorf("%s: expected a number not '%c'", isoPeriod, des.Byte())
			}
			nComponents++

			if err != nil {
				return err
			}

			prevFraction = fraction
			previous = des
		}
	}

	if nComponents == 0 {
		return fmt.Errorf("%s: expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", isoPeriod)
	}

	return nil
}

//-------------------------------------------------------------------------------------------------

type itemState int

const (
	Unready itemState = iota
	Armed
	Set
)

func (i itemState) testAndSet(integer, fraction int32, des designator, result *Period32, value *int32, original string) (itemState, error) {
	switch i {
	case Unready:
		return i, fmt.Errorf("%s: '%c' designator cannot occur here", original, des.Byte())
	case Set:
		return i, fmt.Errorf("%s: '%c' designator cannot occur more than once", original, des.Byte())
	}

	*value = integer
	if integer != 0 || fraction != 0 {
		result.fraction = fraction
		result.lastField = des
	}
	return Set, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(str, original string, isHMS bool) (int32, int32, designator, string, error) {
	i := scanDigits(str)
	switch i {
	case noDigitsFound:
		return 0, 0, 0, "", fmt.Errorf("%s: expected a number but found '%c'", original, str[0])
	case stringIsAllDigits:
		return 0, 0, 0, "", fmt.Errorf("%s: missing designator at the end", original)
	}

	des, err := asDesignator(str[i], isHMS)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("%s: %w", original, err)
	}

	integer, fraction, err := parseDecimalNumber(str[:i], original, des)
	if integer > math.MaxInt32 {
		return 0, 0, 0, "", fmt.Errorf("%s: integer overflow occurred in %s", original, des.field())
	}
	return int32(integer), int32(fraction), des, str[i+1:], err
}

const (
	maxFractionDigits = 9
	trailingZeros     = "000000000" // nine zeros
)

// Fixed-point one decimal place
func parseDecimalNumber(number, original string, des designator) (integer, fraction int64, err error) {
	dec := strings.IndexByte(number, '.')
	if dec < 0 {
		dec = strings.IndexByte(number, ',')
	}

	if dec >= 0 {
		integer, err = strconv.ParseInt(number[:dec], 10, 64)
		if err == nil {
			number = number[dec+1:]
			if len(number) > 0 {
				number = (number + trailingZeros)[:maxFractionDigits]
				fraction, err = strconv.ParseInt(number, 10, 64)
				//fraction *= pow10(maxFractionDigits - 1 - countZeros(number))
			}
		}
	} else {
		integer, err = strconv.ParseInt(number, 10, 64)
	}

	if err != nil {
		return 0, 0, fmt.Errorf("%s: expected a number but found '%c'", original, des)
	}

	return integer, fraction, err
}

// scanDigits finds the index of the first non-digit character after some digits.
func scanDigits(s string) int {
	foundSomeDigits := false
	for i, c := range s {
		if !isDigit(c) {
			if foundSomeDigits {
				return i // index of the next non-digit character
			} else {
				return noDigitsFound
			}
		} else {
			foundSomeDigits = true
		}
	}
	return stringIsAllDigits
}

const (
	noDigitsFound     = -1
	stringIsAllDigits = -2
)

func isDigit(c rune) bool {
	return ('0' <= c && c <= '9') || c == '.' || c == ','
}
