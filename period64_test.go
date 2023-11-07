// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/govalues/decimal"
	. "github.com/onsi/gomega"
	"math"
	"strings"
	"testing"
	"time"
)

// shorthand functions

func dec(i int64, s int) decimal.Decimal {
	return decimal.MustNew(i, s)
}

func decI(i int64) decimal.Decimal {
	return decimal.MustNew(i, 0)
}

func decS(s string) decimal.Decimal {
	return decimal.MustParse(s)
}

func add(a, b decimal.Decimal) decimal.Decimal {
	sum, err := a.Add(b)
	if err != nil {
		panic(err)
	}
	return sum
}

//-------------------------------------------------------------------------------------------------

func TestNewHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	const largeInt = math.MaxInt32

	cases := []struct {
		period                  Period64
		hours, minutes, seconds int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period64{seconds: decI(1)}, seconds: 1},
		{period: Period64{minutes: decI(1)}, minutes: 1},
		{period: Period64{hours: decI(1)}, hours: 1},

		{period: Period64{hours: decI(3), minutes: decI(4), seconds: decI(5)}, hours: 3, minutes: 4, seconds: 5},
		{period: Period64{hours: decI(largeInt), minutes: decI(largeInt), seconds: decI(largeInt)}, hours: largeInt, minutes: largeInt, seconds: largeInt},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %dh %dm %ds", i, c.hours, c.minutes, c.seconds), func(t *testing.T) {
			pp := NewHMS(c.hours, c.minutes, c.seconds)
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.Hours()).To(Equal(decimal.MustNew(int64(c.hours), 0)), info(i, c.period))
			g.Expect(pp.HoursInt()).To(Equal(c.hours), info(i, c.period))
			g.Expect(pp.Minutes()).To(Equal(decimal.MustNew(int64(c.minutes), 0)), info(i, c.period))
			g.Expect(pp.MinutesInt()).To(Equal(c.minutes), info(i, c.period))
			g.Expect(pp.Seconds()).To(Equal(decimal.MustNew(int64(c.seconds), 0)), info(i, c.period))
			g.Expect(pp.SecondsInt()).To(Equal(c.seconds), info(i, c.period))

			pn := NewHMS(-c.hours, -c.minutes, -c.seconds)
			en := c.period.Negate()
			g.Expect(pn).To(Equal(en), info(i, en))
			g.Expect(pn.Hours()).To(Equal(decimal.MustNew(int64(-c.hours), 0)), info(i, c.period))
			g.Expect(pn.HoursInt()).To(Equal(-c.hours), info(i, en))
			g.Expect(pn.Minutes()).To(Equal(decimal.MustNew(int64(-c.minutes), 0)), info(i, c.period))
			g.Expect(pn.MinutesInt()).To(Equal(-c.minutes), info(i, en))
			g.Expect(pn.Seconds()).To(Equal(decimal.MustNew(int64(-c.seconds), 0)), info(i, c.period))
			g.Expect(pn.SecondsInt()).To(Equal(-c.seconds), info(i, en))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewYMWD(t *testing.T) {
	g := NewGomegaWithT(t)

	const largeInt = math.MaxInt32

	cases := []struct {
		period                     Period64
		years, months, weeks, days int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period64{days: decI(1)}, days: 1},
		{period: Period64{weeks: decI(1)}, weeks: 1},
		{period: Period64{months: decI(1)}, months: 1},
		{period: Period64{years: decI(1)}, years: 1},

		{period: Period64{years: decI(100), months: decI(222), weeks: decI(404), days: decI(700)}, years: 100, months: 222, weeks: 404, days: 700},
		{period: Period64{years: decI(largeInt), months: decI(largeInt), weeks: decI(largeInt), days: decI(largeInt)}, years: largeInt, months: largeInt, weeks: largeInt, days: largeInt},
	}
	for i, c := range cases {
		pp := NewYMWD(c.years, c.months, c.weeks, c.days)
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(decimal.MustNew(int64(c.years), 0)), info(i, c.period))
		g.Expect(pp.YearsInt()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(decimal.MustNew(int64(c.months), 0)), info(i, c.period))
		g.Expect(pp.MonthsInt()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Weeks()).To(Equal(decimal.MustNew(int64(c.weeks), 0)), info(i, c.period))
		g.Expect(pp.WeeksInt()).To(Equal(c.weeks), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(decimal.MustNew(int64(c.days), 0)), info(i, c.period))
		g.Expect(pp.DaysInt()).To(Equal(c.days), info(i, c.period))

		pn := NewYMWD(-c.years, -c.months, -c.weeks, -c.days)
		en := c.period.Negate()
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Years()).To(Equal(decimal.MustNew(int64(-c.years), 0)), info(i, en))
		g.Expect(pn.YearsInt()).To(Equal(-c.years), info(i, en))
		g.Expect(pn.Months()).To(Equal(decimal.MustNew(int64(-c.months), 0)), info(i, en))
		g.Expect(pn.MonthsInt()).To(Equal(-c.months), info(i, en))
		g.Expect(pn.Weeks()).To(Equal(decimal.MustNew(int64(-c.weeks), 0)), info(i, en))
		g.Expect(pn.WeeksInt()).To(Equal(-c.weeks), info(i, en))
		g.Expect(pn.Days()).To(Equal(decimal.MustNew(int64(-c.days), 0)), info(i, en))
		g.Expect(pn.DaysInt()).To(Equal(-c.days), info(i, en))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewDecimal(t *testing.T) {
	g := NewGomegaWithT(t)

	var (
		largeInt = decI(math.MaxInt32)
		one      = decI(1)
	)

	cases := []struct {
		period                     Period64
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
	}{
		{}, // zero case

		{period: Period64{seconds: one}, seconds: one},
		{period: Period64{minutes: one}, minutes: one},
		{period: Period64{hours: one}, hours: one},
		{period: Period64{days: one}, days: one},
		{period: Period64{weeks: one}, weeks: one},
		{period: Period64{months: one}, months: one},
		{period: Period64{years: one}, years: one},

		{period: Period64{years: largeInt, months: largeInt, weeks: largeInt, days: largeInt, hours: largeInt, minutes: largeInt, seconds: largeInt},
			years: largeInt, months: largeInt, weeks: largeInt, days: largeInt, hours: largeInt, minutes: largeInt, seconds: largeInt},
	}
	for i, c := range cases {
		pp, err := NewDecimal(c.years, c.months, c.weeks, c.days, c.hours, c.minutes, c.seconds)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Weeks()).To(Equal(c.weeks), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))
		g.Expect(pp.Hours()).To(Equal(c.hours), info(i, c.period))
		g.Expect(pp.Minutes()).To(Equal(c.minutes), info(i, c.period))
		g.Expect(pp.Seconds()).To(Equal(c.seconds), info(i, c.period))
	}
}

