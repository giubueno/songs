package songs

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHandler func(w http.ResponseWriter, r *http.Request)

func serverMock(mock MockHandler) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/search", mock)

	return httptest.NewServer(handler)
}

func songsMock(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("../../test/songs.json")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = w.Write(content)
}

func mock500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func mock404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func mock401(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func TestFindSongsByArtist(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		wantErr    bool
		quantity   int
		mock       MockHandler
	}{
		{"Searching with empty string", "", false, 0, songsMock},
		{"Searching songs from Dan Torres", "Dan Torres", false, 10, songsMock},
		{"Searching songs with API error", "Dan Torres", true, 0, mock500},
		{"Searching songs with API not found", "Dan Torres", true, 0, mock404},
		{"Searching songs with API unauthorized", "Dan Torres", true, 0, mock401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := serverMock(tt.mock)
			defer srv.Close()
			repo := NewRepository(srv.URL+"/search", "accessToken")
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
