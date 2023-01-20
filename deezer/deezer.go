package deezer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// Base URL for all the API methods
// https://api.deezer.com/version/service/id/method/?parameters
const (
	defaultBaseURL = "https://api.deezer.com/"
)

type LoggingRT struct {
	next http.RoundTripper
	out  io.Writer
}

func NewLoggingRT(next http.RoundTripper, out io.Writer) *LoggingRT {
	return &LoggingRT{
		next,
		out,
	}
}

func (rt *LoggingRT) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	defer func(begin time.Time) {
		var statusCode int
		if resp != nil {
			statusCode = resp.StatusCode
		}
		fmt.Fprintf(rt.out, "method=%s host=%s error=%v status_code=%d took=%s\n",
			req.Method, req.URL, err, statusCode, time.Since(begin))
	}(time.Now())

	return rt.next.RoundTrip(req)
}

type ListOptions struct {
	Index int `url:"index,omitempty"`
	Limit int `url:"limit,omitempty"`
}

func addOptions(s string, opt *ListOptions) (string, error) {

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// Client manages communication with the Deezer API.
type Client struct {
	// net/http Client used to communicate with the API.
	client *http.Client
	// BaseURL for the API endpoints.
	BaseURL *url.URL

	// A service per collection of endpoint resources.
	Albums  *AlbumService
	Artists *ArtistService
}

// NewClient returns a new Deezer API client.
func NewClient(client *http.Client) *Client {

	url, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  client,
		BaseURL: url,
	}

	if client == nil {
		c.client = &http.Client{}
	}

	// Register services.
	c.Albums = &AlbumService{client: c}
	c.Artists = &ArtistService{client: c}

	return c
}

type ClientOp func(*Client) error

func SetBaseURL(bu string) ClientOp {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}
		c.BaseURL = u
		return nil
	}
}

func New(client *http.Client, opts ...ClientOp) (*Client, error) {
	c := NewClient(client)

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// NewRequest creates a request. It takes care of using the BaseURL and setup common headers.
// Setting up the BaseURL allows to seamlessly use the loopback interface when using the test server.
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(url)
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
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.DoRequestWithClient(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decodingErr := json.NewDecoder(resp.Body).Decode(v)
	if errors.Is(decodingErr, io.EOF) {
		decodingErr = nil
	}
	if decodingErr != nil {
		err = decodingErr
	}

	// Notice we return the custom response but use the http.Reponse (resp) to decode into the v interface.
	return resp, err
}

// DoRequestWithClient submits an HTTP request using the specified client.
func (c *Client) DoRequestWithClient(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = checkResponse(resp)
	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}

	response := newResponse(resp)

	return response, nil
}

// Response is a wrapper for the standard http.Response.
// Typical use cases are: add pagination data, rate limits, etc.
type Response struct {
	*http.Response

	// Pagination
	Total    int
	PrevPage string
	NextPage string
}

// newResponse creates a new Response from a provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// TODO: Add pagination, response rate limits, etc.
	return response
}

type ErrorResponse struct {
	// HTTP response that cause this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// Error returned by the API on successful HTTP request
	APIError *APIError `json:"error"`

	// Carries any other error up the chain
	Carrier error
}

// checkResponse inspect the repsonse status code for HTTP errors and returns them as errors if present,
// given an http.Client's Do() does not returns an error on non-2xx status codes.
// view: https://golang.org/pkg/net/http/#Client.Do
func checkResponse(r *http.Response) error {

	// As Deezer's API returns its errors in the response body
	// within a 2xx status code, we need to inspect the body
	// to check for errors, view: https://developers.deezer.com/api/errors.
	// Reading the http.Response Body to inspect it will consume the Body, as it's an io.ReadCloser.
	// Thus we need to set the content back after we finish reading it, but set it as an io.ReadCloser.
	// We make the content an io.Reader with a bytes.Buffer.
	// We make it an io.ReadCloser with io.NopCloser which wraps an io.Reader and returns an io.ReadCloser.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	bufBody := bytes.NewBuffer(body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	errorResponse := &ErrorResponse{Response: r}

	if c := r.StatusCode; c >= 200 && c <= 299 {
		err = json.NewDecoder(bufBody).Decode(errorResponse)
		if err != nil {
			errorResponse.Message = string(body)
			errorResponse.Carrier = err
			return errorResponse
		}
		if errorResponse.APIError != nil {
			errorResponse.Carrier = err
			return errorResponse
		}
		return nil
	}

	errorResponse.Message = string(body)

	return errorResponse

}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Message, e.Carrier)
}

type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
