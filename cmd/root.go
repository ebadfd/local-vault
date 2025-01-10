package cmd

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func newRootCmd(version string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "local-vault",
		Short: "A CLI tool to manage environment variables in a local encrypted vault.",
		Long: `local-vault is a CLI tool that stores environment variables in an encrypted vault.
It supports importing and exporting environment variables for projects and applications.`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newImportCommand())
	cmd.AddCommand(newDumpCommand())
	cmd.AddCommand(newProjectCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newInitCommand())

	return cmd
}

func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}
	return nil
}
