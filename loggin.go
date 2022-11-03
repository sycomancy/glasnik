package main

import (
	"context"
	"time"

	"github.com/damjackk/njufetch/types"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next AdsFetcher
}

func NewLoggingService(next AdsFetcher) AdsFetcher {
	return &loggingService{
		next: next,
	}
}

func (l *loggingService) FetchAds(ctx context.Context, url string) ([]types.AddPageResponse, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"url":       url,
		}).Info("fetchAds")
	}(time.Now())

	return l.next.FetchAds(ctx, url)
}
