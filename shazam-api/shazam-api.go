package shazam_api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pseudoelement/go-tg-music-bot/types"
)

type ShazamApiService struct {
	apiToken    string
	apiEndpoint string
}

func NewShazamApiService() (*ShazamApiService, error) {
	chat := &ShazamApiService{
		apiEndpoint: "https://shazam.p.rapidapi.com",
	}
	token, err := chat.GetApiToken()
	if err != nil {
		return nil, err
	}
	chat.apiToken = token

	return chat, nil
}

func (srv *ShazamApiService) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("SHAZAM_API_KEY")
	if !ok {
		return "", errors.New("SHAZAM_API_KEY doesn't exist!")
	}

	return token, nil
}

func (srv *ShazamApiService) QuerySimilarSongs(text string, isRetry bool) (string, error) {
	url := "https://shazam.p.rapidapi.com/shazam-songs/list-similarities?id=track-similarities-id-66398418&locale=en-US"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "c6a8fab7c8msh048c4df1ac026bep1830fcjsn9ac136edb848")
	req.Header.Add("X-RapidAPI-Host", "shazam.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	return string(body), nil
}

var _ types.MusicApiService = (*ShazamApiService)(nil)
