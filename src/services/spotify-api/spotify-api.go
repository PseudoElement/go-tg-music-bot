package spotify_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pseudoelement/go-tg-music-bot/src/common/api"
	app_types "github.com/pseudoelement/go-tg-music-bot/src/common/types"
	app_utils "github.com/pseudoelement/go-tg-music-bot/src/common/utils"
)

type SpotifyApi struct {
	apiToken          string
	tokenExpirationMS int
}

func NewSpotifyApi() *SpotifyApi {
	time.Now()
	api := &SpotifyApi{
		tokenExpirationMS: 0,
	}
	return api
}

func (sa *SpotifyApi) isTokenExpired() bool {
	return sa.tokenExpirationMS < int(time.Now().UnixMilli())
}

func (sa *SpotifyApi) fetchNewApiToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
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
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	sa.apiToken = tokenResponse.AccessToken
	sa.tokenExpirationMS = int(time.Now().UnixMilli()) + tokenResponse.ExpiresIn*1_000

	return nil
}

func (sa *SpotifyApi) makeSearchRequest(videoName string) (string, error) {
	bearerToken := fmt.Sprintf("Bearer %s", sa.apiToken)
	headers := map[string]string{"Authorization": bearerToken}
	params := map[string]string{"q": videoName, "limit": "1", "type": "track"}

	res, err := api.Get[SpotifySearchResponse]("https://api.spotify.com/v1/search", params, headers)
	if err != nil {
		return "", app_utils.Error(err.Error(), "QueryLinkByVideoName")
	} else if len(res.Tracks.Items) == 0 {
		return "", app_utils.EmptyApiResponse()
	}

	return res.Tracks.Items[0].Album.ExternalUrls.Spotify, nil
}

func (sa *SpotifyApi) QueryLinkByVideoName(videoName string) (string, error) {
	if sa.isTokenExpired() {
		err := sa.fetchNewApiToken()
		if err != nil {
			return "", err
		}
	}

	link, err := sa.makeSearchRequest(videoName)
	if err != nil {
		return "", err
	}

	return link, nil
}

var _ app_types.MusicLinkSearcher = (*SpotifyApi)(nil)
