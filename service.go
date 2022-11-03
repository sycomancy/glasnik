package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
	"github.com/sycomancy/glasnik/types"
)

// AdsFetcher is an interface that can fetch a ads
type AdsFetcher interface {
	FetchAds(context.Context, *infra.IncognitoClient, string) ([]types.AdsPageResponse, error)
}

// priceFetcher implements an interface
type adsFetcher struct{}

func (a *adsFetcher) FetchAds(ctx context.Context, ic *infra.IncognitoClient, url string) ([]types.AdsPageResponse, error) {
	result, err := njuskalo.Fetch(url, ic)
	if err != nil {
		fmt.Print(result)
	}

	logrus.WithFields(logrus.Fields{
		"count": len(result),
		"url":   url,
	}).Info("results from njuskalo")

	return result, nil
}
