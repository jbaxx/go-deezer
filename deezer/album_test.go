package deezer

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetAlbum(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/album/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"id": 44132881}`)
	})

	album, _, err := client.Albums.Get("44132881")
	if err != nil {
		t.Errorf("Album.Get return error: %v", err)
	}

	want := &Album{ID: 44132881}
	if !reflect.DeepEqual(album, want) {
		t.Errorf("Album.Get got: %#v, want %#v", album, want)
	}

}

func TestGetAlbumRaw(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/album/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"id": 44132881}`)
	})

	albumRaw, _, err := client.Albums.GetRaw("44132881")
	if err != nil {
		t.Errorf("Album.Get return error: %v", err)
	}

	want := []byte(`{"id": 44132881}`)
	if !bytes.Equal(albumRaw, want) {
		t.Errorf("Album.GetRaw got: %#v, want %#v", albumRaw, want)
	}

	_, _, err = client.Albums.GetRaw("\n")

	if err == nil {
		t.Errorf("bad options err = nil, want error")
	}

	// client.URL.Path = ""
	// f := func() (*Response, error) {
	// 	got, resp, err := client.Albums.GetRaw(44132881)
	// 	if got != nil {
	// 		t.Errorf("want nil, got: %v", got)
	// 	}
	// 	return resp, err
	// }
	// resp, err := f()
	// if resp != nil {
	// 	t.Errorf("client.URL.Path='' resp = %#v, want nil", resp)
	// }

}
