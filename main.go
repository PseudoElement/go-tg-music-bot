package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/pseudoelement/go-tg-music-bot/ai"
	shazam_api "github.com/pseudoelement/go-tg-music-bot/shazam-api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		panic("BOT_TOKEN doesn't exist!")
	}

	useChatGPT := needUseChatGPT()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s!\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 100
	updates := bot.GetUpdatesChan(u)

	var chatGPT *ai.ChatGPT
	if useChatGPT {
		chatGPT, err = ai.NewChatGPTService()
		if err != nil {
			panic(err)
		}
	}

	shazamApi, err := shazam_api.NewShazamApiService()
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			fmt.Println("Name: ", update.Message.From.FirstName)
			// fmt.Println("Video from TG: ", update.Message.Video)
			// fmt.Println("Photo from TG: ", update.Message.Photo)
			msg := handleQuery(update, shazamApi, chatGPT, useChatGPT)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func needUseChatGPT() bool {
	useAiEnv := os.Getenv("USE_AI")
	var useChatGPT bool
	if len(useAiEnv) == 0 {
		useChatGPT = false
	} else {
		boolean, err := strconv.ParseBool(useAiEnv)
		if err != nil {
			fmt.Println("Incorrect format of .env var USE_AI!")
			useChatGPT = false
		}
		useChatGPT = boolean
	}

	return useChatGPT
}

func handleQuery(update tgbotapi.Update, shazamApi *shazam_api.ShazamApiService, chatGPT *ai.ChatGPT, useChatGPT bool) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var response string
	var err error
	if useChatGPT {
		response, err = chatGPT.QuerySimilarSongs(update.Message.Text, false)
	} else {
		response, err = shazamApi.QuerySimilarSongs(update.Message.Text, false)
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Некорректный запрос, попробуй еще раз! Текст ошибки(для разработчика): %s", err.Error())
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
	}

	return msg
}
