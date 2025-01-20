package proxy

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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
	// Basic proxy implementation
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	// Copy body using io.Copy
	_, _ = io.Copy(w, resp.Body)
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
