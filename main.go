package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	curr     uint64
}

// Create a new backend instance
func newBackend(rawurl string) *Backend {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		log.Fatalf("Failed to parse URL %s: %v", rawurl, err)
	}
	return &Backend{
		URL:          parsedURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}
}

// Get the next alive backend in round-robin
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

// Serve incoming requests to available backend
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	if backend != nil {
		backend.ReverseProxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service unavailable: no healthy backends", http.StatusServiceUnavailable)
	}
}

// Ping backends to check their health
func healthCheck(lb *LoadBalancer) {
	for {
		for _, backend := range lb.backends {
			resp, err := http.Get(backend.URL.String())
			if err != nil || resp.StatusCode != 200 {
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
	lb := &LoadBalancer{
		backends: []*Backend{
			newBackend("http://localhost:9001"),
			newBackend("http://localhost:9002"),
		},
	}

	go healthCheck(lb)

	log.Println("🔃 Load balancer started on :8080")
	if err := http.ListenAndServe(":8080", lb); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}
