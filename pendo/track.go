package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	trackEndpoint      = "https://data.pendo-dev.pendo-dev.com/data/track"
	integrationKey     = "99c7b9e5-e6fd-4601-ab64-cbb81750f324"
)

type TrackEvent struct {
	Type       string                 `json:"type"`
	Event      string                 `json:"event"`
	VisitorID  string                 `json:"visitorId"`
	AccountID  string                 `json:"accountId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Track sends a server-side track event to Pendo.
// visitorId and accountId should be set to actual user/account identifiers when available.
func Track(event string, visitorId string, accountId string, properties map[string]interface{}) {
	trackEvent := TrackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorId,
		AccountID:  accountId,
		Timestamp:  time.Now().UnixMilli(),
		Properties: properties,
	}

	body, err := json.Marshal(trackEvent)
	if err != nil {
		fmt.Printf("pendo.Track: failed to marshal event %q: %v\n", event, err)
		return
	}

	req, err := http.NewRequest("POST", trackEndpoint, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("pendo.Track: failed to create request for event %q: %v\n", event, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", integrationKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("pendo.Track: failed to send event %q: %v\n", event, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("pendo.Track: event %q returned status %d\n", event, resp.StatusCode)
	}
}
