package cmd

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const (
	importName         = "import"
	importShort        = "Import an existing environment file into the vault."
	importNumberOfArgs = 1 // file name arg
)

type importCommand struct {
	// flags
	Project string `validate:"required"`
	App     string `validate:"required"`
	Env     string `validate:"required"`
}

func defaultImportOptions() *importCommand {
	return &importCommand{}
}

func newImportCommand() *cobra.Command {
	c := defaultImportOptions()

	cmd := &cobra.Command{
		Use:   c.Name(),
		Short: c.Short(),
		Args:  cobra.ExactArgs(importNumberOfArgs),
		RunE:  c.run,
	}

	cmd.Flags().StringVarP(&c.Project, "project", "p", c.Project, "project name")
	cmd.Flags().StringVarP(&c.App, "app", "a", c.App, "application name")
	cmd.Flags().StringVarP(&c.Env, "env", "e", c.Env, "envrionment name")

	return cmd
}

func (o *importCommand) run(cmd *cobra.Command, args []string) error {
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

	fmt.Println(file)
	fmt.Println(o)

	return nil
}

func (o *importCommand) parseArgs(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("file argument is required")
	}

	return args[0], nil
}

func (c *importCommand) Name() string {
	return importName
}

func (c *importCommand) Short() string {
	return importShort
}
