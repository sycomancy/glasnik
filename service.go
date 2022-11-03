package main

import (
	"context"

	"github.com/damjackk/njufetch/types"
)

// AdsFetcher is an interface that can fetch a ads
type AdsFetcher interface {
	FetchAds(context.Context, string) ([]types.AddPageResponse, error)
}

// priceFetcher implements an interface
type adsFetcher struct{}

func (a *adsFetcher) FetchAds(ctx context.Context, url string) ([]types.AddPageResponse, error) {
	return MockAdsFetcher(ctx, url)
}

var adsMock = []types.AddPageResponse{{
	Id:    "1",
	Title: "Naslov1",
	Link:  "http://url",
	Price: 12323,
}}

func MockAdsFetcher(ctx context.Context, url string) ([]types.AddPageResponse, error) {
	// mimick http roundtrip
	//time.Sleep(2 * time.Second)
	return adsMock, nil
}
