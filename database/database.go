package database

import (
	"bookstore-api/models"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"fmt"
)

var DB *gorm.DB
func ConnectAndMigrate(dsn string) *gorm.DB {
	var err error
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return nil
	}
	fmt.Println("Database connection established successfully")
	DB.AutoMigrate(&models.Book{}, &models.User{})
	fmt.Println("Database migration completed successfully")
	return DB
}	