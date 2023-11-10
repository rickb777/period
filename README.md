# period

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/period)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/period)](https://goreportcard.com/report/github.com/rickb777/period)
[![Issues](https://img.shields.io/github/issues/rickb777/period.svg)](https://github.com/rickb777/period/issues)

Package `period` has types that represent ISO-8601 periods of time.

The two core types are 

 * `ISOString` - an ISO-8601 string
 * `Period` - a struct with the seven numbers years, months, weeks, days, hours, minutes and seconds.

These two can be converted to the other.

`Period` also allows various calculations to be made. Its fields each hold up to 19 digits precision.

## Status

The basic API exists but may yet change.
