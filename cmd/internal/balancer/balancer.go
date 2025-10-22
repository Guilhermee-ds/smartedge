package balancer

import "smartedge/cmd/internal/backend"

type Balancer interface {
	Next() *backend.Backend
}
