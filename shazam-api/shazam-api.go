package shazam_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pseudoelement/go-tg-music-bot/types"
	"github.com/pseudoelement/go-tg-music-bot/utils"
)

type ShazamApiService struct {
	apiToken    string
	apiHost     string
	apiEndpoint string
}

func NewShazamApiService() (*ShazamApiService, error) {
	chat := &ShazamApiService{
		apiEndpoint: "https://shazam.p.rapidapi.com",
		apiHost:     "shazam.p.rapidapi.com",
	}
	token, err := chat.GetApiToken()
	if err != nil {
		return nil, err
	}
	chat.apiToken = token

	return chat, nil
}

func (srv *ShazamApiService) makeGetRequest(subUrl string, params map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", srv.apiEndpoint, subUrl)
	req, _ := http.NewRequest("GET", url, nil)

	queryParams := req.URL.Query()
	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

	req.Header.Add("X-RapidAPI-Key", srv.apiToken)
	req.Header.Add("X-RapidAPI-Host", "shazam.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, nil
}

func (srv *ShazamApiService) querySongId(songName string) (string, error) {
	p := map[string]string{"term": songName}
	resBytes, err := srv.makeGetRequest("search", p)
	var searchResponse SearchQueryResponse
	if err = json.Unmarshal(resBytes, &searchResponse); err != nil {
		return "", utils.Error(err.Error(), "querySongShazamId")
	}

	return searchResponse.Tracks.Hits[0].Track.Key, nil
}

func (srv *ShazamApiService) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("SHAZAM_API_KEY")
	if !ok {
		return "", errors.New("SHAZAM_API_KEY doesn't exist!")
	}

	return token, nil
}

func (srv *ShazamApiService) QuerySimilarSongs(songName string, isRetry bool) (string, error) {
	songId, err := srv.querySongId(songName)
	if err != nil {
		return "", utils.Error(err.Error(), "QuerySimilarSongs")
	}

	validSongId := fmt.Sprintf("track-similarities-id-%s", songId)
	p := map[string]string{"id": validSongId}
	resBytes, err := srv.makeGetRequest("shazam-songs/list-similarities", p)
	var similaritiesResponse ListSimilaritiesResponse
	if err = json.Unmarshal(resBytes, &similaritiesResponse); err != nil {
		return "", utils.Error(err.Error(), "QuerySimilarSongs")
	}

	var list string
	var count int
	for _, value := range similaritiesResponse.Resources.ShazamSongs {
		count++
		str := fmt.Sprintf("%v. %s - %s\n", count, value.Attributes.Artist, value.Attributes.Title)
		list += str
	}

	return list, nil
}

var _ types.MusicApiService = (*ShazamApiService)(nil)
