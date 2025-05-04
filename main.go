package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	Weight       int
}

type BackendConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

type LoadBalancer struct {
	backends []*Backend
	curr     uint64
}

// Load backends from JSON config
func loadBackends(path string) []*Backend {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config: %v", err)
	}
	defer file.Close()

	var backendConfigs []BackendConfig
	if err := json.NewDecoder(file).Decode(&backendConfigs); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	var backends []*Backend
	for _, cfg := range backendConfigs {
		parsedURL, err := url.Parse(cfg.URL)
		if err != nil {
			log.Printf("Invalid backend URL %s: %v", cfg.URL, err)
			continue
		}
		for i := 0; i < cfg.Weight; i++ {
			backends = append(backends, &Backend{
				URL:          parsedURL,
				Alive:        true,
				Weight:       cfg.Weight,
				ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
			})
		}
	}
	return backends
}

// Select next alive backend (weighted round-robin)
func (lb *LoadBalancer) getNextBackend() *Backend {
	total := uint64(len(lb.backends))
	for i := uint64(0); i < total; i++ {
		idx := atomic.AddUint64(&lb.curr, 1) % total
		if lb.backends[idx].Alive {
			return lb.backends[idx]
		}
	}
	return nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	if backend != nil {
		backend.ReverseProxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
	}
}

// Periodically check backend health
func healthCheck(lb *LoadBalancer) {
	for {
		for _, backend := range lb.backends {
			resp, err := http.Get(backend.URL.String())
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Printf("❌ %s marked as dead", backend.URL)
				backend.Alive = false
			} else {
				log.Printf("✅ %s is alive", backend.URL)
				backend.Alive = true
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	backends := loadBackends("config/backends.json")
	if len(backends) == 0 {
		log.Fatal("No valid backends loaded")
	}

	lb := &LoadBalancer{backends: backends}
	go healthCheck(lb)

	log.Println("🔃 Load balancer started on :8080")
	if err := http.ListenAndServe(":8080", lb); err != nil {
		log.Fatalf("Load balancer failed: %v", err)
	}
}
