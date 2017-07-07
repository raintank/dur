package dur

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

//var absoluteTimeFormats = []string{"15:04 20060102", "20060102", "01/02/06"}

var absoluteTimeFormats = []string{"15:04 20060102", "20060102", "01/02/06", "01/02/2006"}
var errUnknownDateTimeFormat = errors.New("parse error. unknown DateTime format")
var errUnknownTimeFormat = errors.New("parse error. unknown Time format")

// ParseDateTime parses a format string to a unix timestamp, or error otherwise.
// 'loc' is the timezone to use for interpretation (when applicable)
// 'now' is a reference, in case a relative specification is given.
// 'def' is a default in case an empty specification is given.
func ParseDateTime(s string, loc *time.Location, now time.Time, def uint32) (uint32, error) {
	now = now.In(loc)
	switch s {
	case "":
		return uint32(def), nil
	case "now":
		return uint32(now.Unix()), nil
	case "today":
		year, month, day := now.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "midnight":
		year, month, day := now.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "noon":
		year, month, day := now.Date()
		return uint32(time.Date(year, month, day, 12, 0, 0, 0, loc).Unix()), nil
	case "teatime":
		year, month, day := now.Date()
		return uint32(time.Date(year, month, day, 16, 0, 0, 0, loc).Unix()), nil
	case "yesterday":
		base := now.AddDate(0, 0, -1)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "tomorrow":
		base := now.AddDate(0, 0, 1)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	}

	// try negative relative duration offset
	if s[0] == '-' {
		dur, err := ParseNDuration(s[1:])
		if err == nil {
			return uint32(now.Add(-time.Duration(dur) * time.Second).Unix()), nil
		}
		return 0, err
	}

	// not documented, though supported by graphite so we must do the same.
	// https://github.com/raintank/metrictank/issues/673
	if strings.HasPrefix(s, "now-") {
		dur, err := ParseNDuration(s[4:])
		if err == nil {
			return uint32(now.Add(-time.Duration(dur) * time.Second).Unix()), nil
		}
		return 0, err
	}
	if strings.HasPrefix(s, "now+") {
		dur, err := ParseNDuration(s[4:])
		if err == nil {
			return uint32(now.Add(time.Duration(dur) * time.Second).Unix()), nil
		}
		return 0, err
	}

	// if it's a plain integer, interpret it as a unix timestamp
	// except if it's a series of numbers that looks like YYYYMMDD,
	// which will be processed further down.
	// if it's not a plain integer, we proceed trying more things
	if len(s) != 8 {
		i, err := strconv.Atoi(s)
		if err == nil {
			return uint32(i), nil
		}
	}

	// try positive relative duration offset like 5s or 5h30min
	// (we already covered negative offsets higher up)
	// since this also accepts a plain number, just like above,
	// we should only do this if the input is not 8 chars
	if len(s) != 8 {
		dur, err := ParseNDuration(s)
		if err == nil {
			return uint32(now.Add(-time.Duration(dur) * time.Second).Unix()), nil
		}
	}

	// try remaining absolute formats.
	// they are either in the following shape: [<time> ]<date>
	// where:
	// * time can be like noon, midnight, teatime, <int>am, <int>AM, <int>pm, <int>PM, or HH:MM
	// * date can be like YYYYMMDD, MM/DD/YY or MM/DD/YYYY
	// or: <monthname> <num>, which we'll try first

	// Go can't parse _ in date strings. this is for HH:MM_YYMMDD
	s = strings.Replace(s, "_", " ", 1)

	base, err := time.ParseInLocation("January 2", s, loc)
	if err == nil {
		y, _, _ := now.Date()
		_, m, d := base.Date()
		return uint32(time.Date(y, m, d, 0, 0, 0, 0, loc).Unix()), nil
	}

	var ts, ds string
	split := strings.Fields(s)

	switch {
	case len(split) == 1:
		ds = s
	case len(split) == 2:
		ts, ds = split[0], split[1]
	case len(split) > 2:
		return 0, errUnknownDateTimeFormat
	}

	// first we need to set our "now" to the right date.

dateStringSwitch:
	switch ds {
	case "today":
		base = now
	case "yesterday":
		base = now.AddDate(0, 0, -1)
	case "tomorrow":
		base = now.AddDate(0, 0, 1)
	case "monday":
		base = RewindToWeekday(now, time.Monday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "tuesday":
		base = RewindToWeekday(now, time.Tuesday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "wednesday":
		base = RewindToWeekday(now, time.Wednesday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "thursday":
		base = RewindToWeekday(now, time.Thursday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "friday":
		base = RewindToWeekday(now, time.Friday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "saturday":
		base = RewindToWeekday(now, time.Saturday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	case "sunday":
		base = RewindToWeekday(now, time.Sunday)
		year, month, day := base.Date()
		return uint32(time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()), nil
	default:
		for _, format := range absoluteTimeFormats {
			base, err = time.ParseInLocation(format, ds, loc)
			if err == nil {
				break dateStringSwitch
			}
		}
		return 0, errUnknownDateTimeFormat
	}
	if ts == "" {
		return uint32(base.Unix()), nil
	}

	hour, minute, err := ParseTime(ts)
	if err != nil {
		return 0, err
	}

	year, month, day := base.Date()
	return uint32(time.Date(year, month, day, hour, minute, 0, 0, loc).Unix()), nil
}

// ParseTime parses a time and returns hours and minutes
func ParseTime(s string) (hour, minute int, err error) {
	switch s {
	case "midnight":
		return 0, 0, nil
	case "noon":
		return 12, 0, nil
	case "teatime":
		return 16, 0, nil
	}
	if strings.HasSuffix(s, "am") || strings.HasSuffix(s, "AM") {
		hour, err := strconv.Atoi(s[:len(s)-2])
		return hour, 0, err
	}
	if strings.HasSuffix(s, "pm") || strings.HasSuffix(s, "PM") {
		hour, err := strconv.Atoi(s[:len(s)-2])
		return hour + 12, 0, err
	}

	parts := strings.Split(s, ":")

	if len(parts) != 2 {
		return 0, 0, errUnknownTimeFormat
	}

	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errUnknownTimeFormat
	}

	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errUnknownTimeFormat
	}

	return hour, minute, nil
}

// RewindToWeekday moves a datetime back to the last occurence of the given weekday (potentially that day without needing to seek back)
// while retaining hour/minute/second values.
func RewindToWeekday(t time.Time, day time.Weekday) time.Time {
	for t.Weekday() != day {
		t = t.AddDate(0, 0, -1)
	}
	return t
}
