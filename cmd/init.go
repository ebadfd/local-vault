package cmd

import (
	"fmt"

	"github.com/ebadfd/local-vault/internal/lib"
	"github.com/ebadfd/local-vault/pkg/config"
	"github.com/spf13/cobra"
)

const (
	initName  = "init"
	initShort = "init the config"
)

type initCommand struct {
}

func defaultinitCommandOptions() *initCommand {
	return &initCommand{}
}

func newInitCommand() *cobra.Command {
	c := defaultinitCommandOptions()

	cmd := &cobra.Command{
		Use:   c.Name(),
		Short: c.Short(),
		RunE:  c.run,
	}

	return cmd
}

func (o *initCommand) run(cmd *cobra.Command, args []string) error {
	c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	db, err := lib.GetDb(c.GetFullDbPath())

	if err != nil {
		return err
	}

	defer db.Close()

	if err := lib.InitDbSchema(db); err != nil {
		return fmt.Errorf("Error initializing schema: %v", err)
	}

	return nil
}

func (c *initCommand) Name() string {
	return initName
}

func (c *initCommand) Short() string {
	return initShort
}
