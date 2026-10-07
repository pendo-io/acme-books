// Package pendo sends server-side Track Events to Pendo's Track API
// (POST <data host>/data/track).
//
// Configuration is read from the environment each time an event is sent:
//
//	PENDO_TRACK_EVENT_SECRET  The subscription's Track Event shared secret, sent
//	                          as the x-pendo-integration-key header. Required:
//	                          while it is unset, events are skipped.
//	PENDO_DATA_HOST           Optional. Pendo data host including the scheme,
//	                          https://data.pendo.io by default. Set it when the
//	                          subscription is hosted in another region.
//
// TODO: neither value is configured for this app yet, so no events reach Pendo
// until PENDO_TRACK_EVENT_SECRET is provided. Set it in the process environment
// or a secret manager, not in .env: that file is committed to the repository.
package pendo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultDataHost = "https://data.pendo.io"

// maxStringLength bounds string property values so that an event's properties
// stay within Pendo's 512-byte limit.
const maxStringLength = 100

// Properties is the metadata sent with a Track Event.
type Properties map[string]any

// trackEvent is the body of a Track API request.
type trackEvent struct {
	Type       string        `json:"type"`
	Event      string        `json:"event"`
	VisitorID  string        `json:"visitorId"`
	AccountID  string        `json:"accountId"`
	Timestamp  int64         `json:"timestamp"`
	Properties Properties    `json:"properties,omitempty"`
	Context    *trackContext `json:"context,omitempty"`
}

// trackContext describes the request that triggered the event.
type trackContext struct {
	UserAgent string `json:"userAgent,omitempty"`
	URL       string `json:"url,omitempty"`
}

var (
	client       = &http.Client{Timeout: 10 * time.Second}
	warnDisabled sync.Once
)

// Track sends event to Pendo for the caller that made request r (r may be nil
// for events that no request triggered). It returns immediately: the event is
// delivered in the background and failures are only logged, so tracking never
// delays or breaks the request being served.
func Track(r *http.Request, event string, properties Properties) {
	secret := os.Getenv("PENDO_TRACK_EVENT_SECRET")
	if secret == "" {
		warnDisabled.Do(func() {
			log.Println("pendo: PENDO_TRACK_EVENT_SECRET is not set, so track events are not sent")
		})
		return
	}

	// TODO: this API has no authentication, sessions or accounts yet, so
	// requests carry nothing that identifies the caller and every event is
	// attributed to the "system" visitor and account. Resolve the real visitor
	// and account from r here once requests are authenticated.
	visitorID, accountID := "system", "system"

	var body bytes.Buffer
	encoder := json.NewEncoder(&body)
	// Keep &, < and > literal (filters are query strings) rather than escaping
	// them, which would eat into the properties size limit.
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(trackEvent{
		Type:       "track",
		Event:      event,
		VisitorID:  visitorID,
		AccountID:  accountID,
		Timestamp:  time.Now().UnixMilli(),
		Properties: truncate(properties),
		Context:    requestContext(r),
	})
	if err != nil {
		log.Printf("pendo: encoding %q event: %v", event, err)
		return
	}

	go send(trackURL(), secret, event, body.Bytes())
}

func send(url, secret, event string, body []byte) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("pendo: sending %q event: %v", event, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pendo-integration-key", secret)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("pendo: sending %q event: %v", event, err)
		return
	}
	defer resp.Body.Close()
	// Drain the body so the connection can be reused.
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("pendo: sending %q event: unexpected response %s", event, resp.Status)
	}
}

func trackURL() string {
	host := os.Getenv("PENDO_DATA_HOST")
	if host == "" {
		host = defaultDataHost
	}
	return strings.TrimSuffix(host, "/") + "/data/track"
}

func requestContext(r *http.Request) *trackContext {
	if r == nil {
		return nil
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return &trackContext{
		UserAgent: r.UserAgent(),
		URL:       scheme + "://" + r.Host + r.URL.RequestURI(),
	}
}

// truncate returns a copy of properties with string values cut to
// maxStringLength bytes.
func truncate(properties Properties) Properties {
	out := make(Properties, len(properties))
	for name, value := range properties {
		if s, ok := value.(string); ok && len(s) > maxStringLength {
			// ToValidUTF8 drops a rune split by the cut.
			value = strings.ToValidUTF8(s[:maxStringLength], "")
		}
		out[name] = value
	}
	return out
}
