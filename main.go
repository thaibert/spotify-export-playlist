package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func main() {

	// var albums []string
	// var artists []string
	// var songs []string

	// define the necessary info (user, playlistID and auth token) on the command line
	args := os.Args[1:] // skip the ./main part
	user := args[0]
	playlist := args[1]
	authToken := args[2]

	//
	//
	//- getTracks()
	// tracks.add(newTracks)
	// if (size == 100) {
	// 	offset = 100
	// 	while (size == 100) {
	// 		getTracks(urlWithOffset)
	// 		tracks.add(newTracks)
	// 		offset += 100
	// 	}
	// }
	//
	//

	offset := 0
	fields := "items(track(name%2Cartists(name)%2Calbum(name)))%2Ctotal" // extract song name, album name, artists only. Also total number of songs
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks?fields=%s&limit=100&offset=%v", user, playlist, fields, offset)
	body := getAPIEndpoint(url, authToken)

	json := string(body)
	// total := gjson.Get(json, "total")
	currentAmount := int(gjson.Get(json, "items.#").Num)

	for i := 0; i < currentAmount; i++ {
		track := gjson.Get(json, fmt.Sprintf("items.%v.track", i)).String()

		album := gjson.Get(track, "album.name").Str
		songName := gjson.Get(track, "name").Str

		var artists []string
		artistArray := gjson.Get(track, "artists")
		for artistIndex := 0; artistIndex < len(artistArray.Array()); artistIndex++ {
			currentArtist := gjson.Get(artistArray.String(), fmt.Sprintf("%v.name", artistIndex))
			artists = append(artists, currentArtist.String())
		}
		fmt.Printf("%s   -   %s     in    %s \n", artists, songName, album)

	}
	// tracks = append(tracksForOutput, items...)

	fmt.Println("Export successful")
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
