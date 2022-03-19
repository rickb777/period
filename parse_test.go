package period

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestParseErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value     string
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
		// integer overflow
		{"P2147483648Y", false, ": integer overflow occurred in years", "P2147483648Y"},
		{"P2147483648M", false, ": integer overflow occurred in months", "P2147483648M"},
		{"P2147483648W", false, ": integer overflow occurred in weeks", "P2147483648W"},
		{"P2147483648D", false, ": integer overflow occurred in days", "P2147483648D"},
		{"PT2147483648H", false, ": integer overflow occurred in hours", "PT2147483648H"},
		{"PT2147483648M", false, ": integer overflow occurred in minutes", "PT2147483648M"},
		{"PT2147483648S", false, ": integer overflow occurred in seconds", "PT2147483648S"},
	}
	for i, c := range cases {
		p := Period32{}
		ep := p.Parse(c.value)
		g.Expect(ep).To(HaveOccurred(), info(i, c.value))
		g.Expect(ep.Error()).To(Equal(c.expvalue+c.expected), info(i, c.value))

		en := p.Parse("-" + c.value)
		g.Expect(en).To(HaveOccurred(), info(i, c.value))
		if c.expvalue != "" {
			g.Expect(en.Error()).To(Equal("-"+c.expvalue+c.expected), info(i, c.value))
		} else {
			g.Expect(en.Error()).To(Equal(c.expected), info(i, c.value))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestParsePeriodVerbatim(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		reversed string
		period   Period32
	}{
		// zero
		{"P0D", "P0D", Period32{}},
		// special zero cases: parse is not identity when reversed
		{"P0", "P0D", Period32{}},
		{"P0Y", "P0D", Period32{}},
		{"P0M", "P0D", Period32{}},
		{"P0W", "P0D", Period32{}},
		{"PT0H", "P0D", Period32{}},
		{"PT0M", "P0D", Period32{}},
		{"PT0S", "P0D", Period32{}},

		// ones
		{"P1Y", "P1Y", Period32{years: 1, lastField: Year}},
		{"P1M", "P1M", Period32{months: 1, lastField: Month}},
		{"P1W", "P1W", Period32{weeks: 1, lastField: Week}},
		{"P1D", "P1D", Period32{days: 1, lastField: Day}},
		{"PT1H", "PT1H", Period32{hours: 1, lastField: Hour}},
		{"PT1M", "PT1M", Period32{minutes: 1, lastField: Minute}},
		{"PT1S", "PT1S", Period32{seconds: 1, lastField: Second}},
		{"P1Y1M1W1DT1H1M1.1S", "P1Y1M1W1DT1H1M1.1S", Period32{years: 1, months: 1, weeks: 1, days: 1, hours: 1, minutes: 1, seconds: 1, fraction: 100_000_000, lastField: Second}},
		{"P1Y1M1W1DT1H1.1M", "P1Y1M1W1DT1H1.1M", Period32{years: 1, months: 1, weeks: 1, days: 1, hours: 1, minutes: 1, fraction: 100_000_000, lastField: Minute}},
		{"P1Y1M1W1DT1.1H", "P1Y1M1W1DT1.1H", Period32{years: 1, months: 1, weeks: 1, days: 1, hours: 1, fraction: 100_000_000, lastField: Hour}},
		{"P1Y1M1W1.1D", "P1Y1M1W1.1D", Period32{years: 1, months: 1, weeks: 1, days: 1, fraction: 100_000_000, lastField: Day}},
		{"P1Y1M1.1W", "P1Y1M1.1W", Period32{years: 1, months: 1, weeks: 1, fraction: 100_000_000, lastField: Week}},
		{"P1Y1.1M", "P1Y1.1M", Period32{years: 1, months: 1, fraction: 100_000_000, lastField: Month}},

		// smallest
		{"PT0.000000001S", "PT0.000000001S", Period32{lastField: Second, fraction: 1}},
		{"PT0.00000001S", "PT0.00000001S", Period32{lastField: Second, fraction: 10}},
		{"PT0.0000001S", "PT0.0000001S", Period32{lastField: Second, fraction: 100}},
		{"PT0.000001S", "PT0.000001S", Period32{lastField: Second, fraction: 1000}},
		{"PT0.00001S", "PT0.00001S", Period32{lastField: Second, fraction: 10_000}},
		{"PT0.0001S", "PT0.0001S", Period32{lastField: Second, fraction: 100_000}},
		{"PT0.001S", "PT0.001S", Period32{lastField: Second, fraction: 1000_000}},
		{"PT0.01S", "PT0.01S", Period32{lastField: Second, fraction: 10_000_000}},
		{"PT0.1S", "PT0.1S", Period32{lastField: Second, fraction: 100_000_000}},
		{"PT0.1M", "PT0.1M", Period32{lastField: Minute, fraction: 100_000_000}},
		{"PT0.1H", "PT0.1H", Period32{lastField: Hour, fraction: 100_000_000}},
		{"P0.1D", "P0.1D", Period32{lastField: Day, fraction: 100_000_000}},
		{"P0.1W", "P0.1W", Period32{lastField: Week, fraction: 100_000_000}},
		{"P0.1M", "P0.1M", Period32{lastField: Month, fraction: 100_000_000}},
		{"P0.1Y", "P0.1Y", Period32{lastField: Year, fraction: 100_000_000}},

		// fraction overflow
		{"PT0.00000000123S", "PT0.000000001S", Period32{lastField: Second, fraction: 1}},
		{"PT0.00000001234S", "PT0.000000012S", Period32{lastField: Second, fraction: 12}},
		{"PT0.00000012345S", "PT0.000000123S", Period32{lastField: Second, fraction: 123}},
		{"PT0.00000123456S", "PT0.000001234S", Period32{lastField: Second, fraction: 1234}},
		{"PT0.00001234567S", "PT0.000012345S", Period32{lastField: Second, fraction: 12_345}},
		{"PT0.00012345678S", "PT0.000123456S", Period32{lastField: Second, fraction: 123_456}},
		{"PT0.00123456789S", "PT0.001234567S", Period32{lastField: Second, fraction: 1234_567}},
		{"PT0.01234567890S", "PT0.012345678S", Period32{lastField: Second, fraction: 12_345_678}},
		{"PT0.12345678901S", "PT0.123456789S", Period32{lastField: Second, fraction: 123_456_789}},

		// largest
		{"PT2147483647S", "PT2147483647S", Period32{lastField: Second, seconds: 2147483647}},
		{"PT2147483647M", "PT2147483647M", Period32{lastField: Minute, minutes: 2147483647}},
		{"PT2147483647H", "PT2147483647H", Period32{lastField: Hour, hours: 2147483647}},
		{"P2147483647D", "P2147483647D", Period32{lastField: Day, days: 2147483647}},
		{"P2147483647W", "P2147483647W", Period32{lastField: Week, weeks: 2147483647}},
		{"P2147483647M", "P2147483647M", Period32{lastField: Month, months: 2147483647}},
		{"P2147483647Y", "P2147483647Y", Period32{lastField: Year, years: 2147483647}},
	}
	for i, c := range cases {
		p := Period32{}
		err := p.Parse(c.value)
		s := info(i, c.value)
		g.Expect(err).NotTo(HaveOccurred(), s)
		expectValid(t, p, s)
		g.Expect(p).To(Equal(c.period), s)
		// reversal is usually expected to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), s+" reversed")
	}
}
