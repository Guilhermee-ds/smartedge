package balancer

import (
	"smartedge/cmd/internal/backend"
	"sync/atomic"
)

type RoundRobin struct {
	manager *backend.Manager
	counter uint64
}

func NewRoundRobin(m *backend.Manager) *RoundRobin {
	return &RoundRobin{manager: m}
}

func (rr *RoundRobin) Next() *backend.Backend {
	backends := rr.manager.GetAlive()
	if len(backends) == 0 {
		return nil
	}
	idx := atomic.AddUint64(&rr.counter, 1)
	return backends[int(idx)%len(backends)]
}

// GeoAffinity: prioriza backend baseado no IP do cliente
func (rr *RoundRobin) NextWithGeo(clientIP string) *backend.Backend {
	backends := rr.manager.GetAlive()
	if len(backends) == 0 {
		return nil
	}
	// Futuro: calcular distância IP → backend
	idx := atomic.AddUint64(&rr.counter, 1)
	return backends[int(idx)%len(backends)]
}
