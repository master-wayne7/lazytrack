package cmd

import (
	"fmt"
	"time"

	"github.com/master-wayne7/lazytrack/notification"
	"github.com/master-wayne7/lazytrack/store"
	"github.com/spf13/cobra"
)

// NewReminderCmd creates the reminder command
func NewReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reminder",
		Short: "Check for pending goals and show late reminders",
		Long: `Check for pending goals and show late reminders.

This command checks all your habits and shows notifications for:
- Pending goals that haven't been reached yet
- Late reminders when it's getting close to 8 PM

Examples:
  lazytrack reminder          # Check all pending goals
  lazytrack reminder --late   # Show late reminder only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lateOnly, _ := cmd.Flags().GetBool("late")
			return runReminder(lateOnly)
		},
	}

	cmd.Flags().BoolP("late", "l", false, "Show late reminder only (after 8 PM)")
	return cmd
}

// runReminder handles the reminder command execution
func runReminder(lateOnly bool) error {
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

	// Get current time
	now := time.Now()
	currentHour := now.Hour()

	// Check if it's late (after 8 PM)
	isLate := currentHour >= 20

	var pendingHabits []string
	var pendingHabitsWithProgress []string

	// Check each habit for pending goals
	for _, habit := range habits {
		if habit.DailyGoal == 0 {
			continue // Skip habits without goals
		}

		// Get today's logs for this habit
		today := now.Truncate(24 * time.Hour)
		tomorrow := today.AddDate(0, 0, 1)

		logs, err := store.GetLogsByHabit(habit.Name, today, tomorrow)
		if err != nil {
			continue // Skip if we can't get logs
		}

		// Calculate current progress
		var currentProgress int
		if habit.GoalType == "count" {
			for _, log := range logs {
				currentProgress += log.Count
			}
		} else {
			// For duration-based habits, convert to minutes
			for _, log := range logs {
				if log.Duration != "" {
					// Parse duration and add to progress
					// This is a simplified version - you might want to use the parser
					currentProgress += 30 // Default 30 minutes for now
				}
			}
		}

		// Check if goal is not reached
		if currentProgress < habit.DailyGoal {
			pendingHabits = append(pendingHabits, habit.Name)
			pendingHabitsWithProgress = append(pendingHabitsWithProgress,
				fmt.Sprintf("%s (%d/%d)", habit.Name, currentProgress, habit.DailyGoal))
		}
	}

	// Show appropriate notifications
	if len(pendingHabits) > 0 {
		if lateOnly && isLate {
			// Show late reminder
			if notification.IsNotificationEnabled() {
				if err := notification.ShowLateReminder(pendingHabits); err != nil {
					fmt.Printf("⚠️  Late reminder notification failed: %v\n", err)
				}
			}
			fmt.Printf("🌙 Late reminder: You still have pending goals: %s\n", joinHabits(pendingHabits))
		} else if !lateOnly {
			// Show general reminder for each pending habit
			for _, habitName := range pendingHabits {
				habit, err := store.GetHabitByName(habitName)
				if err != nil {
					continue
				}

				// Get progress for this specific habit
				today := now.Truncate(24 * time.Hour)
				tomorrow := today.AddDate(0, 0, 1)
				logs, err := store.GetLogsByHabit(habitName, today, tomorrow)
				if err != nil {
					continue
				}

				var currentProgress int
				if habit.GoalType == "count" {
					for _, log := range logs {
						currentProgress += log.Count
					}
				} else {
					for _, log := range logs {
						if log.Duration != "" {
							currentProgress += 30 // Simplified
						}
					}
				}

				if notification.IsNotificationEnabled() {
					if err := notification.ShowGoalReminder(habitName, currentProgress, habit.DailyGoal, habit.GoalType); err != nil {
						fmt.Printf("⚠️  Goal reminder notification failed: %v\n", err)
					}
				}
			}
			fmt.Printf("📋 Pending goals: %s\n", joinHabits(pendingHabitsWithProgress))
		}
	} else {
		if !lateOnly {
			fmt.Println("✅ All goals completed for today!")
		}
	}

	return nil
}

// joinHabits joins habit names with commas (duplicate of sound package, but needed here)
func joinHabits(habits []string) string {
	if len(habits) == 0 {
		return "none"
	}
	if len(habits) == 1 {
		return habits[0]
	}
	if len(habits) == 2 {
		return habits[0] + " and " + habits[1]
	}

	result := ""
	for i, habit := range habits[:len(habits)-1] {
		if i > 0 {
			result += ", "
		}
		result += habit
	}
	result += " and " + habits[len(habits)-1]
	return result
}
