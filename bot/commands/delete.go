package commands

import (
	dbconfig "github.com/nathanjcook/discordbotgo/config"
)

// Add Function To Allow Admins To Delete Microservice Details From The PostGres Database By Microservice Name
// Example: !gobot delete <microservice_name>
// Name: = cmdsplit[2] from microservice_handler
func Delete(name string) (string, string) {
	// Define Microservice Struct For GORM Database Intergration
	type Microservice struct {
		MicroserviceId      int    `gorm:"column:microservice_id;unique;primaryKey;autoIncrement"`
		MicroserviceName    string `gorm:"column:microservice_name;size:25;"`
		MicroserviceUrl     string `gorm:"column:microservice_url;"`
		MicroserviceTimeout int    `gorm:"column:microservice_timeout;size:4;"`
	}
	var query Microservice
	var msg string
	var title string
	// Deploying an SQL Query To Select All Rows Where Microservice Name iS Equal To The Inputted Name Made On Discord
	result := dbconfig.DB.Where("microservice_name = ?", name).Find(&query)
	// If Rows Affected Is Greater Than Zero Then The Microservice Record Exists
	if result.RowsAffected > 0 {
		//Deleting The Microservice Details Into The Postgres Database Where microservice_name iS Equal To The Inputted Name Made On Discord
		dbconfig.DB.Where("microservice_name = ?", name).Delete(&Microservice{})
		title = "Delete Command"
		msg = "Microservice: " + name + " Has Been Deleted"
		return title, msg
		// Error Handling When Rows Affected Is Equal To 0: No Microservice With That Name Exists
	} else {
		title = "Delete Command Error"
		msg = "Bot Name Does Not Exist"
		return title, msg
	}
}
