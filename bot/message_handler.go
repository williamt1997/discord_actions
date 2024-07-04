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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var message string
	var help_message string
	var title string
	var msg string
	var is_help bool

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
			//Call Function Add_Handler Function With User Permission Values And cmdsplit slice
			title, msg = Add_Handler(int(adminCheck), cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			//If Second Item In cmdsplit slice is equal to delete
		} else if cmdsplit[1] == "delete" {
			//Call Function Delete_Handler Function With User Permission Values And cmdsplit slice
			title, msg = Delete_Handler(int(adminCheck), cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			//If Second Item In cmdsplit slice is equal to help
		} else if cmdsplit[1] == "help" {
			//Call Function Help_Handler function With cmdsplit slice
			title, msg, is_help = Help_Handler(cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + title + ": " + msg
			help_message = "<@" + m.Author.ID + ">" + " HELP!!!" + commands.AddTitle + commands.AddMsg + commands.DeleteTitle + commands.DeleteMsg + commands.InfoTitle + commands.InfoMsg + commands.MicroserviceTitle + commands.MicroserviceMsg
			//If Second Item In cmdsplit slice is equal to info
		} else if cmdsplit[1] == "info" {
			//Call Function Info_Handler function With cmdsplit slice
			title, msg = Info_Handler(cmdsplit)
			message = "<@" + m.Author.ID + ">" + " " + "***" + title + "*** \n\n" + msg
			//If Second Item In cmdsplit slice is equal to anything other than 'add' 'delete' 'info' 'help', assume user is trying to call microservice
		} else {
			//Perform SQL Query For Microservices Table Where Microservice_name is equal to the second word in the discord message/command
			host := dbconfig.DB.Table("microservices").Where("microservice_name = ?", string(cmdsplit[1])).Scan(&query)
			//Check Instances Where Name Not Found In Database: If this happens return info details
			if host.RowsAffected < 1 {
				titles, msg := commands.Info()
				title = "Microservice " + cmdsplit[1] + " Does Not Exist"
				message = "<@" + m.Author.ID + ">" + " " + "***" + title + "*** \n\n" + "***" + titles + "***\n" + msg
			} else {
				title, msg = Microservice_Handler(query, cmdsplit, messageContent)
				message = "<@" + m.Author.ID + ">" + " " + "***" + title + "*** \n\n" + msg
				//Discord Has a Limit oF 2000 Characters For Users And Bots.
				//Will Return An Error Specifying That The HTTP Response Cannot Be Sent To Discord Chat Due To This Limit
				if utf8.RuneCountInString(message) >= 2000 {
					title := "Microservice " + cmdsplit[1] + " Error"
					msg := "Response Exceeded 2000 Characters! Report This Microservice To An Admin To Review"
					message = "<@" + m.Author.ID + ">" + " " + "***" + title + "*** \n\n" + msg
				}

			}
		}
		if is_help {
			//Send Microservice Help Message To Discord Chat to the channel identified by m.ChannelID
			_, _ = s.ChannelMessageSend(m.ChannelID, help_message)
		} else {
			//Send Message To Discord Chat to the channel identified by m.ChannelID
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}

func Add_Handler(adminCheck int, cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If Admin Tries To Use Add Command With Less Than Or Greater Than Three Variables
	if len(cmdsplit) < 5 || len(cmdsplit) > 6 {
		title := "Add Command Error"
		msg := "Invalid Amount Of Args Provided"
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

func Delete_Handler(adminCheck int, cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If Admin Tries To Use Delete Command With Less Than Or Greater Than One Variable
	if len(cmdsplit) < 3 || len(cmdsplit) > 3 {
		title := "Delete Command Error"
		msg := "Invalid Amount Of Args Provided"
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

func Help_Handler(cmdsplit []string) (string, string, bool) {
	var title string
	var msg string
	var is_help bool
	//Error prevention: Return Error Message If User Tries To Use Help Command With Added Unnecessary Variables
	if len(cmdsplit) > 2 {
		title := "Help Command Error"
		msg := "Invalid Amount Of Args Provided"
		return title, msg, is_help
	} else {
		is_help = true
		//Return String Variables Title & Message From Returned Strings From Help Function If Error Occured And is_help bool
		return title, msg, is_help
	}
}

func Info_Handler(cmdsplit []string) (string, string) {
	var title string
	var msg string
	//Error prevention: Return Error Message If User Tries To Use Info Command With Added Unnecessary Variables
	if len(cmdsplit) > 2 {
		title := "Info Command Error"
		msg := "Invalid Amount Of Args Provided"
		return title, msg
	} else {
		//Call Info Function
		title, msg = commands.Info()
		//Return String Variables Title & Message From Returned Strings From Info Function
		return title, msg
	}
}

func Microservice_Handler(query Microservice, cmdsplit []string, messageContent string) (string, string) {
	//Create HTTP Response Context With Timeout set to microservice_timeout column details
	resp_timeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(query.MicroserviceTimeout))
	//Cancel Context When Surrounding Function Returns
	defer cancel()
	var title string
	var msg string
	//If User Only Inputs !gobot <microservice_name> then present error as endpoint is required
	if len(cmdsplit) < 3 {
		title = "Microservice Command Error"
		msg = "Invalid Amount Of Args Provided"
		return title, msg
	} else {
		//Call Body_Parser Function To Convert messageContent[HTTP Request Body] into JSON Format
		txt, str := Body_Parser(messageContent)
		//str is used as an error identifyer if error occurs during the body_parser function
		//If str is not empty then return error message from returned str
		if str != "" {
			title = "Pre Microservice JSON Body Error"
			msg = str
			return title, msg
		} else {
			// Create a new buffer with the values of txt [Returned from body_parser].
			body := bytes.NewBuffer(txt)
			//Setting The POST Request URL with the HOST Url and Endpoint
			urls := (query.MicroserviceUrl + "/api/" + cmdsplit[2])
			//Sending a POST request with to the microservice endpoint with the http request method as post and including the body processed by body_parser
			//Adding optional headers application/json to allow json data to be parsed
			req, err := http.NewRequestWithContext(resp_timeout, http.MethodPost, urls, body)
			req.Header.Set("Content-Type", "application/json")
			// Check And Handle Errors Whilst Making The Post Request
			if err != nil {
				title = cmdsplit[1] + " error"
				msg = "Error Connecting To Microservice"
				return title, msg
			}
			// If response fails or exceeds the timeout context then return timeout message
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				title = cmdsplit[1] + " error"
				msg = "Timeout"
				return title, msg
			} else {
				// If response returns a status code of 404 then... [This suggests that the users inputted endpoint is either unreachable or the endpoint doesn't exist]
				if res.StatusCode == 404 {
					// If endpoint is equal to help then just return general error message. This is done to remove the risk of an infinite loop as the next error message calls upon the help section to return the microservices help details
					if cmdsplit[2] == "help" {
						title = cmdsplit[1] + " No Help"
						msg = "The Microservice " + cmdsplit[1] + "Does Not Have A Help Section! Report This To An Admin"
						return title, msg
					} else {
						//Call Function Get Help And Return Help Details For The Specified Microservice
						title = cmdsplit[1] + " Endpoint Not Found"
						helper, txt := commands.Get_Help((query.MicroserviceUrl + "/api/help"))
						if txt != "" {
							msg = txt
						} else {
							msg = Body_Reader(helper)
						}
					}
					return title, msg
				} else {
					// Closing The Response Body After Reading
					defer res.Body.Close()
					// Reading The Response Body
					body, err := io.ReadAll(res.Body)
					// Check And Handle Errors When Reading The Response Body Fails
					if err != nil {
						title = cmdsplit[1] + "error"
						msg = "Error Reading Response Body"
						return title, msg
					} else {
						title = cmdsplit[1]
						//Call Body Reader Function To Convert Response Body From JSON/Byte to String Format
						msg = Body_Reader(body)
						//Return title and msg as to send the microservice responsse on the discord chat
						return title, msg
					}
				}
			}
		}
	}
}

// This function will be called a new shard connects
func onConnect(s *discordgo.Session, evt *discordgo.Connect) {
	fmt.Printf("[INFO] Shard #%v connected.\n", s.ShardID)
}
