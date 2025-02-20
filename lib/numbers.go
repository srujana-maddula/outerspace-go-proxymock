package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

// GetMathFact fetches a random math fact
func (c *NumbersClient) GetMathFact() (*MathFact, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/random/math?json", c.baseURL))
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
