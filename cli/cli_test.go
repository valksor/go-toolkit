package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestSetupRootCmd(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	SetupRootCmd(cmd, RootOptions{
		ToolName: "test",
		VersionInfo: func() string {
			return "test version 1.0.0"
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	})

	if !cmd.SilenceErrors {
		t.Error("SilenceErrors not set")
	}
	if !cmd.SilenceUsage {
		t.Error("SilenceUsage not set")
	}

	// Check if version command was added
	versionCmd := cmd.Commands()[0]
	if versionCmd == nil || versionCmd.Use != "version" {
		t.Error("Version command not added")
	}
}

func TestAddStandardFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	flags := AddStandardFlags(cmd)

	verbose, _ := flags.GetBool("verbose")
	quiet, _ := flags.GetBool("quiet")
	noColor, _ := flags.GetBool("no-color")

	if verbose {
		t.Error("verbose should be false by default")
	}
	if quiet {
		t.Error("quiet should be false by default")
	}
	if noColor {
		t.Error("no-color should be false by default")
	}
}

func TestGetStandardFlags(t *testing.T) {
	flags := GetStandardFlags()

	if flags == nil {
		t.Fatal("GetStandardFlags() returned nil")
	}

	if flags.Verbose {
		t.Error("Verbose should be false initially")
	}
}

func TestDefaultPersistentPreRun(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	globalFlags.Verbose = true

	preRun := defaultPersistentPreRun("test", nil)
	err := preRun(cmd, nil)
	if err != nil {
		t.Errorf("defaultPersistentPreRun() error = %v", err)
	}
}

func TestAddCommandGroup(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	AddCommandGroup(cmd, "workflow", "Workflow Commands:")

	groups := cmd.Groups()
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}

	if groups[0].ID != "workflow" {
		t.Errorf("Group ID = %v, want 'workflow'", groups[0].ID)
	}

	if groups[0].Title != "Workflow Commands:" {
		t.Errorf("Group Title = %v, want 'Workflow Commands:'", groups[0].Title)
	}

	// Adding the same group again should be idempotent
	AddCommandGroup(cmd, "workflow", "Workflow Commands:")
	if len(cmd.Groups()) != 1 {
		t.Errorf("Adding duplicate group should not create a new group")
	}
}

func TestSetCommandGroup(t *testing.T) {
	subCmd := &cobra.Command{Use: "status"}

	SetCommandGroup(subCmd, "workflow")

	if subCmd.GroupID != "workflow" {
		t.Errorf("GroupID = %v, want 'workflow'", subCmd.GroupID)
	}
}

func TestAddCommandGroups(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}

	statusCmd := &cobra.Command{Use: "status"}
	listCmd := &cobra.Command{Use: "list"}

	AddCommandGroups(rootCmd,
		CommandGroup{
			ID:       "task",
			Title:    "Task Commands:",
			Commands: []*cobra.Command{statusCmd, listCmd},
		},
	)

	// Check groups
	groups := rootCmd.Groups()
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}

	// Check commands were added
	commands := rootCmd.Commands()
	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Check commands have group IDs
	for _, cmd := range commands {
		if cmd.GroupID != "task" {
			t.Errorf("Command %s GroupID = %v, want 'task'", cmd.Use, cmd.GroupID)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	SetupRootCmd(cmd, RootOptions{
		ToolName: "test",
		VersionInfo: func() string {
			return "test v1.0.0\n  Commit: abc123\n  Built:  2024-01-15\n  Go:     go1.21.5"
		},
	})

	// Find version command
	var versionCmd *cobra.Command
	for _, c := range cmd.Commands() {
		if c.Use == "version" {
			versionCmd = c

			break
		}
	}

	if versionCmd == nil {
		t.Fatal("Version command not found")
	}

	// Execute version command
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)
	versionCmd.SetArgs([]string{})

	if versionCmd.Run != nil {
		versionCmd.Run(versionCmd, []string{})
	} else if versionCmd.RunE != nil {
		err := versionCmd.RunE(versionCmd, []string{})
		if err != nil {
			t.Errorf("version command error = %v", err)
		}
	}

	output := buf.String()
	if !strings.Contains(output, "test v1.0.0") {
		t.Errorf("Version output doesn't contain version info, got: %s", output)
	}
}
