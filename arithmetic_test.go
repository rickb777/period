// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/govalues/decimal"
	"github.com/rickb777/expect"
)

var filter = cmpopts.IgnoreUnexported(Period{})

func Test_Add_Subtract(t *testing.T) {
	cases := []struct {
		one, two        ISOString
		sum, difference ISOString
	}{
		// simple cases
		{"P0D", "P0D", "P0D", "P0D"},
		{"P1Y", "P1Y", "P2Y", "P0D"},
		{"P1M", "P1M", "P2M", "P0D"},
		{"P1W", "P1W", "P2W", "P0D"},
		{"P1D", "P1D", "P2D", "P0D"},
		{"PT1H", "PT1H", "PT2H", "P0D"},
		{"PT1M", "PT1M", "PT2M", "P0D"},
		{"PT1S", "PT1S", "PT2S", "P0D"},

		{"-P0D", "-P0D", "-P0D", "P0D"},
		{"-P1Y", "-P1Y", "-P2Y", "P0D"},
		{"-P1M", "-P1M", "-P2M", "P0D"},
		{"-P1W", "-P1W", "-P2W", "P0D"},
		{"-P1D", "-P1D", "-P2D", "P0D"},
		{"-PT1H", "-PT1H", "-PT2H", "P0D"},
		{"-PT1M", "-PT1M", "-PT2M", "P0D"},
		{"-PT1S", "-PT1S", "-PT2S", "P0D"},

		{"P0Y", "P1Y", "P1Y", "-P1Y"},
		{"P1Y", "P1M", "P1Y1M", "P1Y-1M"},
		{"P1M", "P1W", "P1M1W", "P1M-1W"},
		{"P1W", "P1D", "P1W1D", "P1W-1D"},
		{"P1D", "PT1H", "P1DT1H", "P1DT-1H"},
		{"PT1H", "PT1M", "PT1H1M", "PT1H-1M"},
		{"PT1M", "PT1S", "PT1M1S", "PT1M-1S"},

		{"P7Y6M5W2DT6H4M2S", "P1Y2M3W2DT3H2M1S", "P8Y8M8W4DT9H6M3S", "P6Y4M2WT3H2M1S"},
		{"P3Y3M3W3DT3H3M3S", "-P3Y3M3W3DT3H3M3S", "P0D", "P6Y6M6W6DT6H6M6S"},
		{"P3Y3M3W3DT3H3M3S", "-P2Y2M2W2DT2H2M2S", "P1Y1M1W1DT1H1M1S", "P5Y5M5W5DT5H5M5S"},
		{"P1Y1M1W1DT1H1M1S", "-P2Y2M2W2DT2H2M2S", "-P1Y1M1W1DT1H1M1S", "P3Y3M3W3DT3H3M3S"},
		{"P1Y2M3W4D", "PT5H6M7S", "P1Y2M3W4DT5H6M7S", "P1Y2M3W4DT-5H-6M-7S"},

		// cases needing borrow/carry
		{"PT16M40S", "PT1000S", "PT33M20S", "P0D"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s %s", i, c.one, c.two), func(t *testing.T) {
			a := MustParse(c.one)
			b := MustParse(c.two)

			s, err := a.Add(b)
			expect.Error(err).Not().ToHaveOccurred(t)
			expect.Any(s).Info("%d %s + %s = %s", i, c.one, c.two, s).Using(filter).ToBe(t, MustParse(c.sum))

			d, err := a.Subtract(b)
			expect.Error(err).Not().ToHaveOccurred(t)
			expect.Any(d).Info("%d %s + %s = %s", i, c.one, c.two, s).Using(filter).ToBe(t, MustParse(c.difference))
		})
	}
}

