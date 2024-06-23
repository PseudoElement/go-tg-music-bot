package spotify_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	app_types "github.com/pseudoelement/go-tg-music-bot/src/common/types"
)

type SpotifyApi struct {
	apiTokenTTL int16
}

func NewSpotifyApi() *SpotifyApi {
	api := &SpotifyApi{
		apiTokenTTL: 3600,
	}
	return api
}

func (sa *SpotifyApi) fetchTokenInfo() (TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		panic("Missed SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET in .env file!")
	}

	req.SetBasicAuth(clientId, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("failed to get token: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return TokenResponse{}, err
	}

	return tokenResponse, nil
}

func (sa *SpotifyApi) QueryLinkByVideoName(videoName string) (string, error) {
	return "string", nil
}

var _ app_types.MusicLinkSearcher = (*SpotifyApi)(nil)
