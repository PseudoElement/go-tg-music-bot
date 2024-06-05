package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pseudoelement/go-tg-music-bot/types"
	"github.com/pseudoelement/go-tg-music-bot/utils"
)

type ChatGPT struct {
	client      *resty.Client
	apiToken    string
	apiEndpoint string
	retryCount  int
}

func NewChatGPTService() (*ChatGPT, error) {
	client := resty.New()
	chat := &ChatGPT{
		client:      client,
		apiEndpoint: "https://api.openai.com/v1/chat/completions",
		retryCount:  0,
	}
	token, err := chat.GetApiToken()
	if err != nil {
		return nil, err
	}
	chat.apiToken = token

	return chat, nil
}

func (c *ChatGPT) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("CHAT_GPT_TOKEN")
	if !ok {
		return "", errors.New("CHAT_GPT_TOKEN doesn't exist!")
	}

	return token, nil
}

func (c *ChatGPT) formatMessageToChatGPT(msg string) string {
	startPart := fmt.Sprintf("Give me ten similar songs on this one - %v", msg)
	lastPart := `
        Give me the answer in such format: 
        1. Song1.
        2. Song2.
        3. Song3.
        ...
        10. Song10!!!
        etc.
        !!! signs are required me to parse this list from full answer.
        Where Song1, Song2, Song3 ... are suggested song names, that you found.
    `
	if strings.HasSuffix(msg, ".") || strings.HasSuffix(msg, "?") || strings.HasSuffix(msg, "!") {
		msg = startPart + msg + ". " + lastPart
	} else {
		msg += startPart + msg + lastPart
	}
	return msg
}

func (c *ChatGPT) getSongsListFromResponse(data map[string]interface{}) (string, error) {
	// Extract the content from the JSON response
	content, ok := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", utils.Error("Can't parse content to string", "getSongsListFromResponse")
	}

	listStartIndex := utils.IndexOfSubstring(content, "1.")
	listEndIndex := utils.IndexOfSubstring(content, "!!!")

	if listStartIndex == -1 || listEndIndex == -1 {
		return "", utils.Error("Invalid ChatGPT response type!", "getSongsListFromResponse")
	}

	list := content[listStartIndex:listEndIndex]

	return list, nil
}

func (c *ChatGPT) QuerySimilarSongs(msg string, isRetry bool) (string, error) {
	c.retryCount++

	msg = c.formatMessageToChatGPT(msg)
	response, err := c.client.R().
		SetAuthToken(c.apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": msg}},
			"max_tokens": 10_000,
		}).
		Post(c.apiEndpoint)

	if err != nil {
		return "", utils.Error(err.Error(), "c.client.R")
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", utils.UnmarshalError(err.Error(), "MakeQuery")
	}

	songsList, err := c.getSongsListFromResponse(data)
	if err != nil {
		if errors.Is(err, utils.InvalidAiResponseFormat()) && c.retryCount < 3 {
			newResponse, err := c.QuerySimilarSongs(msg, true)
			if err != nil {
				return "", err
			}
			return newResponse, nil
		}
		return "", err
	}

	fmt.Println("==========================")
	fmt.Println(songsList)
	fmt.Println("==========================")

	return songsList, nil
}

func (srv *ChatGPT) QuerySongByKeyWords(msg string) (string, error) {
	return "", utils.MethodNotImplemented()
}

var _ types.MusicApiService = (*ChatGPT)(nil)
