package repository

import (
	"account/pkg/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBHandler interface {
	GetDB() (*gorm.DB, error)
}

type gormDBHandler struct {
	config config.DBConfig
}

func (dbHandler *gormDBHandler) GetDB() (*gorm.DB, error) {
	fmt.Println("DB Connection String:"+ dbHandler.config.Address())

	db, err := gorm.Open(postgres.Open(dbHandler.config.Address()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping. Error %v", err)
	}

	return db, nil
}

func NewDBHandler(config config.DBConfig) DBHandler {
	return &gormDBHandler{
		config: config,
	}
}
