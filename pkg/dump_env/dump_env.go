package dumpenv

import (
	"database/sql"
	"fmt"
	"os/exec"

	"github.com/ebadfd/local-vault/internal/lib"
	"github.com/ebadfd/local-vault/pkg/config"
)

func CreateDump(projectName, appName, env, file string) error {
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

	encryptFile, err := getEncryptedFile(db, projectId, appName, env)

	if err != nil {
		return err
	}

	err = decrypt(encryptFile, file, *c)

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

func getEncryptedFile(db *sql.DB, projectId int, appName, env string) (string, error) {
	query := `SELECT encrypted_file FROM apps WHERE project_id = ? AND name = ? AND env = ? LIMIT 1;`
	var encryptedFile string

	err := db.QueryRow(query, projectId, appName, env).Scan(&encryptedFile)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no app found with the specified projectId, appName, and env")
		}
		return "", fmt.Errorf("failed to query encrypted file: %w", err)
	}

	return encryptedFile, nil
}

func decrypt(encrypted_file, file string, config config.Config) error {
	cmd := exec.Command("gpg", "--yes", "--decrypt", "--output", file, "--recipient", *config.Recipient, encrypted_file)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}

	return nil
}
