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

// Period32 holds a period of time as a set of integers, one for each field in the ISO-8601
// period, and additional information to track any fraction.
// The precision is almost unlimited (int32 is used for all fields for calculations). Fractions
// can hold up to 9 decimal places, therefore the finest grain is one nanosecond.
//
// Instances are immutable.
type Period32 struct {
	// always positive values
	years, months, weeks, days, hours, minutes, seconds, fraction int32

	// true if the period is negative
	neg bool

	// the fraction applies to this field; no other fields to the right can be non-zero
	lastField designator

	// ISO-8601 representation
	//s string
}

//-------------------------------------------------------------------------------------------------

// Period converts the period to ISO-8601 string form.
func (p32 Period32) Period() Period {
	return Period(p32.String())
}

// String converts the period to ISO-8601 form.
func (p32 Period32) String() string {
	neg := p32.neg
	p32.neg = false

	if p32 == (Period32{}) {
		return "P0D"
	}

	buf := &strings.Builder{}
	if neg {
		buf.WriteByte('-')
	}

	buf.WriteByte('P')

	if writeField32(buf, p32.years, p32.fraction, Year, p32.lastField) {
		return buf.String()
	}
	if writeField32(buf, p32.months, p32.fraction, Month, p32.lastField) {
		return buf.String()
	}
	if writeField32(buf, p32.weeks, p32.fraction, Week, p32.lastField) {
		return buf.String()
	}
	if writeField32(buf, p32.days, p32.fraction, Day, p32.lastField) {
		return buf.String()
	}

	buf.WriteByte('T')

	if writeField32(buf, p32.hours, p32.fraction, Hour, p32.lastField) {
		return buf.String()
	}
	if writeField32(buf, p32.minutes, p32.fraction, Minute, p32.lastField) {
		return buf.String()
	}
	writeField32(buf, p32.seconds, p32.fraction, Second, p32.lastField)

	return buf.String()
}

func writeField32(w usefulWriter, field, fraction int32, fieldDesignator, lastField designator) bool {
	if field != 0 || (fraction != 0 && fieldDesignator == lastField) {
		fmt.Fprintf(w, "%d", field)
		if fieldDesignator == lastField {
			writeFraction(w, fraction)
		}
		w.WriteByte(fieldDesignator.Byte())
	}
	return fieldDesignator == lastField
}

func writeFraction(w usefulWriter, fraction int32) {
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

		s := strconv.FormatInt(int64(fraction), 10)
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
func (p32 Period32) Abs() Period32 {
	p32.neg = false
	return p32
}

// Negate returns a period with the sign changed.
func (p32 Period32) Negate() Period32 {
	p32.neg = !p32.neg
	return p32
}

// IsZero returns true if applied to a period of zero-length.
func (p32 Period32) IsZero() bool {
	p32.neg = false
	return p32 == Period32{}
}

// Sign returns 1 if the period is positive, -1 if it is negative, or zero otherwise.
func (p32 Period32) Sign() int {
	switch {
	case p32.neg:
		return -1
	case p32 != Period32{}:
		return 1
	default:
		return 0
	}
}

// IsNegative returns true if the period is negative.
func (p32 Period32) IsNegative() bool {
	return p32.neg
}

// IsPositive returns true if the period is positive or zero.
func (p32 Period32) IsPositive() bool {
	return !p32.neg
}

// isValid returns true if all fields with the period are consistent with each other.
func (p32 Period32) isValid() bool {
	if p32.years < 0 || p32.months < 0 || p32.weeks < 0 || p32.days < 0 || p32.hours < 0 || p32.minutes < 0 || p32.seconds < 0 {
		return false
	}

	switch p32.lastField {
	case Second:
		return p32.fraction == 0
	case Minute:
		return p32.seconds == 0
	case Hour:
		return p32.seconds == 0 && p32.minutes == 0
	case Day:
		return p32.seconds == 0 && p32.minutes == 0 && p32.hours == 0
	case Week:
		return p32.seconds == 0 && p32.minutes == 0 && p32.hours == 0 && p32.days == 0
	case Month:
		return p32.seconds == 0 && p32.minutes == 0 && p32.hours == 0 && p32.days == 0 && p32.weeks == 0
	case Year:
		return p32.seconds == 0 && p32.minutes == 0 && p32.hours == 0 && p32.days == 0 && p32.weeks == 0 && p32.months == 0
	}

	return true
}
