// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

// Period holds a period of time and provides conversion to/from ISO-8601 representations.
// Therefore there are seven fields: years, months, weeks, days, hours, minutes, and seconds.
//
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
type Period string

// String converts the period to ISO-8601 form.
func (p Period) String() string {
	return string(p)
}
