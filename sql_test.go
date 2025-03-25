// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"database/sql/driver"
	"fmt"
	"github.com/rickb777/expect"
	"testing"
)

func TestPeriodScan(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected Period
	}{
		{[]byte("P1Y3M"), MustParse("P1Y3M")},
		{"P1Y3M", MustParse("P1Y3M")},
		{[]byte("P48M"), MustParse("P48M")},
		{"P48M", MustParse("P48M")},
		{"P1YT-1S", MustParse("P1YT-1S")},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			r := new(Period)
			e := r.Scan(c.v)
			expect.Error(e).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			expect.Any(*r).ToBe(t, c.expected)

			var d driver.Valuer = *r

			q, e := d.Value()
			expect.Error(e).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			expect.String(q.(string)).ToBe(t, c.expected.String())
		})
	}
}

func TestPeriodScan_nil_value(t *testing.T) {
	r := new(Period)
	e := r.Scan(nil)
	expect.Error(e).Not().ToHaveOccurred(t)
}

func TestPeriodScan_problem_type(t *testing.T) {
	r := new(Period)
	e := r.Scan(1)
	expect.Error(e).ToContain(t, "not a meaningful period")
}
