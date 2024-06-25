package spotify_api

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	// In seconds from now
	ExpiresIn int `json:"expires_in"`
}

type SpotifySearchResponse struct {
	Tracks struct {
		Items []struct {
			Album struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"album"`
		} `json:"items"`
	} `json:"tracks"`
}
