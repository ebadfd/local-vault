package project

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ebadfd/local-vault/internal/lib"
	"github.com/ebadfd/local-vault/pkg/config"

	"github.com/manifoldco/promptui"
)

func ListProjects() error {
	c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	db, err := lib.GetDb(c.GetFullDbPath())

	if err != nil {
		return err
	}

	defer db.Close()

	projects, err := projects(db)

	if err != nil {
		return err
	}

	project, err := projectSelector(projects)

	if err != nil {
		return err
	}

	apps, err := projectApps(db, project.ID)

	if err != nil {
		return err
	}

	app, err := appsSelector(apps)

	if err != nil {
		return err
	}

	envs, err := projectEnv(db, project.ID, app)

	if err != nil {
		return err
	}

	RenderList(envs, fmt.Sprintf("Availible environments for %s", project.Name))

	return nil
}

func projects(db *sql.DB) ([]Project, error) {
	var projects []Project
	query := `SELECT id, name FROM projects;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return projects, nil
}

func projectApps(db *sql.DB, projectId int) ([]string, error) {
	var apps []string
	query := `SELECT DISTINCT name FROM apps WHERE project_id = ?;`

	rows, err := db.Query(query, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan app name: %w", err)
		}
		apps = append(apps, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return apps, nil
}

func projectEnv(db *sql.DB, projectId int, app string) ([]string, error) {
	var envs []string
	query := `SELECT DISTINCT env FROM apps WHERE project_id = ? and name = ?;`

	rows, err := db.Query(query, projectId, app)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve env names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var env string
		if err := rows.Scan(&env); err != nil {
			return nil, fmt.Errorf("failed to scan env name: %w", err)
		}
		envs = append(envs, env)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return envs, nil
}

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
