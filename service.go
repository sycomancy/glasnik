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
	ProcessRequest(context.Context, *infra.IncognitoClient, types.RequestData) (types.RequestResult, error)
}

// priceFetcher implements an interface
type adsFetcher struct{}

func (a *adsFetcher) ProcessRequest(ctx context.Context, ic *infra.IncognitoClient, request types.RequestData) (types.RequestResult, error) {
	response := types.RequestResult{}
	requestId := ctx.Value("requestID").(int)

	if request.CallbackURL != "" {
		response := types.RequestResult{
			Status:      "success",
			CallbackURL: request.CallbackURL,
			RequestID:   requestId,
		}

		go func(request types.RequestData) {
			result, err := njuskalo.Fetch(request.Filter, ic)
			if err != nil {
				// TODO what to do here?
				return
			}

			response.Data = result

			infra.Dispatch(response)
		}(request)

		return response, nil
	}
	result, err := njuskalo.Fetch(request.Filter, ic)

	if err != nil {
		fmt.Print(result)
	}

	logrus.WithFields(logrus.Fields{
		"count":  len(result),
		"filter": request.Filter,
	}).Info("got results from njuskalo")

	response = types.RequestResult{
		Data:        result,
		CallbackURL: request.CallbackURL,
		Status:      "success",
		RequestID:   requestId,
	}

	return response, nil
}
