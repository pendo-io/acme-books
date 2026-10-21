// Package pendo sends server-side Track Events to Pendo's Track Event API.
//
// Configuration is read from the environment, which main loads from .env
// before the server starts:
//   - PENDO_TRACK_EVENT_SECRET: the Track Event shared secret, sent in the
//     x-pendo-integration-key header. Events are not sent while it is unset.
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

	// The API has no user or account context of its own, so callers can name
	// the Pendo visitor and account in these request headers. Events are
	// attributed to "system" when a header is missing.
	visitorIDHeader = "X-Visitor-Id"
	accountIDHeader = "X-Account-Id"
	systemID        = "system"

	// maxStringLen caps string property values, in bytes, to keep the
	// properties payload within Pendo's 512-byte limit.
	maxStringLen = 128

	requestTimeout = 5 * time.Second
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
	loadConfigOnce sync.Once
	integrationKey string
	trackURL       string

	httpClient = &http.Client{Timeout: requestTimeout}
)

// loadConfig reads the configuration on first use instead of at package
// initialization, which runs before main has loaded .env.
func loadConfig() {
	integrationKey = os.Getenv("PENDO_TRACK_EVENT_SECRET")
	if integrationKey == "" {
		log.Println("pendo: PENDO_TRACK_EVENT_SECRET is not set; Track Events will not be sent")
	}

	dataHost := os.Getenv("PENDO_DATA_HOST")
	if dataHost == "" {
		dataHost = defaultDataHost
	}
	trackURL = "https://" + dataHost + "/data/track"
}

// Track records event for the visitor and account named by the X-Visitor-Id
// and X-Account-Id headers of r. The event is sent in the background: failures
// are logged and never affect the response to r.
func Track(r *http.Request, event string, properties Properties) {
	loadConfigOnce.Do(loadConfig)
	if integrationKey == "" {
		return
	}

	body, err := json.Marshal(trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  headerOrSystem(r, visitorIDHeader),
		AccountID:  headerOrSystem(r, accountIDHeader),
		Timestamp:  time.Now().UnixMilli(),
		Properties: truncateStrings(properties),
	})
	if err != nil {
		log.Printf("pendo: encoding %s event: %v", event, err)
		return
	}

	go send(event, body)
}

// send posts an encoded event to the Track Event API.
func send(event string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, trackURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("pendo: creating %s request: %v", event, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", integrationKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("pendo: sending %s event: %v", event, err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("pendo: %s event rejected: %s", event, resp.Status)
	}
}

func headerOrSystem(r *http.Request, header string) string {
	if value := r.Header.Get(header); value != "" {
		return value
	}
	return systemID
}

// truncateStrings returns a copy of properties in which string values longer
// than maxStringLen bytes are cut short, without splitting a UTF-8 character.
func truncateStrings(properties Properties) Properties {
	truncated := make(Properties, len(properties))
	for key, value := range properties {
		if s, ok := value.(string); ok && len(s) > maxStringLen {
			end := maxStringLen
			for end > 0 && !utf8.RuneStart(s[end]) {
				end--
			}
			value = s[:end]
		}
		truncated[key] = value
	}
	return truncated
}
