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
		period                  Period
		hours, minutes, seconds int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period{seconds: decI(1)}, seconds: 1},
		{period: Period{minutes: decI(1)}, minutes: 1},
		{period: Period{hours: decI(1)}, hours: 1},

		{period: Period{hours: decI(3), minutes: decI(4), seconds: decI(5)}, hours: 3, minutes: 4, seconds: 5},
		{period: Period{hours: decI(largeInt), minutes: decI(largeInt), seconds: decI(largeInt)}, hours: largeInt, minutes: largeInt, seconds: largeInt},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %dh %dm %ds", i, c.hours, c.minutes, c.seconds), func(t *testing.T) {
			pp := NewHMS(c.hours, c.minutes, c.seconds)
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.HoursDecimal()).To(Equal(decimal.MustNew(int64(c.hours), 0)), info(i, c.period))
			g.Expect(pp.Hours()).To(Equal(c.hours), info(i, c.period))
			g.Expect(pp.MinutesDecimal()).To(Equal(decimal.MustNew(int64(c.minutes), 0)), info(i, c.period))
			g.Expect(pp.Minutes()).To(Equal(c.minutes), info(i, c.period))
			g.Expect(pp.SecondsDecimal()).To(Equal(decimal.MustNew(int64(c.seconds), 0)), info(i, c.period))
			g.Expect(pp.Seconds()).To(Equal(c.seconds), info(i, c.period))

			pn := NewHMS(-c.hours, -c.minutes, -c.seconds)
			en := c.period.Negate()
			g.Expect(pn).To(Equal(en), info(i, en))
			g.Expect(pn.HoursDecimal()).To(Equal(decimal.MustNew(int64(-c.hours), 0)), info(i, c.period))
			g.Expect(pn.Hours()).To(Equal(-c.hours), info(i, en))
			g.Expect(pn.MinutesDecimal()).To(Equal(decimal.MustNew(int64(-c.minutes), 0)), info(i, c.period))
			g.Expect(pn.Minutes()).To(Equal(-c.minutes), info(i, en))
			g.Expect(pn.SecondsDecimal()).To(Equal(decimal.MustNew(int64(-c.seconds), 0)), info(i, c.period))
			g.Expect(pn.Seconds()).To(Equal(-c.seconds), info(i, en))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewYMD(t *testing.T) {
	g := NewGomegaWithT(t)

	const largeInt = math.MaxInt32

	cases := []struct {
		period              Period
		years, months, days int
	}{
		{}, // zero case

		{period: Period{days: decI(1)}, days: 1},
		{period: Period{months: decI(1)}, months: 1},
		{period: Period{years: decI(1)}, years: 1},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.period), func(t *testing.T) {
			pp := NewYMD(c.years, c.months, c.days)
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
			g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
			g.Expect(pp.Weeks()).To(Equal(0), info(i, c.period))
			g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))
			g.Expect(pp.DaysIncWeeks()).To(Equal(c.days), info(i, c.period))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewYMWD(t *testing.T) {
	g := NewGomegaWithT(t)

	const largeInt = math.MaxInt32

	cases := []struct {
		period                     Period
		years, months, weeks, days int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period{days: decI(1)}, days: 1},
		{period: Period{weeks: decI(1)}, weeks: 1},
		{period: Period{months: decI(1)}, months: 1},
		{period: Period{years: decI(1)}, years: 1},

		{period: Period{years: decI(100), months: decI(222), weeks: decI(404), days: decI(700)}, years: 100, months: 222, weeks: 404, days: 700},
		{period: Period{years: decI(largeInt), months: decI(largeInt), weeks: decI(largeInt), days: decI(largeInt)}, years: largeInt, months: largeInt, weeks: largeInt, days: largeInt},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.period), func(t *testing.T) {
			pp := NewYMWD(c.years, c.months, c.weeks, c.days)
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.YearsDecimal()).To(Equal(decimal.MustNew(int64(c.years), 0)), info(i, c.period))
			g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
			g.Expect(pp.MonthsDecimal()).To(Equal(decimal.MustNew(int64(c.months), 0)), info(i, c.period))
			g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
			g.Expect(pp.WeeksDecimal()).To(Equal(decimal.MustNew(int64(c.weeks), 0)), info(i, c.period))
			g.Expect(pp.Weeks()).To(Equal(c.weeks), info(i, c.period))
			g.Expect(pp.DaysDecimal()).To(Equal(decimal.MustNew(int64(c.days), 0)), info(i, c.period))
			g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))
			g.Expect(pp.DaysIncWeeksDecimal()).To(Equal(decimal.MustNew(int64(7*c.weeks+c.days), 0)), info(i, c.period))
			g.Expect(pp.DaysIncWeeks()).To(Equal(7*c.weeks+c.days), info(i, c.period))

			pn := NewYMWD(-c.years, -c.months, -c.weeks, -c.days)
			en := c.period.Negate()
			g.Expect(pn).To(Equal(en), info(i, en))
			g.Expect(pn.YearsDecimal()).To(Equal(decimal.MustNew(int64(-c.years), 0)), info(i, en))
			g.Expect(pn.Years()).To(Equal(-c.years), info(i, en))
			g.Expect(pn.MonthsDecimal()).To(Equal(decimal.MustNew(int64(-c.months), 0)), info(i, en))
			g.Expect(pn.Months()).To(Equal(-c.months), info(i, en))
			g.Expect(pn.WeeksDecimal()).To(Equal(decimal.MustNew(int64(-c.weeks), 0)), info(i, en))
			g.Expect(pn.Weeks()).To(Equal(-c.weeks), info(i, en))
			g.Expect(pn.DaysDecimal()).To(Equal(decimal.MustNew(int64(-c.days), 0)), info(i, en))
			g.Expect(pn.Days()).To(Equal(-c.days), info(i, en))
			g.Expect(pn.DaysIncWeeksDecimal()).To(Equal(decimal.MustNew(int64(-7*c.weeks-c.days), 0)), info(i, en))
			g.Expect(pn.DaysIncWeeks()).To(Equal(-7*c.weeks-c.days), info(i, en))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestSetGet(t *testing.T) {
	g := NewGomegaWithT(t)

	var (
		two   = decI(2)
		three = decI(3)
		four  = decI(4)
		five  = decI(5)
		six   = decI(6)
		seven = decI(7)
		ten   = decI(10)
	)

	p0 := New(1, 2, 3, 4, 5, 6, 7)

	cases := []struct {
		field                      Designator
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
	}{
		{field: Year, years: ten, months: two, weeks: three, days: four, hours: five, minutes: six, seconds: seven},
		{field: Month, years: one, months: ten, weeks: three, days: four, hours: five, minutes: six, seconds: seven},
		{field: Week, years: one, months: two, weeks: ten, days: four, hours: five, minutes: six, seconds: seven},
		{field: Day, years: one, months: two, weeks: three, days: ten, hours: five, minutes: six, seconds: seven},
		{field: Hour, years: one, months: two, weeks: three, days: four, hours: ten, minutes: six, seconds: seven},
		{field: Minute, years: one, months: two, weeks: three, days: four, hours: five, minutes: ten, seconds: seven},
		{field: Second, years: one, months: two, weeks: three, days: four, hours: five, minutes: six, seconds: ten},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %v", i, c.field.Byte()), func(t *testing.T) {
			p1, err := p0.SetField(ten, c.field)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(p1.YearsDecimal()).To(Equal(c.years))
			g.Expect(p1.MonthsDecimal()).To(Equal(c.months))
			g.Expect(p1.WeeksDecimal()).To(Equal(c.weeks))
			g.Expect(p1.DaysDecimal()).To(Equal(c.days))
			g.Expect(p1.HoursDecimal()).To(Equal(c.hours))
			g.Expect(p1.MinutesDecimal()).To(Equal(c.minutes))
			g.Expect(p1.SecondsDecimal()).To(Equal(c.seconds))

			p2 := p0.SetInt(10, c.field)
			g.Expect(p2.Years()).To(BeEquivalentTo(c.years.Coef()))
			g.Expect(p2.Months()).To(BeEquivalentTo(c.months.Coef()))
			g.Expect(p2.Weeks()).To(BeEquivalentTo(c.weeks.Coef()))
			g.Expect(p2.Days()).To(BeEquivalentTo(c.days.Coef()))
			g.Expect(p2.Hours()).To(BeEquivalentTo(c.hours.Coef()))
			g.Expect(p2.Minutes()).To(BeEquivalentTo(c.minutes.Coef()))
			g.Expect(p2.Seconds()).To(BeEquivalentTo(c.seconds.Coef()))

			d0 := p0.GetField(c.field)
			g.Expect(d0).NotTo(Equal(ten))

			d1 := p1.GetField(c.field)
			g.Expect(d1).To(Equal(ten))

			v0 := p0.GetInt(c.field)
			g.Expect(v0).NotTo(Equal(10))

			v2 := p2.GetInt(c.field)
			g.Expect(v2).To(Equal(10))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewDecimal(t *testing.T) {
	g := NewGomegaWithT(t)

	var (
		largeInt = decI(math.MaxInt64)
		smallInt = dec(1, decimal.MaxScale)
	)

	cases := []struct {
		period                     Period
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
	}{
		{}, // zero case

		{period: Period{seconds: one}, seconds: one},
		{period: Period{minutes: one}, minutes: one},
		{period: Period{hours: one}, hours: one},
		{period: Period{days: one}, days: one},
		{period: Period{weeks: one}, weeks: one},
		{period: Period{months: one}, months: one},
		{period: Period{years: one}, years: one},

		{
			period: Period{years: largeInt, months: largeInt, weeks: largeInt, days: largeInt, hours: largeInt, minutes: largeInt, seconds: largeInt},
			years:  largeInt, months: largeInt, weeks: largeInt, days: largeInt, hours: largeInt, minutes: largeInt, seconds: largeInt,
		},
		{
			period: Period{years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: smallInt},
			years:  decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: smallInt,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.period), func(t *testing.T) {
			pp := MustNewDecimal(c.years, c.months, c.weeks, c.days, c.hours, c.minutes, c.seconds)
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.YearsDecimal()).To(Equal(c.years), info(i, c.period))
			g.Expect(pp.MonthsDecimal()).To(Equal(c.months), info(i, c.period))
			g.Expect(pp.WeeksDecimal()).To(Equal(c.weeks), info(i, c.period))
			g.Expect(pp.DaysDecimal()).To(Equal(c.days), info(i, c.period))
			g.Expect(pp.HoursDecimal()).To(Equal(c.hours), info(i, c.period))
			g.Expect(pp.MinutesDecimal()).To(Equal(c.minutes), info(i, c.period))
			g.Expect(pp.SecondsDecimal()).To(Equal(c.seconds), info(i, c.period))
		})
	}
}

