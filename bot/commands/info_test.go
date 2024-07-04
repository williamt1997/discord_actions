package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	dbconfig "github.com/nathanjcook/discordbotgo/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Microservice struct {
	MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
	MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
	MicroserviceUrl     string `gorm:"column:microservice_url;"`
	MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
}

func setupTestDBInfo() {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load(".env")
		if err != nil {
			zap.L().Panic("Error loading .env file:", zap.Error(err))
		}
	}
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		dbname,
		port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	dbconfig.DB = db

	err = db.AutoMigrate(&Microservice{})
	if err != nil {
		zap.L().Panic("Error Migrating To Database .env file:", zap.Error(err))
	}
}

func TestInfoMS(t *testing.T) {
	setupTestDBInfo()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "tester_test",
		MicroserviceUrl:     "http://localhost:3007",
		MicroserviceTimeout: 70,
	})

	title, msg := Info()
	titleDontWant := "Info Command Null"
	msgDontWant := "No Microservices Available"

	if titleDontWant == title {
		t.Errorf("\n\nError: Info Failing To Get Microservice Data:\nWhat We Wanted: Info Command\nWhat We Got: %q", title)
	} else if msgDontWant == msg {
		t.Errorf("\n\nError: Info Failing To Get Microservice Data:\nWhat We Wanted: All Microservices\nWhat We Got: %q", msg)
	} else {
		Delete("tester_test")
	}
}
