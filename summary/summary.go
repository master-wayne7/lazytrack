package summary

import (
	"fmt"
	"strings"
	"time"

	"github.com/master-wayne7/lazytrack/parser"
	"github.com/master-wayne7/lazytrack/types"
)

// CalculateWeeklySummary calculates a weekly summary for all habits
func CalculateWeeklySummary(habits []types.Habit, logsByHabit map[string][]types.Log) types.WeeklySummary {
	startDate := getWeekStart(time.Now())
	endDate := startDate.AddDate(0, 0, 7)

	var summaries []types.Summary
	totalTime := 0.0

	for _, habit := range habits {
		logs := logsByHabit[habit.Name]
		summary := calculateHabitSummary(habit, logs, startDate, endDate)
		summaries = append(summaries, summary)
		totalTime += summary.TotalTime
	}

	return types.WeeklySummary{
		StartDate: startDate,
		EndDate:   endDate,
		Habits:    summaries,
		TotalTime: totalTime,
	}
}

// calculateHabitSummary calculates summary for a single habit
func calculateHabitSummary(habit types.Habit, logs []types.Log, startDate, endDate time.Time) types.Summary {
	var totalTime float64
	var totalCount int
	var streak int

	// Filter logs for the week
	var weekLogs []types.Log
	for _, log := range logs {
		if log.LoggedAt.After(startDate) && log.LoggedAt.Before(endDate) {
			weekLogs = append(weekLogs, log)
		}
	}

	// Calculate totals
	for _, log := range weekLogs {
		if habit.GoalType == "count" {
			totalCount += log.Count
		} else {
			if log.Duration != "" {
				duration, err := parser.ParseDuration(log.Duration)
				if err == nil {
					totalTime += parser.GetTotalHours(duration)
				}
			}
		}
	}

	// Calculate goal progress
	var goalProgress float64
	if habit.DailyGoal > 0 {
		if habit.GoalType == "count" {
			goalProgress = float64(totalCount) / float64(habit.DailyGoal*7) * 100
		} else {
			goalProgress = totalTime / float64(habit.DailyGoal*7) * 100
		}
	}

	// Calculate streak (simplified)
	streak = calculateStreak(weekLogs)

	// Generate bar chart
	barChart := generateBarChart(habit, totalTime, totalCount, habit.DailyGoal)

	return types.Summary{
		HabitName:    habit.Name,
		Emoji:        habit.Emoji,
		TotalTime:    totalTime,
		TotalCount:   totalCount,
		GoalProgress: goalProgress,
		Streak:       streak,
		BarChart:     barChart,
	}
}

// generateBarChart creates a visual bar chart
func generateBarChart(habit types.Habit, totalTime float64, totalCount int, dailyGoal int) string {
	var value float64
	var maxValue float64

	if habit.GoalType == "count" {
		value = float64(totalCount)
		maxValue = float64(dailyGoal * 7) // Weekly goal
	} else {
		value = totalTime
		maxValue = float64(dailyGoal * 7) // Weekly goal in hours
	}

	if maxValue == 0 {
		maxValue = 10 // Default max for visualization
	}

	// Create bar chart with max 20 characters
	barLength := int((value / maxValue) * 20)
	if barLength > 20 {
		barLength = 20
	}
	if barLength < 0 {
		barLength = 0
	}

	bar := strings.Repeat("█", barLength)
	if barLength < 20 {
		bar += strings.Repeat("░", 20-barLength)
	}

	return bar
}

// calculateStreak calculates the current streak (simplified)
func calculateStreak(logs []types.Log) int {
	if len(logs) == 0 {
		return 0
	}

	// Sort logs by date (most recent first)
	// For simplicity, just count unique days with logs
	uniqueDays := make(map[string]bool)
	for _, log := range logs {
		day := log.LoggedAt.Format("2006-01-02")
		uniqueDays[day] = true
	}

	return len(uniqueDays)
}

