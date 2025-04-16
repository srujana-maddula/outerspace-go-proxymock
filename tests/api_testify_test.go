package tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"outerspace-go/lib"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (suite *APITestSuite) SetupTest() {
	// Setup verbose logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Create clients
	spaceClient := lib.NewSpaceXClient()
	numbersClient := lib.NewNumbersClient()

	// Create router with all handlers
	mux := http.NewServeMux()

	// Add debug logging wrapper for each handler
	debugHandler := func(pattern string, handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Debug().
				Str("pattern", pattern).
				Str("request_path", r.URL.Path).
				Str("request_url", r.URL.String()).
				Msg("Handler matched")
			handler(w, r)
		}
	}

	// Register handlers with debug logging
	mux.HandleFunc("/api/rockets", debugHandler("/api/rockets", lib.HandleListRockets(spaceClient)))
	mux.HandleFunc("/api/rocket", debugHandler("/api/rocket", lib.HandleRocket(spaceClient)))
	mux.HandleFunc("/api/latest-launch", debugHandler("/api/latest-launch", lib.HandleLatestLaunch(spaceClient)))
	mux.HandleFunc("/api/numbers", debugHandler("/api/numbers", lib.HandleNumbers(numbersClient)))
	mux.HandleFunc("/", debugHandler("/", lib.HandleRoot()))

	// Create test server
	suite.server = httptest.NewServer(mux)
	log.Info().Msgf("Test server started at %s", suite.server.URL)
}

func (suite *APITestSuite) TearDownTest() {
	// Stop test server
	if suite.server != nil {
		suite.server.Close()
		log.Info().Msg("Test server stopped")
	}
}

// RRPair represents the structure of the recorded request/response pairs
type RRPair struct {
	HTTP struct {
		Req struct {
			Method string `json:"method"`
			URI    string `json:"uri"`
		} `json:"req"`
		Res struct {
			StatusCode int    `json:"statusCode"`
			BodyBase64 string `json:"bodyBase64"`
		} `json:"res"`
	} `json:"http"`
}

func (suite *APITestSuite) TestRecordedAPIs() {
	// Check if the directory exists
	rrpairsDir := "../proxymock/recording/localhost"
	_, err := os.Stat(rrpairsDir)
	if os.IsNotExist(err) {
		suite.T().Fatalf("Directory %s does not exist", rrpairsDir)
	}
	assert.NoError(suite.T(), err, "Failed to check directory: %s", rrpairsDir)

	// Read all files in the directory
	files, err := filepath.Glob(filepath.Join(rrpairsDir, "*.json"))
	assert.NoError(suite.T(), err, "Failed to read rrpairs directory")

	// Check that we have files to test
	assert.Greater(suite.T(), len(files), 0, "No test files found in %s", rrpairsDir)
	log.Info().Msgf("Found %d files to test", len(files))

	for _, file := range files {
		log.Info().Msgf("Processing file: %s", filepath.Base(file))

		// Read and parse the RRPair file
		data, err := os.ReadFile(file)
		assert.NoError(suite.T(), err, "Failed to read file: %s", file)

		var rrpair RRPair
		err = json.Unmarshal(data, &rrpair)
		assert.NoError(suite.T(), err, "Failed to parse JSON from file: %s", file)

		// Create a subtest for this API request
		suite.Run("URL_"+rrpair.HTTP.Req.URI, func() {
			log.Info().
				Str("uri", rrpair.HTTP.Req.URI).
				Str("method", rrpair.HTTP.Req.Method).
				Msg("Testing endpoint")

			// Make the request to our test server
			req, err := http.NewRequest(rrpair.HTTP.Req.Method, suite.server.URL+rrpair.HTTP.Req.URI, nil)
			assert.NoError(suite.T(), err, "Failed to create request")

			// Add debug logging for the request
			log.Debug().
				Str("request_url", req.URL.String()).
				Str("request_path", req.URL.Path).
				Msg("Making request")

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(suite.T(), err, "Failed to make request")
			defer resp.Body.Close()

			// Check status code
			assert.Equal(suite.T(), rrpair.HTTP.Res.StatusCode, resp.StatusCode,
				"Status code mismatch")

			// Read and compare response bodies
			actualBody, err := io.ReadAll(resp.Body)
			assert.NoError(suite.T(), err, "Failed to read response body")

			// Log the actual response for debugging
			log.Debug().
				RawJSON("response", actualBody).
				Msg("Received response")

			// Decode the expected base64 body
			expectedBodyBytes, err := base64.StdEncoding.DecodeString(rrpair.HTTP.Res.BodyBase64)
			assert.NoError(suite.T(), err, "Failed to decode base64 body")

			// Log the expected response for debugging
			log.Debug().
				RawJSON("expected", expectedBodyBytes).
				Msg("Expected response")

			// For JSON responses, compare the parsed structures instead of raw strings
			if resp.Header.Get("Content-Type") == "application/json" {
				var expectedJSON, actualJSON interface{}

				err = json.Unmarshal(expectedBodyBytes, &expectedJSON)
				assert.NoError(suite.T(), err, "Failed to parse expected JSON")

				err = json.Unmarshal(actualBody, &actualJSON)
				assert.NoError(suite.T(), err, "Failed to parse actual JSON")

				// Compare the full structure
				assert.Equal(suite.T(), expectedJSON, actualJSON,
					"Response body mismatch")
			}
		})
	}
}

// Run the test suite
func TestAPISuite(t *testing.T) {
	// Enable verbose output for tests
	if testing.Verbose() {
		fmt.Println("Running tests in verbose mode")
	}
	suite.Run(t, new(APITestSuite))
}