func TestNewDecimal_error1(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                     Period
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
	}{
		{
			period: Period{years: dec(1, 1), months: dec(2, 1), weeks: dec(3, 1), days: dec(4, 1), hours: dec(5, 1), minutes: dec(6, 1), seconds: dec(7, 1)},
			years:  dec(1, 1), months: dec(2, 1), weeks: dec(3, 1), days: dec(4, 1), hours: dec(5, 1), minutes: dec(6, 1), seconds: dec(7, 1),
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.period), func(t *testing.T) {
			pp, err := NewDecimal(c.years, c.months, c.weeks, c.days, c.hours, c.minutes, c.seconds)
			g.Expect(err).To(HaveOccurred())
			g.Expect(err.Error()).To(ContainSubstring("found YMWD/HMS fractions in P0.1Y0.2M0.3W0.4DT0.5H0.6M0.7S"))
			g.Expect(pp).To(Equal(c.period), info(i, c.period))
			g.Expect(pp.YearsDecimal()).To(Equal(c.years), info(i, c.period))
			g.Expect(pp.MonthsDecimal()).To(Equal(c.months), info(i, c.period))
			g.Expect(pp.WeeksDecimal()).To(Equal(c.weeks), info(i, c.period))
			g.Expect(pp.DaysDecimal()).To(Equal(c.days), info(i, c.period))
			g.Expect(pp.HoursDecimal()).To(Equal(c.hours), info(i, c.period))
			g.Expect(pp.MinutesDecimal()).To(Equal(c.minutes), info(i, c.period))
			g.Expect(pp.SecondsDecimal()).To(Equal(c.seconds), info(i, c.period))
		})
	}
}

