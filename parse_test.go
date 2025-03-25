// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
	"github.com/rickb777/expect"
	"math"
	"testing"
)

func TestParseErrors(t *testing.T) {
	cases := []struct {
		value    ISOString
		expected string
		expvalue string
	}{
		{"", "cannot parse a blank string as a period", ""},
		{`P000`, `: missing designator at the end`, "P000"},
		{`PT1`, `: missing designator at the end`, "PT1"},
		{"XY", ": expected 'P' period mark at the start", "XY"},
		{"PxY", ": expected a number but found 'x'", "PxY"},
		{"PxW", ": expected a number but found 'x'", "PxW"},
		{"PxD", ": expected a number but found 'x'", "PxD"},
		{"PTxH", ": expected a number but found 'x'", "PTxH"},
		{"PTxM", ": expected a number but found 'x'", "PTxM"},
		{"PTxS", ": expected a number but found 'x'", "PTxS"},
		{"PT1A", ": expected a designator Y, M, W, D, H, or S not 'A'", "PT1A"},
		{"P1HT1M", ": 'H' designator cannot occur here", "P1HT1M"},
		{"PT1Y", ": 'Y' designator cannot occur here", "PT1Y"},
		{"P1S", ": 'S' designator cannot occur here", "P1S"},
		{"P1D2D", ": 'D' designator cannot occur more than once", "P1D2D"},
		{"PT1HT1S", ": 'T' designator cannot occur more than once", "PT1HT1S"},
		{"P0.1YT0.1S", ": 'Y' & 'S' only the last field can have a fraction", "P0.1YT0.1S"},
		{"P", ": expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", "P"},
		{"P92233720368547758071Y", ": number invalid or out of range", "P92233720368547758071Y"},
		{"P1.1.1Y", ": number invalid or out of range", "P1.1.1Y"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			_, ep := Parse(c.value)
			expect.Error(ep).Info("%d %v", i, c.value).ToHaveOccurred(t)
			expect.Error(ep).Info("%d %v", i, c.value).ToContain(t, c.expvalue+c.expected)

			_, en := Parse("-" + c.value)
			expect.Error(en).Info("%d %v", i, c.value).ToHaveOccurred(t)
			if c.expvalue != "" {
				expect.Error(en).Info("%d %v", i, c.value).ToContain(t, "-"+c.expvalue+c.expected)
			} else {
				expect.Error(en).Info("%d %v", i, c.value).ToContain(t, c.expected)
			}

			expect.Func(func() { MustParse(c.value) }).ToPanic(t)
		})
	}
}

//-------------------------------------------------------------------------------------------------

var (
	one    = decimal.One
	negOne = decimal.One.Neg()
)

