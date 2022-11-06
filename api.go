package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr      string
	svc             AdsFetcher
	ic              *infra.IncognitoClient
	asyncDispatcher *infra.Dispatcher
}

func NewJSONAPIServer(listenAddr string, svc AdsFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr:      listenAddr,
		svc:             svc,
		ic:              infra.NewIncognitoClient(nil),
		asyncDispatcher: infra.NewDispatcher(),
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/api/request-njuska", makeHTTPHandlerFunc(s.handleFetch))
	http.ListenAndServe(s.listenAddr, nil)
}

func makeHTTPHandlerFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "requestID", rand.Intn(10000000))
		if err := apiFunc(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleFetch(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestData types.Request
	if r.Method != http.MethodPost {
		return fmt.Errorf("unsupported method")
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return fmt.Errorf("unable to decode body")
	}

	requestId := ctx.Value("requestID").(int)

	request := types.RequestData{
		Filter:      requestData.Filter,
		Token:       requestData.Token,
		CallbackURL: requestData.CallbackURL,
		RequestID:   requestId,
	}

	if requestData.CallbackURL != "" {
		response := types.RequestResult{
			Status:      "success",
			CallbackURL: requestData.CallbackURL,
			RequestID:   requestId,
		}

		go func(request types.RequestData) {
			result, err := s.svc.ProcessRequest(ctx, s.ic, request)

			if err != nil {
				// TODO what to do here?
				return
			}

			s.asyncDispatcher.Dispatch(result)
		}(request)

		return writeJSON(w, http.StatusAccepted, response)
	}

	response, err := s.svc.ProcessRequest(ctx, s.ic, request)

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
