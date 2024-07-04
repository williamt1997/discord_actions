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

var TitleError = "Add Command Error"

var CommonErrMessage = "\n\nTitle We Wanted: %q\nWhat We Got: %q\n\nMessage We Wanted: %q\nWhat We Got:%q"

var ExistingMS = "http://localhost:8081"

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

	err = db.AutoMigrate(&Microservice{})
	if err != nil {
		zap.L().Panic("Error Migrating To Database .env file:", zap.Error(err))
	}
}

func TestAddMSNameAlreadyExists(t *testing.T) {
	setupTestDBAdd()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "existing_service",
		MicroserviceUrl:     "http://localhost:3007",
		MicroserviceTimeout: 70,
	})

	title, msg := Add("existing_service", ExistingMS, "50")
	msgWant := "Microservice Name AND Microservice URL Must Be Unique"

	if TitleError != title || msgWant != msg {
		t.Errorf(CommonErrMessage, TitleError, title, msgWant, msg)
	}
}

func TestAddMSHostURLAlreadyExists(t *testing.T) {
	setupTestDBAdd()

	dbconfig.DB.Create(&Microservice{
		MicroserviceName:    "existing_service",
		MicroserviceUrl:     ExistingMS,
		MicroserviceTimeout: 70,
	})

	title, msg := Add("new_service", ExistingMS, "50")

	msgWant := "Microservice Name AND Microservice URL Must Be Unique"

	if TitleError != title || msgWant != msg {
		t.Errorf(CommonErrMessage, TitleError, title, msgWant, msg)
	}
	Delete("existing_service")

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

	title, msg := Add("New_service_test", ExistingMS, "50")
	titleWant := "Add Command"
	msgWant := "Microservice: New_service_test Added To Server"

	if titleWant != title || msgWant != msg {
		t.Errorf(CommonErrMessage, titleWant, title, msgWant, msg)
	}
	Delete("testname_5")
	Delete("New_service_test")
}
