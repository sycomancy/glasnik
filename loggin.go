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

func (l *loggingService) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) (types.RequestResult, error) {
	logrus.WithFields(logrus.Fields{
		"filter":    request.Filter,
		"requestID": ctx.Value("requestID"),
	}).Info("request start")

	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"filter":    request.Filter,
		}).Info("request end")
	}(time.Now())

	return l.next.ProcessRequest(ctx, ic, request)
}
