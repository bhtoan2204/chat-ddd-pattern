package grpc

import (
	"context"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/shared/config"
	baseserver "wechat-clone/core/shared/pkg/server"
	"wechat-clone/core/shared/pkg/stackErr"

	"golang.org/x/sync/errgroup"
	grpcsdk "google.golang.org/grpc"
)

type Server struct {
	cfg            *config.Config
	moduleBuilders []ModuleBuilder
	moduleServers  []GRPCServer
}

type Option func(*Server)

func WithModuleBuilders(builders ...ModuleBuilder) Option {
	return func(s *Server) {
		s.moduleBuilders = append(s.moduleBuilders, builders...)
	}
}

func NewServer(cfg *config.Config, opts ...Option) *Server {
	s := &Server{cfg: cfg}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Enabled() bool {
	return s != nil && len(s.moduleBuilders) > 0 && len(s.ports()) > 0
}

func (s *Server) Start(ctx context.Context, appCtx *appCtx.AppContext) error {
	if !s.Enabled() {
		return nil
	}

	if err := s.buildModuleServers(ctx, appCtx); err != nil {
		return stackErr.Error(err)
	}

	group, groupCtx := errgroup.WithContext(ctx)
	for _, port := range s.ports() {
		port := port
		group.Go(func() error {
			srv, err := baseserver.New(port)
			if err != nil {
				return stackErr.Error(err)
			}
			return stackErr.Error(s.serveWithServer(groupCtx, srv))
		})
	}

	return stackErr.Error(group.Wait())
}

func (s *Server) StartWithServer(ctx context.Context, appCtx *appCtx.AppContext, srv *baseserver.Server) error {
	if s == nil || len(s.moduleBuilders) == 0 {
		return nil
	}
	if err := s.buildModuleServers(ctx, appCtx); err != nil {
		return stackErr.Error(err)
	}

	return stackErr.Error(s.serveWithServer(ctx, srv))
}

func (s *Server) serveWithServer(ctx context.Context, srv *baseserver.Server) error {
	grpcServer := grpcsdk.NewServer()
	for _, moduleServer := range s.moduleServers {
		moduleServer.Register(grpcServer)
	}

	return stackErr.Error(srv.ServeGRPC(ctx, grpcServer))
}

func (s *Server) ports() []int {
	if s == nil || s.cfg == nil {
		return nil
	}

	seen := map[int]struct{}{}
	ports := make([]int, 0, 2)
	for _, port := range []int{s.cfg.GrpcConfig.PaymentGRPCPort, s.cfg.GrpcConfig.LedgerGRPCPort} {
		if port <= 0 {
			continue
		}
		if _, ok := seen[port]; ok {
			continue
		}
		seen[port] = struct{}{}
		ports = append(ports, port)
	}
	return ports
}

func (s *Server) buildModuleServers(ctx context.Context, appCtx *appCtx.AppContext) error {
	if len(s.moduleBuilders) == 0 {
		s.moduleServers = nil
		return nil
	}

	servers, err := BuildModuleServers(ctx, appCtx, s.moduleBuilders...)
	if err != nil {
		return stackErr.Error(err)
	}
	s.moduleServers = servers
	return nil
}
