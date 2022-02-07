// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	//. "github.com/onsi/gomega"
	"testing"
)

func Test_String(t *testing.T) {
	//g := NewGomegaWithT(t)

	cases := map[Period]Period64{
		// note: the negative cases are also covered (see below)

		"P0D": Period64{},

		// ones
		"P1Y":  {years: 1, lastField: Year},
		"P1M":  {months: 1, lastField: Month},
		"P1W":  {weeks: 1, lastField: Week},
		"P1D":  {days: 1, lastField: Day},
		"PT1H": {hours: 1, lastField: Hour},
		"PT1M": {minutes: 1, lastField: Minute},
		"PT1S": {seconds: 1, lastField: Second},

		// smallest fraction
		"P0.000000001Y":  {fraction: 1, lastField: Year},
		"P0.000000001M":  {fraction: 1, lastField: Month},
		"P0.000000001W":  {fraction: 1, lastField: Week},
		"P0.000000001D":  {fraction: 1, lastField: Day},
		"PT0.000000001H": {fraction: 1, lastField: Hour},
		"PT0.000000001M": {fraction: 1, lastField: Minute},
		"PT0.000000001S": {fraction: 1, lastField: Second},

		// 1 + smallest
		"P1.000000001Y":  {years: 1, fraction: 1, lastField: Year},
		"P1.000000001M":  {months: 1, fraction: 1, lastField: Month},
		"P1.000000001W":  {weeks: 1, fraction: 1, lastField: Week},
		"P1.000000001D":  {days: 1, fraction: 1, lastField: Day},
		"PT1.000000001H": {hours: 1, fraction: 1, lastField: Hour},
		"PT1.000000001M": {minutes: 1, fraction: 1, lastField: Minute},
		"PT1.000000001S": {seconds: 1, fraction: 1, lastField: Second},

		// other fractions
		"P0.00000001Y": {fraction: 10, lastField: Year},
		"P0.0000001Y":  {fraction: 100, lastField: Year},
		"P0.000001Y":   {fraction: 1000, lastField: Year},
		"P0.00001Y":    {fraction: 10000, lastField: Year},
		"P0.0001Y":     {fraction: 100000, lastField: Year},
		"P0.001Y":      {fraction: 1000000, lastField: Year},
		"P0.01Y":       {fraction: 10000000, lastField: Year},
		"P0.1Y":        {fraction: 100000000, lastField: Year},

		"P3Y":   {years: 3, lastField: Year},
		"P6M":   {months: 6, lastField: Month},
		"P5W":   {weeks: 5, lastField: Week},
		"P4D":   {days: 4, lastField: Day},
		"PT12H": {hours: 12, lastField: Hour},
		"PT30M": {minutes: 30, lastField: Minute},
		"PT5S":  {seconds: 5, lastField: Second},

		"P3.9Y": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Year, fraction: 900000000},
		"P3Y6.9M": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Month, fraction: 900000000},
		"P3Y6M2.9W": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Week, fraction: 900000000},
		"P3Y6M2W39.9D": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Day, fraction: 900000000},
		"P3Y6M2W39DT1.9H": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Hour, fraction: 900000000},
		"P3Y6M2W39DT1H2.9M": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Minute, fraction: 900000000},
		"P3Y6M2W39DT1H2M4.9S": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Second, fraction: 900000000},
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

func Test_Period64_Sign_Abs_etc(t *testing.T) {
	z := Period64{}
	neg := Period64{years: 1, months: 2, weeks: 3, days: 4, hours: 5, minutes: 6, seconds: 7, fraction: 8, neg: true}
	pos := Period64{years: 1, months: 2, weeks: 3, days: 4, hours: 5, minutes: 6, seconds: 7, fraction: 8, neg: false}

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

func Test_Period64_IsValid_false(t *testing.T) {
	//g := NewGomegaWithT(t)

	cases := []Period64{
		{years: -1},
		{months: -1},
		{weeks: -1},
		{days: -1},
		{hours: -1},
		{minutes: -1},
		{seconds: -1},

		{months: 1, lastField: Year},
		{weeks: 1, lastField: Year},
		{days: 1, lastField: Year},
		{hours: 1, lastField: Year},
		{minutes: 1, lastField: Year},
		{seconds: 1, lastField: Year},

		{weeks: 1, lastField: Month},
		{days: 1, lastField: Month},
		{hours: 1, lastField: Month},
		{minutes: 1, lastField: Month},
		{seconds: 1, lastField: Month},

		{days: 1, lastField: Week},
		{hours: 1, lastField: Week},
		{minutes: 1, lastField: Week},
		{seconds: 1, lastField: Week},

		{hours: 1, lastField: Day},
		{minutes: 1, lastField: Day},
		{seconds: 1, lastField: Day},

		{minutes: 1, lastField: Hour},
		{seconds: 1, lastField: Hour},

		{seconds: 1, lastField: Minute},
	}

	for _, p64 := range cases {
		if p64.isValid() {
			t.Errorf("expected invalid for %+v", p64)
		}
	}
}
