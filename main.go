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
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", user, playlist)

	// define HTTP client and prepare new request
	// https://developer.spotify.com/web-api/authorization-guide/
	client := &http.Client{}
	req, errMakeRequest := http.NewRequest("GET", url, nil)
	if errMakeRequest != nil {
		panic(errMakeRequest)
	}

	authHeader := fmt.Sprintf("Bearer %s", authToken)
	req.Header.Add("Authorization", authHeader)

	resp, errGetResponse := client.Do(req)
	if errGetResponse != nil {
		panic(errGetResponse)
	}
	defer resp.Body.Close()

	// declare body as a byte array of the response body
	body, errReadBody := ioutil.ReadAll(resp.Body)
	if errReadBody != nil {
		panic(errReadBody)
	}

	ioutil.WriteFile("dump.json", body, 0600)
	fmt.Println("Export successful")
}