func TestNewDecimal_error2(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		years, months, weeks, days decimal.Decimal
		hours, minutes, seconds    decimal.Decimal
		message                    string
	}{
		// year fraction
		{
			years: dec(1, 1), months: decI(2), weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found Y fractions in P0.1Y2M",
		},
		{
			years: dec(1, 1), months: decimal.Zero, weeks: decI(2), days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found Y fractions in P0.1Y2W",
		},
		{
			years: dec(1, 1), months: decimal.Zero, weeks: decimal.Zero, days: decI(2), hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found Y fractions in P0.1Y2D",
		},
		{
			years: dec(1, 1), months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decI(2), minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found Y fractions in P0.1YT2H",
		},
		{
			years: dec(1, 1), months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decI(2), seconds: decimal.Zero,
			message: "found Y fractions in P0.1YT2M",
		},
		{
			years: dec(1, 1), months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decI(2),
			message: "found Y fractions in P0.1YT2S",
		},
		// month fraction
		{
			years: decimal.Zero, months: dec(1, 1), weeks: decI(2), days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found M fractions in P0.1M2W",
		},
		{
			years: decimal.Zero, months: dec(1, 1), weeks: decimal.Zero, days: decI(2), hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found M fractions in P0.1M2D",
		},
		{
			years: decimal.Zero, months: dec(1, 1), weeks: decimal.Zero, days: decimal.Zero, hours: decI(2), minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found M fractions in P0.1MT2H",
		},
		{
			years: decimal.Zero, months: dec(1, 1), weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decI(2), seconds: decimal.Zero,
			message: "found M fractions in P0.1MT2M",
		},
		{
			years: decimal.Zero, months: dec(1, 1), weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decI(2),
			message: "found M fractions in P0.1MT2S",
		},
		// week fraction
		{
			years: decimal.Zero, months: decimal.Zero, weeks: dec(1, 1), days: decI(2), hours: decimal.Zero, minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found W fractions in P0.1W2D",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: dec(1, 1), days: decimal.Zero, hours: decI(2), minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found W fractions in P0.1WT2H",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: dec(1, 1), days: decimal.Zero, hours: decimal.Zero, minutes: decI(2), seconds: decimal.Zero,
			message: "found W fractions in P0.1WT2M",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: dec(1, 1), days: decimal.Zero, hours: decimal.Zero, minutes: decimal.Zero, seconds: decI(2),
			message: "found W fractions in P0.1WT2S",
		},
		// day fraction
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: dec(1, 1), hours: decI(2), minutes: decimal.Zero, seconds: decimal.Zero,
			message: "found D fractions in P0.1DT2H",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: dec(1, 1), hours: decimal.Zero, minutes: decI(2), seconds: decimal.Zero,
			message: "found D fractions in P0.1DT2M",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: dec(1, 1), hours: decimal.Zero, minutes: decimal.Zero, seconds: decI(2),
			message: "found D fractions in P0.1DT2S",
		},
		// hour fraction
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: dec(1, 1), minutes: decI(2), seconds: decimal.Zero,
			message: "found H fractions in PT0.1H2M",
		},
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: dec(1, 1), minutes: decimal.Zero, seconds: decI(2),
			message: "found H fractions in PT0.1H2S",
		},
		// minute fraction
		{
			years: decimal.Zero, months: decimal.Zero, weeks: decimal.Zero, days: decimal.Zero, hours: decimal.Zero, minutes: dec(1, 1), seconds: decI(2),
			message: "found M fractions in PT0.1M2S",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.message), func(t *testing.T) {
			_, err := NewDecimal(c.years, c.months, c.weeks, c.days, c.hours, c.minutes, c.seconds)
			g.Expect(err).To(HaveOccurred())
			g.Expect(err.Error()).To(ContainSubstring(c.message))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewOf(t *testing.T) {
	// note: the negative cases are also covered (see below)

	// HMS tests
	testNewOf(t, 1, time.Nanosecond, Period{seconds: dec(1, 9)})
	testNewOf(t, 2, time.Microsecond, Period{seconds: dec(1, 6)})
	testNewOf(t, 3, time.Millisecond, Period{seconds: dec(1, 3)})
	testNewOf(t, 4, 100*time.Millisecond, Period{seconds: dec(1, 1)})
	testNewOf(t, 5, time.Second, Period{seconds: one})
	testNewOf(t, 6, time.Minute, Period{seconds: decI(60)})
	testNewOf(t, 7, time.Hour, Period{seconds: decI(3600)})
	testNewOf(t, 8, time.Hour+time.Minute+time.Second, Period{seconds: decI(3661)})
	testNewOf(t, 9, time.Duration(math.MaxInt64), Period{seconds: dec(math.MaxInt64, 9)})
}

func testNewOf(t *testing.T, i int, source time.Duration, expected Period) {
	t.Helper()
	testNewOf1(t, i, source, expected)
	testNewOf1(t, i, -source, expected.Negate())
}

func testNewOf1(t *testing.T, i int, source time.Duration, expected Period) {
	t.Helper()
	g := NewGomegaWithT(t)

	n := NewOf(source)
	rev, _ := expected.Duration()
	info := fmt.Sprintf("%d: source %v expected %+v rev %v", i, source, expected, rev)
	g.Expect(n).To(Equal(expected), info)
	g.Expect(rev).To(Equal(source), info)
}

//-------------------------------------------------------------------------------------------------

func TestBetween(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		// note: the negative cases are also covered (see below)

		{now, now, Period{}},

		// simple positive date calculations
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 1), Period{seconds: dec(1, 3)}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(4), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(1), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(4), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(3), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(4), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{weeks: decI(4), days: decI(3), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{weeks: decI(26), hours: decI(1), minutes: decI(1), seconds: dec(1001, 3)}},

		// less than one month
		{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{weeks: decI(4), days: decI(2)}},
		{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{weeks: decI(3), days: decI(6)}}, // non-leap
		{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{weeks: decI(4)}},                // leap year
		{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{weeks: decI(4), days: decI(2)}},
		{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{weeks: decI(4), days: decI(1)}},
		{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{weeks: decI(4), days: decI(2)}},
		{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{weeks: decI(4), days: decI(1)}},

		// BST drops an hour at the daylight-saving transition
		{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 1), Period{weeks: decI(26), minutes: decI(1), seconds: dec(1001, 3)}},

		// daytime only
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{seconds: dec(5, 1)}},
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{hours: decI(2), minutes: decI(1), seconds: dec(35, 1)}},
		{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{hours: decI(2), minutes: decI(1), seconds: dec(25, 1)}},

		// different dates and times
		{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{weeks: decI(16), days: decI(6), hours: decI(4), minutes: decI(6), seconds: decI(7)}},
		{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{weeks: decI(16), days: decI(6), hours: decI(3), minutes: decI(6), seconds: decI(7)}},

		// earlier month in later year
		{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{weeks: decI(2), days: decI(5), hours: decI(5), minutes: decI(6), seconds: decI(7)}},
		{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{weeks: decI(47), days: decI(3), hours: decI(18), minutes: decI(53), seconds: dec(525, 1)}},

		// larger ranges
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{weeks: decI(417), days: decI(2), seconds: decI(1)}},
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2017, 12, 21, 0, 0, 2, 0), Period{weeks: decI(468), days: decI(0), seconds: decI(1)}},
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2017, 12, 22, 0, 0, 2, 0), Period{weeks: decI(468), days: decI(1), seconds: decI(1)}},
		{utc(2009, 1, 1, 10, 10, 10, 00), utc(2017, 12, 23, 5, 5, 5, 5), Period{weeks: decI(468), days: decI(1), hours: decI(18), minutes: decI(54), seconds: dec(55005, 3)}},
		{utc(1900, 1, 1, 0, 0, 1, 0), utc(2009, 12, 31, 0, 0, 2, 0), Period{weeks: decI(5739), days: decI(3), seconds: decI(1)}},

		{japan(2021, 3, 1, 0, 0, 0, 0), japan(2021, 9, 7, 0, 0, 0, 0), Period{weeks: decI(27), days: decI(1)}},
		{japan(2021, 3, 1, 0, 0, 0, 0), utc(2021, 9, 7, 0, 0, 0, 0), Period{weeks: decI(27), days: decI(1), hours: decI(9)}},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			pp := Between(c.a, c.b).Normalise(false)
			g.Expect(pp).To(Equal(c.expected), info(i, c.expected))

			pn := Between(c.b, c.a).Normalise(false)
			en := c.expected.Negate()
			g.Expect(pn).To(Equal(en), info(i, en))
		})
	}
}

