package deezer

import (
	"context"
	"fmt"
	"net/http"
)

// AlbumService defines a client to interface with the Deezer Album service
type AlbumService struct {
	client *Client
}

// Get fetches an Album given an album id.
func (a *AlbumService) Get(ctx context.Context, id string) (*Album, *Response, error) {

	url := fmt.Sprintf("album/%s", id)

	req, err := a.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	album := new(Album)
	resp, err := a.client.Do(ctx, req, album)
	if err != nil {
		return nil, resp, err
	}

	return album, resp, err

}

// List fetches an Album given an album id using pagination.
func (a *AlbumService) ListTracks(ctx context.Context, id string, opt *ListOptions) (*Tracks, *Response, error) {

	u := fmt.Sprintf("album/%s/tracks", id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := a.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tracks := new(Tracks)
	resp, err := a.client.Do(ctx, req, tracks)
	if err != nil {
		return nil, resp, err
	}

	return tracks, resp, err

}

// Album represents an artist's album.
type Album struct {
	ID                    int             `json:"id,omitempty"`
	Title                 string          `json:"title,omitempty"`
	Upc                   string          `json:"upc,omitempty"`
	Link                  string          `json:"link,omitempty"`
	Share                 string          `json:"share,omitempty"`
	Cover                 string          `json:"cover,omitempty"`
	CoverSmall            string          `json:"cover_small,omitempty"`
	CoverMedium           string          `json:"cover_medium,omitempty"`
	CoverBig              string          `json:"cover_big,omitempty"`
	CoverXl               string          `json:"cover_xl,omitempty"`
	Md5Image              string          `json:"md5_image,omitempty"`
	GenreID               int             `json:"genre_id,omitempty"`
	Genres                *Genres         `json:"genres,omitempty"`
	Label                 string          `json:"label,omitempty"`
	NbTracks              int             `json:"nb_tracks,omitempty"`
	Duration              int             `json:"duration,omitempty"`
	Fans                  int             `json:"fans,omitempty"`
	Rating                int             `json:"rating,omitempty"`
	ReleaseDate           string          `json:"release_date,omitempty"`
	RecordType            string          `json:"record_type,omitempty"`
	Available             bool            `json:"available,omitempty"`
	Tracklist             string          `json:"tracklist,omitempty"`
	ExplicitLyrics        bool            `json:"explicit_lyrics,omitempty"`
	ExplicitContentLyrics int             `json:"explicit_content_lyrics,omitempty"`
	ExplicitContentCover  int             `json:"explicit_content_cover,omitempty"`
	Contributors          []*Contributors `json:"contributors,omitempty"`
	Artist                *Artist         `json:"artist,omitempty"`
	Type                  string          `json:"type,omitempty"`
	Tracks                *Tracks         `json:"tracks,omitempty"`
}

// Genres is for Album
type Genres struct {
	Data []*GenresData `json:"data,omitempty"`
}

// GenresData is for Album
type GenresData struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
	Type    string `json:"type,omitempty"`
}

// Controbutors is for Album
type Contributors struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Link          string `json:"link,omitempty"`
	Share         string `json:"share,omitempty"`
	Picture       string `json:"picture,omitempty"`
	PictureSmall  string `json:"picture_small,omitempty"`
	PictureMedium string `json:"picture_medium,omitempty"`
	PictureBig    string `json:"picture_big,omitempty"`
	PictureXl     string `json:"picture_xl,omitempty"`
	Radio         bool   `json:"radio,omitempty"`
	Tracklist     string `json:"tracklist,omitempty"`
	Type          string `json:"type,omitempty"`
	Role          string `json:"role,omitempty"`
}

// TracksData is for Album
type TracksData struct {
	ID                    int     `json:"id,omitempty"`
	Readable              bool    `json:"readable,omitempty"`
	Title                 string  `json:"title,omitempty"`
	TitleShort            string  `json:"title_short,omitempty"`
	TitleVersion          string  `json:"title_version,omitempty"`
	Link                  string  `json:"link,omitempty"`
	Duration              int     `json:"duration,omitempty"`
	Rank                  int     `json:"rank,omitempty"`
	ExplicitLyrics        bool    `json:"explicit_lyrics,omitempty"`
	ExplicitContentLyrics int     `json:"explicit_content_lyrics,omitempty"`
	ExplicitContentCover  int     `json:"explicit_content_cover,omitempty"`
	Preview               string  `json:"preview,omitempty"`
	Md5Image              string  `json:"md5_image,omitempty"`
	Artist                *Artist `json:"artist,omitempty"`
	Type                  string  `json:"type,omitempty"`
}

// Tracks is for Album
type Tracks struct {
	Data  []*Track `json:"data,omitempty"`
	Total int      `json:"total"`
	Prev  string   `json:"prev"`
	Next  string   `json:"next"`
}
