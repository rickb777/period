// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	. "github.com/onsi/gomega"
	bigdecimal "github.com/shopspring/decimal"
	"strings"
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

		"P0D": Period64{},

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

		//"P3Y":   {years: 3, lastField: Year},
		//"P6M":   {months: 6, lastField: Month},
		//"P5W":   {weeks: 5, lastField: Week},
		//"P4D":   {days: 4, lastField: Day},
		//"PT12H": {hours: 12, lastField: Hour},
		//"PT30M": {minutes: 30, lastField: Minute},
		//"PT5S":  {seconds: 5, lastField: Second},

		"P3.9Y":              {years: decS("3.9")},
		"P3Y6.9M":            {years: decI(3), months: decS("6.9")},
		"P3Y6M2.9W":          {years: decI(3), months: decI(6), weeks: decS("2.9")},
		"P3Y6M2W4.9D":        {years: decI(3), months: decI(6), weeks: decI(2), days: decS("4.9")},
		"P3Y6M2W4DT1.9H":     {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decS("1.9")},
		"P3Y6M2W4DT1H5.9M":   {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decI(1), minutes: decS("5.9")},
		"P3Y6M2W4DT1H5M7.9S": {years: decI(3), months: decI(6), weeks: decI(2), days: decI(4), hours: decI(1), minutes: decI(5), seconds: decS("7.9")},
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
			if sn != ne {
				t.Errorf("-ve got %s, expected %s", sn, ne)
			}
		}
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
	z := Period64{}
	neg := Period64{years: decI(1), months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: true}
	pos := Period64{years: decI(1), months: decI(2), weeks: decI(3), days: decI(4), hours: decI(5), minutes: decI(6), seconds: decI(7), neg: false}

	a := neg.Abs()
	if a != pos {
		t.Errorf("Abs() failed %+v", a)
	}

	if neg.Sign() != -1 {
		t.Errorf("Sign() -1 failed")
	}

	if pos.Sign() != 1 {
		t.Errorf("Sign() 1 failed")
	}

	if z.Sign() != 0 {
		t.Errorf("Sign() 0 failed")
	}

	if !z.IsZero() {
		t.Errorf("IsZero() 0 failed")
	}

	if pos.IsZero() {
		t.Errorf("IsZero() 1 failed")
	}

	if neg.IsZero() {
		t.Errorf("IsZero() -1 failed")
	}

	if !pos.IsPositive() {
		t.Errorf("+ve IsPositive() failed")
	}

	if pos.IsNegative() {
		t.Errorf("+ve IsNegative() failed")
	}

	if !neg.IsNegative() {
		t.Errorf("-ve IsNegative() failed")
	}

	if neg.IsPositive() {
		t.Errorf("-ve IsPositive() failed")
	}
}

//func Test_Period_IsValid_false(t *testing.T) {
//	if !(Period64{}).isValid() {
//		t.Errorf("expected valid for P0D")
//	}
//
//	cases := []Period64{
//		{years: -1},
//		{months: -1},
//		{weeks: -1},
//		{days: -1},
//		{hours: -1},
//		{minutes: -1},
//		{seconds: -1},
//
//		{months: 1, lastField: Year},
//		{weeks: 1, lastField: Year},
//		{days: 1, lastField: Year},
//		{hours: 1, lastField: Year},
//		{minutes: 1, lastField: Year},
//		{seconds: 1, lastField: Year},
//
//		{weeks: 1, lastField: Month},
//		{days: 1, lastField: Month},
//		{hours: 1, lastField: Month},
//		{minutes: 1, lastField: Month},
//		{seconds: 1, lastField: Month},
//
//		{days: 1, lastField: Week},
//		{hours: 1, lastField: Week},
//		{minutes: 1, lastField: Week},
//		{seconds: 1, lastField: Week},
//
//		{hours: 1, lastField: Day},
//		{minutes: 1, lastField: Day},
//		{seconds: 1, lastField: Day},
//
//		{minutes: 1, lastField: Hour},
//		{seconds: 1, lastField: Hour},
//
//		{seconds: 1, lastField: Minute},
//	}
//
//	for _, p64 := range cases {
//		if p64.isValid() {
//			t.Errorf("expected invalid for %+v", p64)
//		}
//	}
//}

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

func nospace(s string) string {
	b := new(strings.Builder)
	for _, r := range s {
		if r != ' ' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