//-------------------------------------------------------------------------------------------------

func Test_Period64_Sign_Abs_etc(t *testing.T) {
	g := NewGomegaWithT(t)

	z := Zero
	neg := Period{years: one, months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: true}
	pos := Period{years: one, months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: false}

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

var (
	london *time.Location // UTC + 1 hour during summer
	tokyo  *time.Location // UTC + 1 hour during summer
)

func init() {
	london = mustLoadLocation("Europe/London")
	tokyo = mustLoadLocation("Asia/Tokyo")
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load %s: %v", name, err))
	}
	return loc
}

func info(i int, m ...interface{}) string {
	if s, ok := m[0].(string); ok {
		m[0] = i
		return fmt.Sprintf("%d "+s, m...)
	}
	return fmt.Sprintf("%d %v", i, m[0])
}

func utc(year int, month time.Month, day, hour, min, sec, msec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, msec*int(time.Millisecond), time.UTC)
}

func bst(year int, month time.Month, day, hour, min, sec, msec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, msec*int(time.Millisecond), london)
}

func japan(year int, month time.Month, day, hour, min, sec, msec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, msec*int(time.Millisecond), tokyo)
}

//-------------------------------------------------------------------------------------------------

func Test_OnlyYMWD(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "P1Y2M3D"},
		{"-P6Y5M4DT3H2M1S", "-P6Y5M4D"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyYMWD()
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

func Test_OnlyHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "PT4H5M6S"},
		{"-P6Y5M4DT3H2M1S", "-PT3H2M1S"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyHMS()
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}
