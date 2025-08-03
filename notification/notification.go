package notification

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// ShowNotification shows a popup notification
func ShowNotification(title, message string) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		return showMacNotification(title, message)
	case "linux":
		return showLinuxNotification(title, message)
	case "windows":
		return showWindowsNotification(title, message)
	default:
		// Fallback: just print a message
		fmt.Printf("📢 %s: %s\n", title, message)
		return nil
	}
}

// showMacNotification shows notification on macOS
func showMacNotification(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// showLinuxNotification shows notification on Linux
func showLinuxNotification(title, message string) error {
	// Try different notification systems
	notifiers := []string{"notify-send", "zenity", "kdialog"}

	for _, notifier := range notifiers {
		if _, err := exec.LookPath(notifier); err == nil {
			switch notifier {
			case "notify-send":
				cmd := exec.Command("notify-send", title, message)
				return cmd.Run()
			case "zenity":
				cmd := exec.Command("zenity", "--info", "--title", title, "--text", message)
				return cmd.Run()
			case "kdialog":
				cmd := exec.Command("kdialog", "--title", title, "--msgbox", message)
				return cmd.Run()
			}
		}
	}

	// Fallback: just print a message
	fmt.Printf("📢 %s: %s\n", title, message)
	return nil
}

// showWindowsNotification shows notification on Windows
func showWindowsNotification(title, message string) error {
	// Use PowerShell to show a Windows notification
	psScript := fmt.Sprintf(`
		Add-Type -AssemblyName System.Windows.Forms
		$notification = New-Object System.Windows.Forms.NotifyIcon
		$notification.Icon = [System.Drawing.SystemIcons]::Information
		$notification.Visible = $true
		$notification.ShowBalloonTip(5000, "%s", "%s", [System.Windows.Forms.ToolTipIcon]::Info)
		Start-Sleep -Seconds 6
		$notification.Dispose()
	`, title, message)

	cmd := exec.Command("powershell", "-Command", psScript)
	return cmd.Run()
}

// ShowGoalReminder shows a reminder for pending goals
func ShowGoalReminder(habitName string, currentProgress, goal int, goalType string) error {
	var message string
	if goalType == "count" {
		message = fmt.Sprintf("You've completed %d/%d %s today. Don't forget to reach your goal!", currentProgress, goal, habitName)
	} else {
		message = fmt.Sprintf("You've logged %.1f hours of %s today. Keep going!", float64(currentProgress)/60.0, habitName)
	}

	return ShowNotification("LazyTrack Reminder", message)
}

// ShowLateReminder shows a reminder when it's getting late
func ShowLateReminder(pendingHabits []string) error {
	message := fmt.Sprintf("It's getting late! You still have pending goals: %s", joinHabits(pendingHabits))
	return ShowNotification("LazyTrack Late Reminder", message)
}

// joinHabits joins habit names with commas
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

// IsNotificationEnabled checks if notifications are enabled
func IsNotificationEnabled() bool {
	// Check if notifications are disabled via environment variable
	if os.Getenv("LAZYTRACK_NOTIFICATIONS_DISABLED") == "1" {
		return false
	}

	// For now, always return true
	// In the future, this could check a config file
	return true
}

// ShouldShowLateReminder checks if it's time to show the late reminder
func ShouldShowLateReminder() bool {
	now := time.Now()
	currentHour := now.Hour()

	// Show late reminder after 8 PM (20:00)
	return currentHour >= 20
}
