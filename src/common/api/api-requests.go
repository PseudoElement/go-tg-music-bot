package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	app_utils "github.com/pseudoelement/go-tg-music-bot/src/common/utils"
)

func Get[T any](url string, params map[string]string, headers map[string]string) (response T, e error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	queryParams := req.URL.Query()

	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	resBytes, _ := io.ReadAll(res.Body)
	// fmt.Println("\nGET_RESPONSE_STRING ============> ", string(resBytes))

	res_struct := new(T)
	if err := json.Unmarshal(resBytes, &res_struct); err != nil {
		return *res_struct, err
	} else if string(resBytes) == "{}" {
		return *res_struct, app_utils.EmptyApiResponse()
	}

	defer res.Body.Close()

	return *res_struct, nil
}

func Post[T any](url string, body interface{}, headers map[string]string) (T, error) {
	client := &http.Client{}

	var bodyBuffer bytes.Buffer
	err := json.NewEncoder(&bodyBuffer).Encode(body)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, url, &bodyBuffer)

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	bytes, err := io.ReadAll(resp.Body)

	var response T
	if err := json.Unmarshal(bytes, &response); err != nil {
		return response, err
	}

	defer resp.Body.Close()
	return response, nil
}
