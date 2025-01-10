package cmd

import (
	"fmt"

	dumpenv "github.com/ebadfd/local-vault/pkg/dump_env"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const (
	dumpName         = "dump"
	dumpShort        = "Dump environment variables from the vault to a file."
	dumpNumberOfArgs = 1 // file name arg
)

type dumpCommand struct {
	// flags
	Project string `validate:"required"`
	App     string `validate:"required"`
	Env     string `validate:"required"`
}

func defaultDumpCommandOptions() *dumpCommand {
	return &dumpCommand{}
}

func newDumpCommand() *cobra.Command {
	c := defaultDumpCommandOptions()

	cmd := &cobra.Command{
		Use:   c.Name(),
		Short: c.Short(),
		Args:  cobra.ExactArgs(dumpNumberOfArgs),
		RunE:  c.run,
	}

	cmd.Flags().StringVarP(&c.Project, "project", "p", c.Project, "project name")
	cmd.Flags().StringVarP(&c.App, "app", "a", c.App, "application name")
	cmd.Flags().StringVarP(&c.Env, "env", "e", c.Env, "envrionment name")

	return cmd
}

func (o *dumpCommand) run(cmd *cobra.Command, args []string) error {
	file, err := o.parseArgs(args)
	if err != nil {
		return err
	}

	err = validate.Struct(o)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		fmt.Println("Validation errors:")
		for _, ve := range err.(validator.ValidationErrors) {
			fmt.Printf("  - Field: %s\n", ve.Field())
			fmt.Printf("    Issue: Validation '%s' failed\n", ve.Tag())
			if ve.Param() != "" {
				fmt.Printf("    Expected: %s\n", ve.Param())
			}
			fmt.Printf("    Actual: %v\n\n", ve.Value())
		}

		return fmt.Errorf("validation failed: %d errors found", len(err.(validator.ValidationErrors)))
	}

	err = dumpenv.CreateDump(o.Project, o.App, o.Env, file)

	if err != nil {
		return err
	}

	return nil
}

func (o *dumpCommand) parseArgs(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("file argument is required")
	}

	return args[0], nil
}

func (c *dumpCommand) Name() string {
	return dumpName
}

func (c *dumpCommand) Short() string {
	return dumpShort
}
