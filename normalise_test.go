// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	"math"
	"strings"
	"testing"
	"time"
)

func Test_Normalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		input     Period
		precise   Period
		imprecise Period
	}{
		// note: the negative cases are also covered (see below)

		{input: "P0D", precise: "P0D", imprecise: "P0D"},

		// ones unchanged
		{input: "P1Y", precise: "P1Y", imprecise: "P1Y"},
		{input: "P1M", precise: "P1M", imprecise: "P1M"},
		{input: "P1W", precise: "P1W", imprecise: "P1W"},
		{input: "P1D", precise: "P1D", imprecise: "P1D"},
		{input: "PT1H", precise: "PT1H", imprecise: "PT1H"},
		{input: "PT1M", precise: "PT1M", imprecise: "PT1M"},
		{input: "PT1S", precise: "PT1S", imprecise: "PT1S"},

		{input: "P11Y", precise: "P11Y", imprecise: "P11Y"},
		{input: "P24M", precise: "P2Y", imprecise: "P2Y"},
		{input: "P10W", precise: "P10W", imprecise: "P10W"},
		{input: "P14D", precise: "P2W", imprecise: "P2W"},
		{input: "PT48H", precise: "PT48H", imprecise: "P2D"},
		{input: "PT120M", precise: "PT2H", imprecise: "PT2H"},
		{input: "PT120S", precise: "PT2M", imprecise: "PT2M"},

		// big fraction changed
		{input: "P1.1Y", precise: "P1.1Y", imprecise: "P1.1Y"},
		{input: "P0.1Y", precise: "P0.1Y", imprecise: "P0.1Y"},
		{input: "P0.1M", precise: "P0.1M", imprecise: "P0.1M"},
		{input: "P0.1W", precise: "P0.1W", imprecise: "P0.1W"},
		{input: "P0.1D", precise: "P0.1D", imprecise: "P0.1D"},
		{input: "PT0.1H", precise: "PT0.1H", imprecise: "PT0.1H"},
		{input: "PT0.1M", precise: "PT0.1M", imprecise: "PT0.1M"},
		{input: "PT0.1S", precise: "PT0.1S", imprecise: "PT0.1S"},

		// small fraction unchanged
		{input: "P0.000000001Y", precise: "P0.000000001Y", imprecise: "P0.000000001Y"},
		{input: "P0.000000001M", precise: "P0.000000001M", imprecise: "P0.000000001M"},
		{input: "P0.000000001W", precise: "P0.000000001W", imprecise: "P0.000000001W"},
		{input: "P0.000000001D", precise: "P0.000000001D", imprecise: "P0.000000001D"},
		{input: "PT0.000000001H", precise: "PT0.000000001H", imprecise: "PT0.000000001H"},
		{input: "PT0.000000001M", precise: "PT0.000000001M", imprecise: "PT0.000000001M"},
		{input: "PT0.000000001S", precise: "PT0.000000001S", imprecise: "PT0.000000001S"},

		// small overflow disregarded
		{input: "PT60.0005S", precise: "PT60.0005S", imprecise: "PT60.0005S"},

		// no change
		{input: "PT26H", precise: "PT26H", imprecise: "P1DT2H"},
		{input: "PT26.1H", precise: "PT26.1H", imprecise: "P1DT2.1H"},
		{input: "P5.3W", precise: "P5.3W", imprecise: "P5.3W"},

		{input: "PT65.5S", precise: "PT1M5.5S", imprecise: "PT1M5.5S"},
		{input: "PT3601.1S", precise: "PT1H1.1S", imprecise: "PT1H1.1S"},
		{input: "PT3661.1S", precise: "PT1H1M1.1S", imprecise: "PT1H1M1.1S"},
		{input: "P9D", precise: "P1W2D", imprecise: "P1W2D"},
		{input: "P14M", precise: "P1Y2M", imprecise: "P1Y2M"},
		{input: "PT26.1H", precise: "PT26.1H", imprecise: "P1DT2.1H"},
		{input: "P366.1D", precise: "P52W2.1D", imprecise: "P52W2.1D"},
		{input: "PT1440M", precise: "PT24H", imprecise: "P1D"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.precise), func(t *testing.T) {
			p1 := MustParse(c.input)
			sp1 := p1.Normalise(true)
			g.Expect(sp1.Period()).To(Equal(c.precise), "precise +ve case")

			if !p1.IsZero() {
				sp1n := p1.Negate().Normalise(true)
				g.Expect(sp1n.Period()).To(Equal("-"+c.precise), "precise -ve case")
			}

			p2 := MustParse(c.input)
			sp2 := p2.Normalise(false)
			g.Expect(sp2.Period()).To(Equal(c.imprecise), "imprecise +ve case")

			if !p2.IsZero() {
				sp2n := p2.Negate().Normalise(false)
				g.Expect(sp2n.Period()).To(Equal("-"+c.imprecise), "imprecise -ve case")
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_NormaliseDaysToYears(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		input    Period
		expected Period
	}{
		// note: the negative cases are also covered (see below)

		{input: "P0D", expected: "P0D"},

		// ones unchanged
		{input: "P1Y", expected: "P1Y"},
		{input: "P1M", expected: "P1M"},
		{input: "P1W", expected: "P1W"},
		{input: "P1D", expected: "P1D"},
		{input: "PT1H", expected: "PT1H"},
		{input: "PT1M", expected: "PT1M"},
		{input: "PT1S", expected: "PT1S"},

		{input: "P365D", expected: "P365D"},
		{input: "P366D", expected: "P1Y0.7575D"},
		{input: "P367D", expected: "P1Y1.7575D"},
		{input: "P1461D", expected: "P4Y0.03D"},
		{input: "P1469D", expected: "P4Y1W1.03D"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			p1 := MustParse(c.input)
			sp1 := p1.NormaliseDaysToYears()
			g.Expect(sp1.Period()).To(Equal(c.expected), "+ve case")

			if !p1.IsZero() {
				sp1n := p1.Negate().NormaliseDaysToYears()
				g.Expect(sp1n.Period()).To(Equal("-"+c.expected), "-ve case")
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_Simplify(t *testing.T) {
	g := NewGomegaWithT(t)

	var extremeMinSec = Period(fmt.Sprintf("PT%dM0.%s1S", math.MaxInt64, strings.Repeat("0", 18)))

	cases := []struct {
		input     Period
		precise   Period
		imprecise Period
	}{
		// note: the negative cases are also covered (see below)

		{input: "P0D", precise: "P0D", imprecise: "P0D"},

		// ones unchanged
		{input: "P1Y", precise: "P1Y", imprecise: "P1Y"},
		{input: "P1M", precise: "P1M", imprecise: "P1M"},
		{input: "P1W", precise: "P1W", imprecise: "P1W"},
		{input: "P1D", precise: "P1D", imprecise: "P1D"},
		{input: "PT1H", precise: "PT1H", imprecise: "PT1H"},
		{input: "PT1M", precise: "PT1M", imprecise: "PT1M"},
		{input: "PT1S", precise: "PT1S", imprecise: "PT1S"},

		// these are already simple
		{input: "P3Y", precise: "P3Y", imprecise: "P3Y"},
		{input: "P3M", precise: "P3M", imprecise: "P3M"},
		{input: "P3W", precise: "P3W", imprecise: "P3W"},
		{input: "P3D", precise: "P3D", imprecise: "P3D"},
		{input: "PT3H", precise: "PT3H", imprecise: "PT3H"},
		{input: "PT3M", precise: "PT3M", imprecise: "PT3M"},
		{input: "PT3S", precise: "PT3S", imprecise: "PT3S"},

		// simplified where possible
		{input: "P2Y1M", precise: "P25M", imprecise: "P25M"},
		{input: "P2W1D", precise: "P15D", imprecise: "P15D"},
		{input: "P2DT1H", precise: "P2DT1H", imprecise: "PT49H"},
		{input: "P1DT48H", precise: "P1DT48H", imprecise: "P3D"},
		{input: "P1DT23H", precise: "P1DT23H", imprecise: "PT47H"},
		{input: "PT3H120M", precise: "PT5H", imprecise: "PT5H"},
		{input: "PT3M120S", precise: "PT5M", imprecise: "PT5M"},

		// no change because of rounding issues
		{input: "P0.083333333Y", precise: "P0.083333333Y", imprecise: "P0.083333333Y"},
		{input: "P0.1Y", precise: "P0.1Y", imprecise: "P0.1Y"},
		{input: "P0.11111Y", precise: "P0.11111Y", imprecise: "P0.11111Y"},
		{input: "P0.1M", precise: "P0.1M", imprecise: "P0.1M"},

		{input: "P0.1W", precise: "P0.7D", imprecise: "P0.7D"},
		{input: "P0.1D", precise: "P0.1D", imprecise: "P0.1D"},

		{input: "PT0.1H", precise: "PT6M", imprecise: "PT6M"},
		{input: "PT0.1M", precise: "PT6S", imprecise: "PT6S"},
		{input: "PT0.1S", precise: "PT0.1S", imprecise: "PT0.1S"},
		{input: "PT0.05H", precise: "PT3M", imprecise: "PT3M"},
		{input: "PT0.05M", precise: "PT3S", imprecise: "PT3S"},
		{input: "PT0.05S", precise: "PT0.05S", imprecise: "PT0.05S"},

		// small overflow disregarded
		{input: "PT60.0005S", precise: "PT60.0005S", imprecise: "PT60.0005S"},

		// because of addition overflow, this input is unaltered
		{input: extremeMinSec, precise: extremeMinSec, imprecise: extremeMinSec},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.precise), func(t *testing.T) {
			p1 := MustParse(c.input)
			sp1 := p1.Simplify(true)
			g.Expect(sp1.Period()).To(Equal(c.precise), "precise +ve case")

			if !p1.IsZero() {
				sp1n := p1.Negate().Simplify(true)
				g.Expect(sp1n.Period()).To(Equal("-"+c.precise), "precise -ve case")
			}

			p2 := MustParse(c.input)
			sp2 := p2.Simplify(false)
			g.Expect(sp2.Period()).To(Equal(c.imprecise), "imprecise +ve case")

			if !p2.IsZero() {
				sp2n := p2.Negate().Simplify(false)
				g.Expect(sp2n.Period()).To(Equal("-"+c.imprecise), "imprecise -ve case")
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_NormaliseSign(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		expected Period
		input    Period64
	}{
		{expected: "P0D"},

		{expected: "P1Y", input: Period64{years: one}},
		{expected: "P1M", input: Period64{months: one}},
		{expected: "P1W", input: Period64{weeks: one}},
		{expected: "P1D", input: Period64{days: one}},
		{expected: "PT1H", input: Period64{hours: one}},
		{expected: "PT1M", input: Period64{minutes: one}},
		{expected: "PT1S", input: Period64{seconds: one}},

		{expected: "P1Y", input: Period64{years: negOne, neg: true}},
		{expected: "P1M", input: Period64{months: negOne, neg: true}},
		{expected: "P1W", input: Period64{weeks: negOne, neg: true}},
		{expected: "P1D", input: Period64{days: negOne, neg: true}},
		{expected: "PT1H", input: Period64{hours: negOne, neg: true}},
		{expected: "PT1M", input: Period64{minutes: negOne, neg: true}},
		{expected: "PT1S", input: Period64{seconds: negOne, neg: true}},

		{expected: "-P1Y", input: Period64{years: negOne, neg: false}},
		{expected: "-P1M", input: Period64{months: negOne, neg: false}},
		{expected: "-P1W", input: Period64{weeks: negOne, neg: false}},
		{expected: "-P1D", input: Period64{days: negOne, neg: false}},
		{expected: "-PT1H", input: Period64{hours: negOne, neg: false}},
		{expected: "-PT1M", input: Period64{minutes: negOne, neg: false}},
		{expected: "-PT1S", input: Period64{seconds: negOne, neg: false}},

		{expected: "-P1Y", input: Period64{years: one, neg: true}},
		{expected: "-P1M", input: Period64{months: one, neg: true}},
		{expected: "-P1W", input: Period64{weeks: one, neg: true}},
		{expected: "-P1D", input: Period64{days: one, neg: true}},
		{expected: "-PT1H", input: Period64{hours: one, neg: true}},
		{expected: "-PT1M", input: Period64{minutes: one, neg: true}},
		{expected: "-PT1S", input: Period64{seconds: one, neg: true}},

		// complex normalisation
		{expected: "-PT1M-1S", input: Period64{minutes: one, seconds: negOne, neg: true}},
		{expected: "-PT1M1S", input: Period64{minutes: one, seconds: one, neg: true}},
		{expected: "PT1M-1S", input: Period64{minutes: one, seconds: negOne}},

		{expected: "PT1H1M1S", input: Period64{hours: negOne, minutes: negOne, seconds: negOne, neg: true}}, // 111
		{expected: "PT1H1M-1S", input: Period64{hours: negOne, minutes: negOne, seconds: one, neg: true}},   // 110
		{expected: "PT1H-1M1S", input: Period64{hours: negOne, minutes: one, seconds: negOne, neg: true}},   // 101
		{expected: "PT1H-1M-1S", input: Period64{hours: negOne, minutes: one, seconds: one, neg: true}},     // 100
		{expected: "-PT1H-1M-1S", input: Period64{hours: one, minutes: negOne, seconds: negOne, neg: true}}, // 011
		{expected: "-PT1H-1M1S", input: Period64{hours: one, minutes: negOne, seconds: one, neg: true}},     // 010
		{expected: "-PT1H1M-1S", input: Period64{hours: one, minutes: one, seconds: negOne, neg: true}},     // 001
		{expected: "-PT1H1M1S", input: Period64{hours: one, minutes: one, seconds: one, neg: true}},         // 000

		{expected: "-PT1H1M1S", input: Period64{hours: negOne, minutes: negOne, seconds: negOne}}, // 111
		{expected: "-PT1H1M-1S", input: Period64{hours: negOne, minutes: negOne, seconds: one}},   // 110
		{expected: "-PT1H-1M1S", input: Period64{hours: negOne, minutes: one, seconds: negOne}},   // 101
		{expected: "-PT1H-1M-1S", input: Period64{hours: negOne, minutes: one, seconds: one}},     // 100
		{expected: "PT1H-1M-1S", input: Period64{hours: one, minutes: negOne, seconds: negOne}},   // 011
		{expected: "PT1H-1M1S", input: Period64{hours: one, minutes: negOne, seconds: one}},       // 010
		{expected: "PT1H1M-1S", input: Period64{hours: one, minutes: one, seconds: negOne}},       // 001
		{expected: "PT1H1M1S", input: Period64{hours: one, minutes: one, seconds: one}},           // 000
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			sp1 := c.input.NormaliseSign()
			g.Expect(sp1.Period()).To(Equal(c.expected))
		})
	}
}

//-------------------------------------------------------------------------------------------------

const oneDay = 24 * time.Hour
const oneMonthApprox = daysPerMonthE6 * secondsPerDay * time.Microsecond // 30.436875 days
const oneYearApprox = 31556952 * time.Second                             // 365.2425 days

func TestPeriodToDuration(t *testing.T) {
	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", time.Duration(0), true},

		{"PT1S", 1 * time.Second, true},
		{"PT0.1S", 100 * time.Millisecond, true},
		{"PT0.001S", time.Millisecond, true},
		{"PT0.000001S", time.Microsecond, true},
		{"PT0.000000001S", time.Nanosecond, true},
		{"PT3276S", 3276 * time.Second, true},

		{"PT1M", 60 * time.Second, true},
		{"PT0.1M", 6 * time.Second, true},
		{"PT0.0001M", 6 * time.Millisecond, true},
		{"PT0.0000001M", 6 * time.Microsecond, true},
		{"PT0.0000000001M", 6 * time.Nanosecond, true},
		{"PT3276M", 3276 * time.Minute, true},

		{"PT1H", 3600 * time.Second, true},
		{"PT0.1H", 360 * time.Second, true},
		{"PT0.01H", 36 * time.Second, true},
		{"PT0.00001H", 36 * time.Millisecond, true},
		{"PT0.00000001H", 36 * time.Microsecond, true},
		{"PT0.00000000001H", 36 * time.Nanosecond, true},
		{"PT3220H", 3220 * time.Hour, true},

		{"P1D", 24 * time.Hour, false},

		// days, months and years conversions are never precise
		{"P0.1D", 144 * time.Minute, false},
		{"P10000D", 10000 * 24 * time.Hour, false},
		{"P1W", 168 * time.Hour, false},
		{"P0.1W", 16*time.Hour + 48*time.Minute, false},
		{"P10000W", 10000 * 7 * 24 * time.Hour, false},
		{"P1M", oneMonthApprox, false},
		{"P0.1M", oneMonthApprox / 10, false},
		{"P1000M", 1000 * oneMonthApprox, false},
		{"P1Y", oneYearApprox, false},
		{"P0.1Y", oneYearApprox / 10, false},
		{"P100Y", 100 * oneYearApprox, false},
		// long second spans
		{"PT86400S", 86400 * time.Second, true},
	}

	for i, c := range cases {
		testPeriodToDuration(t, i, c.value, c.duration, c.precise)
		testPeriodToDuration(t, i, "-"+c.value, -c.duration, c.precise)
	}
}

func testPeriodToDuration(t *testing.T, i int, value string, duration time.Duration, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)
	hint := info(i, "%s %s %v", value, duration, precise)
	pp := MustParse(value)
	d1, prec := pp.Duration()
	g.Expect(d1).To(Equal(duration), hint)
	g.Expect(prec).To(Equal(precise), hint)
	d2 := pp.DurationApprox()
	g.Expect(d2).To(Equal(duration), hint)
}
