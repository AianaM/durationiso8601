package durationiso8601

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(t time.Time, s string) (time.Duration, error) {
	re := regexp.MustCompile(`(?mi)^(?P<sign>\+|-)?P((?P<year>\d+(\.\d+)?)Y)?((?P<month>\d+(\.\d+)?)M)?((?P<week>\d+(\.\d+)?)W)?((?P<day>\d+(\.\d+)?)D)?(T((?P<hour>\d+(\.\d+)?)H)?((?P<minute>\d+(\.\d+)?)M)?((?P<second>\d+(\.\d+)?)S)?)?$`)
	matches := re.FindStringSubmatch(s)
	matchesLen := len(matches)
	if matchesLen == 0 {
		return time.Duration(0), fmt.Errorf("invalid duration format: %s", s)
	}

	periodComponents := struct {
		year  int
		month int
		week  int
		day   int
	}{}
	timeComponents := struct {
		hour   int
		minute int
		second int
	}{}
	// An optional sign character ( + or - ), which represents positive or negative duration. Default is positive
	positive := true

	for i, name := range re.SubexpNames() {
		if name == "sign" && matches[i] == "-" {
			positive = false
			continue
		}
		if name == "" || name == "sign" || matches[i] == "" {
			continue
		}
		if strings.Contains(matches[i], ".") {
			if num, err := strconv.ParseFloat(matches[i], 32); err != nil {
				return 0, fmt.Errorf("invalid duration format: %s, %s: %w", s, name, err)
			} else {
				switch name {
				case "year":
					periodComponents.month = int(num * 12) // Convert years to months
				case "month":
					return 0, fmt.Errorf("fractional months are not supported: %s, %s: %w", s, name, err)
				case "week":
					periodComponents.day += int(num * 7) // Convert weeks to days
				case "day":
					timeComponents.second = int(num * 24 * 60 * 60) // Convert days to seconds
				case "hour":
					timeComponents.second = int(num * 60 * 60) // Convert hours to seconds
				case "minute":
					timeComponents.second = int(num * 60) // Convert minutes to seconds
				case "second":
					timeComponents.second = int(num) // Seconds remain as is
				}
			}
		} else if num, err := strconv.Atoi(matches[i]); err == nil {
			switch name {
			case "year":
				periodComponents.year = num
			case "month":
				periodComponents.month = num
			case "week":
				periodComponents.week = num
			case "day":
				periodComponents.day = num
			case "hour":
				timeComponents.hour = num
			case "minute":
				timeComponents.minute = num
			case "second":
				timeComponents.second = num
			}
		} else {
			return 0, fmt.Errorf("invalid duration format: %s, %s: %w", s, name, err)
		}
	}

	t2 := t
	if periodComponents.year > 0 || periodComponents.month > 0 || periodComponents.week > 0 || periodComponents.day > 0 {
		if positive {
			t2 = t2.AddDate(periodComponents.year, periodComponents.month, periodComponents.week*7+periodComponents.day)
		} else {
			t2 = t2.AddDate(-periodComponents.year, -periodComponents.month, -(periodComponents.week*7 + periodComponents.day))
		}
	}
	if timeComponents.hour > 0 || timeComponents.minute > 0 || timeComponents.second > 0 {
		duration := time.Duration(timeComponents.hour)*time.Hour + time.Duration(timeComponents.minute)*time.Minute + time.Duration(timeComponents.second)*time.Second
		if positive {
			t2 = t2.Add(duration)
		} else {
			t2 = t2.Add(-duration)
		}
	}

	return t2.Sub(t), nil
}