func Test_AddTo(t *testing.T) {
	const millisec = 1000000
	const second = 1000 * millisec
	const minute = 60 * second
	const hour = 60 * minute

	est := mustLoadLocation("America/New_York")
	aest := mustLoadLocation("Australia/Sydney")

	times := []time.Time{
		// A conveniently round number but with non-zero nanoseconds (14 July 2017 @ 2:40am UTC)
		time.Unix(1500000000, 1).UTC(),
		// This specific time fails for EST due behaviour of time.Time.AddDate
		time.Date(2020, 11, 1, 1, 0, 0, 0, est),
		// Testing to ensure that adding zero with a DST day does not blow up due behaviour of time.Time.AddDate
		time.Date(2024, 4, 6, 15, 0, 0, 0, time.UTC).In(aest),
	}

	for _, t0 := range times {
		cases := []struct {
			value   string
			result  time.Time
			precise bool
		}{
			// precise cases
			{value: "P0D", result: t0, precise: true},
			{value: "PT1S", result: t0.Add(second), precise: true},
			{value: "PT0.1S", result: t0.Add(100 * millisec), precise: true},
			{value: "-PT0.1S", result: t0.Add(-100 * millisec), precise: true},
			{value: "PT1H-0.1S", result: t0.Add(3599900 * millisec), precise: true},
			{value: "PT3276S", result: t0.Add(3276 * second), precise: true},
			{value: "PT1M", result: t0.Add(60 * second), precise: true},
			{value: "PT0.1M", result: t0.Add(6 * second), precise: true},
			{value: "PT3276M", result: t0.Add(3276 * minute), precise: true},
			{value: "PT1H", result: t0.Add(hour), precise: true},
			{value: "PT0.1H", result: t0.Add(6 * minute), precise: true},
			{value: "PT3276H", result: t0.Add(3276 * hour), precise: true},
			{value: "P1D", result: t0.AddDate(0, 0, 1), precise: true},
			{value: "P3276D", result: t0.AddDate(0, 0, 3276), precise: true},
			{value: "P1M", result: t0.AddDate(0, 1, 0), precise: true},
			{value: "P3276M", result: t0.AddDate(0, 3276, 0), precise: true},
			{value: "P1Y", result: t0.AddDate(1, 0, 0), precise: true},
			{value: "-P1Y", result: t0.AddDate(-1, 0, 0), precise: true},
			{value: "P3276Y", result: t0.AddDate(3276, 0, 0), precise: true},   // near the upper limit of range
			{value: "-P3276Y", result: t0.AddDate(-3276, 0, 0), precise: true}, // near the lower limit of range
			{value: "P1DT1M", result: t0.AddDate(0, 0, 1).Add(minute), precise: true},
			{value: "P1DT-1M", result: t0.AddDate(0, 0, 1).Add(-minute), precise: true},
			{value: "P1MT1M", result: t0.AddDate(0, 1, 0).Add(minute), precise: true},
			{value: "P1YT-1S", result: t0.AddDate(1, 0, 0).Add(-second), precise: true},

			// approximate cases
			{value: "P0.1D", result: t0.Add(144 * minute), precise: false},
			{value: "-P0.1D", result: t0.Add(-144 * minute), precise: false},
			{value: "P0.1M", result: t0.Add(oneMonthApprox / 10), precise: false},
			{value: "P0.1Y", result: t0.Add(oneYearApprox / 10), precise: false},
		}
		for i, c := range cases {
			t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
				t1, prec := MustParse(c.value).AddTo(t0)

				hint := info(i, "value=%s t0=%s, t1=%s, exp=%s", c.value,
					t0.Format(time.RFC3339Nano), t1.Format(time.RFC3339Nano), c.result.Format(time.RFC3339Nano))
				expect.Bool(t1.Equal(c.result)).Info(hint).ToBeTrue(t)
				expect.Bool(prec).Info(hint).ToBe(t, c.precise)
			})
		}
	}
}

