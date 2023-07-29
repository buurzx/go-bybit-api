package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buurzx/go-bybit-api/ws"
)

func main() {
	sigCH := make(chan os.Signal, 1)
	signal.Notify(sigCH, os.Interrupt, syscall.SIGTERM)

	cfg := &ws.Configuration{
		Addr:          ws.HostPerpetualtTestnet,
		AutoReconnect: true,
		DebugMode:     true,
	}

	bbWS := ws.New(cfg)
	// bbWS.Subscribe(ws.WSKLine, "1", "BTCUSDT")
	err := bbWS.SubscribeToMany(
		"orderbook.1.BTCUSDT",
		"publicTrade.BTCUSDT",
		"orderbook.1.ETHUSDT",
	)
	if err != nil {
		panic(err)
	}

	bbWS.On(ws.WSKLine, handleKLine)

	bbWS.Start()

	<-sigCH
	fmt.Println("Shutdown ...")

	bbWS.Close()
	<-time.After(time.Second * 1)
}

func handleKLine(symbol string, data ws.KLine) {
	log.Printf("handleKLine %v/%#v \n", symbol, data)
}
