package database

import (
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AllModels returns a slice of all models to be migrated
func AllModels() []interface{} {
	return []interface{}{
		&model.User{},
		&model.Log{},
		&model.Verification{},
		&model.BlacklistedToken{},
	}
}

func migrateDB(db *gorm.DB) error {
	// Enable uuid-ossp extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		logrus.Fatalf("Failed to enable uuid-ossp extension: %v", err)
		return err
	}

	if err := db.AutoMigrate(AllModels()...); err != nil {
		logrus.Errorln("Migration error:", err)
		return err
	}
	return nil
}
