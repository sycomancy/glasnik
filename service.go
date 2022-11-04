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
	ProcessRequest(context.Context, *infra.IncognitoClient, types.RequestData) ([]types.AdsPageResponse, error)
}

// priceFetcher implements an interface
type adsFetcher struct{}

func (a *adsFetcher) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) ([]types.AdsPageResponse, error) {
	result, err := njuskalo.Fetch(request.Filter, ic)
	if err != nil {
		fmt.Print(result)
	}

	logrus.WithFields(logrus.Fields{
		"count":  len(result),
		"filter": request.Filter,
	}).Info("results from njuskalo")

	return result, nil
}
