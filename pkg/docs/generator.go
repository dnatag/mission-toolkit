package docs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FlagSchema represents a command flag
type FlagSchema struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Usage     string `json:"usage"`
	Type      string `json:"type"`
	Default   string `json:"default,omitempty"`
}

// CommandSchema represents a CLI command structure
type CommandSchema struct {
	Name        string          `json:"name"`
	Use         string          `json:"use"`
	Short       string          `json:"short"`
	Long        string          `json:"long,omitempty"`
	Flags       []FlagSchema    `json:"flags,omitempty"`
	Subcommands []CommandSchema `json:"subcommands,omitempty"`
}

// CLISchema represents the complete CLI structure
type CLISchema struct {
	Commands []CommandSchema `json:"commands"`
}

// TraverseCommands recursively extracts schema from commands
func TraverseCommands(cmds []*cobra.Command) []CommandSchema {
	var schemas []CommandSchema
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		schemas = append(schemas, TraverseCommand(cmd))
	}
	return schemas
}

// TraverseCommand extracts schema from a single command
func TraverseCommand(cmd *cobra.Command) CommandSchema {
	schema := CommandSchema{
		Name:  cmd.Name(),
		Use:   cmd.Use,
		Short: cmd.Short,
		Long:  cmd.Long,
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flag := FlagSchema{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Usage:     f.Usage,
			Type:      f.Value.Type(),
			Default:   f.DefValue,
		}
		schema.Flags = append(schema.Flags, flag)
	})

	if cmd.HasSubCommands() {
		schema.Subcommands = TraverseCommands(cmd.Commands())
	}

	return schema
}

// GenerateSchema creates a CLISchema from a root command
func GenerateSchema(rootCmd *cobra.Command) CLISchema {
	return CLISchema{
		Commands: TraverseCommands(rootCmd.Commands()),
	}
}

// GenerateMarkdown generates cli-reference.md content from a root command
func GenerateMarkdown(rootCmd *cobra.Command) string {
	var sb strings.Builder

	sb.WriteString("# Mission Toolkit CLI Reference\n\n")
	sb.WriteString("<!-- AUTO-GENERATED: Do not edit manually. Run 'm init' to regenerate. -->\n\n")

	schema := GenerateSchema(rootCmd)

	for _, cmd := range schema.Commands {
		writeCommandMarkdown(&sb, cmd, "")
	}

	return sb.String()
}

func writeCommandMarkdown(sb *strings.Builder, cmd CommandSchema, prefix string) {
	fullName := prefix + cmd.Name
	if prefix != "" {
		fullName = prefix + " " + cmd.Name
	}

	sb.WriteString(fmt.Sprintf("## `m %s`\n\n", fullName))
	sb.WriteString(fmt.Sprintf("%s\n\n", cmd.Short))

	if cmd.Long != "" && cmd.Long != cmd.Short {
		sb.WriteString(fmt.Sprintf("%s\n\n", cmd.Long))
	}

	sb.WriteString(fmt.Sprintf("**Usage:** `m %s`\n\n", fullName))

	if len(cmd.Flags) > 0 {
		sb.WriteString("**Flags:**\n")
		for _, flag := range cmd.Flags {
			if flag.Shorthand != "" {
				sb.WriteString(fmt.Sprintf("- `-%s, --%s`: %s", flag.Shorthand, flag.Name, flag.Usage))
			} else {
				sb.WriteString(fmt.Sprintf("- `--%s`: %s", flag.Name, flag.Usage))
			}
			if flag.Default != "" && flag.Default != "false" && flag.Default != "[]" {
				sb.WriteString(fmt.Sprintf(" (default: `%s`)", flag.Default))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	for _, sub := range cmd.Subcommands {
		writeCommandMarkdown(sb, sub, fullName)
	}
}

// GenerateCondensedMarkdown generates a condensed CLI reference (50-100 lines)
func GenerateCondensedMarkdown(rootCmd *cobra.Command) string {
	var sb strings.Builder

	sb.WriteString("# CLI Reference (Condensed)\n\n")

	// Core commands only
	coreCommands := []string{"analyze", "mission", "diagnosis", "backlog", "checkpoint", "log", "check"}

	schema := GenerateSchema(rootCmd)

	for _, cmdName := range coreCommands {
		for _, cmd := range schema.Commands {
			if cmd.Name == cmdName {
				writeCondensedCommand(&sb, cmd, "")
				break
			}
		}
	}

	return sb.String()
}

func writeCondensedCommand(sb *strings.Builder, cmd CommandSchema, prefix string) {
	fullName := prefix + cmd.Name
	if prefix != "" {
		fullName = prefix + " " + cmd.Name
	}

	sb.WriteString(fmt.Sprintf("**`m %s`** - %s\n", fullName, cmd.Short))

	// Only show flags for leaf commands (no subcommands)
	if len(cmd.Subcommands) == 0 && len(cmd.Flags) > 0 {
		sb.WriteString("  Flags: ")
		flagNames := []string{}
		for _, flag := range cmd.Flags {
			if flag.Shorthand != "" {
				flagNames = append(flagNames, fmt.Sprintf("`-%s/--%s`", flag.Shorthand, flag.Name))
			} else {
				flagNames = append(flagNames, fmt.Sprintf("`--%s`", flag.Name))
			}
		}
		sb.WriteString(strings.Join(flagNames, ", "))
		sb.WriteString("\n")
	}

	// Recurse for subcommands
	for _, sub := range cmd.Subcommands {
		writeCondensedCommand(sb, sub, fullName)
	}

	sb.WriteString("\n")
}
