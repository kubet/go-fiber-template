package database

import (
	"fmt"
	"log"
	"os"

	"template/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	//   // Set the maximum number of open connections.
	//   sqlDB.SetMaxOpenConns(50)

	//   // Set the maximum number of idle connections.
	//   sqlDB.SetMaxIdleConns(10)

	//   // Set the maximum lifetime of a connection.
	//   sqlDB.SetConnMaxLifetime(time.Minute * 5)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.User{})

	DB = Dbinstance{
		Db: db,
	}
}
