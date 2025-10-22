package main

import (
	"log"
	"net/http"
	"time"

	"smartedge/cmd/internal/backend"
	"smartedge/cmd/internal/balancer"
	"smartedge/cmd/internal/metrics"
	"smartedge/cmd/internal/proxy"
)

func main() {
	manager := backend.NewManager([]string{
		"http://localhost:8081",
		"http://localhost:8082",
	})

	// EWMA adaptativo
	bal := balancer.NewEWMABalancer(manager, 0.3)

	lb := proxy.NewLoadBalancer(bal)

	// Métricas Prometheus
	metrics.Init()

	// Health check periódico
	go func() {
		for {
			manager.HealthCheck()
			time.Sleep(10 * time.Second)
		}
	}()

	// Descoberta automática via Consul
	go func() {
		for {
			backends := backend.DiscoverBackends()
			manager.ReloadBackends(backends)
			time.Sleep(10 * time.Second)
		}
	}()

	// API hot reload
	http.HandleFunc("/api/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		manager.ReloadBackends([]string{
			"http://localhost:8081",
			"http://localhost:8082",
		})
		w.Write([]byte("✅ Backends recarregados"))
	})

	http.Handle("/metrics", metrics.Handler())

	http.Handle("/", lb)

	log.Println("🚀 SmartEdge iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
