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
	http.HandleFunc("/api", makeHTTPHandlerFunc(s.handleFetchAds))
	http.HandleFunc("/api/fetch-njuska", makeHTTPHandlerFunc(s.handleFetchAdsPOST))
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

func (s *JSONAPIServer) handleFetchAds(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// TODO(sycomancy): parse url from request
	url := r.URL.Query().Get("filter")
	data, err := s.svc.FetchAds(ctx, s.ic, url)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, &data)
}

func (s *JSONAPIServer) handleFetchAdsPOST(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestData types.RequestData
	if r.Method != http.MethodPost {
		return fmt.Errorf("unsupported method")
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return fmt.Errorf("unable to decode body")
	}
	fmt.Print(requestData.Filter)
	return nil
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
