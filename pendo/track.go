package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TODO: Set the PENDO_INTEGRATION_KEY environment variable to your Pendo
// server-side Track Event secret (integration key) to enable tracking.
// TODO: If your Pendo subscription uses a custom data host, set the
// PENDO_DATA_HOST environment variable (e.g. "data.eu.pendo.io").
// The default is "data.pendo.io".

// TrackEvent represents a Pendo server-side track event payload.
type TrackEvent struct {
	Type       string                 `json:"type"`
	Event      string                 `json:"event"`
	VisitorID  string                 `json:"visitorId"`
	AccountID  string                 `json:"accountId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Context    *TrackContext          `json:"context,omitempty"`
}

// TrackContext provides optional request context for the event.
type TrackContext struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	URL       string `json:"url,omitempty"`
}

// Track sends a track event to the Pendo data API.
// The call is made asynchronously so it does not block the request handler.
func Track(event string, visitorID string, accountID string, properties map[string]interface{}, ctx *TrackContext) {
	integrationKey := os.Getenv("PENDO_INTEGRATION_KEY")
	if integrationKey == "" {
		// TODO: PENDO_INTEGRATION_KEY is required for server-side tracking.
		fmt.Println("[Pendo] PENDO_INTEGRATION_KEY not set; skipping track event:", event)
		return
	}

	dataHost := os.Getenv("PENDO_DATA_HOST")
	if dataHost == "" {
		dataHost = "data.pendo.io"
	}

	payload := TrackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: properties,
		Context:    ctx,
	}

	go func() {
		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("[Pendo] Failed to marshal track event:", err)
			return
		}

		url := fmt.Sprintf("https://%s/data/track", dataHost)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("[Pendo] Failed to create request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-pendo-integration-key", integrationKey)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[Pendo] Failed to send track event:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			fmt.Printf("[Pendo] Track event returned status %d for event: %s\n", resp.StatusCode, event)
		}
	}()
}

// ContextFromRequest extracts a TrackContext from an HTTP request.
func ContextFromRequest(r *http.Request) *TrackContext {
	if r == nil {
		return nil
	}
	return &TrackContext{
		IP:        r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
		URL:       r.URL.String(),
	}
}
