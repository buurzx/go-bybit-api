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
		DebugMode:     false,
	}

	bbWS := ws.New(cfg)
	bbWS.Subscribe(ws.WSKLine, "1", "BTCUSDT")
	bbWS.Subscribe(ws.WSKPublicTrade, "", "BTCUSDT")

	bbWS.On(ws.WSKLine, handleKLine)

	// bbWS.Start()
	bbWS.StartRAW(handleRaw)

	bbWS.Subscribe(ws.WSKLine, "5", "ETHUSDT")

	<-time.After(time.Second * 2)
	bbWS.Unsubscribe(ws.WSKLine, "1", "BTCUSDT")

	<-sigCH
	fmt.Println("Shutdown ...")

	bbWS.Close()
	<-time.After(time.Second * 1)
}

func handleKLine(symbol string, data ws.KLine) {
	log.Printf("handleKLine %v/%#v \n", symbol, data)
}

func handleRaw(messageType int, data []byte) {
	log.Printf("raw stream %v \n", string(data))
}
