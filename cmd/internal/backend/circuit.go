package backend

import (
	"log"
	"sync"
	"time"
)

type CircuitBreaker struct {
	failures     int
	state        string // "closed", "open", "half-open"
	lastFailTime time.Time
	mu           sync.Mutex
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{state: "closed"}
}

func (cb *CircuitBreaker) Report(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case "open":
		if time.Since(cb.lastFailTime) > 10*time.Second {
			cb.state = "half-open"
		}
	case "half-open":
		if success {
			cb.state = "closed"
			cb.failures = 0
		} else {
			cb.state = "open"
			cb.lastFailTime = time.Now()
			log.Println("⚡ Circuit breaker reaberto!")
		}
	case "closed":
		if !success {
			cb.failures++
			if cb.failures >= 3 {
				cb.state = "open"
				cb.lastFailTime = time.Now()
				log.Println("⚡ Circuit breaker aberto!")
			}
		} else {
			cb.failures = 0
		}
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state != "open"
}
