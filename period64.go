// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Period64 holds a period of time as a set of integers, one for each field in the ISO-8601
// period, and additional information to track any fraction.
// The precision is almost unlimited (int64 is used for all fields for calculations). Fractions
// can hold up to 9 decimal places, therefore the finest grain is one nanosecond.
//
// Instances are immutable.
type Period64 struct {
	// always positive values
	years, months, weeks, days, hours, minutes, seconds, fraction int64

	// true if the period is negative
	neg bool

	// the fraction applies to this field; no other fields to the right can be non-zero
	lastField designator

	// ISO-8601 representation
	s string
}

//-------------------------------------------------------------------------------------------------

// Period converts the period to ISO-8601 string form.
func (p64 Period64) Period() Period {
	if p64.s != "" {
		return Period(p64.s)
	}
	return Period(p64.String())
}

// String converts the period to ISO-8601 form.
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
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] != '0' {
				s = s[:i+1]
				break
			}
		}
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

// Negate returns a period with the sign changed.
func (p64 Period64) Negate() Period64 {
	p64.neg = !p64.neg
	return p64
}

// IsZero returns true if applied to a period of zero-length.
func (p64 Period64) IsZero() bool {
	p64.neg = false
	return p64 == Period64{}
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
func (p64 Period64) Sign() int {
	switch {
	case p64.neg:
		return -1
	case p64 != Period64{}:
		return 1
	default:
		return 0
	}
}

// IsNegative returns true if the period is negative.
func (p64 Period64) IsNegative() bool {
	return p64.neg
}

// IsPositive returns true if the period is positive or zero.
func (p64 Period64) IsPositive() bool {
	return !p64.neg
}

// isValid returns true if all fields with the period are consistent with each other.
func (p64 Period64) isValid() bool {
	if p64.years < 0 || p64.months < 0 || p64.weeks < 0 || p64.days < 0 || p64.hours < 0 || p64.minutes < 0 || p64.seconds < 0 {
		return false
	}

	switch p64.lastField {
	case Second:
		return p64.fraction == 0
	case Minute:
		return p64.seconds == 0
	case Hour:
		return p64.seconds == 0 && p64.minutes == 0
	case Day:
		return p64.seconds == 0 && p64.minutes == 0 && p64.hours == 0
	case Week:
		return p64.seconds == 0 && p64.minutes == 0 && p64.hours == 0 && p64.days == 0
	case Month:
		return p64.seconds == 0 && p64.minutes == 0 && p64.hours == 0 && p64.days == 0 && p64.weeks == 0
	case Year:
		return p64.seconds == 0 && p64.minutes == 0 && p64.hours == 0 && p64.days == 0 && p64.weeks == 0 && p64.months == 0
	}

	return true
}
