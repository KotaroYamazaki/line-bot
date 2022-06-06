package googlemap

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

type Client struct {
	client *maps.Client
	key    string
}

func New(key string) (*Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
		key:    key,
	}, nil
}

var PlaceTypeCafe = maps.PlaceTypeCafe

type TextSearchRequest struct {
	Query   string
	OpenNow bool
	Type    maps.PlaceType
}

type TextSearchResponse struct {
	PhotoURL string
	Name     string
	Addr     string
	URL      string
	Rating   float32
}

func (c *Client) TextSearch(ctx context.Context, req *TextSearchRequest) ([]TextSearchResponse, error) {
	placeResp, err := c.client.TextSearch(ctx, &maps.TextSearchRequest{
		Query:    req.Query,
		OpenNow:  req.OpenNow,
		Language: "ja",
		Type:     req.Type,
	})
	if err != nil {
		return nil, err
	}
	searchResp := make([]TextSearchResponse, 0, len(placeResp.Results))
	for _, v := range placeResp.Results {
		var ref string
		if len(v.Photos) > 0 {
			ref = v.Photos[0].PhotoReference
		}
		searchResp = append(searchResp, TextSearchResponse{
			PhotoURL: fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", ref, c.key),
			Name:     v.Name,
			Addr:     v.FormattedAddress,
			Rating:   v.Rating,
			URL: fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f&quary_place_id=%s",
				v.Geometry.Location.Lat,
				v.Geometry.Location.Lng,
				v.PlaceID),
		})
	}
	return searchResp, nil
}
