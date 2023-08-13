package main

import (
	"fmt"
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

	eventCH := make(chan string)
	doneCH := bbWS.StartRAW(eventCH)

	go handleRaw(eventCH)

	bbWS.Subscribe(ws.WSKLine, "5", "ETHUSDT")

	<-time.After(time.Second * 2)
	bbWS.Unsubscribe(ws.WSKLine, "1", "BTCUSDT")

	<-sigCH

	fmt.Println("Shutdown ...")

	doneCH <- struct{}{}
	bbWS.Close()

	<-time.After(time.Second * 1)
}

func handleRaw(eventCH chan string) {
	for event := range eventCH {
		fmt.Println("Event: ", event)
	}
}
