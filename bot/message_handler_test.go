package bot

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
	dbconfig "github.com/nathanjcook/discordbotgo/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var WantAddErrorTitle = "Add Command Error"
var WantDeleteErrorTitle = "Delete Command Error"
var WantHelpErrorTitle = "Help Command Error"
var WantInfoErrorTitle = "Info Command Error"
var WantMicroserviceErrorTitle = "Microservice Command Error"

var AdminMessage = "Only Admins Can"
var ArgsMessage = "Invalid Amount Of Args Provided"

var CommonErrMessage = "\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q"
var HelpErrMessage = "\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q: IsHelp=%v"
var MsErrMessage = "\n\nTitle Wanted: %q, Title Recieved: %q: Msg Recieved %q"

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
		zap.L().Panic("failed to connect database:", zap.Error(err))
	}
	dbconfig.DB = db

	err = db.AutoMigrate(&Microservice{})
	if err != nil {
		zap.L().Panic("Error Migrating To Database .env file:", zap.Error(err))
	}
}
func TestAddHandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot add test test 55", " ")
	title, msg := AddHandler(0, cmdsplit)

	if WantAddErrorTitle != title || !strings.Contains(msg, AdminMessage) {
		t.Errorf(CommonErrMessage, WantAddErrorTitle, title, AdminMessage, msg)
	}
}

func TestAddArgsValError(t *testing.T) {
	var inputs = []string{"!gobot add test test 55 a b c d e f g", "!gobot add test"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := AddHandler(123, cmdsplit)

		if WantAddErrorTitle != title || ArgsMessage != msg {
			t.Errorf(CommonErrMessage, WantAddErrorTitle, title, ArgsMessage, msg)
		}
	}
}

func TestAddArgsAddExistingIS(t *testing.T) {
	var inputs = []string{"!gobot add add http://localhost:3002 60", "!gobot add delete http://localhost:3002 60", "!gobot add info http://localhost:3002 60", "!gobot add help http://localhost:3002 60"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		wantExisting := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"
		title, msg := AddHandler(123, cmdsplit)

		if WantAddErrorTitle != title || wantExisting != msg {
			t.Errorf(CommonErrMessage, WantAddErrorTitle, title, wantExisting, msg)
		}
	}
}

func TestAddHandlerNameTooLarge(t *testing.T) {
	cmdsplit := strings.Split("!gobot add aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa test 50", " ")
	title, msg := AddHandler(123, cmdsplit)
	wantToolarge := "Microservice Name Cannot Be Larger Than 25 Characters"

	if WantAddErrorTitle != title || !strings.Contains(msg, wantToolarge) {
		t.Errorf(CommonErrMessage, WantAddErrorTitle, title, wantToolarge, msg)
	}
}

func TestAddHandlerNameBadUrl(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:8081/api/help")
	cmdsplit := strings.Split("!gobot add test nonexist_234232 50", " ")
	title, msg := AddHandler(123, cmdsplit)
	wantCantconnect := "Error Connecting To Microservice"

	if WantAddErrorTitle != title || !strings.Contains(msg, wantCantconnect) {
		t.Errorf(CommonErrMessage, WantAddErrorTitle, title, wantCantconnect, msg)
	}
}

func TestAddHandlerNameBadTimeOutFormatl(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:3002")
	cmdsplit := strings.Split("!gobot add tester http://localhost:3002 A", " ")
	title, msg := AddHandler(123, cmdsplit)
	wantBadtimeout := "Timeout Is In An Incorrect Format"

	if WantAddErrorTitle != title || !strings.Contains(msg, wantBadtimeout) {
		t.Errorf(CommonErrMessage, WantAddErrorTitle, title, wantBadtimeout, msg)
	}
}

func TestDeleteHandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot delete test", " ")
	title, msg := DeleteHandler(0, cmdsplit)

	if WantDeleteErrorTitle != title || !strings.Contains(msg, AdminMessage) {
		t.Errorf(CommonErrMessage, WantDeleteErrorTitle, title, AdminMessage, msg)
	}
}

func TestDeleteArgsValError(t *testing.T) {
	var inputs = []string{"!gobot delete test test", "!gobot delete"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := DeleteHandler(123, cmdsplit)
		if WantDeleteErrorTitle != title || ArgsMessage != msg {
			t.Errorf(CommonErrMessage, WantDeleteErrorTitle, title, ArgsMessage, msg)
		}
	}
}

func TestHelpArgsValError(t *testing.T) {
	var inputs = []string{"!gobot help test", "!gobot help test dfs ds"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg, isHelp := HelpHandler(cmdsplit)
		if WantHelpErrorTitle != title || ArgsMessage != msg {
			t.Errorf(HelpErrMessage, WantHelpErrorTitle, title, ArgsMessage, msg, isHelp)
		}
	}
}

func TestHelpHandlerHelpNotReturned(t *testing.T) {
	cmdsplit := strings.Split("!gobot help", " ")
	title, msg, isHelp := HelpHandler(cmdsplit)

	wantTitle := ""
	wantMsg := ""
	wantHelp := true

	if wantTitle != title || wantMsg != msg || wantHelp != isHelp {
		t.Errorf(HelpErrMessage, wantTitle, title, wantMsg, msg, wantHelp)
	}
}

func TestInfoArgsValError(t *testing.T) {
	var inputs = []string{"!gobot info test", "!gobot help test dfs ds"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := InfoHandler(cmdsplit)
		if WantInfoErrorTitle != title || ArgsMessage != msg {
			t.Errorf(CommonErrMessage, WantInfoErrorTitle, title, ArgsMessage, msg)
		}
	}
}

type MicroserviceTest struct {
	MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
	MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
	MicroserviceUrl     string `gorm:"column:microservice_url;"`
	MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
}

var QueryTest Microservice

func TestMicroserviceHandlerTooFewValues(t *testing.T) {
	messageContentTest := "!gobot alcCalc"
	cmdsplit := strings.Split(messageContentTest, " ")
	title, msg := MicroserviceHandler(QueryTest, cmdsplit, messageContentTest)

	if WantMicroserviceErrorTitle != title || ArgsMessage != msg {
		t.Errorf(CommonErrMessage, WantMicroserviceErrorTitle, title, ArgsMessage, msg)
	}
}

func TestMicroserviceHandlerBadVariable(t *testing.T) {
	messageContentTest := "!gobot alcCalc calculate -test"
	cmdsplit := strings.Split(messageContentTest, " ")
	title, msg := MicroserviceHandler(QueryTest, cmdsplit, messageContentTest)

	wantTitle := "Pre Microservice JSON Body Error"

	if wantTitle != title {
		t.Errorf(MsErrMessage, wantTitle, title, msg)
	}
}
