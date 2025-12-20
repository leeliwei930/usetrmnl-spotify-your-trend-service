package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leeliwei930/usetrmnl_spotify_service/config"
)

// AgentClient handles communication with the Python agent service
type AgentClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAgentClient creates a new agent client using configuration
func NewAgentClient() *AgentClient {
	cfg := config.Get()
	return &AgentClient{
		BaseURL: cfg.AgentBaseURL,
		Client: &http.Client{
			Timeout: time.Duration(cfg.AgentClientTimeout) * time.Second,
		},
	}
}

// BehindTheLyricsRequest represents the request payload for the agent
type BehindTheLyricsRequest struct {
	SearchInput struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Cover  string `json:"cover"`
	} `json:"searchInput"`
}

// BehindTheLyricsResponse represents the response from the agent
type BehindTheLyricsResponse struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Cover     string `json:"cover"`
	SummaryEn string `json:"summary_en"`
	SummaryZh string `json:"summary_zh"`
}

// GetBehindTheLyrics calls the Python agent to get lyrics summary
func (ac *AgentClient) GetBehindTheLyrics(title, artist, album, coverURL string) (*BehindTheLyricsResponse, error) {
	// Create context with timeout (from configuration)
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.AgentRequestTimeout)*time.Second)
	defer cancel()

	// Prepare request payload
	reqPayload := BehindTheLyricsRequest{}
	reqPayload.SearchInput.Title = title
	reqPayload.SearchInput.Artist = artist
	reqPayload.SearchInput.Album = album
	reqPayload.SearchInput.Cover = coverURL

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with context
	url := fmt.Sprintf("%s/usetrmnl/agent/behind-the-lyrics", ac.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := ac.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var agentResp BehindTheLyricsResponse
	if err := json.Unmarshal(body, &agentResp); err != nil {
		return nil, fmt.Errorf("failed to parse agent response: %w", err)
	}

	return &agentResp, nil
}