func Test_Mul(t *testing.T) {
	cases := []struct {
		input    ISOString
		factor   decimal.Decimal
		expected ISOString
	}{
		{input: "P0D", factor: decI(2), expected: "P0D"},
		{input: "P1D", factor: decI(2), expected: "P2D"},
		{input: "P1D", factor: decI(0), expected: "P0D"},
		{input: "P1D", factor: decI(365), expected: "P365D"},
		{input: "P1M", factor: decI(2), expected: "P2M"},
		{input: "P1M", factor: decI(12), expected: "P12M"},
		{input: "P1Y", factor: decI(2), expected: "P2Y"},
		{input: "P1Y", factor: decI(0), expected: "P0Y"},
		{input: "PT1H", factor: decI(2), expected: "PT2H"},
		{input: "PT1M", factor: decI(2), expected: "PT2M"},
		{input: "PT1S", factor: decI(2), expected: "PT2S"},
		{input: "P1D", factor: dec(5, 1), expected: "P0.5D"},
		{input: "P1M", factor: dec(5, 1), expected: "P0.5M"},
		{input: "P1Y", factor: dec(5, 1), expected: "P0.5Y"},
		{input: "PT1H", factor: dec(5, 1), expected: "PT0.5H"},
		{input: "PT1H", factor: dec(1, 1), expected: "PT0.1H"},
		{input: "PT1M", factor: dec(5, 1), expected: "PT0.5M"},
		{input: "PT1S", factor: dec(5, 1), expected: "PT0.5S"},
		{input: "PT0.000001S", factor: dec(1, 6), expected: "PT0.000000000001S"},
		{input: "-PT0.000001S", factor: dec(1, 6), expected: "-PT0.000000000001S"},
		{input: "PT1H", factor: dec(2777778, 9), expected: "PT0.002777778H"}, // 1 / 3600
		{input: "PT1M", factor: decI(60), expected: "PT60M"},
		{input: "PT1S", factor: decI(60), expected: "PT60S"},
		{input: "PT1S", factor: decI(86400), expected: "PT86400S"},
		{input: "PT1S", factor: decI(86400000), expected: "PT86400000S"},
		{input: "P365.2425D", factor: decI(10), expected: "P3652.425D"},
		{input: "P1Y2M3W4DT5H6M7S", factor: dec(2, 0), expected: "P2Y4M6W8DT10H12M14S"},
		{input: "P2Y4M6W8DT10H12M14S", factor: dec(-5, 1), expected: "-P1Y2M3W4DT5H6M7S"},
		{input: "-P2Y4M6W8DT10H12M14S", factor: dec(5, 1), expected: "-P1Y2M3W4DT5H6M7S"},
		{input: "-P2Y4M6W8DT10H12M14S", factor: dec(-5, 1), expected: "P1Y2M3W4DT5H6M7S"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.input), func(t *testing.T) {
			s, err := MustParse(c.input).Mul(c.factor)
			expect.Error(err).Not().ToHaveOccurred(t)
			expect.Any(s).Info("%d %s * %s -> %s", i, c.input, c.factor, c.expected).Using(filter).ToBe(t, MustParse(c.expected))
		})
	}
}

