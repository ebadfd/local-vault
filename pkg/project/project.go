package project

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ebadfd/local-vault/internal/lib"
	"github.com/ebadfd/local-vault/pkg/config"

	"github.com/manifoldco/promptui"
)

func CreateProject() error {
	c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	db, err := lib.GetDb(c.GetFullDbPath())

	if err != nil {
		return err
	}

	defer db.Close()

	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("project name must have more than 3 characters")
		}

		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM projects WHERE name = ? LIMIT 1);`
		err := db.QueryRow(query, input).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if project exists: %w", err)
		}

		if exists {
			return errors.New("project with this name already exists")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Project Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	err = createProject(db, result)
	return err
}

func createProject(db *sql.DB, name string) error {
	query := `INSERT INTO projects (name) VALUES (?);`
	_, err := db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}
