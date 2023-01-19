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

	url := fmt.Sprintf("artist/%v", id)

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
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Share         string `json:"share"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	NbAlbum       int    `json:"nb_album"`
	NbFan         int    `json:"nb_fan"`
	Radio         bool   `json:"radio"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}
