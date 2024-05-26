package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		panic("BOT_TOKEN doesn't exists!")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	fmt.Printf("Authorized on account %s!\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 100

	updates := bot.GetUpdatesChan(u)

	// chatGPT, err := ai.NewChatGPTModule()
	// if err != nil {
	// 	msg := "Error in ai.NewChatGPTModule() - " + err.Error()
	// 	panic(msg)
	// }

	for update := range updates {
		if update.Message != nil {
			fmt.Println("Name: ", update.Message.From.FirstName)
			// fmt.Println("Video from TG: ", update.Message.Video)
			// fmt.Println("Photo from TG: ", update.Message.Photo)

			// response, err := chatGPT.MakeQuery(update.Message.Text, false)

			// var msg tgbotapi.MessageConfig
			// if err != nil {
			// 	msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			// } else {
			// 	msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
			// }

			response := fakeRequest()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)

			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func fakeRequest() string {
	apiKey := os.Getenv("CHAT_GPT_TOKEN")
	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": "Hi can you tell me what is the factorial of 10?"}},
			"max_tokens": 50,
		}).
		Post("https://api.openai.com/v1/chat/completions")
	fmt.Println("AFTER REQUEST  ========")

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()
	fmt.Println("BODY ====> ", body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return ""
	}

	fmt.Println("DATA ß=======> ", data)
	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	fmt.Println(content)
	return content
}
