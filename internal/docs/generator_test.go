package docs

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestTraverseCommand(t *testing.T) {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "A test command for testing",
	}
	testCmd.Flags().String("name", "default", "A test flag")
	testCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	schema := TraverseCommand(testCmd)

	if schema.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", schema.Name)
	}
	if schema.Short != "Test command" {
		t.Errorf("expected short 'Test command', got '%s'", schema.Short)
	}
	if len(schema.Flags) != 2 {
		t.Errorf("expected 2 flags, got %d", len(schema.Flags))
	}
}

func TestTraverseCommands_HiddenCommands(t *testing.T) {
	visibleCmd := &cobra.Command{Use: "visible", Short: "Visible"}
	hiddenCmd := &cobra.Command{Use: "hidden", Short: "Hidden", Hidden: true}

	schemas := TraverseCommands([]*cobra.Command{visibleCmd, hiddenCmd})

	if len(schemas) != 1 {
		t.Errorf("expected 1 command (hidden excluded), got %d", len(schemas))
	}
	if schemas[0].Name != "visible" {
		t.Errorf("expected 'visible', got '%s'", schemas[0].Name)
	}
}

func TestTraverseCommand_WithSubcommands(t *testing.T) {
	parentCmd := &cobra.Command{Use: "parent", Short: "Parent"}
	childCmd := &cobra.Command{Use: "child", Short: "Child"}
	parentCmd.AddCommand(childCmd)

	schema := TraverseCommand(parentCmd)

	if len(schema.Subcommands) != 1 {
		t.Errorf("expected 1 subcommand, got %d", len(schema.Subcommands))
	}
	if schema.Subcommands[0].Name != "child" {
		t.Errorf("expected subcommand 'child', got '%s'", schema.Subcommands[0].Name)
	}
}

func TestGenerateMarkdown(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	initCmd := &cobra.Command{Use: "init", Short: "Initialize project"}
	initCmd.Flags().String("ai", "", "AI type")
	rootCmd.AddCommand(initCmd)

	md := GenerateMarkdown(rootCmd)

	if !strings.Contains(md, "# Mission Toolkit CLI Reference") {
		t.Error("expected header in markdown")
	}
	if !strings.Contains(md, "AUTO-GENERATED") {
		t.Error("expected auto-generated comment")
	}
	if !strings.Contains(md, "## `m init`") {
		t.Error("expected init command section")
	}
	if !strings.Contains(md, "--ai") {
		t.Error("expected --ai flag")
	}
}

func TestGenerateMarkdown_WithSubcommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	missionCmd := &cobra.Command{Use: "mission", Short: "Mission management"}
	checkCmd := &cobra.Command{Use: "check", Short: "Check mission"}
	missionCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(missionCmd)

	md := GenerateMarkdown(rootCmd)

	if !strings.Contains(md, "## `m mission`") {
		t.Error("expected mission command section")
	}
	if !strings.Contains(md, "## `m mission check`") {
		t.Error("expected mission check subcommand section")
	}
}
