package ws

type KLine struct {
	Start     int64   `json:"start"` // 1539918000
	End       int64   `json:"end"`   // 1539918000
	Open      float64 `json:"open,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Close     float64 `json:"close,string"`
	Volume    float64 `json:"volume,string"`
	Turnover  float64 `json:"turnover,string"` // 0.0013844
	Interval  string  `json:"interval"`        // 1m
	Confirm   bool    `json:"confirm"`
	Timestamp int64   `json:"timestamp"`
}
