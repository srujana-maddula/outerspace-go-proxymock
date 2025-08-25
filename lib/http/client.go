package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents an HTTP client for the outerspace API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new HTTP client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Close is a no-op for HTTP client (for compatibility with gRPC client interface)
func (c *Client) Close() error {
	return nil
}

// Launch represents the latest launch data
type Launch struct {
	FlightNumber int32  `json:"flight_number"`
	MissionName  string `json:"mission_name"`
	DateUtc      string `json:"date_utc"`
	Success      bool   `json:"success"`
	Details      string `json:"details"`
}

// Rocket represents rocket data
type Rocket struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	HeightMeters  float32 `json:"height_meters"`
	MassKg        int32   `json:"mass_kg"`
}

// GetRocketsResponse represents the response from getting all rockets
type GetRocketsResponse struct {
	Rockets []*Rocket `json:"rockets"`
}

// MathFact represents a math fact
type MathFact struct {
	Number int32  `json:"number"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Found  bool   `json:"found"`
}

// GetLatestLaunch calls the HTTP API to get the latest launch
func (c *Client) GetLatestLaunch(ctx context.Context) (*Launch, error) {
	url := fmt.Sprintf("%s/api/latest-launch", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var launch Launch
	if err := json.NewDecoder(resp.Body).Decode(&launch); err != nil {
		return nil, err
	}

	return &launch, nil
}

// GetRocket calls the HTTP API to get a specific rocket
func (c *Client) GetRocket(ctx context.Context, id string) (*Rocket, error) {
	url := fmt.Sprintf("%s/api/rocket?id=%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var rocket Rocket
	if err := json.NewDecoder(resp.Body).Decode(&rocket); err != nil {
		return nil, err
	}

	return &rocket, nil
}

// GetRockets calls the HTTP API to get all rockets
func (c *Client) GetRockets(ctx context.Context) (*GetRocketsResponse, error) {
	url := fmt.Sprintf("%s/api/rockets", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var rockets []*Rocket
	if err := json.NewDecoder(resp.Body).Decode(&rockets); err != nil {
		return nil, err
	}

	return &GetRocketsResponse{Rockets: rockets}, nil
}

// GetMathFact calls the HTTP API to get a math fact
func (c *Client) GetMathFact(ctx context.Context) (*MathFact, error) {
	url := fmt.Sprintf("%s/api/numbers", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var mathFact MathFact
	if err := json.NewDecoder(resp.Body).Decode(&mathFact); err != nil {
		return nil, err
	}

	return &mathFact, nil
}