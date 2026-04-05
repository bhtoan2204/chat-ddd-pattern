package transport

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/infra/proxy"
	stackErr "gateway/pkg/stackErr"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/hashicorp/consul/api"
)

type HTTPTransport struct {
	cfg          *config.Config
	server       *http.Server
	consulClient *api.Client
}

func NewHTTPTransport(ctx context.Context, cfg *config.Config) (*HTTPTransport, error) {
	consulClient, err := proxy.NewConsulClient(ctx, cfg)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return &HTTPTransport{
		cfg:          cfg,
		consulClient: consulClient,
	}, nil
}

func normalizeListenAddr(value string) string {
	if _, _, err := net.SplitHostPort(value); err == nil {
		return value
	}
	if !strings.Contains(value, ":") {
		return ":" + value
	}
	return value
}

func (t *HTTPTransport) Start() error {
	addr := normalizeListenAddr(t.cfg.HTTP.Port)

	// For now, we don't have mesh networking, so we need to proxy the request to the target service
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			services, _, err := t.consulClient.Health().Service("go-socket", "", true, nil)
			if err != nil || len(services) == 0 {
				// Cố tình trỏ đến một host lỗi để ErrorHandler phía dưới bắt được
				req.URL.Scheme = "http"
				req.URL.Host = "service-not-found"
				return
			}
			service := services[0].Service
			targetHost := fmt.Sprintf("%s:%d", service.Address, service.Port)
			req.URL.Scheme = "http"
			req.URL.Host = targetHost
			req.Host = targetHost
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Bad Gateway: Service 'go-socket' is currently unavailable or not found in Consul"}`))
		},
	}
	t.server = &http.Server{
		Addr:    addr,
		Handler: proxy,
	}
	return t.server.ListenAndServe()
}

func (t *HTTPTransport) Stop() error {
	return t.server.Shutdown(context.Background())
}
