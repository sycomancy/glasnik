package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

type tokenEntry struct {
	Token           string `bson:"token,omitempty"`
	RequestCount    string `bson:"requestCount,omitempty"`
	AllowedRequests string `bson:"allowedRequests,omitempty"`
}

type tokenValidatorService struct {
	next AdsFetcher
}

func NewTokenValidatorService(next AdsFetcher) AdsFetcher {
	return &tokenValidatorService{
		next: next,
	}
}

func (t *tokenValidatorService) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) (types.RequestResult, error) {
	invalidToken := t.ValidateToken(request)
	if invalidToken != nil {
		return types.RequestResult{}, invalidToken
	}
	return t.next.ProcessRequest(ctx, ic, request)
}

func (t *tokenValidatorService) ValidateToken(request types.RequestData) error {
	token := request.Token
	if token == "" {
		return fmt.Errorf("missing  required parameter: token")
	}

	tokenEntry := FindToken(request.Token)

	if tokenEntry.Token != token {
		logrus.WithFields(logrus.Fields{
			"token":  token,
			"status": "token not found",
		}).Warn("token validation")
		return fmt.Errorf("invalid token provided %s", token)
	}

	return nil
}

func FindToken(token string) tokenEntry {
	var tokenEntry tokenEntry
	infra.FindDocument("tokens", bson.D{{Key: "token", Value: token}}, &tokenEntry)
	return tokenEntry
}

func UpdateResultsForToken(token string, response types.RequestResult) error {
	tokenEntry := FindToken(token)
	if tokenEntry.Token != "" {
		return nil
	}
	return nil
}
