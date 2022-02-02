package period

// Period holds a period of time and provides conversion to/from ISO-8601 representations.
// Therefore there are six fields: years, months, days, hours, minutes, and seconds.
//
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
type Period string

func (p Period) String() string {
	return string(p)
}
