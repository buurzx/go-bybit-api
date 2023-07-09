package api

import (
	"fmt"
	"strconv"
)

type BaseResult struct {
	RetCode          int         `json:"ret_code"`
	RetMsg           string      `json:"ret_msg"`
	ExtCode          string      `json:"ext_code"`
	Result           interface{} `json:"result"`
	Time             int64       `json:"time"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

// type OHLC struct {
// 	Symbol   string  `json:"symbol"`
// 	Interval string  `json:"interval"`
// 	OpenTime int64   `json:"open_time"`
// 	Open     float64 `json:"open,string"`
// 	High     float64 `json:"high,string"`
// 	Low      float64 `json:"low,string"`
// 	Close    float64 `json:"close,string"`
// 	Volume   float64 `json:"volume,string"`
// 	Turnover float64 `json:"turnover,string"`
// }

// type GetKlineResult struct {
// 	BaseResult
// 	Result []OHLC `json:"result"`
// }

type OHLCLinear struct {
	Symbol   string     `json:"symbol"`
	Category string     `json:"category"`
	List     [][]string `json:"list"`
	Klines   []OHLC
}

// FIXME:
func (o *OHLCLinear) DecodeList() error {
	var ohlcs []OHLC

	for _, ohlc := range o.List {
		openTime, err := strconv.Atoi(ohlc[0])
		if err != nil {
			return fmt.Errorf("decode list open time %v", err)
		}

		open, err := strconv.ParseFloat(ohlc[1], 64)
		if err != nil {
			return fmt.Errorf("decode list open %v", err)
		}

		high, err := strconv.ParseFloat(ohlc[2], 64)
		if err != nil {
			return fmt.Errorf("decode list high %v", err)
		}

		low, err := strconv.ParseFloat(ohlc[3], 64)
		if err != nil {
			return fmt.Errorf("decode list low %v", err)
		}
		close, err := strconv.ParseFloat(ohlc[4], 64)
		if err != nil {
			return fmt.Errorf("decode list close %v", err)
		}

		volume, err := strconv.ParseFloat(ohlc[5], 64)
		if err != nil {
			return fmt.Errorf("decode list volume %v", err)
		}

		turnover, err := strconv.ParseFloat(ohlc[6], 64)
		if err != nil {
			return fmt.Errorf("decode list turnover %v", err)
		}

		ohlc := OHLC{
			OpenTime: int64(openTime),
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
			Turnover: turnover,
		}

		ohlcs = append(ohlcs, ohlc)
	}

	o.Klines = ohlcs
	o.List = nil

	return nil
}

type OHLC struct {
	OpenTime int64   `json:"open_time"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	Volume   float64 `json:"volume"`
	Turnover float64 `json:"turnover"`
}

type GetLinearKlineDirtyResult struct {
	BaseResult
	Result OHLCLinear `json:"result"`
}
