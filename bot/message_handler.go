package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/nathanjcook/discordbotgo/bot/commands"
	dbconfig "github.com/nathanjcook/discordbotgo/config"
	"go.uber.org/zap"
)

// Define Microservice Struct For GORM Database Intergration
type Microservice struct {
	MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
	MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
	MicroserviceUrl     string `gorm:"column:microservice_url;"`
	MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
}

var ArgMsgErr = "Invalid Amount Of Args Provided"
var MessageEnder = "*** \n\n"

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var message string
	var helpMessage string
	var title string
	var msg string
	var isHelp bool

	var query Microservice

	// check if sender is self, and don't reply if true
	if m.Author.ID == BotId {
		return
	}
	// If Discord message starts with !gobot/etc
	if strings.Contains(m.Content, os.Getenv("BOT_PREFIX")) {
		//Split Discord Message By Whitespace And Save Each Item As A Slice Called cmdsplit
		cmdsplit := strings.Split(m.Content, " ")
		//Return User Permission Variable For The Current Channel
		p, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			zap.L().Error("Error", zap.Error(err))
		}
		//Check If User Is An Administrator For The Current Channel
		adminCheck := p & discordgo.PermissionAdministrator
		//Sets the value of messageContent to the message the user sent on discord
		messageContent := m.Content
		//If Second Item In cmdsplit slice is equal to add
		if cmdsplit[1] == "add" {
			//Call Function AddHandler Function With User Permission Values And cmdsplit slice
			title, msg = AddHandler(int(adminCheck), cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			//If Second Item In cmdsplit slice is equal to delete
		} else if cmdsplit[1] == "delete" {
			//Call Function DeleteHandler Function With User Permission Values And cmdsplit slice
			title, msg = DeleteHandler(int(adminCheck), cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			//If Second Item In cmdsplit slice is equal to help
		} else if cmdsplit[1] == "help" {
			//Call Function HelpHandler function With cmdsplit slice
			title, msg, isHelp = HelpHandler(cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			helpMessage = "<@" + m.Author.ID + ">" + " HELP!!!" + commands.AddTitle + commands.AddMsg + commands.DeleteTitle + commands.DeleteMsg + commands.InfoTitle + commands.InfoMsg + commands.MicroserviceTitle + commands.MicroserviceMsg
			//If Second Item In cmdsplit slice is equal to info
		} else if cmdsplit[1] == "info" {
			//Call Function InfoHandler function With cmdsplit slice
			title, msg = InfoHandler(cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + "***" + title + MessageEnder + msg
			//If Second Item In cmdsplit slice is equal to anything other than 'add' 'delete' 'info' 'help', assume user is trying to call microservice
		} else {
			//Perform SQL Query For Microservices Table Where Microservice_name is equal to the second word in the discord message/command
			host := dbconfig.DB.Table("microservices").Where("microservice_name = ?", string(cmdsplit[1])).Scan(&query)
			//Check Instances Where Name Not Found In Database: If this happens return info details
			if host.RowsAffected < 1 {
				titles, msg := commands.Info()
				title = "Microservice " + cmdsplit[1] + " Does Not Exist"
				message = "<@" + m.Author.ID + ">" + " " + "***" + title + MessageEnder + "***" + titles + "***\n" + msg
			} else {
				title, msg = MicroserviceHandler(query, cmdsplit, messageContent)
				message = "<@" + m.Author.ID + ">" + " " + "***" + title + MessageEnder + msg
				//Discord Has a Limit oF 2000 Characters For Users And Bots.
				//Will Return An Error Specifying That The HTTP Response Cannot Be Sent To Discord Chat Due To This Limit
				if utf8.RuneCountInString(message) >= 2000 {
					title := "Microservice " + cmdsplit[1] + " Error"
					msg := "Response Exceeded 2000 Characters! Report This Microservice To An Admin To Review"
					message = "<@" + m.Author.ID + ">" + " " + "***" + title + MessageEnder + msg
				}

			}
		}
		if isHelp {
			//Send Microservice Help Message To Discord Chat to the channel identified by m.ChannelID
			_, _ = s.ChannelMessageSend(m.ChannelID, helpMessage)
		} else {
			//Send Message To Discord Chat to the channel identified by m.ChannelID
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}

func AddHandler(adminCheck int, cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If Admin Tries To Use Add Command With Less Than Or Greater Than Three Variables
	if len(cmdsplit) < 5 || len(cmdsplit) > 6 {
		title := "Add Command Error"
		msg := ArgMsgErr
		return title, msg
	} else if adminCheck == 0 {
		//If User Is Not An Admin Then Return Error Message Identifying That Only Admins Can Add Microservices
		title = "Add Command Error"
		msg = "Only Admins Can Add MicroServices! Please Contact Any Administrators If You Want A Particular Microservice Added"
		return title, msg
	} else {
		//Call Add Function With Variables For <Microservice_Name> <Microservice_Url> <Microservice_Timeout>
		//Return String Variables Title & Message From Returned Strings From Add Function
		title, msg = commands.Add(cmdsplit[2], cmdsplit[3], cmdsplit[4])
		return title, msg
	}
}

func DeleteHandler(adminCheck int, cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If Admin Tries To Use Delete Command With Less Than Or Greater Than One Variable
	if len(cmdsplit) < 3 || len(cmdsplit) > 3 {
		title := "Delete Command Error"
		msg := ArgMsgErr
		return title, msg
		//If User Is Not An Admin Then Return Error Message Identifying That Only Admins Can Delete Microservices
	} else if adminCheck == 0 {
		title = "Delete Command Error"
		msg = "Only Admins Can Delete MicroServices! Please Contact Any Administrators If You Want A Particular Microservice Deleted"
		return title, msg
	} else {
		//Call Delete Function With Variables For <Microservice_Name>
		//Return String Variables Title & Message From Returned Strings From Delete Function
		title, msg = commands.Delete(cmdsplit[2])
		return title, msg
	}
}

func HelpHandler(cmdsplit []string) (string, string, bool) {
	var title string
	var msg string
	var isHelp bool
	//Error prevention: Return Error Message If User Tries To Use Help Command With Added Unnecessary Variables
	if len(cmdsplit) > 2 {
		title := "Help Command Error"
		msg := ArgMsgErr
		return title, msg, isHelp
	} else {
		isHelp = true
		//Return String Variables Title & Message From Returned Strings From Help Function If Error Occured And is_help bool
		return title, msg, isHelp
	}
}

func InfoHandler(cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If User Tries To Use Info Command With Added Unnecessary Variables
	if len(cmdsplit) > 2 {
		title := "Info Command Error"
		msg := ArgMsgErr
		return title, msg
	} else {
		//Call Info Function
		title, msg = commands.Info()
		//Return String Variables Title & Message From Returned Strings From Info Function
		return title, msg
	}
}

func MicroserviceHandler(query Microservice, cmdsplit []string, messageContent string) (string, string) {
	respTimeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(query.MicroserviceTimeout))
	defer cancel()

	var title string
	var msg string

	if len(cmdsplit) < 3 {
		title = "Microservice Command Error"
		msg = ArgMsgErr
		return title, msg
	}

	txt, str := BodyParser(messageContent)

	if str != "" {
		title = "Pre Microservice JSON Body Error"
		msg = str
		return title, msg
	}

	body := bytes.NewBuffer(txt)
	urls := (query.MicroserviceUrl + "/api/" + cmdsplit[2])
	req, err := http.NewRequestWithContext(respTimeout, http.MethodPost, urls, body)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		title = cmdsplit[1] + " error"
		msg = "Error Connecting To Microservice"
		return title, msg
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		title = cmdsplit[1] + " error"
		msg = "Timeout"
		return title, msg
	}

	if res.StatusCode == 404 && cmdsplit[2] == "help" {
		title = cmdsplit[1] + " No Help"
		msg = "The Microservice " + cmdsplit[1] + "Does Not Have A Help Section! Report This To An Admin"
		return title, msg
	} else if res.StatusCode == 404 {
		title = cmdsplit[1] + " Endpoint Not Found"
		helper, txt := commands.Gethelp((query.MicroserviceUrl + "/api/help"))

		if txt != "" {
			msg = txt
		} else {
			msg = BodyReader(helper)
		}
		return title, msg
	}

	defer res.Body.Close()
	bodyRes, errs := io.ReadAll(res.Body)

	if errs != nil {
		title = cmdsplit[1] + "error"
		msg = "Error Reading Response Body"
		return title, msg
	} else {
		title = cmdsplit[1]
		msg = BodyReader(bodyRes)
		return title, msg
	}
}

// This function will be called a new shard connects
func onConnect(s *discordgo.Session, evt *discordgo.Connect) {
	fmt.Printf("[INFO] Shard #%v connected.\n", s.ShardID)
}
