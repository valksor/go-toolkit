// Package cli provides Cobra CLI helpers for Valksor tools.
//
// This package provides common patterns for building CLI tools with Cobra,
// including standard flags, version information, command grouping, and
// user interaction helpers.
//
// Basic usage:
//
//	var rootCmd = &cobra.Command{
//	    Use:   "mytool",
//	    Short: "My tool description",
//	}
//
//	cli.SetupRootCmd(rootCmd, cli.RootOptions{
//	    ToolName: "mytool",
//	    VersionInfo: func() string { return version.Info("mytool") },
//	})
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/valksor/go-toolkit/display"
	"github.com/valksor/go-toolkit/log"
)

// RootOptions configures the root command setup.
type RootOptions struct {
	// ToolName is the name of the CLI tool.
	ToolName string
	// VersionInfo returns formatted version information.
	VersionInfo func() string
	// PersistentPreRun is called before each command.
	// If nil, a default implementation is used that configures logging and colors.
	PersistentPreRun func(cmd *cobra.Command, args []string) error
	// PreRunHook is called before the default PersistentPreRun logic.
	// Use this for custom initialization that needs to happen before logging/colors setup.
	PreRunHook func() error
	// SilenceErrors silences error messages (errors are handled by the caller).
	SilenceErrors bool
	// SilenceUsage silences usage printing on error.
	SilenceUsage bool
}

// StandardFlags holds the standard persistent flag values.
type StandardFlags struct {
	Verbose bool
	Quiet   bool
	NoColor bool
}

// globalFlags stores the standard flag values.
var globalFlags = &StandardFlags{}

// SetupRootCmd configures a root command with standard settings.
func SetupRootCmd(cmd *cobra.Command, opts RootOptions) {
	// Enable prefix matching by default
	cobra.EnablePrefixMatching = true

	// Set silence options
	if opts.SilenceErrors {
		cmd.SilenceErrors = true
	}
	if opts.SilenceUsage {
		cmd.SilenceUsage = true
	}

	// Set default PersistentPreRun if not provided
	if opts.PersistentPreRun == nil {
		cmd.PersistentPreRunE = defaultPersistentPreRun(opts.ToolName, opts.PreRunHook)
	} else {
		cmd.PersistentPreRunE = opts.PersistentPreRun
	}

	// Add version command if version info provided
	if opts.VersionInfo != nil {
		versionCmd := &cobra.Command{
			Use:   "version",
			Short: "Show version information",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println(opts.VersionInfo())
			},
		}
		cmd.AddCommand(versionCmd)
	}
}

// AddStandardFlags adds the standard persistent flags (--verbose, --quiet, --no-color).
// Returns the flag set for further customization if needed.
func AddStandardFlags(cmd *cobra.Command) *pflag.FlagSet {
	cmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVarP(&globalFlags.Quiet, "quiet", "q", false, "Suppress non-essential output")
	cmd.PersistentFlags().BoolVar(&globalFlags.NoColor, "no-color", false, "Disable color output")

	return cmd.PersistentFlags()
}

// GetStandardFlags returns the current values of standard flags.
func GetStandardFlags() *StandardFlags {
	return globalFlags
}

// defaultPersistentPreRun returns a default PersistentPreRunE function.
func defaultPersistentPreRun(toolName string, preRunHook func() error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Call pre-run hook if provided
		if preRunHook != nil {
			if err := preRunHook(); err != nil {
				return err
			}
		}

		// Configure logging
		log.Configure(log.Options{
			Verbose: globalFlags.Verbose,
		})

		// Initialize colors
		display.InitColors(globalFlags.NoColor)

		log.Debug("initialized",
			"tool", toolName,
			"verbose", globalFlags.Verbose,
			"quiet", globalFlags.Quiet,
			"no_color", globalFlags.NoColor,
		)

		return nil
	}
}

// AddCommandGroup adds a command group to the given command.
func AddCommandGroup(cmd *cobra.Command, groupID, title string) {
	// Add group if it doesn't exist
	for _, g := range cmd.Groups() {
		if g.ID == groupID {
			return
		}
	}

	cmd.AddGroup(&cobra.Group{
		ID:    groupID,
		Title: title,
	})
}

// SetCommandGroup assigns a command to a group.
// This is a convenience wrapper around cobra.Command's GroupID field.
func SetCommandGroup(cmd *cobra.Command, groupID string) {
	cmd.GroupID = groupID
}

// CommandGroup represents a group of related commands.
type CommandGroup struct {
	ID       string
	Title    string
	Commands []*cobra.Command
}

// AddCommandGroups adds multiple command groups to the root command.
// Each group includes the group definition and assigns commands to it.
func AddCommandGroups(rootCmd *cobra.Command, groups ...CommandGroup) {
	for _, group := range groups {
		// Add the group
		AddCommandGroup(rootCmd, group.ID, group.Title)

		// Assign commands to the group
		for _, cmd := range group.Commands {
			SetCommandGroup(cmd, group.ID)
			rootCmd.AddCommand(cmd)
		}
	}
}

// ConfirmAction prompts the user for confirmation with a yes/no prompt.
// If skipConfirm is true, it returns true immediately without prompting.
// Returns true if the user confirms (responds with "y" or "yes"), false otherwise.
//
// Example usage:
//
//	if confirmed, err := cli.ConfirmAction("This will delete all files", yesFlag); err != nil {
//	    return err
//	}
//	if !confirmed {
//	    fmt.Println("Operation cancelled")
//	    return nil
//	}
func ConfirmAction(prompt string, skipConfirm bool) (bool, error) {
	if skipConfirm {
		return true, nil
	}

	fmt.Printf("%s\nAre you sure? [y/N]: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes", nil
}
