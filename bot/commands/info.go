package commands

import (
	dbconfig "github.com/nathanjcook/discordbotgo/config"
)

// Function To Get All Microservice Names And General Help Commands
func Info() (string, string) {
	// Define Microservice Struct For GORM Database Intergration
	type Microservice struct {
		MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
		MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
		MicroserviceUrl     string `gorm:"column:microservice_url;"`
		MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
	}
	var names []string
	msg := "\n"
	//Querying The Database To SELECT And Store All Microservice Names Into A Slice
	dbconfig.DB.Model(&Microservice{}).Pluck("microservice_name", &names)
	//If The Length Of The Slice Is Greater Than 0: Showing That At Least One Microservice Exists
	if len(names) > 0 {
		title := "Available Microservices"
		//Looping Through Each Item In The names slice and setting message to the full command needed to access the help endpoint for each microservice
		for i := 0; i < len(names); i++ {
			msg += "!gobot " + names[i] + " help\n\n"
		}
		return title, msg
		//Handling Instances Where No Microservices Are Currently Added To The Postgres Database
	} else {
		title := "Info Command Null"
		return title, "No Microservices Available"
	}
}
