// Copyright 2016 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package period provides functionality for periods of time using ISO-8601 conventions.
// This deals with years, months, weeks/days, hours, minutes and seconds.
//
// Because of the vagaries of calendar systems, the meaning of year lengths, month lengths
// and even day lengths depends on context. So a period is not necessarily a fixed duration
// of time in terms of seconds.
//
// See https://en.wikipedia.org/wiki/ISO_8601#Durations
//
// Example representations:
//
// * "P2Y" is two years;
//
// * "P6M" is six months;
//
// * "P4D" is four days;
//
// * "P1W" is one week (seven days);
//
// * "PT3H" is three hours.
//
// * "PT20M" is twenty minutes.
//
// * "PT30S" is thirty seconds.
//
// * "-PT30S" or "PT-30S" is minus thirty seconds (implies an "earlier" time).
//
// These can be combined, for example:
//
// * "P3Y6M4W1D" is three years, 6 months, 4 weeks and one day.
//
// * "P2DT12H" is 2 days and 12 hours.
//
// * "P1M-1D" is one month minus one day. Mixed signs may not be widely supported elsewhere.
//
// Also, decimal fractions are supported. To comply with the standard,
// only the last non-zero component is allowed to have a fraction.
// For example
//
// * "P2.5Y" or "P2,5Y" is 2.5 years; both notations are allowed.
//
// * "PT12M7.5S" is 12 minutes and 7.5 seconds.
package period
