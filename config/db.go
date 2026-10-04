package config

import (
	"fmt"

	"learn-gin-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	if App == nil {
		LoadConfig()
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		App.DBHost, App.DBUser, App.DBPass, App.DBName, App.DBPort, App.DBSSL,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB_CONNECTION_ERROR: " + err.Error())
	}

	// Migrate the schema
	DB.AutoMigrate(
		&models.User{},
		&models.Article{},
	)
}
