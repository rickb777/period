// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestParseErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    Period
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
			g.Expect(ep).To(HaveOccurred(), info(i, c.value))
			g.Expect(ep.Error()).To(Equal(c.expvalue+c.expected), info(i, c.value))

			_, en := Parse("-" + c.value)
			g.Expect(en).To(HaveOccurred(), info(i, c.value))
			if c.expvalue != "" {
				g.Expect(en.Error()).To(Equal("-"+c.expvalue+c.expected), info(i, c.value))
			} else {
				g.Expect(en.Error()).To(Equal(c.expected), info(i, c.value))
			}

			g.Expect(func() { MustParse(c.value) }).To(Panic())
		})
	}
}

//-------------------------------------------------------------------------------------------------

var (
	one    = decimal.One
	negOne = decimal.One.Neg()
)

func TestParsePeriod(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    Period
		reversed Period
		period   Period64
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
		{"P1Y", "P1Y", Period64{years: one}},
		{"P1M", "P1M", Period64{months: one}},
		{"P1W", "P1W", Period64{weeks: one}},
		{"P1D", "P1D", Period64{days: one}},
		{"PT1H", "PT1H", Period64{hours: one}},
		{"PT1M", "PT1M", Period64{minutes: one}},
		{"PT1S", "PT1S", Period64{seconds: one}},
		{"+PT1S", "PT1S", Period64{seconds: one}},

		// unusual case: treat this as a double negative
		{"-P-1Y", "P1Y", Period64{years: one}},
		{"-P-1M", "P1M", Period64{months: one}},
		{"-P-1W", "P1W", Period64{weeks: one}},
		{"-P-1D", "P1D", Period64{days: one}},
		{"-PT-1H", "PT1H", Period64{hours: one}},
		{"-PT-1M", "PT1M", Period64{minutes: one}},
		{"-PT-1S", "PT1S", Period64{seconds: one}},

		{"-P1Y", "-P1Y", Period64{years: one, neg: true}},
		{"-P1M", "-P1M", Period64{months: one, neg: true}},
		{"-P1W", "-P1W", Period64{weeks: one, neg: true}},
		{"-P1D", "-P1D", Period64{days: one, neg: true}},
		{"-PT1H", "-PT1H", Period64{hours: one, neg: true}},
		{"-PT1M", "-PT1M", Period64{minutes: one, neg: true}},
		{"-PT1S", "-PT1S", Period64{seconds: one, neg: true}},
		{"-PT1S", "-PT1S", Period64{seconds: one, neg: true}},

		{"P-1Y", "-P1Y", Period64{years: one, neg: true}},
		{"P-1M", "-P1M", Period64{months: one, neg: true}},
		{"P-1W", "-P1W", Period64{weeks: one, neg: true}},
		{"P-1D", "-P1D", Period64{days: one, neg: true}},
		{"PT-1H", "-PT1H", Period64{hours: one, neg: true}},
		{"PT-1M", "-PT1M", Period64{minutes: one, neg: true}},
		{"PT-1S", "-PT1S", Period64{seconds: one, neg: true}},
		{"PT-1S", "-PT1S", Period64{seconds: one, neg: true}},

		{"P1Y1M1W1DT1H1M1.111111111S", "P1Y1M1W1DT1H1M1.111111111S", Period64{years: one, months: one, weeks: one, days: one, hours: one, minutes: one, seconds: decS("1.111111111")}},
		//{"-P1Y-1M-1W-1DT-1H-1M-1.111111111S", "P1Y1M1W1DT1H1M1.111111111S", Period64{years: negOne, months: negOne, weeks: negOne, days: negOne, hours: negOne, minutes: negOne, seconds: decS("-1.111111111")}},
		{"P1Y-1M1W-1DT1H-1M1.111111111S", "P1Y-1M1W-1DT1H-1M1.111111111S", Period64{years: one, months: negOne, weeks: one, days: negOne, hours: one, minutes: negOne, seconds: decS("1.111111111")}},
		{"-P1Y-1M1W-1DT1H-1M1.111111111S", "-P1Y-1M1W-1DT1H-1M1.111111111S", Period64{years: one, months: negOne, weeks: one, days: negOne, hours: one, minutes: negOne, seconds: decS("1.111111111"), neg: true}},

		{"P0.0000000000000000001Y", "P0.0000000000000000001Y", Period64{years: dec(1, 19)}},
		{"P0.0000000000000000001M", "P0.0000000000000000001M", Period64{months: dec(1, 19)}},
		{"P0.0000000000000000001W", "P0.0000000000000000001W", Period64{weeks: dec(1, 19)}},
		{"P0.0000000000000000001D", "P0.0000000000000000001D", Period64{days: dec(1, 19)}},
		{"PT0.0000000000000000001H", "PT0.0000000000000000001H", Period64{hours: dec(1, 19)}},
		{"PT0.0000000000000000001M", "PT0.0000000000000000001M", Period64{minutes: dec(1, 19)}},
		{"PT0.0000000000000000001S", "PT0.0000000000000000001S", Period64{seconds: dec(1, 19)}},

		{"P9223372036854775807Y", "P9223372036854775807Y", Period64{years: decI(math.MaxInt64)}},
		{"P9223372036854775807M", "P9223372036854775807M", Period64{months: decI(math.MaxInt64)}},
		{"P9223372036854775807W", "P9223372036854775807W", Period64{weeks: decI(math.MaxInt64)}},
		{"P9223372036854775807D", "P9223372036854775807D", Period64{days: decI(math.MaxInt64)}},
		{"PT9223372036854775807H", "PT9223372036854775807H", Period64{hours: decI(math.MaxInt64)}},
		{"PT9223372036854775807M", "PT9223372036854775807M", Period64{minutes: decI(math.MaxInt64)}},
		{"PT9223372036854775807S", "PT9223372036854775807S", Period64{seconds: decI(math.MaxInt64)}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			p := MustParse(c.value)
			s := info(i, c.value)
			g.Expect(p).To(Equal(c.period), s)
			// reversal is usually expected to be an identity
			g.Expect(p.Period()).To(Equal(c.reversed), s+" reversed")
		})
	}
}
