package cmd

import (
	"fmt"
	"time"

	"github.com/master-wayne7/lazytrack/notification"
	"github.com/master-wayne7/lazytrack/store"
	"github.com/spf13/cobra"
)

// NewDaemonCmd creates the daemon command for automatic reminders
func NewDaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run LazyTrack daemon for automatic late reminders",
		Long: `Run LazyTrack daemon for automatic late reminders.

This command runs in the background and automatically shows late reminders after 8 PM.
It checks every hour for pending goals and shows notifications when appropriate.

Examples:
  lazytrack daemon              # Run daemon in foreground
  lazytrack daemon --background # Run daemon in background`,
		RunE: func(cmd *cobra.Command, args []string) error {
			background, _ := cmd.Flags().GetBool("background")
			return runDaemon(background)
		},
	}

	cmd.Flags().BoolP("background", "b", false, "Run daemon in background")
	return cmd
}

// runDaemon handles the daemon command execution
func runDaemon(background bool) error {
	if background {
		// For now, just run in foreground
		// In the future, this could fork to background
		fmt.Println("🔄 Starting LazyTrack daemon...")
		fmt.Println("📅 Will check for late reminders after 8 PM")
		fmt.Println("⏰ Checking every hour for pending goals")
		fmt.Println("💡 Press Ctrl+C to stop the daemon")
	}

	// Check interval (1 hour)
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Check immediately on startup
	if err := checkAndShowLateReminder(); err != nil {
		fmt.Printf("⚠️  Error checking reminders: %v\n", err)
	}

	fmt.Println("✅ Daemon started successfully!")

	// Run the daemon loop
	for {
		select {
		case <-ticker.C:
			if err := checkAndShowLateReminder(); err != nil {
				fmt.Printf("⚠️  Error checking reminders: %v\n", err)
			}
		}
	}
}

// checkAndShowLateReminder checks if it's late and shows reminders
func checkAndShowLateReminder() error {
	// Check if it's time to show late reminder
	if !notification.ShouldShowLateReminder() {
		return nil // Not late enough yet
	}

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
	var pendingHabits []string

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
		}
	}

	// Show late reminder if there are pending habits
	if len(pendingHabits) > 0 {
		if notification.IsNotificationEnabled() {
			if err := notification.ShowLateReminder(pendingHabits); err != nil {
				return fmt.Errorf("late reminder notification failed: %w", err)
			}
		}
		fmt.Printf("🌙 Late reminder sent for: %s\n", joinHabitsDaemon(pendingHabits))
	}

	return nil
}

// joinHabitsDaemon joins habit names with commas (for daemon)
func joinHabitsDaemon(habits []string) string {
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
