package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) FindPrice(ctx context.Context, key string) (*types.FindPriceResponse, error) {
	resp, err := http.Get(c.url + "?ticker=" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &types.FindPriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
