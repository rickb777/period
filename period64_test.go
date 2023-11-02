// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	bigdecimal "github.com/shopspring/decimal"
	"testing"
	"time"
)

// shorthand function name
func decI(i int64) decimal {
	return decimal{value: i}
}

func decS(s string) decimal {
	d, err := bigdecimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return newDecimal(d)
}

func Test_String(t *testing.T) {
	//g := NewGomegaWithT(t)

	cases := map[Period]Period64{
		// note: the negative cases are also covered (see below)

		"P0D": {},

		// ones
		"P1Y":  {years: one},
		"P1M":  {months: one},
		"P1W":  {weeks: one},
		"P1D":  {days: one},
		"PT1H": {hours: one},
		"PT1M": {minutes: one},
		"PT1S": {seconds: one},

		// small fraction
		"P0.000000001Y":  {years: one.Shift(-9)},
		"P0.000000001M":  {months: one.Shift(-9)},
		"P0.000000001W":  {weeks: one.Shift(-9)},
		"P0.000000001D":  {days: one.Shift(-9)},
		"PT0.000000001H": {hours: one.Shift(-9)},
		"PT0.000000001M": {minutes: one.Shift(-9)},
		"PT0.000000001S": {seconds: one.Shift(-9)},

		// 1 + small
		"P1.000000001Y":  {years: one.Add(one.Shift(-9))},
		"P1.000000001M":  {months: one.Add(one.Shift(-9))},
		"P1.000000001W":  {weeks: one.Add(one.Shift(-9))},
		"P1.000000001D":  {days: one.Add(one.Shift(-9))},
		"PT1.000000001H": {hours: one.Add(one.Shift(-9))},
		"PT1.000000001M": {minutes: one.Add(one.Shift(-9))},
		"PT1.000000001S": {seconds: one.Add(one.Shift(-9))},

		// other fractions
		"P0.00000001Y": {years: one.Shift(-8)},
		"P0.0000001Y":  {years: one.Shift(-7)},
		"P0.000001Y":   {years: one.Shift(-6)},
		"P0.00001Y":    {years: one.Shift(-5)},
		"P0.0001Y":     {years: one.Shift(-4)},
		"P0.001Y":      {years: one.Shift(-3)},
		"P0.01Y":       {years: one.Shift(-2)},
		"P0.1Y":        {years: one.Shift(-1)},

		"P3.9Y":                  {years: decS("3.9")},
		"P3Y6.9M":                {years: decI(3), months: decS("6.9")},
		"P3Y6M2.9W":              {years: decI(3), months: decI(6), weeks: decS("2.9")},
		"P3Y6M2W4.9D":            {years: decI(3), months: decI(6), weeks: decI(2), days: decS("4.9")},
		"P3Y6M2W4DT1.9H":         {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decS("1.9")},
		"P3Y6M2W4DT1H5.9M":       {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decI(1), minutes: decS("5.9")},
		"P3Y6M2W4DT1H5M7.9S":     {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decI(1), minutes: decI(5), seconds: decS("7.9")},
		"-P3Y6M2W4DT1H5M7.9S":    {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decI(1), minutes: decI(5), seconds: decS("7.9"), neg: true},
		"P-3Y6M-2W4DT-1H5M-7.9S": {years: decI(-3), months: decI(6), weeks: decI(-2), days: decI(4), hours: decI(-1), minutes: decI(5), seconds: decS("-7.9")},
	}

	for expected, p64 := range cases {
		sp1 := p64.Period()
		if sp1 != expected {
			t.Errorf("+ve got %s, expected %s", sp1, expected)
		}

		sp2 := p64.String()
		if sp2 != expected.String() {
			t.Errorf("+ve got %s, expected %s", sp2, expected)
		}

		if !p64.IsZero() {
			sn := p64.Negate().Period()
			ne := "-" + expected
			if expected[0] == '-' {
				ne = expected[1:]
			}
			if sn != ne {
				t.Errorf("-ve got %s, expected %s", sn, ne)
			}
		}
	}
}

