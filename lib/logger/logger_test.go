package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Save original values
	origTimeFormat := zerolog.TimeFieldFormat
	defer func() {
		zerolog.TimeFieldFormat = origTimeFormat
	}()

	// Test initialization
	Init()

	// Verify time format was set
	assert.Equal(t, "2006-01-02T15:04:05Z07:00", zerolog.TimeFieldFormat)
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()

	// Basic check that we get a valid logger
	assert.NotNil(t, logger)
}
