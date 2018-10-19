// Package dur is a package to convert string duration and time specifications
// to numbers of seconds and to unix timestamps.
// It aims to support the full specification as outlined in http://graphite.readthedocs.io/en/latest/render_api.html#from-until
// which incorporates the formats defined by `at`.
//
// this package works with the following shorthands:
// Duration : unsigned (positive) number of seconds
// NDuration: like Duration, but non-zero.
package dur

import (
	"errors"
	"fmt"
	"strconv"
)

var errEmpty = errors.New("number cannot be empty")
var errNegative = errors.New("number cannot be negative")
var errNonZero = errors.New("number must be nonzero")
var errUnknownTimeUnit = errors.New("unknown time unit")

// MustParseNDuration parses a format string to a non-zero number of seconds, or panics otherwise
// unit defaults to s if not specified
func MustParseNDuration(desc, s string) uint32 {
	sec, err := ParseNDuration(s)
	if err != nil {
		panic(fmt.Sprintf("%q: %s", desc, s))
	}
	return sec
}

// MustParseDuration parses a format string to a number of seconds, or panics otherwise
// unit defaults to s if not specified
func MustParseDuration(desc, s string) uint32 {
	sec, err := ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("%q: %s", desc, s))
	}
	return sec
}

// ParseNDuration parses a format string to a non-zero number of seconds, or error otherwise
// unit defaults to s if not specified
func ParseNDuration(s string) (uint32, error) {
	i, e := ParseDuration(s)
	if e == nil && i == 0 {
		return 0, errNonZero
	}
	return i, e
}

// ParseDuration parses a format string to a number of seconds, or error otherwise
// valid units are s/sec/secs/second/seconds, m/min/mins/minute/minutes, h/hour/hours, d/day/days, w/week/weeks, mon/month/months, y/year/years
// unit defaults to s if not specified
func ParseDuration(s string) (uint32, error) {
	if s == "" {
		return 0, errEmpty
	}
	if s[0] == '-' {
		return 0, errNegative
	}
	var sum uint32
	for len(s) > 0 {
		var i int
		for i < len(s) && '0' <= s[i] && s[i] <= '9' {
			i++
		}
		var numStr string
		numStr, s = s[:i], s[i:]
		i = 0
		for i < len(s) && (s[i] < '0' || '9' < s[i]) {
			i++
		}
		var unitStr string
		unitStr, s = s[:i], s[i:]

		var units int
		switch unitStr {
		case "", "s", "sec", "secs", "second", "seconds":
			units = 1
		case "m", "min", "mins", "minute", "minutes":
			units = 60
		case "h", "hour", "hours":
			units = 60 * 60
		case "d", "day", "days":
			units = 24 * 60 * 60
		case "w", "week", "weeks":
			units = 7 * 24 * 60 * 60
		case "mon", "month", "months":
			units = 30 * 24 * 60 * 60
		case "y", "year", "years":
			units = 365 * 24 * 60 * 60
		default:
			return 0, errUnknownTimeUnit
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, err
		}
		sum += uint32(num * units)
	}
	return sum, nil
}

var factors = map[string]uint32{
	"y": 60 * 60 * 24 * 365,
	"w": 60 * 60 * 24 * 7,
	"d": 60 * 60 * 24,
	"h": 60 * 60,
	"m": 60,
	"s": 1,
}

// FormatDuration takes a number of seconds and returns a string represention
// that can be parsed by ParseDuration(). Due to its ambiguity, we dont
// use "month" units in the returned value.
func FormatDuration(seconds uint32) string {
	output := ""

	if seconds == 0 {
		return "0s"
	}

	var numUnits uint32
	for seconds > 0 {
		if seconds > factors["y"] {
			numUnits = seconds / factors["y"]
			seconds = seconds - (numUnits * factors["y"])
			output += fmt.Sprintf("%dy", numUnits)
		} else if seconds > factors["w"] {
			numUnits = seconds / factors["w"]
			seconds = seconds - (numUnits * factors["w"])
			output += fmt.Sprintf("%dw", numUnits)
		} else if seconds > factors["d"] {
			numUnits = seconds / factors["d"]
			seconds = seconds - (numUnits * factors["d"])
			output += fmt.Sprintf("%dd", numUnits)
		} else if seconds > factors["h"] {
			numUnits = seconds / factors["h"]
			seconds = seconds - (numUnits * factors["h"])
			output += fmt.Sprintf("%dh", numUnits)
		} else if seconds > factors["m"] {
			numUnits = seconds / factors["m"]
			seconds = seconds - (numUnits * factors["m"])
			output += fmt.Sprintf("%dm", numUnits)
		} else {
			output += fmt.Sprintf("%ds", seconds)
			seconds = 0
		}
	}
	return output
}
