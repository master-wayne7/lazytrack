package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/master-wayne7/lazytrack/parser"
	"github.com/master-wayne7/lazytrack/store"
	"github.com/master-wayne7/lazytrack/summary"
	"github.com/master-wayne7/lazytrack/types"
	"github.com/spf13/cobra"
)

// NewSummaryCmd creates the summary command
func NewSummaryCmd() *cobra.Command {
	var weekly bool
	var daily bool

	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Show habit summary",
		Long: `Show a summary of your habits.

Examples:
  lazytrack summary          # Show weekly summary
  lazytrack summary --daily  # Show daily summary
  lazytrack summary --weekly # Show weekly summary (default)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSummary(weekly, daily)
		},
	}

	cmd.Flags().BoolVarP(&weekly, "weekly", "w", true, "Show weekly summary")
	cmd.Flags().BoolVarP(&daily, "daily", "d", false, "Show daily summary")
	return cmd
}

// runSummary handles the summary command execution
func runSummary(weekly, daily bool) error {
	// Initialize store
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer store.Close()

	// Get all habits
	habits, err := store.GetAllHabits()
	if err != nil {
		return fmt.Errorf("failed to get habits: %w", err)
	}

	if len(habits) == 0 {
		displayEmptyState()
		return nil
	}

	// Get logs for all habits
	logsByHabit := make(map[string][]types.Log)
	for _, habit := range habits {
		var startDate, endDate time.Time
		
		if daily {
			startDate = time.Now().Truncate(24 * time.Hour)
			endDate = startDate.AddDate(0, 0, 1)
		} else {
			startDate = getWeekStart(time.Now())
			endDate = startDate.AddDate(0, 0, 7)
		}

		logs, err := store.GetLogsByHabit(habit.Name, startDate, endDate)
		if err != nil {
			continue // Skip habits with errors
		}
		logsByHabit[habit.Name] = logs
	}

	// Calculate and display summary
	if daily {
		displayDailySummary(habits, logsByHabit)
	} else {
		displayWeeklySummary(habits, logsByHabit)
	}

	return nil
}

// displayWeeklySummary shows the weekly summary
func displayWeeklySummary(habits []types.Habit, logsByHabit map[string][]types.Log) {
	weeklySummary := summary.CalculateWeeklySummary(habits, logsByHabit)
	
	// Display formatted summary
	fmt.Println(summary.FormatSummary(weeklySummary))
	
	// Display motivational message
	motivationalMsg := summary.GetMotivationalMessage(weeklySummary)
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n" + motivationalMsg)
}

// displayDailySummary shows the daily summary
func displayDailySummary(habits []types.Habit, logsByHabit map[string][]types.Log) {
	today := time.Now().Truncate(24 * time.Hour)
	
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("📅 Daily Summary - %s\n", today.Format("Monday, January 2, 2006"))
	cyan.Println(strings.Repeat("=", 50))
	
	var totalTime float64
	var totalCount int
	
	for _, habit := range habits {
		logs := logsByHabit[habit.Name]
		if len(logs) == 0 {
			continue
		}
		
		// Calculate daily totals
		var habitTime float64
		var habitCount int
		
		for _, log := range logs {
			if habit.GoalType == "count" {
				habitCount += log.Count
			} else {
				if log.Duration != "" {
					duration, err := parser.ParseDuration(log.Duration)
					if err == nil {
						habitTime += parser.GetTotalHours(duration)
					}
				}
			}
		}
		
		// Display habit summary
		displayDailyHabitSummary(habit, habitTime, habitCount)
		
		if habit.GoalType == "count" {
			totalCount += habitCount
		} else {
			totalTime += habitTime
		}
	}
	
	// Display totals
	fmt.Println("\n" + strings.Repeat("=", 50))
	if totalTime > 0 {
		fmt.Printf("🎯 Total Time Today: %.1f hours\n", totalTime)
	}
	if totalCount > 0 {
		fmt.Printf("🎯 Total Count Today: %d\n", totalCount)
	}
}

// displayDailyHabitSummary shows a single habit's daily summary
func displayDailyHabitSummary(habit types.Habit, totalTime float64, totalCount int) {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	
	// Emoji and name
	fmt.Printf("%s %s ", habit.Emoji, habit.Name)
	
	// Values
	if habit.GoalType == "count" {
		if totalCount > 0 {
			green.Printf("%dx", totalCount)
		} else {
			fmt.Print("0")
		}
	} else {
		if totalTime > 0 {
			green.Printf("%.1fh", totalTime)
		} else {
			fmt.Print("0")
		}
	}
	
	// Goal progress
	if habit.DailyGoal > 0 {
		var progress float64
		if habit.GoalType == "count" {
			progress = float64(totalCount) / float64(habit.DailyGoal) * 100
		} else {
			progress = totalTime / float64(habit.DailyGoal) * 100
		}
		
		if progress > 0 {
			yellow.Printf(" (%.0f%% of daily goal)", progress)
		}
	}
	
	fmt.Println()
}

// displayEmptyState shows a message when no habits exist
func displayEmptyState() {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("🌟 Welcome to LazyTrack!")
	fmt.Println()
	fmt.Println("You haven't logged any habits yet. Start by logging your first habit:")
	fmt.Println()
	fmt.Println("  lazytrack code 2h          # Log 2 hours of coding")
	fmt.Println("  lazytrack walk 30m         # Log 30 minutes of walking")
	fmt.Println("  lazytrack water 8x         # Log 8 glasses of water")
	fmt.Println("  lazytrack read             # Log default duration (30m)")
	fmt.Println()
	fmt.Println("Then run 'lazytrack summary' to see your progress!")
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