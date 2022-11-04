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

func (l *loggingService) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) ([]types.AdsPageResponse, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"filter":    request.Filter,
		}).Info("fetchAds")
	}(time.Now())

	return l.next.ProcessRequest(ctx, ic, request)
}
