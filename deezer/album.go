package deezer

import (
	"bytes"
	"fmt"
	"net/http"
)

// AlbumService defines a client to interface with the Deezer Album service
type AlbumService struct {
	client *Client
}

// Get fetches an Album given an album id.
func (a *AlbumService) Get(id string) (*Album, *Response, error) {

	url := fmt.Sprintf("album/%v", id)

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

// GetRaw fetches an Album given an album id.
func (a *AlbumService) GetRaw(id string) ([]byte, *Response, error) {

	url := fmt.Sprintf("album/%v", id)

	req, err := a.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer
	resp, err := a.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, err

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
	ID            int    `json:"id"`
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
	ID                    int    `json:"id"`
	Readable              bool   `json:"readable"`
	Title                 string `json:"title"`
	TitleShort            string `json:"title_short"`
	TitleVersion          string `json:"title_version"`
	Link                  string `json:"link"`
	Duration              int    `json:"duration"`
	Rank                  int    `json:"rank"`
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
