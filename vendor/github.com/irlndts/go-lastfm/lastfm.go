package lastfm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	lastfmAPI = "http://ws.audioscrobbler.com/2.0/"
)

// Client ...
type Client struct {
	User *UserService
}

// New returns a new Client.
func New(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key doesn't available")
	}

	return &Client{
		User: &UserService{apiKey: apiKey},
	}, nil
}

func request(params url.Values, resp interface{}) error {
	r, err := http.NewRequest("GET", lastfmAPI+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &resp)
}
