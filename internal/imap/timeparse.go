package imap

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var relativeTimePattern = regexp.MustCompile(`^(\d+)(m|h|d)$`)

func ParseRelativeTime(timeStr string) (string, error) {
	matches := relativeTimePattern.FindStringSubmatch(timeStr)
	if len(matches) != 3 {
		return "", fmt.Errorf("Invalid time format. Use: 30m, 2h, 7d")
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", fmt.Errorf("Invalid time format. Use: 30m, 2h, 7d")
	}

	now := time.Now()
	unit := matches[2]

	var past time.Time
	switch unit {
	case "m":
		past = now.Add(-time.Duration(value) * time.Minute)
	case "h":
		past = now.Add(-time.Duration(value) * time.Hour)
	case "d":
		past = now.Add(-time.Duration(value) * 24 * time.Hour)
	default:
		return "", fmt.Errorf("Invalid time format. Use: 30m, 2h, 7d")
	}

	return FormatIMAPDate(past), nil
}

func FormatIMAPDate(date time.Time) string {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	utc := date.UTC()
	day := utc.Day()
	month := months[int(utc.Month())-1]
	year := utc.Year()

	return fmt.Sprintf("%02d-%s-%04d", day, month, year)
}
