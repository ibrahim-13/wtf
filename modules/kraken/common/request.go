package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	__req_duration = 1 * time.Second
)

var __req_time time.Time = time.Now().Add(-1 * time.Second)
var __req_time_lock sync.Mutex

func GetKrakenStatus() (*KrakenResponse[KrakenDataStatus], error) {
	return executeRequest[KrakenDataStatus](EndpointSystemStatus)
}

func GetKrakenTicker(pairs ...string) (*KrakenResponse[KrakenDataTicker], error) {
	assetPairs := strings.Join(pairs, ",")
	url := fmt.Sprintf("%s?pair=%s", EndpointTicker, assetPairs)
	return executeRequest[KrakenDataTicker](url)
}

func executeRequest[K interface{}](url string) (*KrakenResponse[K], error) {
	__req_time_lock.Lock()
	now := time.Now()
	if now.Add(__req_duration).After(__req_time) {
		<-time.After(__req_duration)
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data := KrakenResponse[K]{}
	__req_time = time.Now()
	__req_time_lock.Unlock()
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
