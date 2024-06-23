package youtube_api

import (
	"context"
	"errors"
	"fmt"
	"os"

	app_types "github.com/pseudoelement/go-tg-music-bot/src/common/types"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeApi struct {
	youtubeSrv *youtube.Service
}

func NewYouTubeApi() *YouTubeApi {
	srv := &YouTubeApi{}
	token, err := srv.GetApiToken()
	if err != nil {
		panic("No YOUTUBE_API_KEY's provided!")
	}
	yts, _ := youtube.NewService(context.Background(), option.WithAPIKey(token))
	srv.youtubeSrv = yts

	return srv
}

func (srv *YouTubeApi) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		return "", errors.New("env variable YOUTUBE_API_KEY doesn't exist!")
	}

	return token, nil
}

func (srv *YouTubeApi) QueryLinkByVideoName(videoName string) (string, error) {
	call := srv.youtubeSrv.Search.List([]string{"snippet"})
	searchList, err := call.MaxResults(5).Type("video").Q(videoName).Do()

	if err != nil || len(searchList.Items) < 1 {
		return "", errors.New("Video not found")
	}

	videoId := searchList.Items[0].Id.VideoId
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	return videoUrl, nil
}

var _ app_types.MusicLinkSearcher = (*YouTubeApi)(nil)
