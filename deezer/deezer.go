package deezer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// AlbumService defines a client to interface with the Deezer Album service
type AlbumService struct {
	client *Client
}

// Response is a wrapper for the standard http.Response.
// Typical use cases are: add pagination data, rate limits, etc.
type Response struct {
	*http.Response
}

// Get fetches an Album given an album id.
func (a *AlbumService) Get(id int) (*Album, *Response, error) {
	i := strconv.Itoa(id)

	url := fmt.Sprintf("album/%v", i)
	fmt.Println("URL: ", url)

	req, err := a.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	album := new(Album)
	resp, err := a.client.Do(req, album)
	if err != nil {
		return nil, resp, err
	}

	return album, resp, err

}

// NewRequest creates a request. It takes care of using the BaseURL and setup common headers.
// Setting up the BaseURL allows to seamlessly use the loopback interface when using the test server.
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	u, err := c.URL.Parse(url)
	if err != nil {
		return nil, err
	}

	// For now, the methods implemented do not require a Body in the request, setting it to nil.
	fmt.Println("FullURL:", u.String())
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

// Album represents an artist's album.
type Album struct {
	ID                    int            `json:"id"`
	Title                 string         `json:"title"`
	Upc                   string         `json:"upc"`
	Link                  string         `json:"link"`
	Share                 string         `json:"share"`
	Cover                 string         `json:"cover"`
	CoverSmall            string         `json:"cover_small"`
	CoverMedium           string         `json:"cover_medium"`
	CoverBig              string         `json:"cover_big"`
	CoverXl               string         `json:"cover_xl"`
	Md5Image              string         `json:"md5_image"`
	GenreID               int            `json:"genre_id"`
	Genres                Genres         `json:"genres"`
	Label                 string         `json:"label"`
	NbTracks              int            `json:"nb_tracks"`
	Duration              int            `json:"duration"`
	Fans                  int            `json:"fans"`
	Rating                int            `json:"rating"`
	ReleaseDate           string         `json:"release_date"`
	RecordType            string         `json:"record_type"`
	Available             bool           `json:"available"`
	Tracklist             string         `json:"tracklist"`
	ExplicitLyrics        bool           `json:"explicit_lyrics"`
	ExplicitContentLyrics int            `json:"explicit_content_lyrics"`
	ExplicitContentCover  int            `json:"explicit_content_cover"`
	Contributors          []Contributors `json:"contributors"`
	Artist                Artist         `json:"artist"`
	Type                  string         `json:"type"`
	Tracks                Tracks         `json:"tracks"`
}

// Genres is for Album
type Genres struct {
	Data []GenresData `json:"data"`
}

// GenresData is for Album
type GenresData struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Type    string `json:"type"`
}

// Controbutors is for Album
type Contributors struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Share         string `json:"share"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	Radio         bool   `json:"radio"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
	Role          string `json:"role"`
}

// Artist is for Album
type Artist struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}

// TracksData is for Album
type TracksData struct {
	ID                    string `json:"id"`
	Readable              bool   `json:"readable"`
	Title                 string `json:"title"`
	TitleShort            string `json:"title_short"`
	TitleVersion          string `json:"title_version"`
	Link                  string `json:"link"`
	Duration              string `json:"duration"`
	Rank                  string `json:"rank"`
	ExplicitLyrics        bool   `json:"explicit_lyrics"`
	ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
	ExplicitContentCover  int    `json:"explicit_content_cover"`
	Preview               string `json:"preview"`
	Md5Image              string `json:"md5_image"`
	Artist                Artist `json:"artist"`
	Type                  string `json:"type"`
}

// Tracks is for Album
type Tracks struct {
	Data []TracksData `json:"data"`
}
