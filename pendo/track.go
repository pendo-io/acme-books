// Package pendo sends server-side Track Events to Pendo's Track API.
//
// Events are only sent when the PENDO_TRACK_EVENT_SECRET environment variable
// holds the Pendo Track Event shared secret. Keep the secret in the environment
// or a secret manager, never in the repository. PENDO_DATA_HOST overrides the
// Pendo data host, which defaults to data.pendo.io.
package pendo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

// SystemID is the visitor and account id for events that cannot be
// attributed to a specific user or account.
const SystemID = "system"

const (
	secretEnvVar    = "PENDO_TRACK_EVENT_SECRET"
	dataHostEnvVar  = "PENDO_DATA_HOST"
	defaultDataHost = "data.pendo.io"

	// maxValueBytes caps string property values so that the properties of an
	// event stay within Pendo's 512 byte limit.
	maxValueBytes = 100
)

// Properties holds the metadata sent with a Track Event.
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
	httpClient        = &http.Client{Timeout: 10 * time.Second}
	warnMissingSecret sync.Once
)

// Track sends a Track Event to Pendo in the background, so that tracking never
// slows down or fails the request that triggered it. Delivery errors are only
// logged. String property values longer than 100 bytes are truncated.
func Track(event, visitorID, accountID string, properties Properties) {
	secret := os.Getenv(secretEnvVar)

	if secret == "" {
		warnMissingSecret.Do(func() {
			fmt.Println("pendo: " + secretEnvVar + " is not set, Track Events will not be sent")
		})
		return
	}

	body, err := json.Marshal(trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: truncateValues(properties),
	})

	if err != nil {
		fmt.Println("pendo:", err)
		return
	}

	go send(event, trackURL(), secret, body)
}

func send(event, url, secret string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	if err != nil {
		fmt.Println("pendo:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", secret)

	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Println("pendo:", err)
		return
	}

	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusMultipleChoices {
		fmt.Printf("pendo: Track Event %q was rejected with status %d\n", event, resp.StatusCode)
	}
}

func trackURL() string {
	host := os.Getenv(dataHostEnvVar)

	if host == "" {
		host = defaultDataHost
	}

	return "https://" + host + "/data/track"
}

// truncateValues returns a copy of properties whose string values are cut to
// at most maxValueBytes, without splitting a UTF-8 encoded character.
func truncateValues(properties Properties) Properties {
	truncated := make(Properties, len(properties))

	for name, value := range properties {
		if s, ok := value.(string); ok && len(s) > maxValueBytes {
			end := maxValueBytes
			for end > 0 && !utf8.RuneStart(s[end]) {
				end--
			}
			value = s[:end]
		}

		truncated[name] = value
	}

	return truncated
}
