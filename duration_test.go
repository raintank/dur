package dur

import "testing"

func TestDuration(t *testing.T) {
	var cases = []struct {
		in  string
		out uint32
		err bool
	}{
		{"", 0, true},
		{"0", 0, false},
		{"-1", 0, true},
		{"1", 1, false},
		{"3600", 3600, false},
		{"1000000000", 1000000000, false},
		{"1us", 0, true},
		{"1ms", 0, true},
		{"1000ms", 0, true},
		{"1m", 60, false},
		{"1min", 60, false},
		{"1h", 3600, false},
		{"1s", 1, false},
		{"2d", 2 * 60 * 60 * 24, false},
		{"10hours", 60 * 60 * 10, false},
		{"7d13h45min21s", 7*24*60*60 + 13*60*60 + 45*60 + 21, false},
		{"01hours", 60 * 60 * 1, false},
		{"2d2d", 4 * 60 * 60 * 24, false},
	}

	for i, c := range cases {
		d, err := ParseDuration(c.in)
		if (err != nil) != c.err {
			t.Fatalf("case %d %q: expected err %t, got err %s", i, c.in, c.err, err)
		}
		if d != c.out {
			t.Fatalf("case %d %q: expected %d, got %d", i, c.in, c.out, d)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	var cases = []struct {
		out string
		in  uint32
	}{
		{"0s", 0},
		{"1s", 1},
		{"33s", 33},
		{"14m", 60 * 14},
		{"23h", 23 * 60 * 60},
		{"4d", 4 * 24 * 60 * 60},
		{"8w", 8 * 7 * 24 * 60 * 60},
		{"52w", 52 * 7 * 24 * 60 * 60},
		{"2y", 2 * 365 * 24 * 60 * 60},
		{"2y51w6d23h59m59s", 2*365*24*60*60 + 51*7*24*60*60 + 6*24*60*60 + 23*60*60 + 59*60 + 59},
		{"1y1w1d1h1m1s", 1*365*24*60*60 + 1*7*24*60*60 + 1*24*60*60 + 1*60*60 + 1*60 + 1},
		{"1y1d1m1s", 1*365*24*60*60 + 1*24*60*60 + 1*60 + 1},
		{"4h30m", 4*60*60 + 30*60},
	}

	for i, c := range cases {
		s := FormatDuration(c.in)
		if s != c.out {
			t.Fatalf("case %d %q: expected %s, got %s", i, c.in, c.out, s)
		}
	}
}
