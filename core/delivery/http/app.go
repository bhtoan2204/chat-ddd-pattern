package http

import (
	"context"
	"go-socket/config"
	"go-socket/constant"
	"go-socket/core/assembly"
	appCtx "go-socket/core/context"
	"go-socket/core/delivery/http/middleware"
	"go-socket/core/shared/infra/idempotency"
	"go-socket/core/shared/pkg/server"
	coreusecase "go-socket/core/usecase"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var _ App = (*Server)(nil)

type App interface {
	Routes(ctx context.Context, appCtx *appCtx.AppContext) *gin.Engine
	Start(ctx context.Context, appCtx *appCtx.AppContext) error
}

type Server struct {
	cfg        *config.Config
	router     *gin.Engine
	httpServer *http.Server
	handler    RoutingHandler
	usecase    coreusecase.Usecase
	appCtx     *appCtx.AppContext
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Routes(ctx context.Context, appCtx *appCtx.AppContext) *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 50 << 20
	r.RedirectTrailingSlash = false
	cache := appCtx.GetCache()
	r.Use(middleware.SetRequestID())
	idemStore := idempotency.NewRedisStore(cache)
	idemManager := idempotency.NewManager(
		idemStore,
		constant.DEFAULT_IDEMPOTENCY_LOCK_TTL,
		constant.DEFAULT_IDEMPOTENCY_DONE_TTL,
	)
	r.Use(middleware.IdempotencyMiddleware(idemManager))
	r.Use(middleware.RateLimitMiddleware(cache))
	r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"error": "something went wrong"}})
	}))

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{
		"*",
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Inside-Token",
	}
	r.Use(cors.New(corsConfig))

	pingHandler := func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"clientIP": ctx.ClientIP(),
			},
		})
	}
	r.GET("/health-check", pingHandler)
	r.HEAD("/health-check", pingHandler)

	s.router = r
	s.appCtx = appCtx
	s.handler = NewRoutingHandler(s.cfg, appCtx, appCtx.GetRedisClient(), s.usecase)

	// public api
	s.registerPublicAPI()
	s.registerPrivateAPI()
	return r
}

func (s *Server) Start(ctx context.Context, appCtx *appCtx.AppContext) error {
	uc, err := s.buildUsecase(appCtx)
	if err != nil {
		return err
	}
	s.usecase = uc
	srv, err := server.New(s.cfg.HttpConfig.Port)
	if err != nil {
		return err
	}

	return srv.ServeHTTPHandler(ctx, s.Routes(ctx, appCtx))
}

func (s *Server) buildUsecase(appContext *appCtx.AppContext) (coreusecase.Usecase, error) {
	return assembly.BuildUsecase(appContext), nil
}

func (s *Server) registerPublicAPI() {
	h := s.handler
	router := s.router.Group("")
	apiV1 := router.Group("/api/v1")
	h.RegisterPublicHandlers(apiV1)
}

func (s *Server) registerPrivateAPI() {
	h := s.handler
	router := s.router.Group("")
	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.AuthenMiddleware(s.appCtx))
	h.RegisterPrivateHandlers(apiV1)
}
