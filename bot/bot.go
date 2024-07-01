package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards"
)

var Mgr *shards.Manager
var BotId string

func Start() {
	var err error

	// Create a new shard manager using the provided bot token.
	Mgr, err = shards.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("[ERROR] Error creating manager,", err)
		return
	}

	// Get details of bot to check if created
	user, err := Mgr.Gateway.User("@me")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	// Set BotId
	BotId = user.ID

	// Register handler for messages
	Mgr.AddHandler(messageCreate)
	// Register handler for sharding
	Mgr.AddHandler(onConnect)

	// Listen for messages in channels bot is member of Listen for DM messages & Listen for Messages
	Mgr.RegisterIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent)

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = Mgr.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Shut down bot
	fmt.Println("[INFO] Stopping shard manager...")
	err = Mgr.Shutdown()
	if err != nil {
		panic(err)
	}
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")
}
