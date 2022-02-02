// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import "strconv"

type designator int

const (
	_ designator = iota
	Year
	Month
	Week
	Day
	Hour
	Minute
	Second
)

func (d designator) Byte() byte {
	switch d {
	case Year:
		return 'Y'
	case Month:
		return 'M'
	case Week:
		return 'W'
	case Day:
		return 'D'
	case Hour:
		return 'H'
	case Minute:
		return 'M'
	case Second:
		return 'S'
	}
	panic(strconv.Itoa(int(d)))
}

func (d designator) IsOneOf(xx ...designator) bool {
	for _, x := range xx {
		if x == d {
			return true
		}
	}
	return false
}

func (d designator) IsNotOneOf(xx ...designator) bool {
	for _, x := range xx {
		if x == d {
			return false
		}
	}
	return true
}
