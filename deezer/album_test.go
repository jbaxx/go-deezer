package deezer

import (
	"context"
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

	album, _, err := client.Albums.Get(context.Background(), "44132881")
	if err != nil {
		t.Errorf("Album.Get return error: %v", err)
	}

	want := &Album{ID: 44132881}
	if !reflect.DeepEqual(album, want) {
		t.Errorf("Album.Get got: %#v, want %#v", album, want)
	}

	_, _, err = client.Albums.Get(context.Background(), "\n")
	if err == nil {
		t.Errorf("bad options err = nil, want error")
	}

}
