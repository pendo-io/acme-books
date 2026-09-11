package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const trackEndpoint = "https://data.pendo.io/data/track"

type trackEvent struct {
	Type       string                 `json:"type"`
	Event      string                 `json:"event"`
	VisitorID  string                 `json:"visitorId"`
	AccountID  string                 `json:"accountId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Context    map[string]string      `json:"context,omitempty"`
}

// Track sends a track event to Pendo asynchronously.
// The integration key is read from the PENDO_INTEGRATION_KEY environment variable.
func Track(event, visitorID, accountID string, properties map[string]interface{}, ctx map[string]string) {
	key := os.Getenv("PENDO_INTEGRATION_KEY")
	if key == "" {
		return
	}

	te := trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: properties,
		Context:    ctx,
	}

	go func() {
		body, err := json.Marshal(te)
		if err != nil {
			fmt.Printf("[Pendo] Failed to marshal track event: %v\n", err)
			return
		}

		req, err := http.NewRequest("POST", trackEndpoint, bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("[Pendo] Failed to create request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-pendo-integration-key", key)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[Pendo] Failed to send track event: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			fmt.Printf("[Pendo] Track event returned status %d\n", resp.StatusCode)
		}
	}()
}
