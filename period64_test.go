// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	//. "github.com/onsi/gomega"
	"testing"
)

func TestPeriodString(t *testing.T) {
	//g := NewGomegaWithT(t)

	cases := map[string]Period64{
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
		"P3Y6M2W39DT1H2M4.9S": {years: 3, months: 6, weeks: 2, days: 39, hours: 1, minutes: 2, seconds: 4,
			lastField: Second, fraction: 900000000},

		//{"P2.5Y", Period64{years: 25}},
		//{"P2.5M", Period64{months: 25}},
		//{"P2.5D", Period64{days: 25}},
		//{"PT2.5H", Period64{hours: 25}},
		//{"PT2.5M", Period64{minutes: 25}},
		//{"PT2.5S", Period64{seconds: 25}},
	}
	for expected, p64 := range cases {
		sp := p64.String()
		if sp != expected {
			t.Errorf("+ve got %s, expected %s", sp, expected)
		}

		if !p64.IsZero() {
			sn := p64.Negate().String()
			ne := "-" + expected
			if sn != ne {
				t.Errorf("-ve got %s, expected %s", sn, ne)
			}
		}
	}
}
