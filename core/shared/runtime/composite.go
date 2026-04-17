package modruntime

import (
	"fmt"
	"sync"

	"go-socket/core/shared/pkg/stackErr"
)

type compositeModule struct {
	modules []Module
}

func NewComposite(modules ...Module) Module {
	filtered := make([]Module, 0, len(modules))
	for _, module := range modules {
		if module != nil {
			filtered = append(filtered, module)
		}
	}
	return &compositeModule{modules: filtered}
}

func (m *compositeModule) Start() error {
	for idx, module := range m.modules {
		if err := module.Start(); err != nil {
			m.stopStarted(idx - 1)
			return stackErr.Error(fmt.Errorf("start runtime %T failed: %w", module, err))
		}
	}
	return nil
}

func (m *compositeModule) Stop() error {
	var (
		firstErr error
		errMu    sync.Mutex
		wg       sync.WaitGroup
	)

	for idx := len(m.modules) - 1; idx >= 0; idx-- {
		module := m.modules[idx]
		wg.Add(1)
		go func(module Module) {
			defer wg.Done()
			if err := module.Stop(); err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = stackErr.Error(fmt.Errorf("stop runtime %T failed: %w", module, err))
				}
				errMu.Unlock()
			}
		}(module)
	}

	wg.Wait()
	return firstErr
}

func (m *compositeModule) stopStarted(lastIdx int) {
	for idx := lastIdx; idx >= 0; idx-- {
		_ = m.modules[idx].Stop()
	}
}
