// Package songs provides the logic necessary to find all the songs of a specific artist.
package songs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	models "github.com/giubueno/songs/pkg/models"
)

type meta struct {
	Status int `json:"status"`
}

type hit struct {
	Type   string      `json:"type"`
	Result models.Song `json:"result"`
}

type response struct {
	Hits []hit `json:"hits"`
}

type result struct {
	Meta     meta     `json:"meta"`
	Response response `json:"response"`
}

func (r result) toSongs() []models.Song {
	songs := make([]models.Song, len(r.Response.Hits))
	for i, hit := range r.Response.Hits {
		songs[i] = hit.Result
	}
	return songs
}

type repository struct {
	BaseURL     string
	AccessToken string
}

func (r repository) FindSongsByArtistName(name string) ([]models.Song, error) {
	if len(name) == 0 {
		return make([]models.Song, 0), nil
	}

	query := url.QueryEscape(name)
	url := fmt.Sprintf("%s?q=%s", r.BaseURL, query)

	client := &http.Client{}

	var emptySlice []models.Song = make([]models.Song, 0)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptySlice, fmt.Errorf("Could not create a request, %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+r.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return emptySlice, fmt.Errorf("Could not fetch Genius, %v", err)
	}

	if resp.StatusCode != 200 {
		return emptySlice, fmt.Errorf("Failed to fetch songs. Status: %s", resp.Status)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var apiResult result
	json.Unmarshal(bodyBytes, &apiResult)

	return apiResult.toSongs(), nil
}

// Repository interface for all implementation of repositories.
type Repository interface {
	FindSongsByArtistName(name string) ([]models.Song, error)
}

func NewRepository(url string, token string) Repository {
	return repository{url, token}
}
