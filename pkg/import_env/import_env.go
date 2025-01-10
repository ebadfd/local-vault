package importenv

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ebadfd/local-vault/internal/lib"
	"github.com/ebadfd/local-vault/pkg/config"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

func CreateImport(projectName, appName, env, file string) error {
	c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	db, err := lib.GetDb(c.GetFullDbPath())

	if err != nil {
		return err
	}

	defer db.Close()

	projectId, err := project(db, projectName)

	if err != nil {
		return err
	}

	encryptFile, err := encrypt(file, projectName, appName, env, *c)

	if err != nil {
		return err
	}

	err = create(db, projectId, appName, env, encryptFile)

	if err != nil {
		return err
	}

	return nil
}

func project(db *sql.DB, name string) (int, error) {
	var id int
	query := `SELECT id FROM projects WHERE name = ? LIMIT 1;`
	err := db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("project with name '%s' not found", name)
		}
		return 0, fmt.Errorf("failed to query project: %w", err)
	}
	return id, nil
}

func create(db *sql.DB, projectId int, appName, env, encryptedFile string) error {
	var existingID int
	var existingEncryptedFile string

	query := `SELECT id, encrypted_file FROM apps WHERE project_id = ? AND env = ? AND name = ?`
	err := db.QueryRow(query, projectId, env, appName).Scan(&existingID, &existingEncryptedFile)

	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO apps (name, project_id, env, encrypted_file) VALUES (?, ?, ?, ?);`
		_, err := db.Exec(insertQuery, appName, projectId, env, encryptedFile)
		if err != nil {
			return fmt.Errorf("failed to create app: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking for existing app: %w", err)
	}

	prompt := promptui.Prompt{
		Label:     "Overwrite resource",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	if strings.ToLower(result) == "y" {
		archiveQuery := `INSERT INTO archive (name, project_id, env, encrypted_file) 
                         VALUES (?, ?, ?, ?)`
		_, err = db.Exec(archiveQuery, appName, projectId, env, existingEncryptedFile)
		if err != nil {
			return fmt.Errorf("failed to archive app: %w", err)
		}

		updateQuery := `UPDATE apps SET encrypted_file = ? WHERE project_id = ? AND env = ? AND name = ?`
		_, err = db.Exec(updateQuery, encryptedFile, projectId, env, appName)
		if err != nil {
			return fmt.Errorf("failed to overwrite app: %w", err)
		}
	}
	return nil
}

func encrypt(file, project, app, env string, config config.Config) (string, error) {
	id := uuid.New()

	outputFile := filepath.Join(config.GetFullDataPath(), fmt.Sprintf("%s-%s-%s-%s_%s", project, app, env, id, file)) + ".gpg"

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %v", file)
	}

	cmd := exec.Command("gpg", "--yes", "--encrypt", "--recipient", *config.Recipient, "--output", outputFile, file)

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file: %w", err)
	}

	return outputFile, nil
}
