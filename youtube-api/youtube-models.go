package youtube_api

type YuoTubeSearchResponse struct {
	Items []struct {
		Id struct {
			Kind    string `json:"kind"`
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelId    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description'`
			ChannelTitle string `json:"channelTitle"`
			PublishTime  string `json:"publishTime"`
		} `json:"snippet"`
	} `json:"items"`
}
