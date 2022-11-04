package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

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

func (t *tokenValidatorService) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) (types.RequestResult, error) {
	invalidToken := t.ValidateToken(request)
	if invalidToken != nil {
		return types.RequestResult{}, invalidToken
	}
	return t.next.ProcessRequest(ctx, ic, request)
}

type tokenEntry struct {
	Token string `bson:"token,omitempty"`
}

func (t *tokenValidatorService) ValidateToken(request types.RequestData) error {
	token := request.Token
	if token == "" {
		return fmt.Errorf("param: token is required")
	}

	var tokenEntry tokenEntry
	infra.FindDocument("tokens", bson.D{{Key: "token", Value: request.Token}}, &tokenEntry)

	if tokenEntry.Token != token {
		logrus.WithFields(logrus.Fields{
			"token":  token,
			"status": "token not found",
		}).Warn("token validation")
		return fmt.Errorf("invalid token provided %s", token)
	}

	return nil
}
