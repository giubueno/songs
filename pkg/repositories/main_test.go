package songs

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const CLIENT_ACCESS_TOKEN = "XXXXXXX"

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/search", songsMock)

	return httptest.NewServer(handler)
}

func songsMock(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("../../test/songs.json")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = w.Write(content)
}

func TestFindSongsByArtist(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	repo := NewRepository(srv.URL+"/search", CLIENT_ACCESS_TOKEN)

	tests := []struct {
		name       string
		artistName string
		wantErr    bool
		quantity   int
	}{

		{"Searching with empty string", "", false, 0},
		{"Searching songs from Dan Torres", "Dan Torres", false, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindSongsByArtistName(tt.artistName)
			if tt.wantErr && err == nil {
				t.Errorf("FindSongsByArtist() was expected to return error but didn't")
			} else if !tt.wantErr && err != nil {
				t.Errorf("FindSongsByArtist() wasn't expected to return error %v", err)
			} else if tt.quantity != len(got) {
				t.Errorf("FindSongsByArtist() returned %v songs, want %v", len(got), tt.quantity)
			}
		})
	}
}
