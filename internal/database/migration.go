package database_

import (
	"github.com/minilikmila/standard-auth-go/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func migrateDB(db *gorm.DB) error {
	// migration
	if err := db.AutoMigrate(&model.User{}, &model.Log{}); err != nil {
		logrus.Errorln("Migration error : - ", err)
	}
	return nil
}
