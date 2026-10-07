// Package pendo sends server-side Track Events to Pendo.
package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

// TODO: Track Events are only sent once PENDO_TRACK_EVENT_SECRET holds the
// Pendo Track Event shared secret (set it in the environment, not in the
// committed .env file). Subscriptions hosted outside Pendo's default US data
// center must also set PENDO_DATA_HOST to their data host, e.g.
// data.eu.pendo.io.
const (
	secretEnv       = "PENDO_TRACK_EVENT_SECRET"
	dataHostEnv     = "PENDO_DATA_HOST"
	defaultDataHost = "data.pendo.io"
)

// systemID is sent as the visitorId and accountId of every event.
//
// TODO: the API has no authentication yet, so requests carry no user or
// account identity. Once they do, send the authenticated user's ID as the
// visitorId and their account's ID as the accountId instead.
const systemID = "system"

// maxPropertyLength caps the length of string property values, keeping events
// within Pendo's 512 byte limit on properties.
const maxPropertyLength = 64

var (
	client      = &http.Client{Timeout: 5 * time.Second}
	warnMissing sync.Once
)

type trackEvent struct {
	Type       string                 `json:"type"`
	Event      string                 `json:"event"`
	VisitorID  string                 `json:"visitorId"`
	AccountID  string                 `json:"accountId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
	Context    trackContext           `json:"context"`
}

type trackContext struct {
	URL       string `json:"url"`
	UserAgent string `json:"userAgent,omitempty"`
}

// Track sends the named Track Event to Pendo for the request r. The event is
// sent in the background, so tracking never delays or fails the request, and
// failures are logged.
func Track(r *http.Request, event string, properties map[string]interface{}) {
	secret := os.Getenv(secretEnv)

	if secret == "" {
		warnMissing.Do(func() {
			fmt.Println("pendo:", secretEnv, "is not set, Track Events will not be sent")
		})
		return
	}

	host := os.Getenv(dataHostEnv)

	if host == "" {
		host = defaultDataHost
	}

	body, err := json.Marshal(trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  systemID,
		AccountID:  systemID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: truncate(properties),
		Context: trackContext{
			URL:       r.URL.RequestURI(),
			UserAgent: r.UserAgent(),
		},
	})

	if err != nil {
		fmt.Printf("pendo: failed to encode %s: %v\n", event, err)
		return
	}

	go send("https://"+host+"/data/track", secret, event, body)
}

// send posts an encoded Track Event to Pendo, logging any failure.
func send(endpoint, secret, event string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))

	if err != nil {
		fmt.Printf("pendo: failed to send %s: %v\n", event, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", secret)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("pendo: failed to send %s: %v\n", event, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("pendo: %s was rejected: %s\n", event, resp.Status)
	}
}

// truncate returns a copy of properties with string values cut down to
// maxPropertyLength characters.
func truncate(properties map[string]interface{}) map[string]interface{} {
	truncated := make(map[string]interface{}, len(properties))

	for name, value := range properties {
		if s, ok := value.(string); ok && utf8.RuneCountInString(s) > maxPropertyLength {
			value = string([]rune(s)[:maxPropertyLength])
		}

		truncated[name] = value
	}

	return truncated
}
