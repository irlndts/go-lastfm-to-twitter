package lastfm

import (
	"net/url"
	"strconv"
)

const (
	getTopArtistsMethod = "user.getTopArtists"
)

// UserService ...
type UserService struct {
	apiKey string
}

// TopArtistsResponse ...
type TopArtistsResponse struct {
	TopArtists struct {
		Artists []struct {
			Name      string `json:"name"`
			Playcount int    `json:"playcount"`
		} `json:"artist"`
		Attributes struct {
			Page       int `json:"page"`
			PerPage    int `json:"perPage"`
			TotalPages int `json:"totalPages"`
			Total      int `json:"total"`
		} `json:"@attr"`
	} `json:"topartists"`
}

// TopArtists ...
type TopArtists struct {
	Artists    []Artist
	Page       int
	PerPage    int
	TotalPages int
	Total      int
}

// Artist ...
type Artist struct {
	Name      string
	Playcount int
}

// topartists converts response to TopArtists
func (r *TopArtistsResponse) topartists() (*TopArtists, error) {
	top := &TopArtists{
		Page:       r.TopArtists.Attributes.Page,
		PerPage:    r.TopArtists.Attributes.PerPage,
		TotalPages: r.TopArtists.Attributes.TotalPages,
		Total:      r.TopArtists.Attributes.Total,
	}

	if top.Total != 0 {
		for _, artist := range r.TopArtists.Artists {
			top.Artists = append(top.Artists, Artist{Name: artist.Name, Playcount: artist.Playcount})
		}
	}

	return top, nil
}

// TopArtists returns a list of top artists
func (s *UserService) TopArtists(user, period string, limit, offset int) (*TopArtists, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("method", getTopArtistsMethod)

	params.Set("user", user)
	params.Set("period", period)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("format", "json")

	var response *TopArtistsResponse
	if err := request(params, &response); err != nil {
		return nil, err
	}

	return response.topartists()
}
