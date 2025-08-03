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

// NewLogCmd creates the log command
func NewLogCmd() *cobra.Command {
	var notes string

	cmd := &cobra.Command{
		Use:   "log [habit] [duration]",
		Short: "Log a habit with optional duration",
		Long: `Log a habit with optional duration.

Examples:
  lazytrack code 2h          # Log 2 hours of coding
  lazytrack walk 30m         # Log 30 minutes of walking
  lazytrack water 8x         # Log 8 glasses of water
  lazytrack read             # Log default duration (30m)`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLog(args, notes)
		},
	}

	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Add notes to the log entry")
	return cmd
}

// runLog handles the log command execution
func runLog(args []string, notes string) error {
	habitName := strings.ToLower(strings.TrimSpace(args[0]))

	// Initialize store
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer store.Close()

	// Get or create habit
	habit, err := store.GetOrCreateHabit(habitName)
	if err != nil {
		return fmt.Errorf("failed to get/create habit: %w", err)
	}

	// Parse duration if provided
	var duration string
	var count int
	var isCountBased bool

	if len(args) > 1 {
		durationInput := strings.TrimSpace(args[1])
		isCountBased = parser.IsCountBased(durationInput)

		if isCountBased {
			// Handle count-based habits
			parsed, err := parser.ParseDuration(durationInput)
			if err != nil {
				return fmt.Errorf("invalid count format: %v", err)
			}
			count = parsed.Hours // We use Hours field for count
		} else {
			// Handle time-based habits
			parsed, err := parser.ParseDuration(durationInput)
			if err != nil {
				return fmt.Errorf("invalid duration format: %v", err)
			}
			duration = parser.FormatDuration(parsed)
		}
	} else {
		// Use habit's default duration and determine if it's count-based
		duration = habit.DefaultDuration
		isCountBased = parser.IsCountBased(duration)

		if isCountBased {
			// Parse the default count
			parsed, err := parser.ParseDuration(duration)
			if err != nil {
				return fmt.Errorf("invalid default count format: %v", err)
			}
			count = parsed.Hours // We use Hours field for count
		}
	}

	// Add log entry
	err = store.AddLog(habit.ID, habit.Name, duration, count, notes)
	if err != nil {
		return fmt.Errorf("failed to add log: %w", err)
	}

	// Success notification removed - only show console output

	// Check if goal is reached (console output only)
	checkAndShowGoalMessage(store, habit)

	// Display success message
	displaySuccessMessage(habit, duration, count, isCountBased)

	return nil
}

// displaySuccessMessage shows a colorful success message
func displaySuccessMessage(habit *types.Habit, duration string, count int, isCountBased bool) {
	// Try color output first, fallback to regular if it fails
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)

	// Use a more robust approach - check if we're in a terminal
	if !color.NoColor {
		green.Print("✅ Logged ")
		cyan.Printf("\"%s\"", habit.Name)

		if isCountBased {
			green.Printf(" for %s", parser.FormatCount(count))
		} else {
			green.Printf(" for %s", duration)
		}

		green.Println()
	} else {
		// Fallback to regular output
		fmt.Print("✅ Logged ")
		fmt.Printf("\"%s\"", habit.Name)

		if isCountBased {
			fmt.Printf(" for %s", parser.FormatCount(count))
		} else {
			fmt.Printf(" for %s", duration)
		}

		fmt.Println()
	}

	// Show emoji and habit name
	fmt.Printf("%s %s\n", habit.Emoji, habit.Name)
}

// checkAndShowGoalMessage checks if today's goal is reached and shows console message only
func checkAndShowGoalMessage(store *store.Store, habit *types.Habit) {
	// Get today's logs for this habit
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)

	logs, err := store.GetLogsByHabit(habit.Name, today, tomorrow)
	if err != nil {
		return // Don't fail if we can't check goals
	}

	// Check if goal is reached
	if summary.IsGoalReached(*habit, logs) {
		// Show goal reached message (console only, no notification)
		if !color.NoColor {
			yellow := color.New(color.FgYellow, color.Bold)
			yellow.Printf("🎉 Goal reached for %s today!\n", habit.Name)
		} else {
			fmt.Printf("🎉 Goal reached for %s today!\n", habit.Name)
		}
	}
}

// LogHabit is a convenience function for logging habits programmatically
func LogHabit(habitName, duration string, notes string) error {
	args := []string{habitName}
	if duration != "" {
		args = append(args, duration)
	}
	return runLog(args, notes)
}
