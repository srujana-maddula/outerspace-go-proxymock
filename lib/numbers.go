package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// NumbersClient handles API calls to Numbers API
type NumbersClient struct {
	baseURL    string
	httpClient *http.Client
}

// RocketSummary provides a simplified view of rocket data
type MathFact struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
	Found  bool   `json:"found"`
	Type   string `json:"type"`
}

// NewNumbersClient creates a new Numbers API client
func NewNumbersClient() *NumbersClient {
	return &NumbersClient{
		baseURL: "http://numbersapi.com",
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Add makeRequest method
func (c *NumbersClient) makeRequest(method, url string) (*http.Response, error) {
	start := time.Now()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Error().
			Str("method", method).
			Str("host", req.URL.Host).
			Dur("latency", duration).
			Err(err).
			Msg("API request failed")
		return nil, err
	}

	// Log X-prefixed headers
	for header, values := range resp.Header {
		if len(header) > 0 && (header[0] == 'X' || header[0] == 'x') {
			log.Info().
				Str("header", header).
				Strs("values", values).
				Msg("X-Header found")
		}
	}

	log.Info().
		Str("method", method).
		Str("host", req.URL.Host).
		Int("status", resp.StatusCode).
		Dur("latency", duration).
		Msg("Outbound")

	return resp, nil
}

// Update GetMathFact to use makeRequest
func (c *NumbersClient) GetMathFact() (*MathFact, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("%s/random/math?json", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mathFact MathFact
	if err := json.NewDecoder(resp.Body).Decode(&mathFact); err != nil {
		return nil, err
	}
	return &mathFact, nil
}
