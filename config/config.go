package config

import (
	"os"
)

const (
	// MaxDiffSize is the maximum number of characters to include in the diff
	// sent to the AI. This is to prevent the API from rejecting large requests.
	MaxDiffSize = 10000

	// EnvKeyGeminiAPI is the environment variable name for the Gemini API key
	EnvKeyGeminiAPI = "GEMINI_API_KEY"
)

// GetAPIKey returns the Gemini API key from environment variables.
// Returns an empty string if the key is not set.
func GetAPIKey() string {
	return os.Getenv(EnvKeyGeminiAPI)
}

// GetMaxDiffSize returns the maximum number of characters to include in the diff.
func GetMaxDiffSize() int {
	return MaxDiffSize
}
