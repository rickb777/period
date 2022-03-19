package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func expectValid(t *testing.T, period Period32, hint interface{}) Period32 {
	t.Helper()
	g := NewGomegaWithT(t)
	info := fmt.Sprintf("%v: invalid: %#v", hint, period)

	// check all the signs are consistent
	nPoz := pos(period.years) + pos(period.months) + pos(period.days) + pos(period.hours) + pos(period.minutes) + pos(period.seconds)
	nNeg := neg(period.years) + neg(period.months) + neg(period.days) + neg(period.hours) + neg(period.minutes) + neg(period.seconds)
	g.Expect(nPoz == 0 || nNeg == 0).To(BeTrue(), info+" inconsistent signs")

	if period.lastField == Year {
		g.Expect(period.months).To(BeZero(), info+" year fraction exists")
		g.Expect(period.weeks).To(BeZero(), info+" year fraction exists")
		g.Expect(period.days).To(BeZero(), info+" year fraction exists")
		g.Expect(period.hours).To(BeZero(), info+" year fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" year fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" year fraction exists")
	}

	if period.lastField == Month {
		g.Expect(period.weeks).To(BeZero(), info+" month fraction exists")
		g.Expect(period.days).To(BeZero(), info+" month fraction exists")
		g.Expect(period.hours).To(BeZero(), info+" month fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" month fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" month fraction exists")
	}

	if period.lastField == Week {
		g.Expect(period.days).To(BeZero(), info+" month fraction exists")
		g.Expect(period.hours).To(BeZero(), info+" month fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" month fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" month fraction exists")
	}

	if period.lastField == Day {
		g.Expect(period.hours).To(BeZero(), info+" day fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" day fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" day fraction exists")
	}

	if period.lastField == Hour {
		g.Expect(period.minutes).To(BeZero(), info+" hour fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" hour fraction exists")
	}

	if period.lastField == Minute {
		g.Expect(period.seconds).To(BeZero(), info+" minute fraction exists")
	}

	return period
}

func pos(i int32) int {
	if i > 0 {
		return 1
	}
	return 0
}

func neg(i int32) int {
	if i < 0 {
		return 1
	}
	return 0
}
