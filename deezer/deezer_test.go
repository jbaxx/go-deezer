package deezer

import (
	"context"
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

	c := NewClient(nil)
	c.BaseURL, _ = url.Parse(server.URL)

	return c, mux, server.Close
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	want := defaultBaseURL
	got := c.BaseURL.String()

	if got != want {
		t.Errorf("NewClient URL is: %v, want: %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

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
	c := NewClient(nil)
	_, err := c.NewRequest(http.MethodGet, ":", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want: %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"key": "value"}`)
	})

	type base struct {
		Key string
	}

	req, _ := client.NewRequest(http.MethodGet, ".", nil)
	body := new(base)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	want := &base{"value"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body: %v, want %v", body, want)
	}

}

func TestDo_withHTTPError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest(http.MethodGet, ".", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Errorf("expected HTTP 400 error got nil")
	}

}

func TestDo_withDeezerAPIError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"error":{"type":"DataException","message":"no data","code":800}}`)
	})

	req, _ := client.NewRequest(http.MethodGet, ".", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Errorf("expected HTTP 400 error got nil")
	}

	want := &APIError{
		Type:    "DataException",
		Message: "no data",
		Code:    800,
	}

	if rerr, ok := err.(*ErrorResponse); ok {
		if got := rerr.APIError; !reflect.DeepEqual(got, want) {
			t.Errorf("ErrorResponse.APIError: %#v, want: %#v", got, want)
		}
	}

}

func TestDo_malformedResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"key": "value"`)
	})

	type base struct {
		Key string
	}

	req, _ := client.NewRequest(http.MethodGet, ".", nil)
	body := new(base)
	_, err := client.Do(context.Background(), req, body)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

}
