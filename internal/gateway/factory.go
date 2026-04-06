package gateway

import (
	"sync"

	"github.com/falser101/gate-trading/pkg/gateapi"
)

// Gate 客户端工厂（用于多用户支持）
type GateClientFactory struct {
	clients map[string]*gateapi.Client
	mu      sync.RWMutex
	baseURL string
}

func NewGateClientFactory(baseURL string) *GateClientFactory {
	return &GateClientFactory{
		clients: make(map[string]*gateapi.Client),
		baseURL: baseURL,
	}
}

func (f *GateClientFactory) GetClient(apiKey, apiSecret string) *gateapi.Client {
	key := apiKey + ":" + apiSecret

	f.mu.RLock()
	if client, ok := f.clients[key]; ok {
		f.mu.RUnlock()
		return client
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	client := gateapi.NewClient(apiKey, apiSecret, f.baseURL)
	f.clients[key] = client
	return client
}
