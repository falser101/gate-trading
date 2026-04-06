package gateway

import (
	"context"
	"sync"

	gateapi "github.com/gate/gateapi-go/v7"
)

// Gate 客户端工厂（用于多用户支持）
type GateClientFactory struct {
	clients map[string]*gateapi.APIClient
	mu      sync.RWMutex
	baseURL string
}

func NewGateClientFactory(baseURL string) *GateClientFactory {
	return &GateClientFactory{
		clients: make(map[string]*gateapi.APIClient),
		baseURL: baseURL,
	}
}

func (f *GateClientFactory) GetClient(apiKey, apiSecret string) *gateapi.APIClient {
	key := apiKey + ":" + apiSecret

	f.mu.RLock()
	if client, ok := f.clients[key]; ok {
		f.mu.RUnlock()
		return client
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	client := f.createClient(apiKey, apiSecret)
	f.clients[key] = client
	return client
}

func (f *GateClientFactory) createClient(apiKey, apiSecret string) *gateapi.APIClient {
	cfg := gateapi.NewConfiguration()
	// 设置基础 URL
	cfg.BasePath = f.baseURL
	// 设置 API Key 和 Secret
	cfg.Key = apiKey
	cfg.Secret = apiSecret

	return gateapi.NewAPIClient(cfg)
}

// 获取带认证的 context
func (f *GateClientFactory) GetContext(apiKey, apiSecret string) context.Context {
	return context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    apiKey,
		Secret: apiSecret,
	})
}
