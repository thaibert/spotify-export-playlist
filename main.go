package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	authToken := ""
	user := ""
	playlist := ""
	offset := 0

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks?limit=100&offset=%v", user, playlist, offset)

	// ioutil.WriteFile("dump.json", body, 0600)
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