func TestParsePeriod(t *testing.T) {
	cases := []struct {
		value    ISOString
		reversed ISOString
		period   Period
	}{
		// zero
		{"P0D", CanonicalZero, Zero},
		{"P0Y", CanonicalZero, Zero},
		{"P0M", CanonicalZero, Zero},
		{"P0W", CanonicalZero, Zero},
		{"PT0H", CanonicalZero, Zero},
		{"PT0M", CanonicalZero, Zero},
		{"PT0S", CanonicalZero, Zero},
		{"-P0D", CanonicalZero, Zero},
		{"-P0Y", CanonicalZero, Zero},
		{"-P0M", CanonicalZero, Zero},
		{"-P0W", CanonicalZero, Zero},
		{"-PT0H", CanonicalZero, Zero},
		{"-PT0M", CanonicalZero, Zero},
		{"-PT0S", CanonicalZero, Zero},
		{"+PT0S", CanonicalZero, Zero},

		// ones
		{"P1Y", "P1Y", Period{years: one}},
		{"P1M", "P1M", Period{months: one}},
		{"P1W", "P1W", Period{weeks: one}},
		{"P1D", "P1D", Period{days: one}},
		{"PT1H", "PT1H", Period{hours: one}},
		{"PT1M", "PT1M", Period{minutes: one}},
		{"PT1S", "PT1S", Period{seconds: one}},
		{"+PT1S", "PT1S", Period{seconds: one}},

		// unusual case: treat this as a double negative
		{"-P-1Y", "P1Y", Period{years: one}},
		{"-P-1M", "P1M", Period{months: one}},
		{"-P-1W", "P1W", Period{weeks: one}},
		{"-P-1D", "P1D", Period{days: one}},
		{"-PT-1H", "PT1H", Period{hours: one}},
		{"-PT-1M", "PT1M", Period{minutes: one}},
		{"-PT-1S", "PT1S", Period{seconds: one}},

		{"-P1Y", "-P1Y", Period{years: one, neg: true}},
		{"-P1M", "-P1M", Period{months: one, neg: true}},
		{"-P1W", "-P1W", Period{weeks: one, neg: true}},
		{"-P1D", "-P1D", Period{days: one, neg: true}},
		{"-PT1H", "-PT1H", Period{hours: one, neg: true}},
		{"-PT1M", "-PT1M", Period{minutes: one, neg: true}},
		{"-PT1S", "-PT1S", Period{seconds: one, neg: true}},
		{"-PT1S", "-PT1S", Period{seconds: one, neg: true}},

		{"P-1Y", "-P1Y", Period{years: one, neg: true}},
		{"P-1M", "-P1M", Period{months: one, neg: true}},
		{"P-1W", "-P1W", Period{weeks: one, neg: true}},
		{"P-1D", "-P1D", Period{days: one, neg: true}},
		{"PT-1H", "-PT1H", Period{hours: one, neg: true}},
		{"PT-1M", "-PT1M", Period{minutes: one, neg: true}},
		{"PT-1S", "-PT1S", Period{seconds: one, neg: true}},
		{"PT-1S", "-PT1S", Period{seconds: one, neg: true}},

		{"P1Y1M1W1DT1H1M1.111111111S", "P1Y1M1W1DT1H1M1.111111111S", Period{years: one, months: one, weeks: one, days: one, hours: one, minutes: one, seconds: decS("1.111111111")}},
		//{"-P1Y-1M-1W-1DT-1H-1M-1.111111111S", "P1Y1M1W1DT1H1M1.111111111S", Period{years: negOne, months: negOne, weeks: negOne, days: negOne, hours: negOne, minutes: negOne, seconds: decS("-1.111111111")}},
		{"P1Y-1M1W-1DT1H-1M1.111111111S", "P1Y-1M1W-1DT1H-1M1.111111111S", Period{years: one, months: negOne, weeks: one, days: negOne, hours: one, minutes: negOne, seconds: decS("1.111111111")}},
		{"-P1Y-1M1W-1DT1H-1M1.111111111S", "-P1Y-1M1W-1DT1H-1M1.111111111S", Period{years: one, months: negOne, weeks: one, days: negOne, hours: one, minutes: negOne, seconds: decS("1.111111111"), neg: true}},

		{"P0.0000000000000000001Y", "P0.0000000000000000001Y", Period{years: dec(1, 19)}},
		{"P0.0000000000000000001M", "P0.0000000000000000001M", Period{months: dec(1, 19)}},
		{"P0.0000000000000000001W", "P0.0000000000000000001W", Period{weeks: dec(1, 19)}},
		{"P0.0000000000000000001D", "P0.0000000000000000001D", Period{days: dec(1, 19)}},
		{"PT0.0000000000000000001H", "PT0.0000000000000000001H", Period{hours: dec(1, 19)}},
		{"PT0.0000000000000000001M", "PT0.0000000000000000001M", Period{minutes: dec(1, 19)}},
		{"PT0.0000000000000000001S", "PT0.0000000000000000001S", Period{seconds: dec(1, 19)}},

		{"P9223372036854775807Y", "P9223372036854775807Y", Period{years: decI(math.MaxInt64)}},
		{"P9223372036854775807M", "P9223372036854775807M", Period{months: decI(math.MaxInt64)}},
		{"P9223372036854775807W", "P9223372036854775807W", Period{weeks: decI(math.MaxInt64)}},
		{"P9223372036854775807D", "P9223372036854775807D", Period{days: decI(math.MaxInt64)}},
		{"PT9223372036854775807H", "PT9223372036854775807H", Period{hours: decI(math.MaxInt64)}},
		{"PT9223372036854775807M", "PT9223372036854775807M", Period{minutes: decI(math.MaxInt64)}},
		{"PT9223372036854775807S", "PT9223372036854775807S", Period{seconds: decI(math.MaxInt64)}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			p := MustParse(c.value)
			s := info(i, c.value)
			expect.Any(p).Info(s).ToBe(t, c.period)
			// reversal is usually expected to be an identity
			expect.Any(p.Period()).Info(s+" reversed").ToBe(t, c.reversed)
		})
	}
}
