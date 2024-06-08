package youtube_api

import (
	"errors"
	"fmt"
	"os"

	"github.com/pseudoelement/go-tg-music-bot/api"
)

type YouTubeApi struct {
	apiToken       string
	apiBearerToken string
	apiEndpoint    string
}

func NewYouTubeApi() *YouTubeApi {
	api := &YouTubeApi{}
	token, bearer, err := api.GetApiTokens()
	if err != nil {
		panic(err)
	}
	api.apiToken = token
	api.apiBearerToken = bearer

	return api
}

func (srv *YouTubeApi) GetApiTokens() (string, string, error) {
	token, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		return "", "", errors.New("YOUTUBE_API_KEY doesn't exist!")
	}

	bearer, ok := os.LookupEnv("YOUTUBE_BEARER_TOKEN")
	if !ok {
		return "", "", errors.New("YOUTUBE_BEARER_TOKEN doesn't exist!")
	}

	return token, bearer, nil
}

// @TODO fix link search
func (srv *YouTubeApi) QueryLinkByVideoName(videoName string) (string, error) {
	p := map[string]string{"part": "snippet", "type": "video", "key": srv.apiToken, "q": videoName}
	h := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", srv.apiBearerToken)}
	res, err := api.Get[YuoTubeSearchResponse]("https://youtube.googleapis.com/youtube/v3/search", p, h)
	if err != nil || len(res.Items) < 1 {
		msg := fmt.Sprintf("Video `%s` not found", videoName)
		return "", errors.New(msg)
	}
	videoId := res.Items[0].Id.VideoId
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	return videoUrl, nil
}
