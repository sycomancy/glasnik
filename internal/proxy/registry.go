package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ProxyInfo struct {
	URL      string
	Username string
	Password string
	Details  ServerDetails
}

type Registry struct {
	proxies []ProxyInfo
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		proxies: make([]ProxyInfo, 0),
	}
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
