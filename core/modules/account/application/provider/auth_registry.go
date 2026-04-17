package provider

import (
	"fmt"
	"strings"
	"sync"
)

type AuthProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]AuthProvider
}

func NewProviderRegistry() *AuthProviderRegistry {
	return &AuthProviderRegistry{
		providers: make(map[string]AuthProvider),
	}
}

func (r *AuthProviderRegistry) Register(provider AuthProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[strings.ToLower(strings.TrimSpace(provider.Name()))] = provider
}

func (r *AuthProviderRegistry) Get(name string) (AuthProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return nil, fmt.Errorf("%v: %s", ErrProviderNotFound, name)
	}
	return provider, nil
}
