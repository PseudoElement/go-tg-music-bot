package youtube_api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeApi struct {
	apiTokens         map[string]YouTubeToken
	selectedTokenName string
	useAnotherToken   bool
	youtubeSrv        *youtube.Service
	mu                sync.Mutex
}

func NewYouTubeApi() *YouTubeApi {
	srv := &YouTubeApi{}
	tokens := srv.GetApiTokens()
	if len(tokens) < 1 {
		panic("No YOUTUBE_API_KEY's provided!")
	}
	srv.apiTokens = tokens
	srv.selectTokenWithMinUseCount()

	return srv
}

func (srv *YouTubeApi) selectTokenWithMinUseCount() {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	var tokenName string
	var minUseCount float64 = math.MaxFloat64
	for name, token := range srv.apiTokens {
		if token.UseCount < minUseCount {
			minUseCount = token.UseCount
			tokenName = name
		}
	}

	srv.selectedTokenName = tokenName
	token := srv.apiTokens[tokenName]
	yts, _ := youtube.NewService(context.Background(), option.WithAPIKey(token.Value))
	srv.youtubeSrv = yts
}

func (srv *YouTubeApi) GetApiTokens() map[string]YouTubeToken {
	tokens := make(map[string]YouTubeToken)
	for i := 1; i <= 10; i++ {
		envName := "YOUTUBE_API_KEY_" + strconv.Itoa(i)
		token, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}
		tokens[envName] = YouTubeToken{
			Value:    token,
			UseCount: 0,
		}
	}

	return tokens
}

func (srv *YouTubeApi) QueryLinkByVideoName(videoName string) (string, error) {
	srv.selectTokenWithMinUseCount()

	call := srv.youtubeSrv.Search.List([]string{"snippet"})
	searchList, err := call.MaxResults(5).Type("video").Q(videoName).Do()

	srv.mu.Lock()
	token, _ := srv.apiTokens[srv.selectedTokenName]
	token.UseCount++
	srv.apiTokens[srv.selectedTokenName] = token
	srv.mu.Unlock()

	if err != nil || len(searchList.Items) < 1 {
		srv.useAnotherToken = true
		return "", errors.New("Video not found")
	}

	videoId := searchList.Items[0].Id.VideoId
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	return videoUrl, nil
}
