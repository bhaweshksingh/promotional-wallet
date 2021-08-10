package repository

import (
	"account/pkg/config"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	rollBackStep = -1
	cutSet       = "file://"
	databaseName = "mysql"
)

func RunMigrations(configFile string) {
	newMigrate, err := newMigrate(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := newMigrate.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		fmt.Println(err)
		return
	}
}

func RollBackMigrations(configFile string) {
	newMigrate, err := newMigrate(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := newMigrate.Steps(rollBackStep); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
	}
}

func newMigrate(configFile string) (*migrate.Migrate, error) {
	cfg := config.NewConfig(configFile)
	dbConfig := cfg.GetDBConfig()

	dbHandler := NewDBHandler(dbConfig)

	gormDB, err := dbHandler.GetDB()
	if err != nil {
		return nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	sourcePath, err := getSourcePath(dbConfig.MigrationPath())
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(sourcePath, databaseName, driver)
}

func getSourcePath(directory string) (string, error) {
	directory = strings.TrimLeft(directory, cutSet)

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", cutSet, absPath), nil
}
