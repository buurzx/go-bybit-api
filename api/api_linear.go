package api

import (
	"net/http"
)

const linearCategory = "linear"

// LinearGetKLine
func (b *ByBit) LinearGetKLine(symbol string, interval string) (query string, resp []byte, result OHLCLinear, err error) {
	var ret GetLinearKlineDirtyResult

	params := map[string]interface{}{}
	params["category"] = linearCategory
	params["symbol"] = symbol
	params["interval"] = interval

	query, resp, err = b.PublicRequest(http.MethodGet, "v5/market/kline", params, &ret)
	if err != nil {
		return
	}

	result = ret.Result
	// FIXME:
	err = result.DecodeList()

	return
}

// ParametrizedLinearGetKLine
func (b *ByBit) ParametrizedLinearGetKLine(symbol string, interval string, start int64, end int64, limit uint) (query string, resp []byte, result OHLCLinear, err error) {
	var ret GetLinearKlineDirtyResult

	params := map[string]interface{}{}
	params["category"] = linearCategory
	params["symbol"] = symbol
	params["interval"] = interval
	params["start"] = start
	params["end"] = end
	params["limit"] = limit

	query, resp, err = b.PublicRequest(http.MethodGet, "v5/market/kline", params, &ret)
	if err != nil {
		return
	}

	result = ret.Result
	// FIXME:
	err = result.DecodeList()

	return
}
