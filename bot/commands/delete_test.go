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

func setupTestDBDelete() {
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

	db.AutoMigrate(&Microservice{})
}

func TestDeleteMockExist(t *testing.T) {
	setupTestDBDelete()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "testname_1",
		MicroserviceUrl:     "http://localhost:3007",
		MicroserviceTimeout: 70,
	})

	title, msg := Delete("testname_1")
	title_want := "Delete Command"
	msg_want := "Microservice: testname_1 Has Been Deleted"

	if title_want != title {
		t.Errorf("\n\nError: Delete Failed For Existing Microservice:\nWhat We Wanted: %q\nWhat We Got: %q", title_want, title)
	}
	if msg_want != msg {
		t.Errorf("\n\nError: Delete Failed For Existing Microservice:\nWhat We Wanted: %q\nWhat We Got: %q", msg_want, msg)
	}
}

func TestDeleteBadInput(t *testing.T) {
	setupTestDBDelete()

	title, msg := Delete("adsadadsadssadsaddsa")
	title_want := "Delete Command Error"
	msg_want := "Bot Name Does Not Exist"

	if title_want != title {
		t.Errorf("\n\nError: Bot still trying to delete non existent bot\nWhat We Wanted: %q\nWhat We Got: %q", title_want, title)
	}
	if msg_want != msg {
		t.Errorf("\n\nError: Bot still trying to delete non existent bot\nWhat We Wanted: %q\nWhat We Got: %q", msg_want, msg)
	}
}
