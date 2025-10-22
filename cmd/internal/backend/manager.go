package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Backend struct {
	URL   string
	Alive bool
	mu    sync.RWMutex
}

type Manager struct {
	backends []*Backend
	mu       sync.RWMutex
}

func NewManager(urls []string) *Manager {
	bs := make([]*Backend, 0, len(urls))
	for _, u := range urls {
		bs = append(bs, &Backend{URL: u, Alive: true})
	}
	return &Manager{backends: bs}
}

func (m *Manager) GetAll() []*Backend {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.backends
}

func (m *Manager) GetAlive() []*Backend {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var alive []*Backend
	for _, b := range m.backends {
		if b.Alive {
			alive = append(alive, b)
		}
	}
	return alive
}

func (m *Manager) HealthCheck() {
	for _, b := range m.backends {
		go func(b *Backend) {
			client := http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get(b.URL + "/health")
			b.mu.Lock()
			if err != nil || resp.StatusCode != http.StatusOK {
				b.Alive = false
				log.Printf("❌ %s offline", b.URL)
			} else {
				b.Alive = true
				log.Printf("✅ %s OK", b.URL)
			}
			if resp != nil {
				resp.Body.Close()
			}
			b.mu.Unlock()
		}(b)
	}
}

func (m *Manager) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(m.GetAll())
	})
}
func (m *Manager) ReloadBackends(urls []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newBackends := make([]*Backend, 0, len(urls))
	for _, u := range urls {
		newBackends = append(newBackends, &Backend{URL: u, Alive: false})
	}
	m.backends = newBackends
	log.Println("🔄 Backends atualizados via /api/reload")
}
