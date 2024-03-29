// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"database/sql/driver"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPeriodScan(t *testing.T) {
	g := NewGomegaWithT(t)

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
			g.Expect(e).NotTo(HaveOccurred())
			g.Expect(*r).To(Equal(c.expected))

			var d driver.Valuer = *r

			q, e := d.Value()
			g.Expect(e).NotTo(HaveOccurred())
			g.Expect(q.(string)).To(Equal(c.expected.String()))
		})
	}
}

func TestPeriodScan_nil_value(t *testing.T) {
	g := NewGomegaWithT(t)
	r := new(Period)
	e := r.Scan(nil)
	g.Expect(e).NotTo(HaveOccurred())
}

func TestPeriodScan_problem_type(t *testing.T) {
	g := NewGomegaWithT(t)
	r := new(Period)
	e := r.Scan(1)
	g.Expect(e).To(HaveOccurred())
	g.Expect(e.Error()).To(ContainSubstring("not a meaningful period"))
}
