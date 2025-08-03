package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/master-wayne7/lazytrack/types"
)

type Store struct {
	dataPath string
	habits   map[string]*types.Habit
	logs     []types.Log
	config   map[string]string
}

// NewStore creates a new store instance
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataPath := filepath.Join(homeDir, ".lazytrack")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &Store{
		dataPath: dataPath,
		habits:   make(map[string]*types.Habit),
		logs:     []types.Log{},
		config:   make(map[string]string),
	}

	if err := store.loadData(); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return store, nil
}

// Close closes the store (saves data)
func (s *Store) Close() error {
	return s.saveData()
}

// loadData loads data from JSON files
func (s *Store) loadData() error {
	// Load habits
	habitsPath := filepath.Join(s.dataPath, "habits.json")
	if data, err := os.ReadFile(habitsPath); err == nil {
		if err := json.Unmarshal(data, &s.habits); err != nil {
			return fmt.Errorf("failed to unmarshal habits: %w", err)
		}
	}

	// Load logs
	logsPath := filepath.Join(s.dataPath, "logs.json")
	if data, err := os.ReadFile(logsPath); err == nil {
		if err := json.Unmarshal(data, &s.logs); err != nil {
			return fmt.Errorf("failed to unmarshal logs: %w", err)
		}
	}

	// Load config
	configPath := filepath.Join(s.dataPath, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &s.config); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return nil
}

// saveData saves data to JSON files
func (s *Store) saveData() error {
	// Save habits
	habitsPath := filepath.Join(s.dataPath, "habits.json")
	if data, err := json.MarshalIndent(s.habits, "", "  "); err == nil {
		if err := os.WriteFile(habitsPath, data, 0644); err != nil {
			return fmt.Errorf("failed to save habits: %w", err)
		}
	}

	// Save logs
	logsPath := filepath.Join(s.dataPath, "logs.json")
	if data, err := json.MarshalIndent(s.logs, "", "  "); err == nil {
		if err := os.WriteFile(logsPath, data, 0644); err != nil {
			return fmt.Errorf("failed to save logs: %w", err)
		}
	}

	// Save config
	configPath := filepath.Join(s.dataPath, "config.json")
	if data, err := json.MarshalIndent(s.config, "", "  "); err == nil {
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}

	return nil
}

// GetOrCreateHabit gets an existing habit or creates a new one
func (s *Store) GetOrCreateHabit(name string) (*types.Habit, error) {
	if habit, exists := s.habits[name]; exists {
		return habit, nil
	}

	// Create new habit with smart defaults
	habit := &types.Habit{
		ID:              len(s.habits) + 1,
		Name:            name,
		Emoji:           getDefaultEmoji(name),
		DefaultDuration: getDefaultDuration(name),
		DailyGoal:       getDefaultGoal(name),
		GoalType:        getDefaultGoalType(name),
		CreatedAt:       time.Now(),
	}

	s.habits[name] = habit
	return habit, nil
}

// GetHabitByName gets a habit by name
func (s *Store) GetHabitByName(name string) (*types.Habit, error) {
	if habit, exists := s.habits[name]; exists {
		return habit, nil
	}
	return nil, fmt.Errorf("habit not found: %s", name)
}

// AddLog adds a new log entry
func (s *Store) AddLog(habitID int, habitName, duration string, count int, notes string) error {
	log := types.Log{
		ID:        len(s.logs) + 1,
		HabitID:   habitID,
		HabitName: habitName,
		Duration:  duration,
		Count:     count,
		LoggedAt:  time.Now(),
		Notes:     notes,
	}

	s.logs = append(s.logs, log)
	return nil
}

// GetLogsByHabit gets logs for a specific habit within a date range
func (s *Store) GetLogsByHabit(habitName string, startDate, endDate time.Time) ([]types.Log, error) {
	var filteredLogs []types.Log

	for _, log := range s.logs {
		if log.HabitName == habitName && log.LoggedAt.After(startDate) && log.LoggedAt.Before(endDate) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs, nil
}

// GetAllHabits gets all habits
func (s *Store) GetAllHabits() ([]types.Habit, error) {
	var habits []types.Habit
	for _, habit := range s.habits {
		habits = append(habits, *habit)
	}
	return habits, nil
}

// UpdateHabit updates a habit's configuration
func (s *Store) UpdateHabit(habit *types.Habit) error {
	s.habits[habit.Name] = habit
	return nil
}

// GetConfig gets a configuration value
func (s *Store) GetConfig(key string) (string, error) {
	if value, exists := s.config[key]; exists {
		return value, nil
	}
	return "", fmt.Errorf("config not found: %s", key)
}

// SetConfig sets a configuration value
func (s *Store) SetConfig(key, value string) error {
	s.config[key] = value
	return nil
}

// getDefaultEmoji returns a default emoji based on habit name
func getDefaultEmoji(name string) string {
	emojiMap := map[string]string{
		"code":     "💻",
		"read":     "📖",
		"walk":     "🚶",
		"run":      "🏃",
		"exercise": "💪",
		"water":    "💧",
		"sleep":    "😴",
		"meditate": "🧘",
		"write":    "✍️",
		"study":    "📚",
		"work":     "💼",
		"gym":      "🏋️",
		"yoga":     "🧘‍♀️",
		"cook":     "👨‍🍳",
		"clean":    "🧹",
		"paint":    "🎨",
		"music":    "🎵",
		"game":     "🎮",
		"social":   "👥",
		"family":   "👨‍👩‍👧‍👦",
	}

	if emoji, exists := emojiMap[name]; exists {
		return emoji
	}
	return "📝"
}

// getDefaultDuration returns a default duration based on habit name
func getDefaultDuration(name string) string {
	countBasedHabits := map[string]bool{
		"water":    true,
		"medicine": true,
		"vitamins": true,
		"pills":    true,
		"steps":    true,
		"pushups":  true,
		"squats":   true,
		"pullups":  true,
	}

	if countBasedHabits[name] {
		return "1x" // Default to 1 count for count-based habits
	}
	return "30m" // Default to 30 minutes for time-based habits
}

// getDefaultGoal returns a default daily goal based on habit name
func getDefaultGoal(name string) int {
	goalMap := map[string]int{
		"water":    8,     // 8 glasses of water
		"medicine": 1,     // 1 dose
		"vitamins": 1,     // 1 dose
		"pills":    1,     // 1 dose
		"steps":    10000, // 10,000 steps
		"pushups":  20,    // 20 pushups
		"squats":   50,    // 50 squats
		"pullups":  10,    // 10 pullups
	}

	if goal, exists := goalMap[name]; exists {
		return goal
	}
	return 0 // No default goal for other habits
}

// getDefaultGoalType returns the default goal type based on habit name
func getDefaultGoalType(name string) string {
	countBasedHabits := map[string]bool{
		"water":    true,
		"medicine": true,
		"vitamins": true,
		"pills":    true,
		"steps":    true,
		"pushups":  true,
		"squats":   true,
		"pullups":  true,
	}

	if countBasedHabits[name] {
		return "count"
	}
	return "duration"
}
