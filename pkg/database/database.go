package database

import (
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"github.com/Adedunmol/zephyr/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	if helpers.EnvConfig.Environment == "test" {
		DB, err = gorm.Open(postgres.Open(helpers.EnvConfig.TestDatabaseUrl), &gorm.Config{TranslateError: true})
	} else {
		DB, err = gorm.Open(postgres.Open(helpers.EnvConfig.DatabaseUrl), &gorm.Config{TranslateError: true})
	}

	if err != nil {
		helpers.Error.Fatal("error connecting to db: %w", err)
	}

	if helpers.EnvConfig.Environment != "test" {
		DB.Logger = logger.Default.LogMode(logger.Info)

		helpers.Info.Println("Running migrations")
	}

	DB.AutoMigrate(&models.User{})
}
