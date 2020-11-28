package songs

import (
	"testing"

	models "github.com/giubueno/songs/pkg/models"
)

func TestGetContent(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		songs      []models.Song
		want       string
	}{
		{"Rendering without songs", "Dan Torres", make([]models.Song, 0), "Dan Torres\n"},
		{"Rendering with a single song", "Dan Torres", []models.Song{{"Pretty Song"}}, "Dan Torres\n1 - Pretty Song\n"},
		{"Rendering with two songs", "Dan Torres", []models.Song{{"Pretty"}, {"Song"}}, "Dan Torres\n1 - Pretty\n2 - Song\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := terminalPresenter{tt.artistName, tt.songs}
			got := p.getContent()
			if got != tt.want {
				t.Errorf("getContent(): %v, want %v", got, tt.want)
			}
		})
	}
}
