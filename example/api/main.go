package main

import (
	"log"
	"net/http"
	"time"

	"github.com/buurzx/go-bybit-api/api"
)

func main() {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	client := api.New(httpClient, "apikey", "secretKey", api.DebugMode(true))

	_, _, result, err := client.LinearGetKLine("BTCUSDT", "5")
	if err != nil {
		log.Fatalf("failed to get klines %v \n", err)
	}

	_ = result
}
