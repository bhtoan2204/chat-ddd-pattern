package http

import (
	"context"
	"go-socket/config"
	"go-socket/core/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

type HTTPHandler interface {
	Start(ctx context.Context)
	RegisterHandlers(ctx context.Context, svcCtx *svc.ServiceContext)
}

type httpHandler struct {
	cfg    *config.Config
	server *rest.Server
}

func NewHTTPHandler(ctx context.Context, cfg *config.Config) (HTTPHandler, error) {
	server := rest.MustNewServer(rest.RestConf{
		Port: cfg.HttpConfig.Port,
	})
	return &httpHandler{
		cfg:    cfg,
		server: server,
	}, nil
}

func (h *httpHandler) RegisterHandlers(ctx context.Context, svcCtx *svc.ServiceContext) {
	h.server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		},
	})
}

func (h *httpHandler) Start(ctx context.Context) {
	h.server.Start()
}
