package common

import "time"

const (
	baseUrl              = "https://api.kraken.com"
	apiVersionV0         = "/0"
	epSystemStatus       = "/public/SystemStatus"
	epTicker             = "/public/Ticker"
	EndpointSystemStatus = baseUrl + apiVersionV0 + epSystemStatus
	TimestampFormat      = time.RFC3339
)

type KrakenStatus = string

type KrakenResponse[K interface{}] struct {
	Error  []string `json:"error"`
	Result K        `json:"result"`
}

type KrakenDataStatus struct {
	Status    KrakenStatus `json:"status"`
	Timestamp string       `json:"timestamp"`
}

const (
	KrakenStatusOnline      KrakenStatus = "online"
	KrakenStatusMaintenance KrakenStatus = "maintenance"
	KrakenStatusCancelOnly  KrakenStatus = "cancel_only"
	KrakenStatusPostOnly    KrakenStatus = "post_only"
)
