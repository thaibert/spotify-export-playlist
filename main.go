package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

func main() {
	unformattedURL := "https://api.spotify.com/v1/users/%s/playlists/%s/tracks?fields=%s&limit=100&offset=%v"
	fields := "items(track(name%2Cartists(name)%2Calbum(name)))%2Ctotal" // extract song name, album name, artists only. Also total number of songs

	var albums []string
	var songs []string
	var artists [][]string

	// define the necessary info (user, playlistID and auth token) on the command line
	args := os.Args[1:] // skip the ./main part
	user := args[0]
	playlist := args[1]
	authToken := args[2]

	// get up to first 100 tracks
	offset := 0
	url := fmt.Sprintf(unformattedURL, user, playlist, fields, offset)

	body := getAPIEndpoint(url, authToken)
	json := string(body)
	total := int(gjson.Get(json, "total").Num) // disregard ugly conversion :(
	trackCount := int(gjson.Get(json, "items.#").Num)

	extractSongData(body, &songs, &albums, &artists)
	offset += trackCount
	// if totalTracks > 100, loop through the playlist with offsets increasing by 100 until end is reached
	for offset < total {
		url := fmt.Sprintf(unformattedURL, user, playlist, fields, offset)
		body := getAPIEndpoint(url, authToken)
		json := string(body)
		extractSongData(body, &songs, &albums, &artists)

		trackCount = int(gjson.Get(json, "items.#").Num)
		offset += trackCount

	}

	// prepare output file with search terms
	nowTime := time.Now().String()
	filename := fmt.Sprintf("%s  %s.txt", nowTime, playlist)
	file, errCreateFile := os.Create(filename)
	if errCreateFile != nil {
		panic(errCreateFile)
	}
	defer file.Close()

	for i := 0; i < total; i++ {
		// split up the array of artists on each individual song into one string
		var artist string
		for aIndex := 0; aIndex < len(artists[i]); aIndex++ {
			artist = artist + " " + artists[i][aIndex]
		}
		song := songs[i]

		searchTerm := fmt.Sprintf("%s %s\n", artist, song)
		fmt.Fprintf(file, searchTerm)
	}

	fmt.Println("Export successful")
}

func extractSongData(body []byte, songArray *[]string, albumArray *[]string, artistArray *[][]string) {
	// takes in the json data from the Spotify API endpoint and appends the
	// song name, album name and artists to their respective arrays.
	json := string(body)

	trackCount := int(gjson.Get(json, "items.#").Num)
	for i := 0; i < trackCount; i++ {
		trackObject := gjson.Get(json, fmt.Sprintf("items.%v.track", i)).String()
		album := gjson.Get(trackObject, "album.name").Str
		songName := gjson.Get(trackObject, "name").Str

		var artistOutput []string
		artistsFromJSON := gjson.Get(trackObject, "artists")
		for index := 0; index < len(artistsFromJSON.Array()); index++ {
			currentArtist := gjson.Get(artistsFromJSON.String(), fmt.Sprintf("%v.name", index))
			artistOutput = append(artistOutput, currentArtist.String())
		}
		*songArray = append(*songArray, songName)
		*albumArray = append(*albumArray, album)
		*artistArray = append(*artistArray, artistOutput)
	}
}

func getAPIEndpoint(url string, authToken string) []byte {
	// define HTTP client and prepare new request
	client := &http.Client{}
	req, errMakeRequest := http.NewRequest("GET", url, nil)
	if errMakeRequest != nil {
		panic(errMakeRequest)
	}

	// add authorization header for the Spotify API
	// https://developer.spotify.com/web-api/authorization-guide/
	authHeader := fmt.Sprintf("Bearer %s", authToken)
	req.Header.Add("Authorization", authHeader)

	resp, errGetResponse := client.Do(req)
	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}
	if errGetResponse != nil {
		panic(errGetResponse)
	}
	defer resp.Body.Close()

	// declare a byte array of the response body
	body, errReadBody := ioutil.ReadAll(resp.Body)
	if errReadBody != nil {
		panic(errReadBody)
	}
	return body
}
