// Package proxy used for forwarding requests to micro services
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
	"go.uber.org/zap"
)

type HTTPProxy struct {
	simulationURL *url.URL
	proxy         *httputil.ReverseProxy
	log           *zap.Logger
}

func NewHTTPProxy(cfg *config.Config, l *zap.Logger) (*HTTPProxy, error) {
	simURL, err := url.Parse(cfg.Services.SimulationServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(simURL)

	proxy.Transport = &http.Transport{
		MaxIdleConns:          cfg.Server.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.Server.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.Server.IdleTimeout,
		ResponseHeaderTimeout: cfg.Server.ResponseHeaderTimeout,
	}

	return &HTTPProxy{
		simulationURL: simURL,
		proxy:         proxy,
		log:           l,
	}, nil
}

func (p *HTTPProxy) ProxyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		p.proxy.ServeHTTP(w, r)

		p.log.Info("Proxy Request",
			zap.String("method", r.Method),
			zap.String("URL path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
