package period

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestDecimalString(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := map[string]decimal{
		"0":           {value: 0},
		"10":          {value: 10},
		"100":         {value: 100},
		"1000":        {value: 1, exp: 3},
		"0.1":         {value: 1, exp: -1},
		"0.001":       {value: 1, exp: -3},
		"0.123456789": {value: 123456789, exp: -9},
		"1.23456789":  {value: 123456789, exp: -8},
		"12.3456789":  {value: 123456789, exp: -7},
		"123.456789":  {value: 123456789, exp: -6},
		"1234.56789":  {value: 123456789, exp: -5},
		"12345.6789":  {value: 123456789, exp: -4},
		"123456.789":  {value: 123456789, exp: -3},
		"1234567.89":  {value: 123456789, exp: -2},
		"12345678.9":  {value: 123456789, exp: -1},
		"123456789":   {value: 123456789, exp: 0},
		"1234567890":  {value: 123456789, exp: 01},
	}

	for s, d := range cases {
		result := d.String()
		g.Expect(result).To(Equal(s))
	}
}

func TestDecimalAdd(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct{ a, b, sum decimal }{
		{a: decimal{value: 0}, b: decimal{value: 0}, sum: decimal{value: 0}},
		{a: decimal{value: 1}, b: decimal{value: 999}, sum: decimal{value: 1, exp: 3}},
		{a: decimal{value: 1}, b: decimal{value: -1}, sum: decimal{value: 0}},
		{a: decimal{value: 1234, exp: -1}, b: decimal{value: 123, exp: -3}, sum: decimal{value: 123523, exp: -3}},
		{a: decimal{value: 12345, exp: -2}, b: decimal{value: 123, exp: -3}, sum: decimal{value: 123573, exp: -3}},
		{a: decimal{value: 12345, exp: -2}, b: decimal{value: -123, exp: -3}, sum: decimal{value: 123327, exp: -3}},
		{a: newFromInt(12340), b: decimal{value: 123, exp: -3}, sum: decimal{value: 12340123, exp: -3}},
	}

	for _, c := range cases {
		r1 := c.a.Add(c.b)
		g.Expect(r1).To(Equal(c.sum), "a+b")

		r2 := c.b.Add(c.a)
		g.Expect(r2).To(Equal(c.sum), "b+a")
	}
}
