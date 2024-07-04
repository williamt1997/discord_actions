package commands

import (
	"bytes"
	"net/http"
	"strconv"

	dbconfig "github.com/nathanjcook/discordbotgo/config"
	"go.uber.org/zap"
)

// Add Function To Allow Admins To Add Microservice Details Into The PostGres Database
// Example: !gobot add <microservice_name> <microservice_url> <microservice_timeout>
// Name: = cmdsplit[2] from microservice_handler
// Url: = cmdsplit[3] from microservice_handler
// Timeout: = cmdsplit[4] from microservice_handler
func Add(name string, url string, timeout string) (string, string) {
	// Define Microservice Struct For GORM Database Intergration
	type Microservice struct {
		MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
		MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
		MicroserviceUrl     string `gorm:"column:microservice_url;"`
		MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
	}
	var query Microservice
	var title string
	var msg string
	// Checks If The Character Length Of The Inputted Microservice Name Exceeds The 25 Character Limit
	if len(name) > 25 {
		title = "Add Command Error"
		msg = "Microservice Name Cannot Be Larger Than 25 Characters"
		return title, msg
	} else {
		//Error Handling Where Users Try To Add A Microservice With The Same Name Of Any Internal Bot Command Names
		if name != "add" && name != "info" && name != "delete" && name != "help" {
			//Checking If Microservice Name Or Microservice URL Already Exists In The Postgres Database
			result := dbconfig.DB.Where("microservice_name = ? OR microservice_url = ?", name, url).Find(&query)
			//If Rows Affected Is Less Than One Then Microservice Is Unique
			if result.RowsAffected < 1 {
				timeout_int, err := strconv.Atoi(timeout)
				//Error Handling To Check Instances Where User Did Not Input Timeout As A Integer
				if err != nil {
					title = "Add Command Error"
					msg = "Timeout Is In An Incorrect Format"
					return title, msg
				} else {
					// Preparing To Make A POST Request To The Microservice ENDPOINT Of Help

					//Setting POST Request BODY As Blank Byte
					body := new(bytes.Buffer)
					//Setting The POST Request URL with the HOST Url and Help Endpoint
					urls := (url + "/api/help")
					//Sending a POST request with an empty body to the microservice's help endpoint to check its availability
					resp, err := http.Post(urls, "application/json", body)
					// Check And Handle Errors Whilst Making The Post Request
					if err != nil {
						title = "Add Command Error"
						msg = "Error Connecting To Microservice"
						zap.L().Error("Error", zap.Error(err))
						return title, msg
					} else {
						// Check if the response status code is less than 400: Successful Request
						if resp.StatusCode < 400 {
							//Creating And Adding The New Microservice Details Into The Postgres Database
							microserviceAdd := Microservice{MicroserviceName: name, MicroserviceUrl: url, MicroserviceTimeout: timeout_int}
							err := dbconfig.DB.Create(&microserviceAdd).Error
							//Error Handling If The Connection To The Database Fails During The Add Query
							if err != nil {
								title = "Add Command Error"
								msg = "Error Connecting To Database"
								return title, msg
							} else {
								title = "Add Command"
								msg = "Microservice: " + name + " Added To Server"
								return title, msg
							}

							//If HTTP Request Returns A Status Code Greater Or Equal To 400 For The Help Endpoint
						} else {
							title = "Add Command Error"
							msg = "Cannot Connect To Microservice Via Selected Host URL"
							return title, msg
						}
					}
				}
				// If Microservice Name Or Microservice URL Is Not Unique: With Exisiting Microservices
			} else {
				title = "Add Command Error"
				msg = "Microservice Name AND Microservice URL Must Be Unique"
				return title, msg
			}
			// If Microservice Name Is Not Unique: With Internal Bot Commands
		} else {
			title = "Add Command Error"
			msg = "Microservice Name Cannot Be The Same As Internal Commands: add, delete, help, info"
			return title, msg
		}
	}
}
