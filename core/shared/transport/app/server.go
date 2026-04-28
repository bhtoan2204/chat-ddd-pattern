package app

import (
	"context"
	"fmt"
	"sync"

	appCtx "wechat-clone/core/context"
	ledgerassembly "wechat-clone/core/modules/ledger/assembly"
	notificationassembly "wechat-clone/core/modules/notification/assembly"
	paymentassembly "wechat-clone/core/modules/payment/assembly"
	relationshipassembly "wechat-clone/core/modules/relationship/assembly"
	roomassembly "wechat-clone/core/modules/room/assembly"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/logging"
	baseserver "wechat-clone/core/shared/pkg/server"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
	grpctransport "wechat-clone/core/shared/transport/grpc"
	httptransport "wechat-clone/core/shared/transport/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -package=app -destination=server_mock.go -source=server.go
type Server interface {
	Start(ctx context.Context, appCtx *appCtx.AppContext) error
	StartWithServer(ctx context.Context, appCtx *appCtx.AppContext, srv *baseserver.Server) error
}

type appServer struct {
	cfg            *config.Config
	httpServer     *httptransport.Server
	httpOptions    []httptransport.Option
	grpcServer     *grpctransport.Server
	grpcOptions    []grpctransport.Option
	moduleRuntimes []modruntime.Module
}

type Option func(*appServer)

func WithHTTPServer(server *httptransport.Server) Option {
	return func(s *appServer) {
		s.httpServer = server
	}
}

func WithHTTPModuleBuilders(builders ...httptransport.ModuleBuilder) Option {
	return func(s *appServer) {
		s.httpOptions = append(s.httpOptions, httptransport.WithModuleBuilders(builders...))
	}
}

func WithGRPCServer(server *grpctransport.Server) Option {
	return func(s *appServer) {
		s.grpcServer = server
	}
}

func WithGRPCModuleBuilders(builders ...grpctransport.ModuleBuilder) Option {
	return func(s *appServer) {
		s.grpcOptions = append(s.grpcOptions, grpctransport.WithModuleBuilders(builders...))
	}
}

func NewServer(cfg *config.Config, opts ...Option) Server {
	s := &appServer{
		cfg: cfg,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.httpServer == nil {
		s.httpServer = httptransport.NewServer(cfg, s.httpOptions...)
	}
	if s.grpcServer == nil {
		s.grpcServer = grpctransport.NewServer(cfg, s.grpcOptions...)
	}
	return s
}

func (s *appServer) Start(ctx context.Context, appContext *appCtx.AppContext) error {
	return s.StartWithServer(ctx, appContext, nil)
}

func (s *appServer) StartWithServer(ctx context.Context, appContext *appCtx.AppContext, srv *baseserver.Server) error {
	if err := s.buildModuleRuntimes(appContext); err != nil {
		return stackErr.Error(err)
	}

	if err := s.startModuleRuntimes(ctx); err != nil {
		return stackErr.Error(err)
	}
	defer s.stopModuleRuntimes(ctx)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if srv != nil {
			return s.httpServer.StartWithServer(groupCtx, appContext, srv)
		}
		return s.httpServer.Start(groupCtx, appContext)
	})

	if s.grpcServer != nil && s.grpcServer.Enabled() {
		group.Go(func() error {
			return s.grpcServer.Start(groupCtx, appContext)
		})
	}

	return stackErr.Error(group.Wait())
}

func (s *appServer) buildModuleRuntimes(appContext *appCtx.AppContext) error {
	notificationRuntime, err := notificationassembly.BuildMessagingRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build notification messaging runtime failed: %w", err))
	}

	ledgerProjectionRuntime, err := ledgerassembly.BuildProjectionRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build ledger messaging runtime failed: %w", err))
	}

	roomProjectionRuntime, err := roomassembly.BuildProjectionRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build room projection runtime failed: %w", err))
	}

	relationshipMessagingRuntime, err := relationshipassembly.BuildMessagingRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build relationship messaging runtime failed: %w", err))
	}

	paymentMessagingRuntime, err := paymentassembly.BuildMessagingRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build payment messaging runtime failed: %w", err))
	}

	paymentTaskRuntime, err := paymentassembly.BuildTaskRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build payment task runtime failed: %w", err))
	}

	paymentCronRuntime, err := paymentassembly.BuildCronRuntime(s.cfg, appContext)
	if err != nil {
		return stackErr.Error(fmt.Errorf("build payment cron runtime failed: %w", err))
	}

	s.moduleRuntimes = []modruntime.Module{
		notificationRuntime,
		roomProjectionRuntime,
		relationshipMessagingRuntime,
		ledgerProjectionRuntime,
		paymentMessagingRuntime,
		paymentTaskRuntime,
		paymentCronRuntime,
	}
	return nil
}

func (s *appServer) startModuleRuntimes(ctx context.Context) error {
	for idx, runtime := range s.moduleRuntimes {
		if err := runtime.Start(); err != nil {
			s.stopStartedRuntimes(ctx, idx-1)
			return stackErr.Error(fmt.Errorf("start module runtime %T failed: %w", runtime, err))
		}
	}
	return nil
}

func (s *appServer) stopStartedRuntimes(ctx context.Context, lastIdx int) {
	for i := lastIdx; i >= 0; i-- {
		if err := s.moduleRuntimes[i].Stop(); err != nil {
			logging.FromContext(ctx).Errorw("Failed to stop module runtime", zap.Error(err))
		}
	}
}

func (s *appServer) stopModuleRuntimes(ctx context.Context) {
	var wg sync.WaitGroup

	for i := len(s.moduleRuntimes) - 1; i >= 0; i-- {
		runtime := s.moduleRuntimes[i]
		wg.Add(1)
		go func(runtime modruntime.Module) {
			defer wg.Done()
			if err := runtime.Stop(); err != nil {
				logging.FromContext(ctx).Errorw("Failed to stop module runtime", zap.Error(err))
			}
		}(runtime)
	}

	wg.Wait()
}
