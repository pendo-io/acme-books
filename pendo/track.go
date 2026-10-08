// Package pendo sends server-side Track Events to Pendo.
//
// Tracking is configured through environment variables:
//
//   - PENDO_TRACK_EVENT_SECRET: the Pendo Track Event shared secret, sent as the
//     x-pendo-integration-key header. Tracking is disabled when it is not set.
//     Do not add it to .env, which is committed.
//   - PENDO_DATA_HOST: the Pendo data host. Defaults to data.pendo.io.
package pendo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	defaultDataHost = "data.pendo.io"

	// The API has no authenticated users or accounts yet, so every event is
	// attributed to the "system" visitor and account.
	systemID = "system"

	// Pendo rejects events whose serialized properties exceed 512 bytes, so
	// string values are capped to keep every event well under that limit.
	maxStringLength = 100
)

// Properties is the metadata sent with a Track Event.
type Properties map[string]any

type trackEvent struct {
	Type       string     `json:"type"`
	Event      string     `json:"event"`
	VisitorID  string     `json:"visitorId"`
	AccountID  string     `json:"accountId"`
	Timestamp  int64      `json:"timestamp"`
	Properties Properties `json:"properties,omitempty"`
}

var (
	httpClient      = &http.Client{Timeout: 5 * time.Second}
	disabledWarning sync.Once
)

// Track sends a Track Event to Pendo. The request is made in the background so
// tracking never delays or fails the API response; failures are only logged.
func Track(event string, properties Properties) {
	secret := os.Getenv("PENDO_TRACK_EVENT_SECRET")
	if secret == "" {
		disabledWarning.Do(func() {
			log.Println("pendo: PENDO_TRACK_EVENT_SECRET is not set, track events are disabled")
		})
		return
	}

	body, err := json.Marshal(trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  systemID,
		AccountID:  systemID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: truncateStrings(properties),
	})
	if err != nil {
		log.Printf("pendo: could not encode %s event: %v", event, err)
		return
	}

	go send(event, secret, body)
}

func send(event string, secret string, body []byte) {
	host := os.Getenv("PENDO_DATA_HOST")
	if host == "" {
		host = defaultDataHost
	}

	req, err := http.NewRequest(http.MethodPost, "https://"+host+"/data/track", bytes.NewReader(body))
	if err != nil {
		log.Printf("pendo: could not send %s event: %v", event, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", secret)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("pendo: could not send %s event: %v", event, err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("pendo: %s event was rejected with status %d", event, resp.StatusCode)
	}
}

// truncateStrings returns a copy of properties with string values cut to
// maxStringLength bytes, without splitting a UTF-8 character.
func truncateStrings(properties Properties) Properties {
	truncated := make(Properties, len(properties))
	for name, value := range properties {
		if s, ok := value.(string); ok && len(s) > maxStringLength {
			cut := maxStringLength
			for cut > 0 && !utf8.RuneStart(s[cut]) {
				cut--
			}
			value = s[:cut]
		}
		truncated[name] = value
	}
	return truncated
}
