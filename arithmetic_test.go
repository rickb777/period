// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPeriodAddSubtract(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one, two        Period
		sum, difference Period
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
		//{"PT16M40S", "-PT1017S", "PT43S"},
		//{"PT17M40S", "-PT1017S", "-PT17S"},
		//{"P3Y3M3W3D", "-P1Y4M", "P1Y11M3W3D"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s %s", i, c.one, c.two), func(t *testing.T) {
			a := MustParse(c.one)
			b := MustParse(c.two)

			s := a.Add(b)
			g.Expect(s).To(Equal(MustParse(c.sum)), info(i, "%s + %s = %s", c.one, c.two, s))

			d := a.Subtract(b)
			g.Expect(d).To(Equal(MustParse(c.difference)), info(i, "%s - %s = %s", c.one, c.two, d))
		})
	}
}

//func expectValid(t *testing.T, period Period32, hint interface{}) Period32 {
//	t.Helper()
//	g := NewGomegaWithT(t)
//	info := fmt.Sprintf("%v: invalid: %#v", hint, period)
//
//	// check all the signs are consistent
//	nPoz := pos(period.years) + pos(period.months) + pos(period.days) + pos(period.hours) + pos(period.minutes) + pos(period.seconds)
//	nNeg := neg(period.years) + neg(period.months) + neg(period.days) + neg(period.hours) + neg(period.minutes) + neg(period.seconds)
//	g.Expect(nPoz == 0 || nNeg == 0).To(BeTrue(), info+" inconsistent signs")
//
//	if period.lastField == Year {
//		g.Expect(period.months).To(BeZero(), info+" year fraction exists")
//		g.Expect(period.weeks).To(BeZero(), info+" year fraction exists")
//		g.Expect(period.days).To(BeZero(), info+" year fraction exists")
//		g.Expect(period.hours).To(BeZero(), info+" year fraction exists")
//		g.Expect(period.minutes).To(BeZero(), info+" year fraction exists")
//		g.Expect(period.seconds).To(BeZero(), info+" year fraction exists")
//	}
//
//	if period.lastField == Month {
//		g.Expect(period.weeks).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.days).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.hours).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.minutes).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.seconds).To(BeZero(), info+" month fraction exists")
//	}
//
//	if period.lastField == Week {
//		g.Expect(period.days).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.hours).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.minutes).To(BeZero(), info+" month fraction exists")
//		g.Expect(period.seconds).To(BeZero(), info+" month fraction exists")
//	}
//
//	if period.lastField == Day {
//		g.Expect(period.hours).To(BeZero(), info+" day fraction exists")
//		g.Expect(period.minutes).To(BeZero(), info+" day fraction exists")
//		g.Expect(period.seconds).To(BeZero(), info+" day fraction exists")
//	}
//
//	if period.lastField == Hour {
//		g.Expect(period.minutes).To(BeZero(), info+" hour fraction exists")
//		g.Expect(period.seconds).To(BeZero(), info+" hour fraction exists")
//	}
//
//	if period.lastField == Minute {
//		g.Expect(period.seconds).To(BeZero(), info+" minute fraction exists")
//	}
//
//	return period
//}
//
//func pos(i int32) int {
//	if i > 0 {
//		return 1
//	}
//	return 0
//}
//
//func neg(i int32) int {
//	if i < 0 {
//		return 1
//	}
//	return 0
//}
