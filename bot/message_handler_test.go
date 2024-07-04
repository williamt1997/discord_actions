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

var Want_Add_Error_Title = "Add Command Error"
var Want_Delete_Error_Title = "Delete Command Error"
var Want_Help_Error_Title = "Help Command Error"
var Want_Info_Error_Title = "Info Command Error"
var Want_Microservice_Error_Title = "Microservice Command Error"

var Admin_Message = "Only Admins Can"
var Args_Message = "Invalid Amount Of Args Provided"

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
func TestAdd_HandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot add test test 55", " ")
	title, msg := Add_Handler(0, cmdsplit)

	if Want_Add_Error_Title != title || !strings.Contains(msg, Admin_Message) {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, Admin_Message, msg)
	}
}

func TestAdd_ArgsValError(t *testing.T) {
	var inputs = []string{"!gobot add test test 55 a b c d e f g", "!gobot add test"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := Add_Handler(123, cmdsplit)

		if Want_Add_Error_Title != title || Args_Message != msg {
			t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, Args_Message, msg)
		}
	}
}

func TestAdd_ArgsAddExistingIS(t *testing.T) {
	var inputs = []string{"!gobot add add http://localhost:3002 60", "!gobot add delete http://localhost:3002 60", "!gobot add info http://localhost:3002 60", "!gobot add help http://localhost:3002 60"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		want_existing := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"
		title, msg := Add_Handler(123, cmdsplit)

		if Want_Add_Error_Title != title || want_existing != msg {
			t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, want_existing, msg)
		}
	}
}

func TestAdd_HandlerNameTooLarge(t *testing.T) {
	cmdsplit := strings.Split("!gobot add aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa test 50", " ")
	title, msg := Add_Handler(0, cmdsplit)
	want_toolarge := "Microservice Name Cannot Be Larger Than 25 Characters"

	if Want_Add_Error_Title != title || !strings.Contains(msg, want_toolarge) {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, want_toolarge, msg)
	}
}

func TestAdd_HandlerNameBadUrl(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:8081/api/help")
	cmdsplit := strings.Split("!gobot add test nonexist_234232 50", " ")
	title, msg := Add_Handler(0, cmdsplit)
	want_cantconnect := "Error Connecting To Microservice"

	if Want_Add_Error_Title != title || !strings.Contains(msg, want_cantconnect) {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, want_cantconnect, msg)
	}
}

func TestAdd_HandlerNameBadTimeOutFormatl(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:3002")
	cmdsplit := strings.Split("!gobot add tester http://localhost:3002 A", " ")
	title, msg := Add_Handler(0, cmdsplit)
	want_badtimeout := "Timeout Is In An Incorrect Format"

	if Want_Add_Error_Title != title || !strings.Contains(msg, want_badtimeout) {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Add_Error_Title, title, want_badtimeout, msg)
	}
}

func TestDelete_HandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot delete test", " ")
	title, msg := Delete_Handler(0, cmdsplit)

	if Want_Delete_Error_Title != title || !strings.Contains(msg, Admin_Message) {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Delete_Error_Title, title, Admin_Message, msg)
	}
}

func TestDelete_ArgsValError(t *testing.T) {
	var inputs = []string{"!gobot delete test test", "!gobot delete"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := Delete_Handler(123, cmdsplit)
		if Want_Delete_Error_Title != title || Args_Message != msg {
			t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Delete_Error_Title, title, Args_Message, msg)
		}
	}
}

func TestHelp_ArgsValError(t *testing.T) {
	var inputs = []string{"!gobot help test", "!gobot help test dfs ds"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg, is_help := Help_Handler(cmdsplit)
		if Want_Help_Error_Title != title || Args_Message != msg {
			t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q: IsHelp=%v", Want_Help_Error_Title, title, Args_Message, msg, is_help)
		}
	}
}

func TestHelp_HandlerHelpNotReturned(t *testing.T) {
	cmdsplit := strings.Split("!gobot help", " ")
	title, msg, is_help := Help_Handler(cmdsplit)

	want_title := ""
	want_msg := ""
	want_help := true

	if want_title != title || want_msg != msg || want_help != is_help {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q: IsHelp=%v", Want_Help_Error_Title, title, Args_Message, msg, is_help)
	}
}

func TestInfo_ArgsValError(t *testing.T) {
	var inputs = []string{"!gobot info test", "!gobot help test dfs ds"}
	for i := 0; i < len(inputs); i++ {
		cmdsplit := strings.Split(inputs[i], " ")
		title, msg := Info_Handler(cmdsplit)
		if Want_Info_Error_Title != title || Args_Message != msg {
			t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Info_Error_Title, title, Args_Message, msg)
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

func TestMicroservice_HandlerTooFewValues(t *testing.T) {
	message_content_test := "!gobot alcCalc"
	cmdsplit := strings.Split(message_content_test, " ")
	title, msg := Microservice_Handler(QueryTest, cmdsplit, message_content_test)

	if Want_Microservice_Error_Title != title || Args_Message != msg {
		t.Errorf("\n\nTitle Wanted: %q, Title Recieved: %q\n\n Message Wanted: %q, Message Recieved: %q", Want_Microservice_Error_Title, title, Args_Message, msg)
	}
}

func TestMicroservice_HandlerBadVariable(t *testing.T) {
	message_content_test := "!gobot alcCalc calculate -test"
	cmdsplit := strings.Split(message_content_test, " ")
	title, msg := Microservice_Handler(QueryTest, cmdsplit, message_content_test)

	want_title := "Pre Microservice JSON Body Error"

	if want_title != title {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting Using Microservice Without An Endpoint\n\n Title Wanted: %q, Title Recieved: %q: Msg Recieved %q", want_title, title, msg)
	}
}
