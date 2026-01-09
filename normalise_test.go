// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/rickb777/expect"
)

func Test_Normalise(t *testing.T) {
	cases := []struct {
		input     ISOString
		precise   ISOString
		imprecise ISOString
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
			expect.String(sp1.Period()).Info("precise +ve case").ToBe(t, c.precise)

			if !p1.IsZero() {
				sp1n := p1.Negate().Normalise(true)
				expect.String(sp1n.Period()).Info("precise -ve case").ToBe(t, "-"+c.precise)
			}

			p2 := MustParse(c.input)
			sp2 := p2.Normalise(false)
			expect.String(sp2.Period()).Info("approximate +ve case").ToBe(t, c.imprecise)

			if !p2.IsZero() {
				sp2n := p2.Negate().Normalise(false)
				expect.String(sp2n.Period()).Info("approximate -ve case").ToBe(t, "-"+c.imprecise)
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_NormaliseDaysToYears(t *testing.T) {
	cases := []struct {
		input    ISOString
		expected ISOString
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
			expect.String(sp1.Period()).Info("+ve case").ToBe(t, c.expected)

			if !p1.IsZero() {
				sp1n := p1.Negate().NormaliseDaysToYears()
				expect.String(sp1n.Period()).Info("-ve case").ToBe(t, "-"+c.expected)
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_SimplifyWeeksToDays(t *testing.T) {
	cases := []struct {
		input    ISOString
		expected ISOString
	}{
		// note: the negative cases are also covered (see below)

		{input: "P0D", expected: "P0D"},

		// ones unchanged
		{input: "P1Y", expected: "P1Y"},
		{input: "P1M", expected: "P1M"},
		{input: "P1D", expected: "P1D"},
		{input: "PT1H", expected: "PT1H"},
		{input: "PT1M", expected: "PT1M"},
		{input: "PT1S", expected: "PT1S"},

		// simplified
		{input: "P1W", expected: "P7D"},
		{input: "P2W1D", expected: "P15D"},
		{input: "P2W-2D", expected: "P12D"},
		{input: "P10W", expected: "P70D"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			p1 := MustParse(c.input)
			sp1 := p1.SimplifyWeeksToDays()
			expect.String(sp1.Period()).Info("precise +ve case").ToBe(t, c.expected)

			if !p1.IsZero() {
				sp1n := p1.Negate().SimplifyWeeksToDays()
				expect.String(sp1n.Period()).Info("precise -ve case").ToBe(t, "-"+c.expected)
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_Simplify(t *testing.T) {
	var extremeMinSec = ISOString(fmt.Sprintf("PT%dM0.%s1S", int64(math.MaxInt64), strings.Repeat("0", 18)))

	cases := []struct {
		input     ISOString
		precise   ISOString
		imprecise ISOString
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

		{input: "P0.1W", precise: "P0.1W", imprecise: "P0.1W"}, // weeks simplification is nuanced
		{input: "P1M0.1W", precise: "P1M0.7D", imprecise: "P1M0.7D"},

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
			expect.String(sp1.Period()).Info("precise +ve case").ToBe(t, c.precise)

			if !p1.IsZero() {
				sp1n := p1.Negate().Simplify(true)
				expect.String(sp1n.Period()).Info("precise -ve case").ToBe(t, "-"+c.precise)
			}

			p2 := MustParse(c.input)
			sp2 := p2.Simplify(false)
			expect.String(sp2.Period()).Info("approximate +ve case").ToBe(t, c.imprecise)

			if !p2.IsZero() {
				sp2n := p2.Negate().Simplify(false)
				expect.String(sp2n.Period()).Info("approximate -ve case").ToBe(t, "-"+c.imprecise)
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_normaliseSign(t *testing.T) {
	cases := []struct {
		expected ISOString
		input    Period
	}{
		{expected: "P0D"},

		{expected: "P1Y", input: Period{years: one}},
		{expected: "P1M", input: Period{months: one}},
		{expected: "P1W", input: Period{weeks: one}},
		{expected: "P1D", input: Period{days: one}},
		{expected: "PT1H", input: Period{hours: one}},
		{expected: "PT1M", input: Period{minutes: one}},
		{expected: "PT1S", input: Period{seconds: one}},

		{expected: "P1Y", input: Period{years: negOne, neg: true}},
		{expected: "P1M", input: Period{months: negOne, neg: true}},
		{expected: "P1W", input: Period{weeks: negOne, neg: true}},
		{expected: "P1D", input: Period{days: negOne, neg: true}},
		{expected: "PT1H", input: Period{hours: negOne, neg: true}},
		{expected: "PT1M", input: Period{minutes: negOne, neg: true}},
		{expected: "PT1S", input: Period{seconds: negOne, neg: true}},

		{expected: "-P1Y", input: Period{years: negOne, neg: false}},
		{expected: "-P1M", input: Period{months: negOne, neg: false}},
		{expected: "-P1W", input: Period{weeks: negOne, neg: false}},
		{expected: "-P1D", input: Period{days: negOne, neg: false}},
		{expected: "-PT1H", input: Period{hours: negOne, neg: false}},
		{expected: "-PT1M", input: Period{minutes: negOne, neg: false}},
		{expected: "-PT1S", input: Period{seconds: negOne, neg: false}},

		{expected: "-P1Y", input: Period{years: one, neg: true}},
		{expected: "-P1M", input: Period{months: one, neg: true}},
		{expected: "-P1W", input: Period{weeks: one, neg: true}},
		{expected: "-P1D", input: Period{days: one, neg: true}},
		{expected: "-PT1H", input: Period{hours: one, neg: true}},
		{expected: "-PT1M", input: Period{minutes: one, neg: true}},
		{expected: "-PT1S", input: Period{seconds: one, neg: true}},

		// complex normalisation
		{expected: "-PT1M-1S", input: Period{minutes: one, seconds: negOne, neg: true}},
		{expected: "-PT1M1S", input: Period{minutes: one, seconds: one, neg: true}},
		{expected: "PT1M-1S", input: Period{minutes: one, seconds: negOne}},

		{expected: "PT1H1M1S", input: Period{hours: negOne, minutes: negOne, seconds: negOne, neg: true}}, // 111
		{expected: "PT1H1M-1S", input: Period{hours: negOne, minutes: negOne, seconds: one, neg: true}},   // 110
		{expected: "PT1H-1M1S", input: Period{hours: negOne, minutes: one, seconds: negOne, neg: true}},   // 101
		{expected: "PT1H-1M-1S", input: Period{hours: negOne, minutes: one, seconds: one, neg: true}},     // 100
		{expected: "-PT1H-1M-1S", input: Period{hours: one, minutes: negOne, seconds: negOne, neg: true}}, // 011
		{expected: "-PT1H-1M1S", input: Period{hours: one, minutes: negOne, seconds: one, neg: true}},     // 010
		{expected: "-PT1H1M-1S", input: Period{hours: one, minutes: one, seconds: negOne, neg: true}},     // 001
		{expected: "-PT1H1M1S", input: Period{hours: one, minutes: one, seconds: one, neg: true}},         // 000

		{expected: "-PT1H1M1S", input: Period{hours: negOne, minutes: negOne, seconds: negOne}}, // 111
		{expected: "-PT1H1M-1S", input: Period{hours: negOne, minutes: negOne, seconds: one}},   // 110
		{expected: "-PT1H-1M1S", input: Period{hours: negOne, minutes: one, seconds: negOne}},   // 101
		{expected: "-PT1H-1M-1S", input: Period{hours: negOne, minutes: one, seconds: one}},     // 100
		{expected: "PT1H-1M-1S", input: Period{hours: one, minutes: negOne, seconds: negOne}},   // 011
		{expected: "PT1H-1M1S", input: Period{hours: one, minutes: negOne, seconds: one}},       // 010
		{expected: "PT1H1M-1S", input: Period{hours: one, minutes: one, seconds: negOne}},       // 001
		{expected: "PT1H1M1S", input: Period{hours: one, minutes: one, seconds: one}},           // 000
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			sp1 := c.input.normaliseSign()
			expect.String(sp1.Period()).ToBe(t, c.expected)
		})
	}
}
