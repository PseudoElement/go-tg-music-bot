package shazam_api

type SearchQueryResponse struct {
	Tracks struct {
		Hits []struct {
			Track struct {
				Layout   string `json:"layout"`
				Type     string `json:"type"`
				Key      string `json:"key"`
				Title    string `json:"title"`
				Subtitle string `json:"subtitle"`
			} `json:"track"`
		} `json:"hits"`
	} `json:"tracks"`
	Artists any `json:"artists"`
}

type GetDetailsResponse struct {
	Data      any `json:"data"`
	Resources struct {
		RelatedTracks map[string]string `json:"related-tracks"`
	} `json:"resources"`
}

type ListSimilaritiesResponse struct {
	Data      any `json:"data"`
	Resources struct {
		ShazamSongs map[string]struct {
			Id         string `json:"id"`
			Attributes struct {
				Title  string `json:"title"`
				Artist string `json:"artist"`
				WebUrl string `json:"webUrl"`
			} `json:"attributes"`
		} `json:"shazam-songs"`
	} `json:"resources"`
}
