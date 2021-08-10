package config

import "fmt"

type DBConfig struct {
	host          string
	port          int
	name          string
	user          string
	password      string
	migrationPath string
}

func (db *DBConfig) Address() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", db.host, db.user, db.password, db.name, db.port)
}
func (db *DBConfig) MigrationPath() string {
	return db.migrationPath
}

func newDBConfig() DBConfig {
	return DBConfig{
		host:          getString("DB_HOST"),
		port:          getInt("DB_PORT"),
		name:          getString("DB_NAME"),
		user:          getString("DB_USER"),
		password:      getString("DB_PASSWORD"),
		migrationPath: getString("MIGRATION_PATH"),
	}
}