func TestNewDecimal_error(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                     Period64
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
	}{
		{period: Period64{years: dec(1, 1), months: dec(2, 1), weeks: dec(3, 1), days: dec(4, 1), hours: dec(5, 1), minutes: dec(6, 1), seconds: dec(7, 1)},
			years: dec(1, 1), months: dec(2, 1), weeks: dec(3, 1), days: dec(4, 1), hours: dec(5, 1), minutes: dec(6, 1), seconds: dec(7, 1)},
	}
	for i, c := range cases {
		pp, err := NewDecimal(c.years, c.months, c.weeks, c.days, c.hours, c.minutes, c.seconds)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("0.1Y0.2M0.3W0.4DT0.5H0.6M0.7S"))
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Weeks()).To(Equal(c.weeks), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))
		g.Expect(pp.Hours()).To(Equal(c.hours), info(i, c.period))
		g.Expect(pp.Minutes()).To(Equal(c.minutes), info(i, c.period))
		g.Expect(pp.Seconds()).To(Equal(c.seconds), info(i, c.period))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewOf(t *testing.T) {
	// note: the negative cases are also covered (see below)

	// HMS tests
	testNewOf(t, 1, time.Nanosecond, Period64{seconds: dec(1, 9)})
	testNewOf(t, 2, time.Microsecond, Period64{seconds: dec(1, 6)})
	testNewOf(t, 3, time.Millisecond, Period64{seconds: dec(1, 3)})
	testNewOf(t, 4, 100*time.Millisecond, Period64{seconds: dec(1, 1)})
	testNewOf(t, 5, time.Second, Period64{seconds: one})
	testNewOf(t, 6, time.Minute, Period64{seconds: decI(60)})
	testNewOf(t, 7, time.Hour, Period64{seconds: decI(3600)})
	testNewOf(t, 8, time.Hour+time.Minute+time.Second, Period64{seconds: decI(3661)})
	testNewOf(t, 9, time.Duration(math.MaxInt64), Period64{seconds: dec(math.MaxInt64, 9)})
}

func testNewOf(t *testing.T, i int, source time.Duration, expected Period64) {
	t.Helper()
	testNewOf1(t, i, source, expected)
	testNewOf1(t, i, -source, expected.Negate())
}

func testNewOf1(t *testing.T, i int, source time.Duration, expected Period64) {
	t.Helper()
	g := NewGomegaWithT(t)

	n := NewOf(source)
	rev, _ := expected.Duration()
	info := fmt.Sprintf("%d: source %v expected %+v rev %v", i, source, expected, rev)
	g.Expect(n).To(Equal(expected), info)
	g.Expect(rev).To(Equal(source), info)
}

