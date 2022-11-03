package main

import (
	"context"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

type tokenValidatorService struct {
	next AdsFetcher
}

func NewTokenValidatorService(next AdsFetcher) AdsFetcher {
	return &tokenValidatorService{
		next: next,
	}
}

func (t *tokenValidatorService) FetchAds(ctx context.Context, ic *infra.IncognitoClient, url string) ([]types.AdsPageResponse, error) {
	return t.next.FetchAds(ctx, ic, url)
}
