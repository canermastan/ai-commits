package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/canermastan/ai-commits/config"
)

const (
	geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"
	httpTimeout    = 30 * time.Second
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response structure from the Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// BuildPrompt creates a prompt for the AI model with the given explanation and diff.
// The prompt instructs the AI to generate a conventional commit message in English
// based on a Turkish explanation and code diff.
func BuildPrompt(explanation, diff string) string {
	return fmt.Sprintf(`You are an expert commit message generator that strictly follows the Conventional Commit format (e.g., feat:, fix:, chore:, etc.).

The user will provide a description of what they did in Turkish, and a partial code diff will also be included.

Your job is to write a single, short, clear, and descriptive commit message in English, starting with the proper Conventional Commit prefix (like fix:, feat:, chore:).

The entire commit message must be in lowercase letters except for acronyms or proper nouns. This includes the commit type prefix and the description. For example:

- chore: add spacing for readability in test.py
- fix: correct typo in login flow
- feat: improve UI responsiveness

Do not include any explanation, do not add any extra text — return only the commit message.

User explanation (in Turkish):
"%s"

Code diff (partial):
%s
`, explanation, diff)
}

// makeGeminiRequest sends a request to the Gemini API and returns the generated text.
// It handles API key validation, request creation, and response parsing.
func makeGeminiRequest(prompt string) (string, error) {
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", geminiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", apiKey)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func CallAI(prompt string) (string, error) {
	return makeGeminiRequest(prompt)
}
