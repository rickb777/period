package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_String(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		expected ISOString
		p64      Period
	}{
		// note: the negative cases are also covered (see below)

		{expected: "P0D", p64: Period{}},

		// ones
		{expected: "P1Y", p64: Period{years: one}},
		{expected: "P1M", p64: Period{months: one}},
		{expected: "P1W", p64: Period{weeks: one}},
		{expected: "P1D", p64: Period{days: one}},
		{expected: "PT1H", p64: Period{hours: one}},
		{expected: "PT1M", p64: Period{minutes: one}},
		{expected: "PT1S", p64: Period{seconds: one}},

		// small fraction
		{expected: "P0.000000001Y", p64: Period{years: dec(1, 9)}},
		{expected: "P0.000000001M", p64: Period{months: dec(1, 9)}},
		{expected: "P0.000000001W", p64: Period{weeks: dec(1, 9)}},
		{expected: "P0.000000001D", p64: Period{days: dec(1, 9)}},
		{expected: "PT0.000000001H", p64: Period{hours: dec(1, 9)}},
		{expected: "PT0.000000001M", p64: Period{minutes: dec(1, 9)}},
		{expected: "PT0.000000001S", p64: Period{seconds: dec(1, 9)}},

		// 1 + small
		{expected: "P1.000000001Y", p64: Period{years: add(one, dec(1, 9))}},
		{expected: "P1.000000001M", p64: Period{months: add(one, dec(1, 9))}},
		{expected: "P1.000000001W", p64: Period{weeks: add(one, dec(1, 9))}},
		{expected: "P1.000000001D", p64: Period{days: add(one, dec(1, 9))}},
		{expected: "PT1.000000001H", p64: Period{hours: add(one, dec(1, 9))}},
		{expected: "PT1.000000001M", p64: Period{minutes: add(one, dec(1, 9))}},
		{expected: "PT1.000000001S", p64: Period{seconds: add(one, dec(1, 9))}},

		// other fractions
		{expected: "P0.00000001Y", p64: Period{years: dec(1, 8)}},
		{expected: "P0.0000001Y", p64: Period{years: dec(1, 7)}},
		{expected: "P0.000001Y", p64: Period{years: dec(1, 6)}},
		{expected: "P0.00001Y", p64: Period{years: dec(1, 5)}},
		{expected: "P0.0001Y", p64: Period{years: dec(1, 4)}},
		{expected: "P0.001Y", p64: Period{years: dec(1, 3)}},
		{expected: "P0.01Y", p64: Period{years: dec(1, 2)}},
		{expected: "P0.1Y", p64: Period{years: dec(1, 1)}},

		{expected: "P3.9Y", p64: Period{years: decS("3.9")}},
		{expected: "P3Y6.9M", p64: Period{years: decI(3), months: decS("6.9")}},
		{expected: "P3Y6M2.9W", p64: Period{years: decI(3), months: decI(6), weeks: decS("2.9")}},
		{expected: "P3Y6M2W4.9D", p64: Period{years: decI(3), months: decI(6), weeks: decI(2), days: decS("4.9")}},
		{expected: "P3Y6M2W4DT1.9H", p64: Period{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decS("1.9")}},
		{expected: "P3Y6M2W4DT1H5.9M", p64: Period{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decS("5.9")}},
		{expected: "P3Y6M2W4DT1H5M7.9S", p64: Period{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decI(5), seconds: decS("7.9")}},
		{expected: "-P3Y6M2W4DT1H5M7.9S", p64: Period{years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: one, minutes: decI(5), seconds: decS("7.9"), neg: true}},
		{expected: "P-3Y6M-2W4DT-1H5M-7.9S", p64: Period{years: decI(-3), months: decI(6), weeks: decI(-2), days: decI(4), hours: decI(-1), minutes: decI(5), seconds: decS("-7.9")}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			// check the normal case
			sp1 := c.p64.Period()
			g.Expect(sp1).To(Equal(c.expected))

			// check the negative case
			if !c.p64.IsZero() {
				sn := c.p64.Negate().Period()
				ne := "-" + c.expected
				if c.expected[0] == '-' {
					ne = c.expected[1:]
				}
				g.Expect(sn).To(Equal(ne))
			}

			// also check WriteTo method is consistent and returns the correct count
			buf := simpleBuffer{}
			n, err := c.p64.WriteTo(&buf)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(n).To(Equal(int64(len(string(buf.bs)))))
			g.Expect(string(buf.bs)).To(Equal(string(sp1)))
		})
	}
}

// simpleBuffer intentionally only has Write method.
type simpleBuffer struct {
	bs []byte
}

func (sb *simpleBuffer) Write(bs []byte) (int, error) {
	sb.bs = append(sb.bs, bs...)
	return len(bs), nil
}

//-------------------------------------------------------------------------------------------------

func Test_Format(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period   string
		expected string
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", "zero"},

		{"P1Y1M1W1D", "1 year, 1 month, 1 week, 1 day"},
		{"PT1H1M1S", "1 hour, 1 minute, 1 second"},
		{"P1Y1M1W1DT1H1M1S", "1 year, 1 month, 1 week, 1 day, 1 hour, 1 minute, 1 second"},
		{"P3Y6M5W4DT2H7M9S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9 seconds"},

		{"P1Y", "1 year"},
		{"P3Y", "3 years"},
		{"P1.1Y", "1.1 years"},
		{"P2.5Y", "2.5 years"},

		{"P1M", "1 month"},
		{"P6M", "6 months"},
		{"P1.1M", "1.1 months"},
		{"P2.5M", "2.5 months"},

		{"P1W", "1 week"},
		{"P10W", "10 weeks"},
		{"P1.1W", "1.1 weeks"},
		{"P1M-1W", "1 month, minus 1 week"},

		{"P1D", "1 day"},
		{"P4D", "4 days"},
		{"P1.1D", "1.1 days"},

		{"PT1H", "1 hour"},
		{"PT8H", "8 hours"},
		{"PT1.1H", "1.1 hours"},
		{"P1DT-1H", "1 day, minus 1 hour"},

		{"PT1M", "1 minute"},
		{"PT6M", "6 minutes"},
		{"PT1.1M", "1.1 minutes"},

		{"PT1S", "1 second"},
		{"PT30S", "30 seconds"},
		{"PT1.1S", "1.1 seconds"},
		{"P1YT-1S", "1 year, minus 1 second"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.period), func(t *testing.T) {
			p := MustParse(c.period)
			sp := p.Format()
			g.Expect(sp).To(Equal(c.expected), info(i, "%s -> %s", p, c.expected))

			en := p.Negate()
			sn := en.Format()
			g.Expect(sn).To(Equal(c.expected), info(i, "%s -> %s", en, c.expected))
		})
	}
}