//-------------------------------------------------------------------------------------------------

func Test_String(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		expected Period
		p64      Period64
	}{
		// note: the negative cases are also covered (see below)

		{expected: "P0D", p64: Period64{}},

		// ones
		{expected: "P1Y", p64: Period64{years: one}},
		{expected: "P1M", p64: Period64{months: one}},
		{expected: "P1W", p64: Period64{weeks: one}},
		{expected: "P1D", p64: Period64{days: one}},
		{expected: "PT1H", p64: Period64{hours: one}},
		{expected: "PT1M", p64: Period64{minutes: one}},
		{expected: "PT1S", p64: Period64{seconds: one}},

		// small fraction
		{expected: "P0.000000001Y", p64: Period64{years: dec(1, 9)}},
		{expected: "P0.000000001M", p64: Period64{months: dec(1, 9)}},
		{expected: "P0.000000001W", p64: Period64{weeks: dec(1, 9)}},
		{expected: "P0.000000001D", p64: Period64{days: dec(1, 9)}},
		{expected: "PT0.000000001H", p64: Period64{hours: dec(1, 9)}},
		{expected: "PT0.000000001M", p64: Period64{minutes: dec(1, 9)}},
		{expected: "PT0.000000001S", p64: Period64{seconds: dec(1, 9)}},

		// 1 + small
		{expected: "P1.000000001Y", p64: Period64{years: add(one, dec(1, 9))}},
		{expected: "P1.000000001M", p64: Period64{months: add(one, dec(1, 9))}},
		{expected: "P1.000000001W", p64: Period64{weeks: add(one, dec(1, 9))}},
		{expected: "P1.000000001D", p64: Period64{days: add(one, dec(1, 9))}},
		{expected: "PT1.000000001H", p64: Period64{hours: add(one, dec(1, 9))}},
		{expected: "PT1.000000001M", p64: Period64{minutes: add(one, dec(1, 9))}},
		{expected: "PT1.000000001S", p64: Period64{seconds: add(one, dec(1, 9))}},

		// other fractions
		{expected: "P0.00000001Y", p64: Period64{years: dec(1, 8)}},
		{expected: "P0.0000001Y", p64: Period64{years: dec(1, 7)}},
		{expected: "P0.000001Y", p64: Period64{years: dec(1, 6)}},
		{expected: "P0.00001Y", p64: Period64{years: dec(1, 5)}},
		{expected: "P0.0001Y", p64: Period64{years: dec(1, 4)}},
		{expected: "P0.001Y", p64: Period64{years: dec(1, 3)}},
		{expected: "P0.01Y", p64: Period64{years: dec(1, 2)}},
		{expected: "P0.1Y", p64: Period64{years: dec(1, 1)}},

		{expected: "P3.9Y", p64: Period64{years: decS("3.9")}},
		{expected: "P3Y6.9M", p64: Period64{years: decI(3), months: decS("6.9")}},
		{expected: "P3Y6M2.9W", p64: Period64{years: decI(3), months: decI(6), weeks: decS("2.9")}},
		{expected: "P3Y6M2W4.9D", p64: Period64{years: decI(3), months: decI(6), weeks: decI(2), days: decS("4.9")}},
		{expected: "P3Y6M2W4DT1.9H", p64: Period64{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decS("1.9")}},
		{expected: "P3Y6M2W4DT1H5.9M", p64: Period64{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decS("5.9")}},
		{expected: "P3Y6M2W4DT1H5M7.9S", p64: Period64{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decI(5), seconds: decS("7.9")}},
		{expected: "-P3Y6M2W4DT1H5M7.9S", p64: Period64{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decI(5), seconds: decS("7.9"), neg: true}},
		{expected: "P-3Y6M-2W4DT-1H5M-7.9S", p64: Period64{years: decI(-3), months: decI(6), weeks: decI(-2), days: decI(4), hours: decI(-1), minutes: decI(5), seconds: decS("-7.9")}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			sp1 := c.p64.Period()
			g.Expect(sp1).To(Equal(c.expected))

			if !c.p64.IsZero() {
				sn := c.p64.Negate().Period()
				ne := "-" + c.expected
				if c.expected[0] == '-' {
					ne = c.expected[1:]
				}
				g.Expect(sn).To(Equal(ne))
			}
		})
	}
}

//-------------------------------------------------------------------------------------------------

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

//-------------------------------------------------------------------------------------------------

func Test_Period64_Sign_Abs_etc(t *testing.T) {
	g := NewGomegaWithT(t)

	z := Zero
	neg := Period64{years: one, months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: true}
	pos := Period64{years: one, months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: false}

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
