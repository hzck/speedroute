package parser

import (
	"testing"
)

// TestParseTime tests all supported time formats into ms.
func TestParseTime(t *testing.T) {
	var data = []struct {
		in  string
		out int
	}{
		{"", 0},                  // 0 value
		{".", 0},                 // 0 value
		{"0", 0},                 // 0 value
		{"55", 55},               // ms
		{"555", 555},             // ms
		{"5555", 5555},           // ms
		{"05555", 5555},          // ms 0 prefix
		{".9994", 999},           // ms rounding down
		{".9995", 1000},          // ms rounding up
		{"1.", 1000},             // secs no ms number
		{"25.0", 25000},          // secs
		{"075.0", 75000},         // secs 0 prefix
		{"100.0", 100000},        // secs 3 numbers
		{"10.5", 10500},          // secs & ms
		{"1:00", 60000},          // mins
		{"02:00", 120000},        // mins 0 prefix
		{"100:10", 6010000},      // mins & secs 3 numbers
		{"100:03.501", 6003501},  // mins, secs & ms
		{"1:00:00", 3600000},     // hours
		{"02:00:00", 7200000},    // hours 0 prefix
		{"100:00:00", 360000000}, // hours 3 numbers
		{"1:02:05.333", 3725333}, // hours, mins, secs & ms
	}
	for _, d := range data {
		parsed, err := parseTime(d.in)
		if err != nil {
			t.Errorf("The parsed time for " + d.in + " is error: " + err.Error())
		}
		if parsed != d.out {
			t.Errorf("The parsed time for "+d.in+" is incorrect: %v != %v", parsed, d.out)
		}
	}
}
