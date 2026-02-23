package http

import (
	"context"
	"go-socket/config"
	"go-socket/constant"
	appCtx "go-socket/core/context"
	accountusecase "go-socket/core/modules/account/application/usecase"
	accountassembly "go-socket/core/modules/account/assembly"
	accounthttp "go-socket/core/modules/account/transport/http"
	roomusecase "go-socket/core/modules/room/application/usecase"
	roomassembly "go-socket/core/modules/room/assembly"
	roomhttp "go-socket/core/modules/room/transport/http"
	"go-socket/core/shared/infra/idempotency"
	"go-socket/core/shared/pkg/server"
	"go-socket/core/shared/transport/http/middleware"
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
	cfg            *config.Config
	router         *gin.Engine
	httpServer     *http.Server
	authUsecase    accountusecase.AuthUsecase
	roomUsecase    roomusecase.RoomUsecase
	messageUsecase roomusecase.MessageUsecase
	appCtx         *appCtx.AppContext
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

	// public api
	s.registerPublicAPI()
	s.registerPrivateAPI()
	return r
}

func (s *Server) Start(ctx context.Context, appCtx *appCtx.AppContext) error {
	if err := s.buildUsecases(appCtx); err != nil {
		return err
	}
	srv, err := server.New(s.cfg.HttpConfig.Port)
	if err != nil {
		return err
	}

	return srv.ServeHTTPHandler(ctx, s.Routes(ctx, appCtx))
}

func (s *Server) buildUsecases(appContext *appCtx.AppContext) error {
	s.authUsecase = accountassembly.BuildAuthUsecase(appContext)
	roomUsecases := roomassembly.BuildUsecases(appContext)
	s.roomUsecase = roomUsecases.Room
	s.messageUsecase = roomUsecases.Message
	return nil
}

func (s *Server) registerPublicAPI() {
	apiV1 := s.router.Group("/api/v1")
	accounthttp.RegisterPublicRoutes(apiV1, s.authUsecase)
}

func (s *Server) registerPrivateAPI() {
	apiV1 := s.router.Group("/api/v1")
	apiV1.Use(middleware.AuthenMiddleware(s.appCtx))
	accounthttp.RegisterPrivateRoutes(apiV1, s.authUsecase)
	roomhttp.RegisterPrivateRoutes(apiV1, s.roomUsecase, s.messageUsecase)
}