func Test_Mul_errors(t *testing.T) {

	cases := []struct {
		input  Period
		factor decimal.Decimal
	}{
		{input: Period{years: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{months: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{weeks: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{days: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{hours: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{minutes: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
		{input: Period{seconds: dec(math.MaxInt64, 0)}, factor: dec(math.MaxInt64, 0)},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.input), func(t *testing.T) {
			_, err := c.input.Mul(c.factor)
			expect.Error(err).ToHaveOccurred(t)
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_TotalDaysApprox(t *testing.T) {

	cases := []struct {
		value      string
		approxDays int
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0},
		{"PT24H", 1},
		{"PT49H", 2},
		{"P1D", 1},
		{"P1M", 30},
		{"P1Y", 365},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td1 := p.TotalDaysApprox()
		expect.Number(td1).Info("%d %v", i, c.value).ToBe(t, c.approxDays)

		td2 := p.Negate().TotalDaysApprox()
		expect.Number(td2).Info("%d %v", i, c.value).ToBe(t, -c.approxDays)
	}
}

//-------------------------------------------------------------------------------------------------

func Test_TotalMonthsApprox(t *testing.T) {
	cases := []struct {
		value        string
		approxMonths int
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0},
		{"P1D", 0},
		{"P30D", 0},
		{"P31D", 1},
		{"P60D", 1},
		{"P62D", 2},
		{"P1M", 1},
		{"P12M", 12},
		{"P2M31D", 3},
		{"P1Y", 12},
		{"P2Y3M", 27},
		{"PT24H", 0},
		{"PT744H", 1},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td1 := p.TotalMonthsApprox()
		expect.Number(td1).Info("%d %v", i, c.value).ToBe(t, c.approxMonths)

		td2 := p.Negate().TotalMonthsApprox()
		expect.Number(td2).Info("%d %v", i, c.value).ToBe(t, -c.approxMonths)
	}
}

//-------------------------------------------------------------------------------------------------

const oneMonthApprox = daysPerMonthE6 * secondsPerDay * time.Microsecond // 30.436875 days
const oneYearApprox = 31556952 * time.Second                             // 365.2425 days

func Test_Duration(t *testing.T) {
	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0, true},

		{"PT1S", 1 * time.Second, true},
		{"PT0.1S", 100 * time.Millisecond, true},
		{"PT0.001S", time.Millisecond, true},
		{"PT0.000001S", time.Microsecond, true},
		{"PT0.000000001S", time.Nanosecond, true},
		{"PT0.0000000001S", 0, false},
		{"PT3276S", 3276 * time.Second, true},

		{"PT1M", 60 * time.Second, true},
		{"PT0.1M", 6 * time.Second, true},
		{"PT0.0001M", 6 * time.Millisecond, true},
		{"PT0.0000001M", 6 * time.Microsecond, true},
		{"PT0.0000000001M", 6 * time.Nanosecond, true},
		{"PT0.00000000001M", 0, false},
		{"PT3276M", 3276 * time.Minute, true},

		{"PT1H", 3600 * time.Second, true},
		{"PT0.1H", 360 * time.Second, true},
		{"PT0.01H", 36 * time.Second, true},
		{"PT0.00001H", 36 * time.Millisecond, true},
		{"PT0.00000001H", 36 * time.Microsecond, true},
		{"PT0.00000000001H", 36 * time.Nanosecond, true},
		{"PT0.0000000000001H", 0, false},
		{"PT3220H", 3220 * time.Hour, true},
		{"PT1H-1M-1S", 3539 * time.Second, true},

		{"P1D", 24 * time.Hour, false},

		// days, months and years conversions are never precise
		{"P0.1D", 144 * time.Minute, false},
		{"P10000D", 10000 * 24 * time.Hour, false},
		{"P1W", 168 * time.Hour, false},
		{"P0.1W", 16*time.Hour + 48*time.Minute, false},
		{"P10000W", 10000 * 7 * 24 * time.Hour, false},
		{"P1M", oneMonthApprox, false},
		{"P0.1M", oneMonthApprox / 10, false},
		{"P3504M", 3504 * oneMonthApprox, false}, // 292 years
		{"P1Y", oneYearApprox, false},
		{"P0.1Y", oneYearApprox / 10, false},
		{"P292Y", 292 * oneYearApprox, false}, // time.Duration represents up to 292 years
		// long second spans
		{"PT86400000S", 86400000 * time.Second, true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			testPeriodToDuration(t, i, c.value, c.duration, c.precise)
			testPeriodToDuration(t, i, "-"+c.value, -c.duration, c.precise)
		})
	}
}

func testPeriodToDuration(t *testing.T, i int, value string, duration time.Duration, precise bool) {
	t.Helper()
	hint := info(i, "%s %s %v", value, duration, precise)
	pp := MustParse(value)
	d1, prec := pp.Duration()
	expect.Number(d1).Info(hint).ToBe(t, duration)
	expect.Bool(prec).Info(hint).ToBe(t, precise)
	d2 := pp.DurationApprox()
	expect.Number(d2).Info(hint).ToBe(t, duration)
}
