package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

type loggingService struct {
	next AdsFetcher
}

func NewLoggingService(next AdsFetcher) AdsFetcher {
	return &loggingService{
		next: next,
	}
}

func (l *loggingService) FetchAds(ctx context.Context, ic *infra.IncognitoClient, url string) ([]types.AdsPageResponse, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"url":       url,
		}).Info("fetchAds")
	}(time.Now())

	return l.next.FetchAds(ctx, ic, url)
}
