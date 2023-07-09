package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	defaultMainnet           = "https://api.bybit.com/"
	defaultDebugMode         = false
	defaultHttpCLientTimeout = 10 * time.Second
)

// Bybit
type ByBit struct {
	baseURL          string
	apiKey           string
	secretKey        string
	serverTimeOffset int64
	client           *http.Client
	debugMode        bool
}

// New instantiate bybit http client
func New(httpClient *http.Client, apiKey string, secretKey string, opts ...Option) *ByBit {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHttpCLientTimeout}
	}

	bb := &ByBit{
		baseURL:   defaultMainnet,
		apiKey:    apiKey,
		secretKey: secretKey,
		client:    httpClient,
		debugMode: defaultDebugMode,
	}

	// Custom options
	for _, opt := range opts {
		opt(bb)
	}

	return bb
}

// PublicRequest
func (b *ByBit) PublicRequest(method string, apiURL string, params map[string]interface{}, result interface{}) (fullURL string, resp []byte, err error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params[k]))
	}

	param := strings.Join(p, "&")
	fullURL = b.baseURL + apiURL
	if param != "" {
		fullURL += "?" + param
	}

	if b.debugMode {
		log.Printf("PublicRequest: %v", fullURL)
	}

	var binBody = bytes.NewReader(make([]byte, 0))

	// get a http request
	var request *http.Request
	request, err = http.NewRequest(method, fullURL, binBody)
	if err != nil {
		return
	}

	var response *http.Response
	response, err = b.client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if b.debugMode {
		log.Printf("PublicRequest: %v", string(resp))
	}

	err = json.Unmarshal(resp, result)

	return
}
