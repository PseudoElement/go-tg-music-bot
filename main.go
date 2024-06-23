package main

import (
	"fmt"

	"github.com/joho/godotenv"
	bot_manager "github.com/pseudoelement/go-tg-music-bot/src/bot-manager"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	botManager := bot_manager.NewBotManager()
	botManager.Broadcast()
}