// getWeekStart returns the start of the current week (Monday)
func getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	daysToSubtract := weekday - 1 // Monday = 1
	return t.AddDate(0, 0, -daysToSubtract).Truncate(24 * time.Hour)
}

// FormatSummary formats the summary for display
func FormatSummary(summary types.WeeklySummary) string {
	var result strings.Builder

	// Header
	result.WriteString("📊 Weekly Summary\n")
	result.WriteString("=" + strings.Repeat("=", 50) + "\n")
	result.WriteString(fmt.Sprintf("📅 %s - %s\n\n", 
		summary.StartDate.Format("Jan 2"), 
		summary.EndDate.AddDate(0, 0, -1).Format("Jan 2")))

	// Habit summaries
	for _, habit := range summary.Habits {
		result.WriteString(formatHabitSummary(habit))
		result.WriteString("\n")
	}

	// Total
	result.WriteString("\n" + strings.Repeat("=", 52) + "\n")
	result.WriteString(fmt.Sprintf("🎯 Total Time: %.1f hours\n", summary.TotalTime))

	return result.String()
}

// formatHabitSummary formats a single habit summary
func formatHabitSummary(summary types.Summary) string {
	var result strings.Builder

	// Emoji and name
	result.WriteString(fmt.Sprintf("%s %s ", summary.Emoji, summary.HabitName))

	// Bar chart
	result.WriteString(summary.BarChart + " ")

	// Values
	if summary.TotalTime > 0 {
		result.WriteString(fmt.Sprintf("%.1fh", summary.TotalTime))
	} else if summary.TotalCount > 0 {
		result.WriteString(fmt.Sprintf("%dx", summary.TotalCount))
	} else {
		result.WriteString("0")
	}

	// Goal progress
	if summary.GoalProgress > 0 {
		result.WriteString(fmt.Sprintf(" (%.0f%% of goal)", summary.GoalProgress))
	}

	// Streak
	if summary.Streak > 0 {
		result.WriteString(fmt.Sprintf(" 🔥 %d day streak", summary.Streak))
	}

	return result.String()
}

// GetMotivationalMessage returns a motivational message based on progress
func GetMotivationalMessage(summary types.WeeklySummary) string {
	var totalProgress float64
	var completedHabits int

	for _, habit := range summary.Habits {
		if habit.GoalProgress > 0 {
			totalProgress += habit.GoalProgress
			completedHabits++
		}
	}

	if completedHabits == 0 {
		return "🌟 Every journey starts with a single step! Log your first habit today!"
	}

	avgProgress := totalProgress / float64(completedHabits)

	switch {
	case avgProgress >= 100:
		return "🎉 Amazing! You're crushing your goals this week!"
	case avgProgress >= 80:
		return "🚀 Great progress! You're so close to your goals!"
	case avgProgress >= 60:
		return "💪 Good work! Keep up the momentum!"
	case avgProgress >= 40:
		return "👍 You're making progress! Every bit counts!"
	case avgProgress >= 20:
		return "🌱 Getting started is the hardest part. You're doing great!"
	default:
		return "🌟 Every small step counts! Keep going!"
	}
}

// CalculateDailyProgress calculates progress for today
func CalculateDailyProgress(habit types.Habit, todayLogs []types.Log) float64 {
	if habit.DailyGoal == 0 {
		return 0
	}

	var total float64
	for _, log := range todayLogs {
		if habit.GoalType == "count" {
			total += float64(log.Count)
		} else {
			if log.Duration != "" {
				duration, err := parser.ParseDuration(log.Duration)
				if err == nil {
					total += parser.GetTotalHours(duration)
				}
			}
		}
	}

	return (total / float64(habit.DailyGoal)) * 100
}

// IsGoalReached checks if today's goal is reached
func IsGoalReached(habit types.Habit, todayLogs []types.Log) bool {
	progress := CalculateDailyProgress(habit, todayLogs)
	return progress >= 100
} 