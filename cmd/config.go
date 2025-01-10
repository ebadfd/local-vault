package cmd

import (
	"fmt"

	"github.com/ebadfd/local-vault/pkg/config"
	"github.com/spf13/cobra"
)

const (
	configName  = "config"
	configShort = "view project config"
)

type configCommand struct {
}

func defaultconfigCommandOptions() *configCommand {
	return &configCommand{}
}

func newConfigCommand() *cobra.Command {
	c := defaultconfigCommandOptions()

	cmd := &cobra.Command{
		Use:   c.Name(),
		Short: c.Short(),
		RunE:  c.run,
	}

	return cmd
}

func (o *configCommand) run(cmd *cobra.Command, args []string) error {
	c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	fmt.Println(c)

	return nil
}

func (c *configCommand) Name() string {
	return configName
}

func (c *configCommand) Short() string {
	return configShort
}
