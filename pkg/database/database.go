package database

import (
	"github.com/Adedunmol/zephyr/pkg/config"
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	var err error

	if config.EnvConfig.Environment == "test" {
		DB, err = gorm.Open(postgres.Open(config.EnvConfig.TestDatabaseUrl), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(config.EnvConfig.DatabaseUrl), &gorm.Config{})
	}

	if err != nil {
		helpers.Error.Fatal("error connecting to db: %w", err)
	}

	if config.EnvConfig.Environment != "test" {
		DB.Logger = logger.Default.LogMode(logger.Info)

		helpers.Info.Println("Running migrations")
	}

	DB.AutoMigrate()
}
