package main

import (
	"context"
	"fmt"

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

func (t *tokenValidatorService) ValidateToken(request types.RequestData) error {
	token := request.Token
	if token == "" {
		return fmt.Errorf("param: token is required")
	}

	return nil
}
