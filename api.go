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
	listenAddr string
	svc        AdsFetcher
	ic         *infra.IncognitoClient
}

func NewJSONAPIServer(listenAddr string, svc AdsFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
		ic:         infra.NewIncognitoClient(nil),
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/api/request-njuska", makeHTTPHandlerFunc(s.handleFetch))
	http.ListenAndServe(s.listenAddr, nil)
}

func makeHTTPHandlerFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", rand.Intn(10000000))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFunc(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleFetch(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestData types.RequestDataDTO
	if r.Method != http.MethodPost {
		return fmt.Errorf("unsupported method")
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return fmt.Errorf("unable to decode body")
	}

	requestId := ctx.Value("requestID").(int)
	data, err := s.svc.ProcessRequest(ctx, s.ic, types.RequestData{
		Filter:      requestData.Filter,
		Token:       requestData.Token,
		CallbackURL: requestData.CallbackURL,
		RequestID:   requestId,
	})

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, data)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
