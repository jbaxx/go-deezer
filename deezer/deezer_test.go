package deezer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setup sets a test server to make test requests
// a multiplexer to register test handlers into
// and a client that knows to to talk with the test server
func setup() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	c := NewClient()
	c.URL, _ = url.Parse(server.URL)

	return c, mux, server.Close
}

func TestNewClient(t *testing.T) {
	c := NewClient()

	want := defaultBaseURL
	got := c.URL.String()

	if got != want {
		t.Errorf("NewClient URL is: %v, want: %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient()

	inURL, outURL := "/algo", defaultBaseURL+"algo"
	req, err := c.NewRequest(http.MethodGet, inURL, nil)
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is: %v, want %v", inURL, got, want)
	}

	if got, want := req.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("NewRequest() Accept header is: %v, want: %v", got, want)
	}

}

func TestNewRequest_BadURL(t *testing.T) {
	c := NewClient()
	_, err := c.NewRequest(http.MethodGet, ":", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
