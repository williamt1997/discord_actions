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

	db.AutoMigrate(&Microservice{}).Error()
}
func TestAdd_HandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot add test test 55", " ")
	title, msg := Add_Handler(0, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Only Admins Can Add MicroServices! Please Contact Any Administrators If You Want A Particular Microservice Added"

	if want_title != title {
		t.Errorf("\n\nOnly Admins Should Be Able To Add Microservices!\nIf This Fails Then This Suggests That The Add Command Is Accessible To Everyone\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nOnly Admins Should Be Able To Add Microservices!\nIf This Fails Then This Suggests That The Add Command Is Accessible To Everyonee\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestDelete_HandlerNotAdmin(t *testing.T) {
	cmdsplit := strings.Split("!gobot delete test", " ")
	title, msg := Delete_Handler(0, cmdsplit)

	want_title := "Delete Command Error"
	want_msg := "Only Admins Can Delete MicroServices! Please Contact Any Administrators If You Want A Particular Microservice Deleted"

	if want_title != title {
		t.Errorf("\n\nOnly Admins Should Be Able To Delete Microservices!\nIf This Fails Then This Suggests That The Delete Command Is Accessible To Everyone\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nOnly Admins Should Be Able To Delete Microservices!\nIf This Fails Then This Suggests That The Delete Command Is Accessible To Everyonee\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerTooManyValues(t *testing.T) {
	cmdsplit := strings.Split("!gobot add test test 55 a b c d e f g", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Invalid Amount Of Args Provided"

	if want_title != title {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestDelete_HandlerTooManyValues(t *testing.T) {
	cmdsplit := strings.Split("!gobot delete test test", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Delete Command Error"
	want_msg := "Invalid Amount Of Args Provided"

	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestHelp_HandlerTooManyValues(t *testing.T) {
	cmdsplit := strings.Split("!gobot help test", " ")
	title, msg, is_help := Help_Handler(cmdsplit)

	want_title := "Help Command Error"
	want_msg := "Invalid Amount Of Args Provided"
	want_help := false

	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
	if want_help != is_help {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Msg Wanted: %v, Msg Recieved: %v", want_help, is_help)
	}
}

func TestHelp_HandlerHelpNotReturned(t *testing.T) {
	cmdsplit := strings.Split("!gobot help", " ")
	title, msg, is_help := Help_Handler(cmdsplit)

	want_title := ""
	want_msg := ""
	want_help := true

	if want_msg != title {
		t.Errorf("\n\nError: Title Should Return As Empty\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Msg Should Return As Empty\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
	if want_help != is_help {
		t.Errorf("\n\nError: is_help bool not returning true and preventing help functionality in messageCreate\n\n Msg Wanted: %v, Msg Recieved: %v", want_help, is_help)
	}
}

func TestAdd_HandlerTooFewValues(t *testing.T) {
	cmdsplit := strings.Split("!gobot add test", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Invalid Amount Of Args Provided"

	if want_title != title {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Few Variables\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Few Variables\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestDelete_HandlerTooFewValues(t *testing.T) {
	cmdsplit := strings.Split("!gobot delete", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Delete Command Error"
	want_msg := "Invalid Amount Of Args Provided"

	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting To Many Variables\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerNameTooLarge(t *testing.T) {
	cmdsplit := strings.Split("!gobot add aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa test 50", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Microservice Name Cannot Be Larger Than 25 Characters"

	if want_title != title {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Providing A Microservice_Name beyond the databases required size/length\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Providing A Microservice_Name beyond the databases required size/length\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerNameBadUrl(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:8081/api/help")
	cmdsplit := strings.Split("!gobot add test nonexist_234232 50", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Error Connecting To Microservice"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That System Is Providing User With Incorrect Error Code [If At All]\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That System Is Providing User With Incorrect Error Code [If At All]\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerNameBadTimeOutFormat(t *testing.T) {
	setupTestDBAdd()
	defer gock.Off()
	gock.New("http://localhost:3002")
	cmdsplit := strings.Split("!gobot add tester http://localhost:3002 A", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Timeout Is In An Incorrect Format"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That Users Are Able To Input Non Interger Values In The Timeout Column\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That Users Are Able To Input Non Interger Values In The Timeout Column\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerUserTriesToAddAdd(t *testing.T) {
	cmdsplit := strings.Split("!gobot add add http://localhost:3002 60", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerUserTriesToAddDelete(t *testing.T) {
	cmdsplit := strings.Split("!gobot add delete http://localhost:3002 60", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerUserTriesToAddHelp(t *testing.T) {
	cmdsplit := strings.Split("!gobot help delete http://localhost:3002 60", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
	}
}

func TestAdd_HandlerUserTriesToAddInfo(t *testing.T) {
	cmdsplit := strings.Split("!gobot help delete http://localhost:3002 60", " ")
	title, msg := Add_Handler(123, cmdsplit)

	want_title := "Add Command Error"
	want_msg := "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"

	if want_title != title {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: Suggests That Users Can Add More Than One Add Command And Cause Multiple Commands To Run At Once\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
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

	want_title := "Microservice Command Error"
	want_msg := "Invalid Amount Of Args Provided"

	if want_title != title {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting Using Microservice Without An Endpoint\n\n Title Wanted: %q, Title Recieved: %q", want_title, title)
	}
	if want_msg != msg {
		t.Errorf("\n\nError: System Has Failed To Prevent User From Inputting Using Microservice Without An Endpoint\n\n Msg Wanted: %q, Msg Recieved: %q", want_msg, msg)
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
