package deezer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Base URL for all the API methods
// https://api.deezer.com/version/service/id/method/?parameters
const (
	defaultBaseURL = "https://api.deezer.com/"
)

// Client manages communication with the Deezer API.
type Client struct {
	client *http.Client
	URL    *url.URL

	Albums *AlbumService
}

// NewClient returns a new Deezer API client.
func NewClient() *Client {
	// TODO: add authentication options
	url, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client: &http.Client{},
		URL:    url,
	}

	c.Albums = &AlbumService{client: c}

	return c
}

// Response is a wrapper for the standard http.Response.
// Typical use cases are: add pagination data, rate limits, etc.
type Response struct {
	*http.Response
}

// NewRequest creates a request. It takes care of using the BaseURL and setup common headers.
// Setting up the BaseURL allows to seamlessly use the loopback interface when using the test server.
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	u, err := c.URL.Parse(url)
	if err != nil {
		return nil, err
	}

	// For now, the methods implemented do not require a Body in the request, setting it to nil.
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Do makes an API request and returns a custom Response. The API response is JSON and will be decoded in the value
// specified by v.
// If v implements the io.Writer interface, the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decodingErr := json.NewDecoder(resp.Body).Decode(v)
		if decodingErr == io.EOF {
			decodingErr = nil
		}
		if decodingErr != nil {
			err = decodingErr
		}
	}

	// Notice we return the custom response but use the http.Reponse to decode into the v interface.
	return response, err
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(client *http.Client, req *http.Request) (*http.Response, error) {
	// req = req.WithContext(ctx)
	return client.Do(req)
}

// newResponse creates a new Response from a provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// TODO: Add pagination, response rate limits, etc.
	return response
}
