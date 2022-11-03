package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        AdsFetcher
}

func NewJSONAPIServer(listenAddr string, svc AdsFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/api", makeHTTPHandlerFunc(s.handleFetchAds))
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
	// TODO(dhudek): parse url from request
	url := r.URL.Query().Get("ticker")
	data, err := s.svc.FetchAds(ctx, url)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, &data)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
