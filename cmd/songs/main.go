package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	presenters "github.com/giubueno/songs/pkg/presenters"
	repositories "github.com/giubueno/songs/pkg/repositories"
)

func getArtistNameFromStdIn() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	var sb strings.Builder

	for j := 0; j < len(output); j++ {
		sb.WriteRune(output[j])
	}

	return sb.String()
}

func getArtistNameFromArgs() string {
	if len(os.Args) < 2 {
		return ""
	}

	var sb strings.Builder

	for i, str := range os.Args[1:] {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(str)
	}
	return sb.String()
}

// Gets the artist name from the command argument list or from a pipe.
func getArtistName() string {

	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error accessing standard input information. More details: %v", err)
		os.Exit(1)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return getArtistNameFromArgs()
	}

	return getArtistNameFromStdIn()
}

// Generates a string with the info about how to use this command.
func printHowToUse() string {
	var sb strings.Builder

	sb.WriteString("Songs\n")
	sb.WriteString("\tUsage: songs [artist name]\n")

	return sb.String()
}

func getURL() string {
	var url string = os.Getenv("GENIUS_API_URL")
	if len(url) == 0 {
		return "https://api.genius.com/search"
	}

	return url
}

func getAccessToken() string {
	accessToken := os.Getenv("CLIENT_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		fmt.Fprintf(os.Stdout, "CLIENT_ACCESS_TOKEN environment variable is not set.")
		os.Exit(1)
	}
	return accessToken
}

func main() {
	url := getURL()
	accessToken := getAccessToken()
	artistName := getArtistName()

	if len(artistName) == 0 {
		fmt.Fprintf(os.Stdout, printHowToUse())
		os.Exit(1)
	}

	repo := repositories.NewRepository(url, accessToken)
	songs, err := repo.FindSongsByArtistName(artistName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error fetching songs. %v\n", err)
		os.Exit(1)
	}

	presenters.NewPresenter(artistName, songs).Render()

	os.Exit(0)
}
