package songs

import (
	"fmt"
	"os"
	"strings"

	models "github.com/giubueno/songs/pkg/models"
)

type terminalPresenter struct {
	artistName string
	songs      []models.Song
}

// "Prints" all the song names as a string.
func (p terminalPresenter) getContent() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", p.artistName))
	for i, song := range p.songs {
		sb.WriteString(fmt.Sprintf("%d - %s\n", i+1, song.FullTitle))
	}

	return sb.String()
}

func (p terminalPresenter) Render() {
	fmt.Fprintf(os.Stdout, p.getContent())
}

// Presenter interface for all presenters.
type Presenter interface {
	Render()
}

// NewPresenter instantiates a new presenter.
func NewPresenter(artistName string, songs []models.Song) Presenter {
	return terminalPresenter{artistName, songs}
}
