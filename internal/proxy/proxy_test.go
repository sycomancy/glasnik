package proxy

import (
	"testing"
	"time"
)

func TestProxyServerAndRegistry(t *testing.T) {
	// Start proxy server
	server := NewProxyServer("localhost", "8081", "testuser", "testpass")
	go func() {
		if err := server.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(1000 * time.Millisecond)

	// Create registry
	registry := NewRegistry()

	// Test registering the proxy
	hosts := []string{"http://localhost:8081/details"}
	err := registry.RegisterProxies(hosts, "testuser", "testpass")
	if err != nil {
		t.Fatalf("Failed to register proxy: %v", err)
	}

	// Verify proxy was registered
	registry.mu.RLock()
	if len(registry.proxies) != 1 {
		t.Errorf("Expected 1 proxy to be registered, got %d", len(registry.proxies))
	}

	// Verify proxy details
	proxy := registry.proxies[0]
	if proxy.Username != "testuser" || proxy.Password != "testpass" {
		t.Error("Proxy credentials don't match")
	}

	if proxy.Details.Hostname == "" {
		t.Error("Proxy hostname is empty")
	}
	t.Log(proxy.Details.IP)
	if proxy.Details.IP == "" {
		t.Error("Proxy IP is empty")
	}
	registry.mu.RUnlock()

	// Test with wrong credentials
	err = registry.RegisterProxies(hosts, "wronguser", "wrongpass")
	if err == nil {
		t.Error("Expected error with wrong credentials, got nil")
	}
}

func TestProxyServerAuthFailure(t *testing.T) {
	// Start proxy server
	server := NewProxyServer("localhost", "8082", "testuser", "testpass")
	go func() {
		if err := server.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create registry
	registry := NewRegistry()

	// Test with wrong credentials
	hosts := []string{"http://localhost:8082/details"}
	err := registry.RegisterProxies(hosts, "wronguser", "wrongpass")
	if err == nil {
		t.Error("Expected authentication error, got nil")
	}
}

func TestProxyServerInvalidHost(t *testing.T) {
	registry := NewRegistry()
	hosts := []string{"http://invalid-host:9999/details"}
	err := registry.RegisterProxies(hosts, "testuser", "testpass")
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}
}
