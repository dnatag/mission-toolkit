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

// TestGenerateMarkdown_UsageLineUsesFullName verifies that the Usage line
// uses the full command path (fullName) rather than just the subcommand name (cmd.Use).
// This is a regression test for the bug where Usage showed "m check" instead of "m mission check".
func TestGenerateMarkdown_UsageLineUsesFullName(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	analyzeCmd := &cobra.Command{Use: "analyze", Short: "Analysis tools"}
	clarityCmd := &cobra.Command{Use: "clarify", Short: "Clarification analysis"} // cmd.Use is just "clarify"
	analyzeCmd.AddCommand(clarityCmd)
	rootCmd.AddCommand(analyzeCmd)

	md := GenerateMarkdown(rootCmd)

	// The heading should use fullName: "## `m analyze clarify`"
	if !strings.Contains(md, "## `m analyze clarify`") {
		t.Error("expected analyze clarity heading with full path")
	}

	// The Usage line should also use fullName: "**Usage:** `m analyze clarify`"
	// NOT cmd.Use which would be "**Usage:** `m clarify`"
	if !strings.Contains(md, "**Usage:** `m analyze clarify`") {
		t.Error("expected Usage line to use full command path 'm analyze clarity', not just 'clarify'")
	}

	// Verify the incorrect format (bug) is NOT present
	if strings.Contains(md, "**Usage:** `m clarify`") {
		t.Error("Usage line should NOT contain subcommand-only path 'm clarify'")
	}
}

func TestGenerateSchema(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	initCmd := &cobra.Command{Use: "init", Short: "Initialize"}
	rootCmd.AddCommand(initCmd)

	schema := GenerateSchema(rootCmd)

	if len(schema.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(schema.Commands))
	}
	if schema.Commands[0].Name != "init" {
		t.Errorf("expected command 'init', got '%s'", schema.Commands[0].Name)
	}
	if schema.Commands[0].Short != "Initialize" {
		t.Errorf("expected short 'Initialize', got '%s'", schema.Commands[0].Short)
	}
}

func TestGenerateCondensedMarkdown(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	analyzeCmd := &cobra.Command{Use: "analyze", Short: "Analysis tools"}
	intentCmd := &cobra.Command{Use: "intent", Short: "Analyze intent"}
	intentCmd.Flags().String("input", "", "User input")
	analyzeCmd.AddCommand(intentCmd)
	rootCmd.AddCommand(analyzeCmd)

	md := GenerateCondensedMarkdown(rootCmd)

	if !strings.Contains(md, "# CLI Reference (Condensed)") {
		t.Error("expected condensed header")
	}
	if !strings.Contains(md, "**`m analyze`**") {
		t.Error("expected analyze command")
	}
	if !strings.Contains(md, "`--input`") {
		t.Error("expected --input flag")
	}
}

func TestTraverseCommand_FlagHandling(t *testing.T) {
	cmd := &cobra.Command{Use: "test", Short: "Test"}
	cmd.Flags().StringP("output", "o", "default.txt", "Output file")
	cmd.Flags().Bool("verbose", false, "Verbose mode")

	schema := TraverseCommand(cmd)

	if len(schema.Flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(schema.Flags))
	}

	// Verify flag with shorthand and default value
	outputFlag := schema.Flags[0]
	if outputFlag.Name != "output" {
		t.Errorf("expected flag name 'output', got '%s'", outputFlag.Name)
	}
	if outputFlag.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got '%s'", outputFlag.Shorthand)
	}
	if outputFlag.Default != "default.txt" {
		t.Errorf("expected default 'default.txt', got '%s'", outputFlag.Default)
	}
	if outputFlag.Usage != "Output file" {
		t.Errorf("expected usage 'Output file', got '%s'", outputFlag.Usage)
	}

	// Verify flag without shorthand
	verboseFlag := schema.Flags[1]
	if verboseFlag.Name != "verbose" {
		t.Errorf("expected flag name 'verbose', got '%s'", verboseFlag.Name)
	}
	if verboseFlag.Shorthand != "" {
		t.Errorf("expected empty shorthand, got '%s'", verboseFlag.Shorthand)
	}
	if verboseFlag.Default != "false" {
		t.Errorf("expected default 'false', got '%s'", verboseFlag.Default)
	}
}

func TestGenerateCondensedMarkdown_FlagFormatting(t *testing.T) {
	rootCmd := &cobra.Command{Use: "m", Short: "Test CLI"}
	checkCmd := &cobra.Command{Use: "check", Short: "Check input"}
	checkCmd.Flags().StringP("input", "i", "", "Input text")
	checkCmd.Flags().Bool("strict", false, "Strict mode")
	rootCmd.AddCommand(checkCmd)

	md := GenerateCondensedMarkdown(rootCmd)

	if !strings.Contains(md, "`-i/--input`") {
		t.Error("expected shorthand flag format '-i/--input'")
	}
	if !strings.Contains(md, "`--strict`") {
		t.Error("expected long-only flag format '--strict'")
	}
}

func TestTraverseCommands_EmptyList(t *testing.T) {
	schemas := TraverseCommands([]*cobra.Command{})

	if len(schemas) != 0 {
		t.Errorf("expected 0 commands, got %d", len(schemas))
	}
}
