package common

type KrakenDataTickerPair struct {
	Ask                        []string `json:"a"`
	Bid                        []string `json:"b"`
	LastTradeClosed            []string `json:"c"`
	Volume                     []string `json:"v"`
	VolumeWeightedAveragePrice []string `json:"p"`
	NumberOfTrades             []int    `json:"t"`
	Low                        []string `json:"l"`
	High                       []string `json:"h"`
	OpeningPriceToday          string   `json:"o"`
}

func (t KrakenDataTickerPair) AskPrice() string {
	return t.Ask[0]
}

func (t KrakenDataTickerPair) AskWholeLotVolume() string {
	return t.Ask[1]
}

func (t KrakenDataTickerPair) AskLotVolume() string {
	return t.Ask[2]
}

func (t KrakenDataTickerPair) BidPrice() string {
	return t.Bid[0]
}

func (t KrakenDataTickerPair) BidWholeLotVolume() string {
	return t.Bid[1]
}

func (t KrakenDataTickerPair) BidLotVolume() string {
	return t.Bid[2]
}

func (t KrakenDataTickerPair) LastTradeClosedPrice() string {
	return t.LastTradeClosed[0]
}

func (t KrakenDataTickerPair) LastTradeClosedLotVolume() string {
	return t.LastTradeClosed[1]
}

func (t KrakenDataTickerPair) VolumeToday() string {
	return t.Volume[0]
}

func (t KrakenDataTickerPair) VolumeLast24Hours() string {
	return t.Volume[1]
}

func (t KrakenDataTickerPair) VolumeWeightedAveragePriceToday() string {
	return t.VolumeWeightedAveragePrice[0]
}

func (t KrakenDataTickerPair) VolumeWeightedAveragePriceLast24Hours() string {
	return t.VolumeWeightedAveragePrice[1]
}

func (t KrakenDataTickerPair) NumberOfTradesToday() int {
	return t.NumberOfTrades[0]
}

func (t KrakenDataTickerPair) NumberOfTradesLast24Hours() int {
	return t.NumberOfTrades[1]
}

func (t KrakenDataTickerPair) LowToday() string {
	return t.Low[0]
}

func (t KrakenDataTickerPair) LowLast24Hours() string {
	return t.Low[1]
}

func (t KrakenDataTickerPair) HighToday() string {
	return t.High[0]
}

func (t KrakenDataTickerPair) HighLast24Hours() string {
	return t.High[1]
}
