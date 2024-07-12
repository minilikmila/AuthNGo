package db

import (
	"embed"
	"fmt"

	"github.com/minilikmila/goAuth/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type Database interface {
	DB() *gorm.DB
}

// This method implements Database interface and any method implements Database interface can satisfy database
type GormDatabase struct {
	db *gorm.DB
}

// Method implementation
func (db GormDatabase) DB() *gorm.DB {
	return db.db
}

func generateDSN(config *config.AppConfig) (dsn string) {
	if config.DB.DSN != "" {
		return config.DB.DSN
	}

	dsn = fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=%s`, config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name, config.DB.SSLMode)
	return
}

func InitDatabase(files embed.FS) (*gorm.DB, error) {
	dsn := generateDSN(&config.Config)
	logrus.Infoln("Url : ", dsn)
	// IF we don't need
	// // If using an external migrations directory
	// var fs fs.FS

	// fs = os.DirFS("migrations") // Assuming "migrations" is your directory and replace 'files' in 'fs' in iofs.New(fs, 'migrations)

	// migrations, err := iofs.New(files, "migrations")
	// if err != nil {
	// 	fmt.Println("while embed migration", err)
	// 	os.Exit(1)
	// }

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "auth.",
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Silent),
		// TranslateError: true, // convert db specific errors into its (Gorm) error types.

	})

	if err != nil {
		logrus.Debug("error encountered when try to open database via gorm.", err)
		logrus.Fatal(err)

	}

	sql, err := db.DB()

	if err != nil {
		logrus.Debug("while initializing sql :\n", err)
		logrus.Fatalln(err)

	}

	// Gorm also automatically ping the db to check its availability.
	if err = sql.Ping(); err != nil {
		fmt.Println("database connection problem.", err)
		logrus.Fatalln(err)

	}

	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(20)

	// migration_driver, err := migrate.NewWithSourceInstance("iofs", migrations, dsn)

	// if err != nil {
	// 	logrus.Debugln("error while creating driver.")
	// 	logrus.Fatalln(err)
	// }

	// if err = migration_driver.Up(); err != nil && err != migrate.ErrNoChange {
	// 	logrus.Debugln("while applying migrations.")
	// 	logrus.Debugln(err)
	// }
	migrationPath := "./migrations" // Path to your migrations folder

	if err := applyMigrations(db, migrationPath); err != nil {
		return nil, err
	}

	fmt.Println("DB successfully connected:")

	return db, nil
}
