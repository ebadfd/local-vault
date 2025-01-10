package cmd

import (
	"fmt"

	"github.com/ebadfd/local-vault/pkg/project"
	"github.com/spf13/cobra"
)

const (
	projectName  = "project"
	projectShort = "Add new project or list all projects."
)

type projectCommand struct {
	// flags
	NewProject bool
}

func defaultProjectCommandOptions() *projectCommand {
	return &projectCommand{
		NewProject: false,
	}
}

func newProjectCommand() *cobra.Command {
	c := defaultProjectCommandOptions()

	cmd := &cobra.Command{
		Use:   c.Name(),
		Short: c.Short(),
		RunE:  c.run,
	}

	cmd.Flags().BoolVarP(&c.NewProject, "new-project", "n", c.NewProject, "Add new project")

	return cmd
}

func (o *projectCommand) run(cmd *cobra.Command, args []string) error {
	if o.NewProject {
		err := project.CreateProject()

		return err
	} else {
		fmt.Println("list all projects")
	}

	return nil
}

func (c *projectCommand) Name() string {
	return projectName
}

func (c *projectCommand) Short() string {
	return projectShort
}
