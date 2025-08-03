package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/master-wayne7/lazytrack/store"
	"github.com/master-wayne7/lazytrack/types"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	var habitName string
	var emoji string
	var goal string
	var goalType string
	var defaultDuration string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure habits and settings",
		Long: `Configure habits and settings.

Examples:
  lazytrack config                    # Interactive configuration
  lazytrack config --habit code --emoji 💻
  lazytrack config --habit water --goal 8 --type count
  lazytrack config --habit read --goal 2 --type duration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(habitName, emoji, goal, goalType, defaultDuration)
		},
	}

	cmd.Flags().StringVarP(&habitName, "habit", "a", "", "Habit name to configure")
	cmd.Flags().StringVarP(&emoji, "emoji", "e", "", "Emoji for the habit")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Daily goal value")
	cmd.Flags().StringVarP(&goalType, "type", "t", "", "Goal type (duration or count)")
	cmd.Flags().StringVarP(&defaultDuration, "duration", "d", "", "Default duration")

	return cmd
}

// runConfig handles the config command execution
func runConfig(habitName, emoji, goal, goalType, defaultDuration string) error {
	// Initialize store
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer store.Close()

	// If no habit specified, run interactive mode
	if habitName == "" {
		return runInteractiveConfig(store)
	}

	// Get or create habit
	habit, err := store.GetOrCreateHabit(habitName)
	if err != nil {
		return fmt.Errorf("failed to get/create habit: %w", err)
	}

	// Update habit configuration
	updated := false

	if emoji != "" {
		habit.Emoji = emoji
		updated = true
	}

	if goal != "" {
		goalValue, err := strconv.Atoi(goal)
		if err != nil {
			return fmt.Errorf("invalid goal value: %s", goal)
		}
		habit.DailyGoal = goalValue
		updated = true
	}

	if goalType != "" {
		if goalType != "duration" && goalType != "count" {
			return fmt.Errorf("invalid goal type: %s (must be 'duration' or 'count')", goalType)
		}
		habit.GoalType = goalType
		updated = true
	}

	if defaultDuration != "" {
		habit.DefaultDuration = defaultDuration
		updated = true
	}

	if updated {
		err = store.UpdateHabit(habit)
		if err != nil {
			return fmt.Errorf("failed to update habit: %w", err)
		}

		green := color.New(color.FgGreen, color.Bold)
		green.Printf("✅ Updated configuration for '%s'\n", habit.Name)
		displayHabitConfig(*habit)
	} else {
		displayHabitConfig(*habit)
	}

	return nil
}

// runInteractiveConfig runs interactive configuration mode
func runInteractiveConfig(store *store.Store) error {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("🔧 LazyTrack Configuration")
	cyan.Println(strings.Repeat("=", 50))

	// Get all habits
	habits, err := store.GetAllHabits()
	if err != nil {
		return fmt.Errorf("failed to get habits: %w", err)
	}

	if len(habits) == 0 {
		cyan.Println("No habits found. Create your first habit:")
		fmt.Println()
		fmt.Println("  lazytrack code 2h")
		fmt.Println("  lazytrack walk 30m")
		fmt.Println()
		return nil
	}

	// Display current configuration
	cyan.Println("Current Habits:")
	fmt.Println()

	for i, habit := range habits {
		fmt.Printf("%d. ", i+1)
		displayHabitConfig(habit)
	}

	// Interactive configuration
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		cyan.Print("Enter habit number to configure (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" || input == "quit" {
			break
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(habits) {
			fmt.Println("❌ Invalid selection. Please try again.")
			continue
		}

		habit := habits[index-1]
		if err := configureHabit(store, &habit); err != nil {
			fmt.Printf("❌ Error configuring habit: %v\n", err)
		} else {
			habits[index-1] = habit // Update the habit in our slice
		}
	}

	return nil
}

// configureHabit interactively configures a single habit
func configureHabit(store *store.Store, habit *types.Habit) error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)

	cyan.Printf("\n🔧 Configuring '%s'\n", habit.Name)
	cyan.Println(strings.Repeat("-", 30))

	reader := bufio.NewReader(os.Stdin)

	// Configure emoji
	fmt.Printf("Current emoji: %s\n", habit.Emoji)
	fmt.Print("New emoji (press Enter to keep current): ")
	emoji, _ := reader.ReadString('\n')
	emoji = strings.TrimSpace(emoji)
	if emoji != "" {
		habit.Emoji = emoji
	}

	// Configure goal type
	fmt.Printf("Current goal type: %s\n", habit.GoalType)
	fmt.Print("Goal type (duration/count, press Enter to keep current): ")
	goalType, _ := reader.ReadString('\n')
	goalType = strings.TrimSpace(goalType)
	if goalType == "duration" || goalType == "count" {
		habit.GoalType = goalType
	}

	// Configure daily goal
	fmt.Printf("Current daily goal: %d\n", habit.DailyGoal)
	fmt.Print("New daily goal (press Enter to keep current): ")
	goalStr, _ := reader.ReadString('\n')
	goalStr = strings.TrimSpace(goalStr)
	if goalStr != "" {
		if goal, err := strconv.Atoi(goalStr); err == nil && goal >= 0 {
			habit.DailyGoal = goal
		} else {
			fmt.Println("❌ Invalid goal value")
		}
	}

	// Configure default duration
	fmt.Printf("Current default duration: %s\n", habit.DefaultDuration)
	fmt.Print("New default duration (e.g., 30m, 1h, press Enter to keep current): ")
	duration, _ := reader.ReadString('\n')
	duration = strings.TrimSpace(duration)
	if duration != "" {
		habit.DefaultDuration = duration
	}

	// Save changes
	err := store.UpdateHabit(habit)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	green.Printf("✅ Updated configuration for '%s'\n", habit.Name)
	displayHabitConfig(*habit)

	return nil
}

// displayHabitConfig displays a habit's current configuration
func displayHabitConfig(habit types.Habit) {
	fmt.Printf("%s %s", habit.Emoji, habit.Name)

	if habit.DailyGoal > 0 {
		fmt.Printf(" (Goal: %d", habit.DailyGoal)
		if habit.GoalType == "count" {
			fmt.Print(" times")
		} else {
			fmt.Print(" hours")
		}
		fmt.Print(")")
	}

	if habit.DefaultDuration != "" {
		fmt.Printf(" [Default: %s]", habit.DefaultDuration)
	}

	fmt.Println()
}
