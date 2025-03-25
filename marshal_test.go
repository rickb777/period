// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/rickb777/expect"
	"testing"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []string{
		"P0D",
		"P1D",
		"P1W",
		"P1M",
		"P1Y",
		"PT1H",
		"PT1M",
		"PT1S",
		"P2Y3M4W5D",
		"-P2Y3M4W5D",
		"P2Y3M4W5DT-1H7M9S",
		"-P2Y3M4W5DT1H7M0.9S",
		"P48M",
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c), func(t *testing.T) {
			period := MustParse(c)
			var p Period
			err := encoder.Encode(&period)
			expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			if err == nil {
				err = decoder.Decode(&p)
				expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
				expect.Any(p).Info("%d %v", i, c).ToBe(t, period)
			}
		})
	}
}

func TestISOStringJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -5, -3, -11, -59, -59), `"-P1111Y4M5W3DT11H59M59S"`},
		{New(-1, -10, -5, -31, -5, -4, -20), `"-P1Y10M5W31DT5H4M20S"`},
		{New(0, 0, 0, 0, 0, 0, 0), `"P0D"`},
		{New(0, 0, 0, 0, 0, 0, 1), `"PT1S"`},
		{New(0, 0, 0, 0, 0, 1, 0), `"PT1M"`},
		{New(0, 0, 0, 0, 1, 0, 0), `"PT1H"`},
		{New(0, 0, 0, 1, 0, 0, 0), `"P1D"`},
		{New(0, 0, 1, 0, 0, 0, 0), `"P1W"`},
		{New(0, 1, 0, 0, 0, 0, 0), `"P1M"`},
		{New(1, 0, 0, 0, 0, 0, 0), `"P1Y"`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			bb, err := json.Marshal(c.value.Period())
			expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			expect.String(bb).Info("%d %v", i, c).ToEqual(t, c.want)
		})
	}
}

func TestPeriodJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -5, -3, -11, -59, -59), `"-P1111Y4M5W3DT11H59M59S"`},
		{New(-1, -10, -5, -31, -5, -4, -20), `"-P1Y10M5W31DT5H4M20S"`},
		{New(0, 0, 0, 0, 0, 0, 0), `"P0D"`},
		{New(0, 0, 0, 0, 0, 0, 1), `"PT1S"`},
		{New(0, 0, 0, 0, 0, 1, 0), `"PT1M"`},
		{New(0, 0, 0, 0, 1, 0, 0), `"PT1H"`},
		{New(0, 0, 0, 1, 0, 0, 0), `"P1D"`},
		{New(0, 0, 1, 0, 0, 0, 0), `"P1W"`},
		{New(0, 1, 0, 0, 0, 0, 0), `"P1M"`},
		{New(1, 0, 0, 0, 0, 0, 0), `"P1Y"`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			var p Period
			bb, err := json.Marshal(c.value)
			expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			expect.String(bb).Info("%d %v", i, c).ToEqual(t, c.want)
			if string(bb) == c.want {
				err = json.Unmarshal(bb, &p)
				expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
				expect.Any(p).Info("%d %v", i, c).ToBe(t, c.value)
			}
		})
	}
}

func TestPeriodTextMarshalling(t *testing.T) {
	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -5, -3, -11, -59, -59), "-P1111Y4M5W3DT11H59M59S"},
		{New(-1, -9, -5, -31, -5, -4, -20), "-P1Y9M5W31DT5H4M20S"},
		{New(0, 0, 0, 0, 0, 0, 0), "P0D"},
		{New(0, 0, 0, 0, 0, 0, 1), "PT1S"},
		{New(0, 0, 0, 0, 0, 1, 0), "PT1M"},
		{New(0, 0, 0, 0, 1, 0, 0), "PT1H"},
		{New(0, 0, 0, 1, 0, 0, 0), "P1D"},
		{New(0, 0, 1, 0, 0, 0, 0), "P1W"},
		{New(0, 1, 0, 0, 0, 0, 0), "P1M"},
		{New(1, 0, 0, 0, 0, 0, 0), "P1Y"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			var p Period
			bb, err := c.value.MarshalText()
			expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
			expect.String(bb).Info("%d %v", i, c).ToEqual(t, c.want)
			if string(bb) == c.want {
				err = p.UnmarshalText(bb)
				expect.Error(err).Info("%d %v", i, c).Not().ToHaveOccurred(t)
				expect.Any(p).Info("%d %v", i, c).ToBe(t, c.value)
			}
		})
	}
}

func TestInvalidPeriodText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{``, `cannot parse a blank string as a period`},
		{`not-a-period`, `not-a-period: expected 'P' period mark at the start`},
		{`P000`, `P000: missing designator at the end`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			var p Period
			err := p.UnmarshalText([]byte(c.value))
			expect.Error(err).Info("%d %v", i, c).ToHaveOccurred(t)
			expect.Error(err).Info("%d %v", i, c).ToContain(t, c.want)
		})
	}
}
