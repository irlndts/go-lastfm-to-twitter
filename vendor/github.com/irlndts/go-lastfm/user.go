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
			Playcount string `json:"playcount"`
		} `json:"artist"`
		Attributes struct {
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			TotalPages string `json:"totalPages"`
			Total      string `json:"total"`
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
	page, err := strconv.Atoi(r.TopArtists.Attributes.Page)
	if err != nil {
		return nil, err
	}
	perPage, err := strconv.Atoi(r.TopArtists.Attributes.PerPage)
	if err != nil {
		return nil, err
	}
	totalPages, err := strconv.Atoi(r.TopArtists.Attributes.TotalPages)
	if err != nil {
		return nil, err
	}
	total, err := strconv.Atoi(r.TopArtists.Attributes.Total)
	if err != nil {
		return nil, err
	}
	top := &TopArtists{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Total:      total,
	}

	if top.Total != 0 {
		for _, artist := range r.TopArtists.Artists {
			playcount, err := strconv.Atoi(artist.Playcount)
			if err != nil {
				return nil, err
			}
			top.Artists = append(top.Artists, Artist{Name: artist.Name, Playcount: playcount})
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
