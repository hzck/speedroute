package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// StoMS is the factor which needs to be multiplies when converting seconds to milliseconds (1000).
const StoMS = 1000

// MtoMS is the factor which needs to be multiplies when converting minutes to milliseconds (60*1000).
const MtoMS = 60 * StoMS

// HtoMS is the factor which needs to be multiplies when converting hours to milliseconds (60*60*1000).
const HtoMS = 60 * MtoMS

func parseTime(time string) (int, error) {
	times := strings.Split(time, ":")
	var h, m, s string
	switch len(times) {
	case 1:
		s = times[0]
	case 2:
		m = times[0]
		s = times[1]
	case 3:
		h = times[0]
		m = times[1]
		s = times[2]
	default:
		return -1, errors.New("Can't parse time " + time)
	}
	result := 0
	if len(h) > 0 {
		hours, err := parseInt(h)
		if err != nil {
			return -1, err
		}
		result += hours * HtoMS
	}
	if len(m) > 0 {
		mins, err := parseInt(m)
		if err != nil {
			return -1, err
		}
		result += mins * MtoMS
	}
	if len(s) > 0 {
		if strings.Contains(s, ".") || len(m) > 0 {
			secs, err := parseSeconds(s)
			if err != nil {
				return -1, err
			}
			result += secs
		} else {
			ms, err := parseInt(s)
			if err != nil {
				return -1, err
			}
			result += ms
		}
	}

	return result, nil
}

func parseInt(time string) (int, error) {
	result, err := strconv.ParseInt(time, 10, 0)
	if err != nil {
		return -1, err
	}
	return int(result), nil
}

func parseSeconds(time string) (int, error) {
	if time == "." {
		return 0, nil
	}
	result, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return -1, err
	}
	return int(round(result * StoMS)), nil
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
