package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FlagSchema represents a command flag
type FlagSchema struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Usage       string `json:"usage"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required,omitempty"`
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

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate CLI documentation schema",
	Long:  `Generate a JSON schema of all CLI commands for documentation purposes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		
		if format != "json" {
			return fmt.Errorf("unsupported format: %s (only 'json' is supported)", format)
		}

		schema := CLISchema{
			Commands: traverseCommands(rootCmd.Commands()),
		}

		output, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling schema: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// traverseCommands recursively extracts schema from commands
func traverseCommands(cmds []*cobra.Command) []CommandSchema {
	var schemas []CommandSchema
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		schemas = append(schemas, traverseCommand(cmd))
	}
	return schemas
}

// traverseCommand extracts schema from a single command
func traverseCommand(cmd *cobra.Command) CommandSchema {
	schema := CommandSchema{
		Name:  cmd.Name(),
		Use:   cmd.Use,
		Short: cmd.Short,
		Long:  cmd.Long,
	}

	// Extract flags
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

	// Recursively process subcommands
	if cmd.HasSubCommands() {
		schema.Subcommands = traverseCommands(cmd.Commands())
	}

	return schema
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().String("format", "json", "Output format (json)")
}
