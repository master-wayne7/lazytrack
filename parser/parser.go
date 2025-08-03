package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/master-wayne7/lazytrack/types"
)

// ParseDuration parses duration strings like "2h", "30m", "1h30m"
func ParseDuration(input string) (types.ParsedDuration, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return types.ParsedDuration{IsValid: false}, fmt.Errorf("empty duration")
	}

	// Handle count-based habits (e.g., "5x", "3 times")
	if strings.HasSuffix(input, "x") || strings.HasSuffix(input, "times") {
		return parseCount(input)
	}

	// Parse time duration
	hours, minutes, err := parseTimeDuration(input)
	if err != nil {
		return types.ParsedDuration{IsValid: false}, err
	}

	return types.ParsedDuration{
		Hours:   hours,
		Minutes: minutes,
		IsValid: true,
	}, nil
}

// parseTimeDuration handles time-based durations
func parseTimeDuration(input string) (int, int, error) {
	// Regex to match patterns like "2h", "30m", "1h30m", "1.5h"
	timeRegex := regexp.MustCompile(`^(\d+(?:\.\d+)?)?h?(\d+)?m?$`)
	matches := timeRegex.FindStringSubmatch(input)

	if len(matches) == 0 {
		return 0, 0, fmt.Errorf("invalid duration format: %s", input)
	}

	hours := 0
	minutes := 0

	// Parse hours
	if len(matches) > 1 && matches[1] != "" {
		if h, err := strconv.ParseFloat(matches[1], 64); err == nil {
			hours = int(h)
			// Handle fractional hours (e.g., "1.5h" = 1h30m)
			if h != float64(int(h)) {
				fractionalMinutes := int((h - float64(int(h))) * 60)
				minutes += fractionalMinutes
			}
		}
	}

	// Parse minutes
	if len(matches) > 2 && matches[2] != "" {
		if m, err := strconv.Atoi(matches[2]); err == nil {
			minutes += m
		}
	}

	// Convert excess minutes to hours
	if minutes >= 60 {
		hours += minutes / 60
		minutes = minutes % 60
	}

	return hours, minutes, nil
}

// parseCount handles count-based durations
func parseCount(input string) (types.ParsedDuration, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	// Handle "x" suffix
	if strings.HasSuffix(input, "x") {
		countStr := strings.TrimSuffix(input, "x")
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			return types.ParsedDuration{
				Hours:   count, // Use hours field for count
				Minutes: 0,
				IsValid: true,
			}, nil
		}
	}

	// Handle "times" suffix
	if strings.HasSuffix(input, "times") {
		countStr := strings.TrimSuffix(input, "times")
		countStr = strings.TrimSpace(countStr)
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			return types.ParsedDuration{
				Hours:   count, // Use hours field for count
				Minutes: 0,
				IsValid: true,
			}, nil
		}
	}

	return types.ParsedDuration{IsValid: false}, fmt.Errorf("invalid count format: %s", input)
}

// FormatDuration formats duration for display
func FormatDuration(duration types.ParsedDuration) string {
	if !duration.IsValid {
		return "0m"
	}

	if duration.Hours > 0 && duration.Minutes > 0 {
		return fmt.Sprintf("%dh%dm", duration.Hours, duration.Minutes)
	} else if duration.Hours > 0 {
		return fmt.Sprintf("%dh", duration.Hours)
	} else {
		return fmt.Sprintf("%dm", duration.Minutes)
	}
}

// FormatCount formats count for display
func FormatCount(count int) string {
	if count == 1 {
		return "1 time"
	}
	return fmt.Sprintf("%d times", count)
}

// IsCountBased checks if the input represents a count-based habit
func IsCountBased(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))
	return strings.HasSuffix(input, "x") || strings.HasSuffix(input, "times")
}

// GetTotalMinutes converts duration to total minutes
func GetTotalMinutes(duration types.ParsedDuration) int {
	return duration.Hours*60 + duration.Minutes
}

// GetTotalHours converts duration to total hours (as float)
func GetTotalHours(duration types.ParsedDuration) float64 {
	return float64(duration.Hours) + float64(duration.Minutes)/60.0
}
