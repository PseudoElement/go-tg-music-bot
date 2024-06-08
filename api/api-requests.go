package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pseudoelement/go-tg-music-bot/utils"
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
	fmt.Println("\nGET_RESPONSE_STRING ============> ", string(resBytes))

	res_struct := new(T)
	if err := json.Unmarshal(resBytes, &res_struct); err != nil {
		return *res_struct, err
	} else if string(resBytes) == "{}" {
		return *res_struct, utils.EmptyApiResponse()
	}

	defer res.Body.Close()

	return *res_struct, nil
}
