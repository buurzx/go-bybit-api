package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/buurzx/go-bybit-api/recws"
	"github.com/chuckpreslar/emission"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	// HeartBeatDuration period for sending ping packet according documentation
	HeartBeatDuration = time.Second * 20
	KeepAliveDuration = time.Second * 60
)

const (
	HostPerpetualReal     = "wss://stream.bybit.com/v5/public/linear"
	HostPerpetualtTestnet = "wss://stream-testnet.bybit.com/v5/public/linear"
)

const (
	WSKLine = "kline"

	WSDisconnected = "disconnected"
)

// TODO make default conf
type Configuration struct {
	Addr          string `json:"addr"`
	Proxy         string `json:"proxy"` // http://127.0.0.1:1081
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
}

type ByBitWS struct {
	cfg    *Configuration
	ctx    context.Context
	cancel context.CancelFunc
	conn   *recws.RecConn
	mu     sync.RWMutex
	Ended  bool

	subscribeCmds []Cmd
	emitter       *emission.Emitter
}

func New(config *Configuration) *ByBitWS {
	b := &ByBitWS{
		cfg:     config,
		emitter: emission.NewEmitter(),
	}

	b.ctx, b.cancel = context.WithCancel(context.Background())

	b.conn = &recws.RecConn{
		KeepAliveTimeout: KeepAliveDuration,
		NonVerbose:       true,
	}

	if config.Proxy != "" {
		proxy, err := url.Parse(config.Proxy)
		if err != nil {
			return nil
		}

		b.conn.Proxy = http.ProxyURL(proxy)
	}

	b.conn.SubscribeHandler = b.subscribeHandler

	return b
}

func (b *ByBitWS) subscribeHandler() error {
	if b.cfg.DebugMode {
		log.Printf("BybitWs subscribeHandler")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, cmd := range b.subscribeCmds {
		err := b.SendCmd(cmd)
		if err != nil {
			log.Printf("BybitWs SendCmd return error: %v", err)
		}
	}

	return nil
}

// IsConnected returns the WebSocket connection state
func (b *ByBitWS) IsConnected() bool {
	return b.conn.IsConnected()
}

// Topic => `kline`
//
// Interval
// 1 3 5 15 30 60 120 240 360 720 minute
// D day
// W week
// M month
//
// CoinPair => BTCUSDT
//
// Subscribe subscribes on ws topic to fetch data
func (b *ByBitWS) Subscribe(topic, interval, coinPair string) {
	arg := strings.Join([]string{topic, interval, coinPair}, ".")

	cmd := Cmd{
		Op:   "subscribe",
		Args: []interface{}{arg},
	}

	b.subscribeCmds = append(b.subscribeCmds, cmd)
	b.SendCmd(cmd)
}

func (b *ByBitWS) SendCmd(cmd Cmd) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return b.Send(string(data))
}

func (b *ByBitWS) Send(msg string) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("BybitWs send error: %v", r)
		}
	}()

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))

	return err
}

func (b *ByBitWS) Start() error {
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(HeartBeatDuration)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			messageType, data, err := b.conn.ReadMessage()
			if err != nil {
				log.Printf("BybitWs Read error, closing connection: %v", err)
				b.conn.Close()
				b.Ended = true
				return
			}

			b.processMessage(messageType, data)
		}
	}()

	return nil
}

func (b *ByBitWS) connect() {
	b.conn.Dial(b.cfg.Addr, nil)
}

func (b *ByBitWS) ping() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("BybitWs ping error: %v", r)
		}
	}()

	if !b.IsConnected() {
		return
	}
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(`{"op":"ping"}`))
	if err != nil {
		log.Printf("BybitWs ping error: %v", err)
	}
}

func (b *ByBitWS) processMessage(messageType int, data []byte) {
	ret := gjson.ParseBytes(data)

	if b.cfg.DebugMode {
		log.Printf("BybitWs %v", string(data))
	}

	if ret.Get("ret_msg").String() == "pong" {
		b.handlePong()
	}

	if topicValue := ret.Get("topic"); topicValue.Exists() {
		topic := topicValue.String()

		if !strings.HasPrefix(topic, WSKLine) {
			return
		}

		// kline.1.BTCUSD
		topicArray := strings.Split(topic, ".")
		if len(topicArray) != 3 {
			return
		}

		symbol := topicArray[2]
		raw := ret.Get("data").Raw

		var data []KLine

		err := json.Unmarshal([]byte(raw), &data)
		if err != nil {
			log.Printf("BybitWs %v", err)
			return
		}

		for _, kline := range data {
			b.processKLine(symbol, kline)
		}
	}
}

func (b *ByBitWS) handlePong() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handlePong error: %v", r)
		}
	}()

	pongHandler := b.conn.PongHandler()

	if pongHandler != nil {
		pongHandler("pong")
	}
	return nil
}

func (b *ByBitWS) CloseAndReconnect() {
	b.conn.CloseAndReconnect()
}

func (b *ByBitWS) Close() {
	if b.conn.IsConnected() {
		b.conn.Close()
	}
}
