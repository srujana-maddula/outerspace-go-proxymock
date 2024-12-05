package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SpaceXClient handles API calls to SpaceX
type SpaceXClient struct {
	baseURL    string
	httpClient *http.Client
}

// Response structures for SpaceX API
type Launch struct {
	FlightNumber int    `json:"flight_number"`
	MissionName  string `json:"name"`
	DateUTC      string `json:"date_utc"`
	Success      bool   `json:"success"`
	Details      string `json:"details"`
}

type Rocket struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Height      struct {
		Meters float64 `json:"meters"`
	} `json:"height"`
	Mass struct {
		Kg int `json:"kg"`
	} `json:"mass"`
}

// RocketSummary provides a simplified view of rocket data
type RocketSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewSpaceXClient creates a new SpaceX API client
func NewSpaceXClient() *SpaceXClient {
	return &SpaceXClient{
		baseURL: "https://api.spacexdata.com/v4",
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetAllRockets fetches all rocket summaries
func (c *SpaceXClient) GetAllRockets() ([]RocketSummary, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/rockets", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rockets []Rocket
	if err := json.NewDecoder(resp.Body).Decode(&rockets); err != nil {
		return nil, err
	}

	// Convert to summaries
	summaries := make([]RocketSummary, len(rockets))
	for i, rocket := range rockets {
		summaries[i] = RocketSummary{
			ID:   rocket.ID,
			Name: rocket.Name,
		}
	}
	return summaries, nil
}

// GetLatestLaunch fetches details of the latest SpaceX launch
func (c *SpaceXClient) GetLatestLaunch() (*Launch, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/launches/latest", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var launch Launch
	if err := json.NewDecoder(resp.Body).Decode(&launch); err != nil {
		return nil, err
	}
	return &launch, nil
}

// GetRocket fetches details of a specific rocket by its ID
func (c *SpaceXClient) GetRocket(rocketID string) (*Rocket, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/rockets/%s", c.baseURL, rocketID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rocket Rocket
	if err := json.NewDecoder(resp.Body).Decode(&rocket); err != nil {
		return nil, err
	}
	return &rocket, nil
}
