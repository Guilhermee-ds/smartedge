package balancer

import (
	"math"
	"smartedge/cmd/internal/backend"
	"sync"
	"time"
)

type ewmaStats struct {
	latency float64
	weight  float64
}

type EWMABalancer struct {
	manager *backend.Manager
	alpha   float64
	stats   map[string]*ewmaStats
	mu      sync.RWMutex
}

func NewEWMABalancer(m *backend.Manager, alpha float64) *EWMABalancer {
	return &EWMABalancer{
		manager: m,
		alpha:   alpha,
		stats:   make(map[string]*ewmaStats),
	}
}

func (e *EWMABalancer) Next() *backend.Backend {
	e.mu.RLock()
	defer e.mu.RUnlock()

	backends := e.manager.GetAlive()
	if len(backends) == 0 {
		return nil
	}

	var best *backend.Backend
	bestScore := math.MaxFloat64
	for _, b := range backends {
		s, ok := e.stats[b.URL]
		if !ok {
			s = &ewmaStats{latency: 100, weight: 1}
			e.stats[b.URL] = s
		}
		if s.latency < bestScore {
			bestScore = s.latency
			best = b
		}
	}
	return best
}

func (e *EWMABalancer) Report(backendURL string, latency time.Duration, success bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	s, ok := e.stats[backendURL]
	if !ok {
		s = &ewmaStats{latency: latency.Seconds() * 1000, weight: 1}
		e.stats[backendURL] = s
	}

	if !success {
		s.latency = s.latency*0.7 + 300
		return
	}

	newVal := latency.Seconds() * 1000
	s.latency = e.alpha*newVal + (1-e.alpha)*s.latency
}
