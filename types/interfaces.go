package types

type MusiApiFields struct {
	ApiToken string
}

type MusicApiService interface {
	GetApiToken() (string, error)
	QuerySimilarSongs(msg string, isRetry bool) (string, error)
	QuerySimilarSongsLinks(msg string) (string, error)
	QuerySongByKeyWords(msg string) (string, error)
	QuerySongByKeyWordsLinks(msg string) (string, error)
}
