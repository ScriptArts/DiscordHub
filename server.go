package main

import (
	"github.com/ScriptArts/DiscordHub/app"
	"github.com/ScriptArts/DiscordHub/models"
	"github.com/ScriptArts/DiscordHub/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Llongfile)
	err := utils.LoadEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// discord setting
	discord, err := utils.GetDiscordClient()
	if err != nil {
		log.Fatalln(err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}

	models.GetDatabase()
	models.Migration()

	discord.AddHandler(app.Handler)

	n := utils.GetNatsConnection()

	// システムが終了させられるまで起動し続ける
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
	n.Close()
}
