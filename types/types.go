package types

import (
	"time"
)

// Habit represents a tracked habit
type Habit struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Emoji       string    `json:"emoji" db:"emoji"`
	DefaultDuration string `json:"default_duration" db:"default_duration"`
	DailyGoal   int       `json:"daily_goal" db:"daily_goal"`
	GoalType    string    `json:"goal_type" db:"goal_type"` // "count" or "duration"
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Log represents a single habit log entry
type Log struct {
	ID        int       `json:"id" db:"id"`
	HabitID   int       `json:"habit_id" db:"habit_id"`
	HabitName string    `json:"habit_name" db:"habit_name"`
	Duration  string    `json:"duration" db:"duration"` // e.g., "30m", "2h"
	Count     int       `json:"count" db:"count"`       // for count-based habits
	LoggedAt  time.Time `json:"logged_at" db:"logged_at"`
	Notes     string    `json:"notes" db:"notes"`
}

// Config represents user configuration
type Config struct {
	SoundEnabled bool              `json:"sound_enabled"`
	DefaultHabits map[string]Habit `json:"default_habits"`
	Theme        string            `json:"theme"` // "default", "dark", "colorful"
}

// Summary represents aggregated habit data
type Summary struct {
	HabitName    string  `json:"habit_name"`
	Emoji        string  `json:"emoji"`
	TotalTime    float64 `json:"total_time"`    // in hours
	TotalCount   int     `json:"total_count"`
	GoalProgress float64 `json:"goal_progress"` // percentage
	Streak       int     `json:"streak"`
	BarChart     string  `json:"bar_chart"`
}

// WeeklySummary represents a week's worth of data
type WeeklySummary struct {
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Habits    []Summary `json:"habits"`
	TotalTime float64    `json:"total_time"`
}

// ParsedDuration represents parsed time duration
type ParsedDuration struct {
	Hours   int
	Minutes int
	IsValid bool
}

// Goal represents a daily goal for a habit
type Goal struct {
	HabitName string `json:"habit_name"`
	Target    int    `json:"target"`
	Type      string `json:"type"` // "count" or "duration"
	Unit      string `json:"unit"` // "times" or "hours"
} 