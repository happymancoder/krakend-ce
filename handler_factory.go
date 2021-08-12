package krakend

import (
	"net/http"
	"strings"

	jose "github.com/devopsfaith/krakend-jose"
	muxjose "github.com/devopsfaith/krakend-jose/mux"
	lua "github.com/devopsfaith/krakend-lua/router/mux"
	metrics "github.com/devopsfaith/krakend-metrics/mux"
	opencensus "github.com/devopsfaith/krakend-opencensus/router/mux"
	krakendrate "github.com/devopsfaith/krakend-ratelimit"
	"github.com/devopsfaith/krakend-ratelimit/juju"
	jujurouter "github.com/devopsfaith/krakend-ratelimit/juju/router"

	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
	"github.com/luraproject/lura/router/mux"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) mux.HandlerFactory {
	handlerFactory := RateLimitHandlerFactory
	handlerFactory = lua.HandlerFactory(logger, handlerFactory, GorillaParamsExtractor)
	handlerFactory = muxjose.HandlerFactory(handlerFactory, GorillaParamsExtractor, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	return handlerFactory
}

type handlerFactory struct{}

// HandlerFactory is the out-of-the-box basic ratelimit handler factory using the default krakend endpoint
// handler for the mux router
var RateLimitHandlerFactory = NewRateLimiterMw(mux.CustomEndpointHandler(mux.NewRequestBuilder(GorillaParamsExtractor)))

// NewRateLimiterMw builds a rate limiting wrapper over the received handler factory.
func NewRateLimiterMw(next mux.HandlerFactory) mux.HandlerFactory {
	return func(remote *config.EndpointConfig, p proxy.Proxy) http.HandlerFunc {
		handlerFunc := next(remote, p)

		cfg := jujurouter.ConfigGetter(remote.ExtraConfig).(jujurouter.Config)
		if cfg == jujurouter.ZeroCfg || (cfg.MaxRate <= 0 && cfg.ClientMaxRate <= 0) {
			return handlerFunc
		}

		if cfg.MaxRate > 0 {
			handlerFunc = NewEndpointRateLimiterMw(juju.NewLimiter(float64(cfg.MaxRate), cfg.MaxRate))(handlerFunc)
		}
		if cfg.ClientMaxRate > 0 {
			switch strings.ToLower(cfg.Strategy) {
			case "header":
				handlerFunc = NewHeaderLimiterMw(cfg.Key, float64(cfg.ClientMaxRate), cfg.ClientMaxRate)(handlerFunc)
			}
		}
		return handlerFunc
	}
}

// EndpointMw is a function that decorates the received handlerFunc with some rateliming logic
type EndpointMw func(http.HandlerFunc) http.HandlerFunc

// NewEndpointRateLimiterMw creates a simple ratelimiter for a given handlerFunc
func NewEndpointRateLimiterMw(tb juju.Limiter) EndpointMw {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !tb.Allow() {
				http.Error(w, krakendrate.ErrLimited.Error(), http.StatusServiceUnavailable)
				return
			}

			next(w, r)
		}
	}
}

// NewHeaderLimiterMw creates a token ratelimiter using the value of a header as a token
func NewHeaderLimiterMw(header string, maxRate float64, capacity int64) EndpointMw {
	return NewTokenLimiterMw(HeaderTokenExtractor(header), juju.NewMemoryStore(maxRate, capacity))
}

// TokenExtractor defines the interface of the functions to use in order to extract a token for each request
type TokenExtractor func(r *http.Request) string

// HeaderTokenExtractor returns a TokenExtractor that looks for the value of the designed header
func HeaderTokenExtractor(header string) TokenExtractor {
	return func(r *http.Request) string { return r.Header.Get(header) }
}

// NewTokenLimiterMw returns a token based ratelimiting endpoint middleware with the received TokenExtractor and LimiterStore
func NewTokenLimiterMw(tokenExtractor TokenExtractor, limiterStore krakendrate.LimiterStore) EndpointMw {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenKey := tokenExtractor(r)
			if tokenKey == "" {
				http.Error(w, krakendrate.ErrLimited.Error(), http.StatusTooManyRequests)
				return
			}
			if !limiterStore(tokenKey).Allow() {
				http.Error(w, krakendrate.ErrLimited.Error(), http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}

func (h handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory) mux.HandlerFactory {
	return NewHandlerFactory(l, m, r)
}
