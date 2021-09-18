package database

import (
	"ambassador/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	var err error

	DB, err = gorm.Open(postgres.Open("postgres://postgres:postgres@db:5432/postgres?sslmode=disable"), &gorm.Config{})

	if err != nil {
		panic("Could not connect with database!")
	}

	return DB
}

func AutoMigrate(DB *gorm.DB) {
	DB.AutoMigrate(models.User{})
}
