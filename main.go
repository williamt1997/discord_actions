package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanjcook/discordbotgo/bot"
	dbconfig "github.com/nathanjcook/discordbotgo/config"
	migration "github.com/nathanjcook/discordbotgo/utils"
	"go.uber.org/zap"
)

var Sugar *zap.Logger

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	// Find .env file
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load(".env")
		if err != nil {
			zap.L().Panic("Error loading .env file:", zap.Error(err))
		}
	}
	migration.BuildDB()
	// Connect to DB on app start up
	dbconfig.Connect()
}

func main() {

	bot.Start()
}
