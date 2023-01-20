package deezer

import (
	"context"
	"fmt"
	"net/http"
)

// ArtistService defines a client to interface with the Deezer Artist service
type ArtistService struct {
	client *Client
}

// Get fetches an Artist given an artist id.
func (a *ArtistService) Get(ctx context.Context, id string) (*Artist, *Response, error) {

	url := fmt.Sprintf("artist/%s", id)

	req, err := a.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	artist := new(Artist)
	resp, err := a.client.Do(ctx, req, artist)
	if err != nil {
		return nil, resp, err
	}

	return artist, resp, err

}

// Artist represents an artist
type Artist struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Link          string `json:"link,omitempty"`
	Share         string `json:"share,omitempty"`
	Picture       string `json:"picture,omitempty"`
	PictureSmall  string `json:"picture_small,omitempty"`
	PictureMedium string `json:"picture_medium,omitempty"`
	PictureBig    string `json:"picture_big,omitempty"`
	PictureXl     string `json:"picture_xl,omitempty"`
	NbAlbum       int    `json:"nb_album,omitempty"`
	NbFan         int    `json:"nb_fan,omitempty"`
	Radio         bool   `json:"radio,omitempty"`
	Tracklist     string `json:"tracklist,omitempty"`
	Type          string `json:"type,omitempty"`
}
