package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"smartedge/cmd/internal/backend"
	"smartedge/cmd/internal/balancer"
	"smartedge/cmd/internal/metrics"
)

type LoadBalancer struct {
	balancer balancer.Balancer
}

func NewLoadBalancer(b balancer.Balancer) *LoadBalancer {
	return &LoadBalancer{balancer: b}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var backendInstance *backend.Backend

	// Escolhe o backend com ou sem GeoAffinity
	if rr, ok := lb.balancer.(*balancer.RoundRobin); ok {
		clientIP := r.RemoteAddr
		backendInstance = rr.NextWithGeo(clientIP)
	} else {
		backendInstance = lb.balancer.Next()
	}

	if backendInstance == nil {
		http.Error(w, "Nenhum backend disponível", http.StatusServiceUnavailable)
		return
	}

	target, err := url.Parse(backendInstance.URL)
	if err != nil {
		http.Error(w, "Backend inválido", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// 🔧 Corrige o host para preservar o backend original
	r.Host = target.Host
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme

	success := true
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, e error) {
		success = false
		log.Printf("❌ Erro no backend %s: %v", backendInstance.URL, e)
		http.Error(rw, "Falha no backend", http.StatusBadGateway)
	}

	// 🔹 Log de requisição com tempo
	log.Printf("➡️  %s → %s%s", r.RemoteAddr, backendInstance.URL, r.URL.Path)

	proxy.ServeHTTP(w, r)
	log.Printf("📥 Nova requisição recebida: %s %s", r.Method, r.URL.Path)

	elapsed := time.Since(start)

	metrics.ObserveRequest(backendInstance.URL, elapsed, success)

	if ewma, ok := lb.balancer.(*balancer.EWMABalancer); ok {
		ewma.Report(backendInstance.URL, elapsed, success)
	}

	if success {
		log.Printf("✅ %s respondeu em %v", backendInstance.URL, elapsed)
	}
}
