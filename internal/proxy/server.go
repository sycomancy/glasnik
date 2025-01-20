package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type ProxyServer struct {
	Host     string
	Port     string
	Username string
	Password string
	server   *http.Server
}

type ServerDetails struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type ProxyInfo struct {
	URL      string        `json:"url"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Details  ServerDetails `json:"details"`
}

func NewProxyServer(host, port, username, password string) *ProxyServer {
	return &ProxyServer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (p *ProxyServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/details", p.handleDetails)
	mux.HandleFunc("/", p.handleProxy)
	p.server = &http.Server{
		Addr:    p.Host + ":" + p.Port,
		Handler: p.basicAuth(mux),
	}

	log.Printf("Proxy server started on %s:%s", p.Host, p.Port)
	return p.server.ListenAndServe()
}

func (p *ProxyServer) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != p.Username || password != p.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (p *ProxyServer) handleDetails(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	ip, _ := p.getOutboundIP()

	details := ServerDetails{
		Hostname: hostname,
		IP:       ip,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (p *ProxyServer) handleProxy(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "missing target URL", http.StatusBadRequest)
		return
	}

	// Create new request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy original headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *ProxyServer) getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func (p *ProxyServer) RegisterWithRegistry(registryURL string) error {
	details := ServerDetails{}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	details.Hostname = hostname

	ip, err := p.getOutboundIP()
	if err != nil {
		return err
	}
	details.IP = ip

	proxyInfo := ProxyInfo{
		URL:      fmt.Sprintf("http://%s:%s", p.Host, p.Port),
		Username: p.Username,
		Password: p.Password,
		Details:  details,
	}

	jsonData, err := json.Marshal(proxyInfo)
	if err != nil {
		return err
	}

	resp, err := http.Post(registryURL+"/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register with registry: %d", resp.StatusCode)
	}

	return nil
}
