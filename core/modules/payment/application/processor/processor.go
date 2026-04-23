package processor

import (
	"context"
	"sync"
	"time"

	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/pkg/logging"

	"go.uber.org/zap"
)

type Processor interface {
	Start() error
	Stop() error
}

type processor struct {
	service  paymentservice.PaymentCommandService
	interval time.Duration

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewProcessor(service paymentservice.PaymentCommandService, interval time.Duration) Processor {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	return &processor{
		service:  service,
		interval: interval,
	}
}

func (p *processor) Start() error {
	if p == nil || p.service == nil {
		return nil
	}
	if p.cancel != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg.Add(1)
	go p.loop(ctx)
	return nil
}

func (p *processor) Stop() error {
	if p == nil || p.cancel == nil {
		return nil
	}
	p.cancel()
	p.wg.Wait()
	p.cancel = nil
	return nil
}

func (p *processor) loop(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		p.runOnce(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (p *processor) runOnce(ctx context.Context) {
	if p == nil || p.service == nil {
		return
	}

	if err := p.service.ProcessPendingWithdrawals(ctx); err != nil {
		logging.FromContext(ctx).Warnw("process pending withdrawals failed", zap.Error(err))
	}
}
