package shazam_api

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/pseudoelement/go-tg-music-bot/src/common/api"
	app_types "github.com/pseudoelement/go-tg-music-bot/src/common/types"
	app_utils "github.com/pseudoelement/go-tg-music-bot/src/common/utils"
	youtube_api "github.com/pseudoelement/go-tg-music-bot/src/services/youtube-api"
)

type ShazamApiService struct {
	apiToken     string
	apiHost      string
	apiEndpoint  string
	musicLinkSrv app_types.MusicLinkSearcher
}

func NewShazamApiService() (*ShazamApiService, error) {
	chat := &ShazamApiService{
		apiEndpoint:  "https://shazam.p.rapidapi.com",
		apiHost:      "shazam.p.rapidapi.com",
		musicLinkSrv: youtube_api.NewYouTubeApi(),
	}
	token, err := chat.GetApiToken()
	if err != nil {
		return nil, err
	}
	chat.apiToken = token

	return chat, nil
}

func (srv *ShazamApiService) querySongInfo(songName string) (SearchQueryResponse, error) {
	p := map[string]string{"term": songName}
	h := map[string]string{"X-RapidAPI-Key": srv.apiToken, "X-RapidAPI-Host": srv.apiHost}
	url := srv.apiEndpoint + "/search"
	searchResponse, err := api.Get[SearchQueryResponse](url, p, h)
	if err != nil {
		return SearchQueryResponse{}, app_utils.SimilarSongsNotFound()
	}

	return searchResponse, nil
}

func (srv *ShazamApiService) querySimilarityList(songName string) (ListSimilaritiesResponse, error) {
	songInfo, err := srv.querySongInfo(songName)
	if err != nil || len(songInfo.Tracks.Hits) < 1 {
		return ListSimilaritiesResponse{}, err
	}

	validSongId := fmt.Sprintf("track-similarities-id-%s", songInfo.Tracks.Hits[0].Track.Key)
	p := map[string]string{"id": validSongId}
	h := map[string]string{"X-RapidAPI-Key": srv.apiToken, "X-RapidAPI-Host": srv.apiHost}
	url := srv.apiEndpoint + "/shazam-songs/list-similarities"
	similaritiesResponse, err := api.Get[ListSimilaritiesResponse](url, p, h)
	if err != nil {
		return ListSimilaritiesResponse{}, errors.New("Similar songs not found!")
	}

	return similaritiesResponse, nil
}

func (srv *ShazamApiService) GetApiToken() (string, error) {
	token, ok := os.LookupEnv("SHAZAM_API_KEY")
	if !ok {
		return "", errors.New("SHAZAM_API_KEY doesn't exist!")
	}

	return token, nil
}

func (srv *ShazamApiService) QuerySimilarSongs(songName string, isRetry bool) (string, error) {
	similaritiesResponse, err := srv.querySimilarityList(songName)
	if err != nil {
		return "", err
	} else if len(similaritiesResponse.Resources.ShazamSongs) < 1 {
		return "", app_utils.SimilarSongsNotFound()
	}

	var list string
	var count int
	for _, value := range similaritiesResponse.Resources.ShazamSongs {
		if count > 9 {
			break
		}
		count++
		str := fmt.Sprintf("%v. %s - %s\n", count, value.Attributes.Artist, value.Attributes.Title)
		list += str
	}

	return list, nil
}

func (srv *ShazamApiService) QuerySongByKeyWords(keyWord string) (string, error) {
	searchResponse, err := srv.querySongInfo(keyWord)
	if err != nil {
		return "", err
	} else if len(searchResponse.Tracks.Hits) < 1 {
		return "", app_utils.SimilarSongsNotFound()
	}

	var list string
	for ind, value := range searchResponse.Tracks.Hits {
		if ind > 9 {
			break
		}
		str := fmt.Sprintf("%v. %s - %s\n", ind+1, value.Track.Subtitle, value.Track.Title)
		list += str
	}

	return list, nil
}

func (srv *ShazamApiService) QuerySimilarSongsLinks(songName string) (string, error) {
	similaritiesResponse, err := srv.querySimilarityList(songName)
	if err != nil {
		return "", err
	} else if len(similaritiesResponse.Resources.ShazamSongs) < 1 {
		return "", app_utils.SimilarSongsNotFound()
	}
	chanCapacity := len(similaritiesResponse.Resources.ShazamSongs)
	listCh := make(chan string, chanCapacity)

	var wg sync.WaitGroup
	var ind int
	for _, value := range similaritiesResponse.Resources.ShazamSongs {
		if ind > 9 {
			break
		}
		ind++
		wg.Add(1)
		go func() {
			fullSongName := value.Attributes.Artist + " - " + value.Attributes.Title
			songRow := srv.getListRow(fullSongName)
			listCh <- songRow
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(listCh)

	if len(listCh) < 1 {
		return "", errors.New("QuerySimilarSongsLinks - List creation error!")
	}

	var list string
	var count int
	for songInfo := range listCh {
		count++
		list += fmt.Sprintf("%v. %s\n\n", count, songInfo)
	}
	return list, nil
}

func (srv *ShazamApiService) QuerySongByKeyWordsLinks(msg string) (string, error) {
	searchResponse, err := srv.querySongInfo(msg)
	if err != nil {
		return "", err
	} else if len(searchResponse.Tracks.Hits) < 1 {
		return "", app_utils.SimilarSongsNotFound()
	}
	chanCapacity := len(searchResponse.Tracks.Hits)
	listCh := make(chan string, chanCapacity)

	var wg sync.WaitGroup
	for ind, value := range searchResponse.Tracks.Hits {
		if ind > 9 {
			break
		}
		wg.Add(1)
		go func() {
			fullSongName := value.Track.Subtitle + " - " + value.Track.Title
			songRow := srv.getListRow(fullSongName)
			listCh <- songRow
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(listCh)

	if len(listCh) < 1 {
		return "", errors.New("QuerySongByKeyWordsLinks - List creation error!")
	}

	var list string
	var count int
	for songInfo := range listCh {
		count++
		list += fmt.Sprintf("%v. %s\n\n", count, songInfo)
	}

	return list, nil
}

func (srv *ShazamApiService) getListRow(song string) string {
	link, err := srv.musicLinkSrv.QueryLinkByVideoName(song)
	if link == "" || err != nil {
		link = "Ссылка не найдена."
	}
	songRow := fmt.Sprintf(`Название -  %s.
	Ссылка - %v`, song, link)
	return songRow
}

var _ app_types.MusicApiService = (*ShazamApiService)(nil)
