package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func applyMigrations(db *gorm.DB, migrationPath string) error {
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationBytes, err := ioutil.ReadFile(filepath.Join(migrationPath, file.Name()))
			if err != nil {
				return err
			}

			migrationSQL := string(migrationBytes)
			if err := db.Exec(migrationSQL).Error; err != nil {
				return err
			}

			fmt.Printf("Applied migration: %s\n", file.Name())
		}
	}

	return nil
}