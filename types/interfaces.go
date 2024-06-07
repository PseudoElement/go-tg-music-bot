package types

type MusiApiFields struct {
	ApiToken string
}

type MusicApiService interface {
	GetApiToken() (string, error)
	QuerySimilarSongs(msg string, isRetry bool, needLinks bool) (string, error)
	QuerySongByKeyWords(msg string, needLinks bool) (string, error)
}
