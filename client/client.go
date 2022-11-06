package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sycomancy/glasnik/types"
)

type Client struct {
	listenAddr  string
	path        string
	serviceAddr string
	Data        chan *types.RequestResult
}

func NewClient(listenAddr string, path string, serviceAddr string) *Client {
	if path == "" {
		path = "/result-hook"
	}

	if listenAddr == "" {
		listenAddr = ":3333"
	}

	return &Client{
		listenAddr:  listenAddr,
		path:        path,
		serviceAddr: serviceAddr,
		Data:        make(chan *types.RequestResult),
	}
}

func (c *Client) Run() {
	http.HandleFunc(c.path, c.handleResultDataReceived)
	fmt.Printf("client listener started on %s %s\n", c.listenAddr, c.path)
	http.ListenAndServe(c.listenAddr, nil)
}

func (c *Client) handleResultDataReceived(w http.ResponseWriter, r *http.Request) {
	var result *types.RequestResult
	err := json.NewDecoder(r.Body).Decode(&result)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte("OKK"))
	c.Data <- result
}

func (c *Client) SendRequest(request *types.Request) (types.RequestResult, error) {
	postBody, err := json.Marshal(request)
	requestBody := bytes.NewBuffer(postBody)
	if err != nil {
		return types.RequestResult{}, err
	}

	resp, err := http.Post(c.serviceAddr, "application/json", requestBody)

	if err != nil {
		return types.RequestResult{}, err
	}

	defer resp.Body.Close()

	var result types.RequestResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.RequestResult{}, err
	}

	return result, nil
}
