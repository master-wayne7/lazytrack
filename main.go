package main

import (
	"fmt"
	"os"

	"github.com/master-wayne7/lazytrack/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lazytrack",
		Short: "A fun CLI-based time/habit tracker",
		Long: `LazyTrack - Your personal productivity companion in the terminal!

Track your habits, view beautiful summaries, and stay motivated with sound feedback.
Perfect for developers, students, and anyone who loves the command line!`,
		Version: "1.0.0",
	}

	// Add subcommands
	rootCmd.AddCommand(cmd.NewLogCmd())
	rootCmd.AddCommand(cmd.NewSummaryCmd())
	rootCmd.AddCommand(cmd.NewConfigCmd())
	rootCmd.AddCommand(cmd.NewReminderCmd())
	rootCmd.AddCommand(cmd.NewDaemonCmd())

	// Set up default behavior for logging habits
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// Handle case where no subcommand is provided (treat as log command)
	if len(os.Args) > 1 && os.Args[1] != "help" && os.Args[1] != "--help" && os.Args[1] != "-h" && os.Args[1] != "version" && os.Args[1] != "--version" && os.Args[1] != "-v" {
		// Check if it's not a known subcommand
		knownCommands := []string{"summary", "config", "reminder", "daemon", "help", "version"}
		isKnownCommand := false
		for _, cmd := range knownCommands {
			if os.Args[1] == cmd {
				isKnownCommand = true
				break
			}
		}

		if !isKnownCommand {
			// Treat as log command
			logCmd := cmd.NewLogCmd()
			logCmd.SetArgs(os.Args[1:])
			if err := logCmd.Execute(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
