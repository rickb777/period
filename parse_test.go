// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	bigdecimal "github.com/shopspring/decimal"
	"testing"
)

func TestParseErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value     Period
		normalise bool
		expected  string
		expvalue  string
	}{
		{"", false, "cannot parse a blank string as a period", ""},
		{`P000`, false, `: missing designator at the end`, "P000"},
		{`PT1`, false, `: missing designator at the end`, "PT1"},
		{"XY", false, ": expected 'P' period mark at the start", "XY"},
		{"PxY", false, ": expected a number but found 'x'", "PxY"},
		{"PxW", false, ": expected a number but found 'x'", "PxW"},
		{"PxD", false, ": expected a number but found 'x'", "PxD"},
		{"PTxH", false, ": expected a number but found 'x'", "PTxH"},
		{"PTxM", false, ": expected a number but found 'x'", "PTxM"},
		{"PTxS", false, ": expected a number but found 'x'", "PTxS"},
		{"PT1A", false, ": expected a designator Y, M, W, D, H, or S not 'A'", "PT1A"},
		{"P1HT1M", false, ": 'H' designator cannot occur here", "P1HT1M"},
		{"PT1Y", false, ": 'Y' designator cannot occur here", "PT1Y"},
		{"P1S", false, ": 'S' designator cannot occur here", "P1S"},
		{"P1D2D", false, ": 'D' designator cannot occur more than once", "P1D2D"},
		{"PT1HT1S", false, ": 'T' designator cannot occur more than once", "PT1HT1S"},
		{"P0.1YT0.1S", false, ": 'Y' & 'S' only the last field can have a fraction", "P0.1YT0.1S"},
		{"P", false, ": expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", "P"},
	}
	for i, c := range cases {
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
	}
}

//-------------------------------------------------------------------------------------------------

var (
	bigOne = bigdecimal.NewFromInt(1)
	one    = decimal{value: 1}
	negOne = decimal{value: -1}
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
