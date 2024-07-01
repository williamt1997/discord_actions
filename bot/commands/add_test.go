package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
	dbconfig "github.com/nathanjcook/discordbotgo/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDBAdd() {
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

func TestAddMSNameAlreadyExists(t *testing.T) {
	setupTestDBAdd()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "existing_service",
		MicroserviceUrl:     "http://localhost:3007",
		MicroserviceTimeout: 70,
	})

	title, msg := Add("existing_service", "http://localhost:8081", "50")
	title_want := "Add Command Error"
	msg_want := "Microservice Name AND Microservice URL Must Be Unique"

	if title_want != title {
		t.Errorf("\n\nError: Failed To Prevent User From Adding A Microservice Name That Already Exists:\nWhat We Wanted: %q\nWhat We Got: %q", title_want, title)
		Delete("existing_service")
	} else if msg_want != msg {
		t.Errorf("\n\nError: Failed To Prevent User From Adding A Microservice Name That Already Exists:\nWhat We Wanted: %q\nWhat We Got: %q", msg_want, msg)
		Delete("existing_service")
	} else {
		Delete("existing_service")
	}
}

func TestAddMSHostURLAlreadyExists(t *testing.T) {
	setupTestDBAdd()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "existing_service",
		MicroserviceUrl:     "http://localhost:8081",
		MicroserviceTimeout: 70,
	})

	title, msg := Add("new_service", "http://localhost:8081", "50")
	title_want := "Add Command Error"
	msg_want := "Microservice Name AND Microservice URL Must Be Unique"

	if title_want != title {
		t.Errorf("\n\nError: Failed To Prevent User From Adding A Microservice Host URL That Already Exists:\nWhat We Wanted: %q\nWhat We Got: %q", title_want, title)
		Delete("existing_service")
	} else if msg_want != msg {
		t.Errorf("\n\nError: Failed To Prevent User From Adding A Microservice Host URL That Already Exists:\nWhat We Wanted: %q\nWhat We Got: %q", msg_want, msg)
		Delete("existing_service")
	} else {
		Delete("existing_service")
	}
}

func TestAddSuccess(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:8081/api/help")

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "testname_5",
		MicroserviceUrl:     "http://localhost:3007",
		MicroserviceTimeout: 70,
	})

	title, msg := Add("New_service_test", "http://localhost:8081", "50")
	title_want := "Add Command"
	msg_want := "Microservice: New_service_test Added To Server"

	if title_want != title {
		t.Errorf("\n\nError: Failed To Add To Database Even If All Conditions Met:\nWhat We Wanted: %q\nWhat We Got: %q", title_want, title)
		Delete("testname_5")
		Delete("New_service_test")
	} else if msg_want != msg {
		t.Errorf("\n\nError: Failed To Add To Database Even If All Conditions Met:\nWhat We Wanted: %q\nWhat We Got: %q", msg_want, msg)
		Delete("testname_5")
		Delete("New_service_test")
	} else {
		Delete("testname_5")
		Delete("New_service_test")
	}
}
