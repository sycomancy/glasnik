package proxy

import (
	"testing"
	"time"
)

func TestProxyServerAndRegistry(t *testing.T) {
	// Start registry server with auth
	registry := NewRegistry("8080", "admin", "secret")
	go func() {
		if err := registry.Start(); err != nil {
			t.Logf("Registry stopped: %v", err)
		}
	}()

	// Give registry time to start
	time.Sleep(100 * time.Millisecond)

	// Start proxy server
	server := NewProxyServer("localhost", "8081", "testuser", "testpass")
	go func() {
		if err := server.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	// Give proxy server time to start
	time.Sleep(100 * time.Millisecond)

	// Test self-registration
	err := server.RegisterWithRegistry("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to register proxy: %v", err)
	}

	// Verify proxy was registered
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if len(registry.proxies) != 1 {
		t.Errorf("Expected 1 proxy to be registered, got %d", len(registry.proxies))
	}

	// Verify proxy details
	if len(registry.proxies) > 0 {
		proxy := registry.proxies[0]
		if proxy.Username != "testuser" || proxy.Password != "testpass" {
			t.Error("Proxy credentials don't match")
		}

		if proxy.Details.Hostname == "" {
			t.Error("Proxy hostname is empty")
		}

		if proxy.Details.IP == "" {
			t.Error("Proxy IP is empty")
		}
		t.Logf("Registered proxy details: %+v", proxy)
	}
}

func TestUnauthenticatedRegistry(t *testing.T) {
	// Start registry server without auth
	registry := NewRegistry("8081", "", "")
	go func() {
		if err := registry.Start(); err != nil {
			t.Logf("Registry stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Start proxy server
	server := NewProxyServer("localhost", "8082", "testuser", "testpass")
	go func() {
		if err := server.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Test registration with unauthenticated registry
	err := server.RegisterWithRegistry("http://localhost:8081")
	if err != nil {
		t.Fatalf("Failed to register with unauthenticated registry: %v", err)
	}
}

func TestAuthenticatedRegistryFailure(t *testing.T) {
	// Start registry server with auth
	registry := NewRegistry("8083", "admin", "secret")
	go func() {
		if err := registry.Start(); err != nil {
			t.Logf("Registry stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Start proxy server
	server := NewProxyServer("localhost", "8084", "testuser", "testpass")
	go func() {
		if err := server.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Test registration without proper auth
	err := server.RegisterWithRegistry("http://localhost:8083")
	if err == nil {
		t.Error("Expected authentication error, got nil")
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
	registry := NewRegistry("8080", "testuser", "testpass")

	// Test with wrong credentials
	hosts := []string{"http://localhost:8082/details"}
	err := registry.RegisterProxies(hosts, "wronguser", "wrongpass")
	if err == nil {
		t.Error("Expected authentication error, got nil")
	}
}

func TestProxyServerInvalidHost(t *testing.T) {
	registry := NewRegistry("8080", "testuser", "testpass")
	hosts := []string{"http://invalid-host:9999/details"}
	err := registry.RegisterProxies(hosts, "testuser", "testpass")
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}
}
