package database_

import (
	"fmt"
	"os"

	"github.com/minilikmila/standard-auth-go/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database interface {
	DB() *gorm.DB
	// Close() error
}

type GormDatabase struct {
	db *gorm.DB
}

func (d *GormDatabase) DB() *gorm.DB {
	return d.db
}

func InitDatabase(config *config.Config) (*GormDatabase, error) {
	fmt.Println("DB- URL : ", config.DatabaseUri)
	db, err := gorm.Open(postgres.Open(config.DatabaseUri), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "genie.",
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Silent),
		// TranslateError: true, // convert db specific errors into its (Gorm) error types.

	})

	if err != nil {
		logrus.Debug("error encountered when try to open database via gorm : ", err)
		logrus.Fatal(err)

	}

	sql, err := db.DB()

	if err != nil {
		logrus.Debug("while initializing sql :\n", err)
		logrus.Fatalln(err)

	}

	// Gorm also automatically ping the db to check its availability.
	if err = sql.Ping(); err != nil {
		fmt.Println("database connection problem : ", err)
		logrus.Fatalln(err)

	}

	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(20)

	if err = migrateDB(db); err != nil {
		os.Exit(1)
	}

	fmt.Println("DB successfully connected:")

	return &GormDatabase{db: db}, nil
}
