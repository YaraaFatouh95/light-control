package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"light-control/models"
)

var DB *gorm.DB

func ConnectDB() {
	connStr := "user=database password=databasepassword dbname=light_control sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = DB.AutoMigrate(
		&models.City{},
		&models.Zone{},
		&models.Luminaire{},
		&models.Command{})

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database connected and migrated successfully")
}
