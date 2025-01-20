package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Registry struct {
	proxies  []ProxyInfo
	mu       sync.RWMutex
	server   *http.Server
	username string
	password string
}

func NewRegistry(port, username, password string) *Registry {
	r := &Registry{
		proxies:  make([]ProxyInfo, 0),
		username: username,
		password: password,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", r.basicAuth(r.handleRegister))

	r.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return r
}

func (r *Registry) Start() error {
	return r.server.ListenAndServe()
}

func (r *Registry) basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Skip auth if credentials are not configured
		if r.username == "" && r.password == "" {
			handler(w, req)
			return
		}

		username, password, ok := req.BasicAuth()
		if !ok || username != r.username || password != r.password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, req)
	}
}

func (r *Registry) handleRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var proxyInfo ProxyInfo
	if err := json.NewDecoder(req.Body).Decode(&proxyInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	r.proxies = append(r.proxies, proxyInfo)
	r.mu.Unlock()

	fmt.Printf("Registered proxy: %+v\n", proxyInfo)

	w.WriteHeader(http.StatusOK)
}

func (r *Registry) RegisterProxies(hosts []string, username, password string) error {
	var wg sync.WaitGroup
	resultChan := make(chan ProxyInfo)
	errorChan := make(chan error)

	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			if details, err := checkProxy(host, username, password); err != nil {
				errorChan <- fmt.Errorf("failed to check proxy %s: %w", host, err)
			} else {
				resultChan <- ProxyInfo{
					URL:      host,
					Username: username,
					Password: password,
					Details:  details,
				}
			}
		}(host)
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	for {
		select {
		case proxy := <-resultChan:
			r.mu.Lock()
			r.proxies = append(r.proxies, proxy)
			r.mu.Unlock()
		case err := <-errorChan:
			return err
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timeout while registering proxies")
		}
	}
}

func checkProxy(host, username, password string) (ServerDetails, error) {
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return ServerDetails{}, err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ServerDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ServerDetails{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var details ServerDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return ServerDetails{}, err
	}

	return details, nil
}
