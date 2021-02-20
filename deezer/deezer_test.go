package deezer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

func TestGetAlbum(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/album/", func(w http.ResponseWriter, r *http.Request) {
		want := "GET"
		got := r.Method
		if got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}
		fmt.Fprintf(w, `{"id": 44132881}`)
	})

	album, _, err := client.Albums.Get(44132881)
	if err != nil {
		t.Errorf("Album.Get return error: %v", err)
	}

	want := &Album{ID: 44132881}
	if !reflect.DeepEqual(album, want) {
		t.Errorf("Album.Get got: %#v, want %#v", album, want)
	}

}
