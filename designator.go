// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import "strconv"

type designator int

const (
	Second designator = iota
	Minute
	Hour
	Day
	Week
	Month
	Year
)

func (d designator) Byte() byte {
	switch d {
	case Second:
		return 'S'
	case Minute:
		return 'M'
	case Hour:
		return 'H'
	case Day:
		return 'D'
	case Week:
		return 'W'
	case Month:
		return 'M'
	case Year:
		return 'Y'
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
