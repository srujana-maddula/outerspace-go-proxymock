package lib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock SpaceX client
type MockSpaceXClient struct {
	mock.Mock
}

func (m *MockSpaceXClient) GetAllRockets() ([]RocketSummary, error) {
	args := m.Called()
	return args.Get(0).([]RocketSummary), args.Error(1)
}

func (m *MockSpaceXClient) GetRocket(id string) (*Rocket, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Rocket), args.Error(1)
}

func (m *MockSpaceXClient) GetLatestLaunch() (*Launch, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Launch), args.Error(1)
}

// Mock Numbers client
type MockNumbersClient struct {
	mock.Mock
}

func (m *MockNumbersClient) GetMathFact() (*MathFact, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MathFact), args.Error(1)
}

func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	HandleRoot()(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var endpoints map[string]string
	json.Unmarshal(body, &endpoints)

	assert.Contains(t, endpoints, "/")
	assert.Contains(t, endpoints, "/api/rockets")
	assert.Contains(t, endpoints, "/api/rocket")
	assert.Contains(t, endpoints, "/api/latest-launch")
	assert.Contains(t, endpoints, "/api/numbers")
}

func TestHandleListRockets(t *testing.T) {
	mockClient := new(MockSpaceXClient)
	mockRockets := []RocketSummary{
		{ID: "123", Name: "Falcon 9"},
		{ID: "456", Name: "Falcon Heavy"},
	}

	mockClient.On("GetAllRockets").Return(mockRockets, nil)

	req := httptest.NewRequest("GET", "/api/rockets", nil)
	w := httptest.NewRecorder()

	handler := HandleListRockets(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var rockets []RocketSummary
	json.Unmarshal(body, &rockets)

	assert.Len(t, rockets, 2)
	assert.Equal(t, "Falcon 9", rockets[0].Name)
	assert.Equal(t, "456", rockets[1].ID)

	mockClient.AssertExpectations(t)
}

func TestHandleRocket(t *testing.T) {
	mockClient := new(MockSpaceXClient)
	mockRocket := &Rocket{
		ID:          "123",
		Name:        "Falcon 9",
		Description: "Orbital rocket",
		Height: struct {
			Meters float64 `json:"meters"`
		}{Meters: 70},
		Mass: struct {
			Kg int `json:"kg"`
		}{Kg: 549054},
	}

	mockClient.On("GetRocket", "123").Return(mockRocket, nil)

	req := httptest.NewRequest("GET", "/api/rocket?id=123", nil)
	w := httptest.NewRecorder()

	handler := HandleRocket(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var rocket Rocket
	json.Unmarshal(body, &rocket)

	assert.Equal(t, "Falcon 9", rocket.Name)
	assert.Equal(t, "123", rocket.ID)

	mockClient.AssertExpectations(t)
}

func TestHandleRocket_MissingID(t *testing.T) {
	mockClient := new(MockSpaceXClient)

	req := httptest.NewRequest("GET", "/api/rocket", nil)
	w := httptest.NewRecorder()

	handler := HandleRocket(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockClient.AssertNotCalled(t, "GetRocket")
}

func TestHandleRocket_Error(t *testing.T) {
	mockClient := new(MockSpaceXClient)
	mockClient.On("GetRocket", "999").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/api/rocket?id=999", nil)
	w := httptest.NewRecorder()

	handler := HandleRocket(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockClient.AssertExpectations(t)
}

func TestHandleLatestLaunch(t *testing.T) {
	mockClient := new(MockSpaceXClient)
	mockLaunch := &Launch{
		FlightNumber: 100,
		MissionName:  "Mission X",
		DateUTC:      "2023-01-01T12:00:00Z",
		Success:      true,
		Details:      "Test mission",
	}

	mockClient.On("GetLatestLaunch").Return(mockLaunch, nil)

	req := httptest.NewRequest("GET", "/api/latest-launch", nil)
	w := httptest.NewRecorder()

	handler := HandleLatestLaunch(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var launch Launch
	json.Unmarshal(body, &launch)

	assert.Equal(t, 100, launch.FlightNumber)

	mockClient.AssertExpectations(t)
}

func TestHandleNumbers(t *testing.T) {
	mockClient := new(MockNumbersClient)
	mockFact := &MathFact{
		Text:   "42 is the meaning of life",
		Number: 42,
		Found:  true,
		Type:   "math",
	}

	mockClient.On("GetMathFact").Return(mockFact, nil)

	req := httptest.NewRequest("GET", "/api/numbers", nil)
	w := httptest.NewRecorder()

	handler := HandleNumbers(mockClient)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var fact MathFact
	json.Unmarshal(body, &fact)

	assert.Equal(t, "42 is the meaning of life", fact.Text)
	assert.Equal(t, 42, fact.Number)

	mockClient.AssertExpectations(t)
}
