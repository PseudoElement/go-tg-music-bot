package youtube_api

import (
	"errors"
	"fmt"
	"os"

	"github.com/pseudoelement/go-tg-music-bot/api"
)

type YouTubeApi struct {
	apiToken    string
	apiEndpoint string
}

func NewYouTubeApi() *YouTubeApi {
	api := &YouTubeApi{}
	token, err := api.GetApiToken()
	if err != nil {
		panic(err)
	}
	api.apiToken = token

	return api
}

func (srv *YouTubeApi) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		return "", errors.New("YOUTUBE_API_KEY doesn't exist!")
	}

	return token, nil
}

// @TODO fix link search
func (srv *YouTubeApi) QueryLinkByVideoName(videoName string) (string, error) {
	p := map[string]string{"part": "snippet", "type": "video", "key": srv.apiToken, "q": videoName}
	h := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", srv.apiToken)}
	res, err := api.Get[YuoTubeSearchResponse]("https://youtube.googleapis.com/youtube/v3/search", p, h)
	fmt.Println("\nQueryLinkByVideoName-Response =========> ", res)
	if err != nil || len(res.Items) < 1 {
		msg := fmt.Sprintf("Video `%s` not found", videoName)
		return "", errors.New(msg)
	}
	videoId := res.Items[0].Id.VideoId
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	return videoUrl, nil
}
