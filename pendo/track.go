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
}

// Track sends a server-side track event to the Pendo Track API.
// The integration key is read from the PENDO_INTEGRATION_KEY environment variable.
// Tracking is performed asynchronously; failures are logged but never
// affect application flow.
func Track(event, visitorID, accountID string, properties map[string]interface{}) {
	integrationKey := os.Getenv("PENDO_INTEGRATION_KEY")
	if integrationKey == "" {
		fmt.Println("pendo: PENDO_INTEGRATION_KEY not set, skipping track event")
		return
	}

	te := trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: properties,
	}

	go func() {
		body, err := json.Marshal(te)
		if err != nil {
			fmt.Printf("pendo: failed to marshal track event %q: %v\n", event, err)
			return
		}

		req, err := http.NewRequest("POST", trackEndpoint, bytes.NewReader(body))
		if err != nil {
			fmt.Printf("pendo: failed to create request for event %q: %v\n", event, err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-pendo-integration-key", integrationKey)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("pendo: failed to send track event %q: %v\n", event, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("pendo: track event %q returned status %d\n", event, resp.StatusCode)
		}
	}()
}
