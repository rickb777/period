// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"strconv"
)

type designator int8

const (
	_ designator = iota
	second
	minute
	hour
	day
	week
	month
	year
)

func asDesignator(d byte, isHMS bool) (designator, error) {
	switch d {
	case 'S':
		return second, nil
	case 'H':
		return hour, nil
	case 'D':
		return day, nil
	case 'W':
		return week, nil
	case 'Y':
		return year, nil
	case 'M':
		if isHMS {
			return minute, nil
		}
		return month, nil
	}
	return 0, fmt.Errorf("expected a designator Y, M, W, D, H, or S not '%c'", d)
}

func (d designator) Byte() byte {
	switch d {
	case second:
		return 'S'
	case minute:
		return 'M'
	case hour:
		return 'H'
	case day:
		return 'D'
	case week:
		return 'W'
	case month:
		return 'M'
	case year:
		return 'Y'
	}
	panic(strconv.Itoa(int(d)))
}

//func (d designator) field() string {
//	switch d {
//	case second:
//		return "seconds"
//	case minute:
//		return "minutes"
//	case hour:
//		return "hours"
//	case Day:
//		return "days"
//	case week:
//		return "weeks"
//	case month:
//		return "months"
//	case year:
//		return "years"
//	}
//	panic(strconv.Itoa(int(d)))
//}
//
//func (d designator) min(other designator) designator {
//	if d < other {
//		return d
//	}
//	return other
//}
//
//func (d designator) IsOneOf(xx ...designator) bool {
//	for _, x := range xx {
//		if x == d {
//			return true
//		}
//	}
//	return false
//}
//
//func (d designator) IsNotOneOf(xx ...designator) bool {
//	for _, x := range xx {
//		if x == d {
//			return false
//		}
//	}
//	return true
//}
