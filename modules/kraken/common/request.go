package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetKrakenStatus() (*KrakenResponse[KrakenDataStatus], error) {
	return executeRequest[KrakenDataStatus](EndpointSystemStatus)
}

func GetKrakenTicker(pairs ...string) (*KrakenResponse[KrakenDataTicker], error) {
	assetPairs := strings.Join(pairs, ",")
	url := fmt.Sprintf("%s?pair=%s", EndpointTicker, assetPairs)
	return executeRequest[KrakenDataTicker](url)
}

func executeRequest[K interface{}](url string) (*KrakenResponse[K], error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data := KrakenResponse[K]{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
