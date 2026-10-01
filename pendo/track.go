package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TODO: Set the PENDO_INTEGRATION_KEY environment variable with your Pendo
// Track Event secret (server-side integration key). Server-side tracking
// will be skipped until this is configured.

const trackEndpoint = "https://data.pendo.io/data/track"

// TrackEvent represents a Pendo server-side track event payload.
type TrackEvent struct {
	Type       string                 `json:"type"`
	Event      string                 `json:"event"`
	VisitorID  string                 `json:"visitorId"`
	AccountID  string                 `json:"accountId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Track sends a server-side track event to the Pendo Track API.
// It runs in a goroutine so it does not block the caller.
// Tracking failures are logged but never break application flow.
func Track(event string, visitorID string, accountID string, properties map[string]interface{}) {
	go func() {
		if err := sendTrackEvent(event, visitorID, accountID, properties); err != nil {
			fmt.Printf("[pendo] failed to send track event %q: %v\n", event, err)
		}
	}()
}

func sendTrackEvent(event string, visitorID string, accountID string, properties map[string]interface{}) error {
	integrationKey := os.Getenv("PENDO_INTEGRATION_KEY")
	if integrationKey == "" {
		// TODO: PENDO_INTEGRATION_KEY is required for server-side tracking.
		// Set this environment variable to your Pendo Track Event secret.
		return fmt.Errorf("PENDO_INTEGRATION_KEY environment variable is not set")
	}

	payload := TrackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: properties,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal track event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, trackEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", integrationKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("pendo API returned status %d", resp.StatusCode)
	}

	return nil
}