func Test_Normalise(t *testing.T) {
	g := NewGomegaWithT(t)

	const (
		both = iota
		precise
		imprecise
	)

	cases := []struct {
		expected Period
		input    Period64
		mode     int
	}{
		// note: the negative cases are also covered (see below)

		{expected: "P0D", mode: both},

		// ones unchanged
		{expected: "P1Y", input: Period64{years: one}, mode: both},
		{expected: "P1M", input: Period64{months: one}, mode: both},
		{expected: "P1W", input: Period64{weeks: one}, mode: both},
		{expected: "P1D", input: Period64{days: one}, mode: both},
		{expected: "PT1H", input: Period64{hours: one}, mode: both},
		{expected: "PT1M", input: Period64{minutes: one}, mode: both},
		{expected: "PT1S", input: Period64{seconds: one}, mode: both},

		// small fraction unchanged
		{expected: "P0.000000001Y", input: Period64{years: one.Shift(-9)}, mode: both},
		{expected: "P0.000000001M", input: Period64{months: one.Shift(-9)}, mode: both},
		{expected: "P0.000000001W", input: Period64{weeks: one.Shift(-9)}, mode: both},
		{expected: "P0.000000001D", input: Period64{days: one.Shift(-9)}, mode: both},
		{expected: "PT0.000000001H", input: Period64{hours: one.Shift(-9)}, mode: both},
		{expected: "PT0.000000001M", input: Period64{minutes: one.Shift(-9)}, mode: both},
		{expected: "PT0.000000001S", input: Period64{seconds: one.Shift(-9)}, mode: both},

		// small overflow discarded
		{expected: "PT60.5S", input: Period64{seconds: decimal{value: 605, exp: -1}}, mode: both},
		{expected: "-PT60.5S", input: Period64{seconds: decimal{value: -605, exp: -1}}, mode: both},

		// no change
		{expected: "PT26.1H", input: Period64{hours: decimal{value: 261, exp: -1}}, mode: precise},
		{expected: "-PT26.1H", input: Period64{hours: decimal{value: -261, exp: -1}}, mode: precise},
		{expected: "P53W", input: Period64{weeks: decimal{value: 53}}, mode: precise},
		{expected: "-P53W", input: Period64{weeks: decimal{value: -53}}, mode: precise},

		// precise normalisation
		{expected: "PT1M5.5S", input: Period64{seconds: decimal{value: 655, exp: -1}}, mode: both},
		{expected: "PT1H1.1S", input: Period64{seconds: decimal{value: 36011, exp: -1}}, mode: both},
		{expected: "PT1H1M1.1S", input: Period64{seconds: decimal{value: 36611, exp: -1}}, mode: both},
		{expected: "P1W2D", input: Period64{days: decimal{value: 9}}, mode: both},
		{expected: "P1Y2M", input: Period64{months: decimal{value: 14}}, mode: both},
		{expected: "-P1Y2M", input: Period64{months: decimal{value: -14}}, mode: both},

		// imprecise normalisation
		{expected: "P1DT2.1H", input: Period64{hours: decimal{value: 261, exp: -1}}, mode: imprecise},
		{expected: "P52W2.1D", input: Period64{days: decimal{value: 3661, exp: -1}}, mode: imprecise},
		{expected: "P54W2DT1H", input: Period64{hours: decimal{value: 9121}}, mode: imprecise},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			if c.mode == precise || c.mode == both {
				sp1 := c.input.Normalise(true)
				g.Expect(sp1.Period()).To(Equal(c.expected), "precise case")
			}

			if c.mode == imprecise || c.mode == both {
				sp1 := c.input.Normalise(false)
				g.Expect(sp1.Period()).To(Equal(c.expected), "imprecise case")
			}
		})
	}
}

func Test_NormaliseSign(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		expected Period
		input    Period64
	}{
		// note: the negative cases are also covered (see below)

		{expected: "P0D"},

		// ones unchanged
		{expected: "P1Y", input: Period64{years: one}},
		{expected: "P1M", input: Period64{months: one}},
		{expected: "P1W", input: Period64{weeks: one}},
		{expected: "P1D", input: Period64{days: one}},
		{expected: "PT1H", input: Period64{hours: one}},
		{expected: "PT1M", input: Period64{minutes: one}},
		{expected: "PT1S", input: Period64{seconds: one}},

		// normalisation
		{expected: "-P1Y", input: Period64{years: negOne}},
		{expected: "-P1M", input: Period64{months: negOne}},
		{expected: "-P1W", input: Period64{weeks: negOne}},
		{expected: "-P1D", input: Period64{days: negOne}},
		{expected: "-PT1H", input: Period64{hours: negOne}},
		{expected: "-PT1M", input: Period64{minutes: negOne}},
		{expected: "-PT1S", input: Period64{seconds: negOne}},

		// complex normalisation
		{expected: "-PT1M1S", input: Period64{minutes: negOne, seconds: negOne}},
		{expected: "-PT1M-1S", input: Period64{minutes: negOne, seconds: one}},
		{expected: "PT1M-1S", input: Period64{minutes: one, seconds: negOne}},

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

//-------------------------------------------------------------------------------------------------

func Test_Period64_Sign_Abs_etc(t *testing.T) {
	g := NewGomegaWithT(t)

	z := Zero
	neg := Period64{years: decI(1), months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: true}
	pos := Period64{years: decI(1), months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: false}

	g.Expect(z.Negate()).To(Equal(z))
	g.Expect(pos.Negate()).To(Equal(neg))
	g.Expect(neg.Negate()).To(Equal(pos))

	g.Expect(z.Abs()).To(Equal(z))
	g.Expect(pos.Abs()).To(Equal(pos))
	g.Expect(neg.Abs()).To(Equal(pos))

	g.Expect(z.Sign()).To(Equal(0))
	g.Expect(pos.Sign()).To(Equal(1))
	g.Expect(neg.Sign()).To(Equal(-1))

	g.Expect(z.IsZero()).To(BeTrue())
	g.Expect(pos.IsZero()).To(BeFalse())
	g.Expect(neg.IsZero()).To(BeFalse())

	g.Expect(z.IsPositive()).To(BeTrue()) // n.b
	g.Expect(pos.IsPositive()).To(BeTrue())
	g.Expect(neg.IsPositive()).To(BeFalse())

	g.Expect(z.IsNegative()).To(BeFalse())
	g.Expect(pos.IsNegative()).To(BeFalse())
	g.Expect(neg.IsNegative()).To(BeTrue())
}

var london *time.Location // UTC + 1 hour during summer

func init() {
	london, _ = time.LoadLocation("Europe/London")
}

func info(i int, m ...interface{}) string {
	if s, ok := m[0].(string); ok {
		m[0] = i
		return fmt.Sprintf("%d "+s, m...)
	}
	return fmt.Sprintf("%d %v", i, m[0])
}
