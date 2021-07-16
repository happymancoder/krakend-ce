package krakend

import (
	//botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	//"net/http"

	jose "github.com/devopsfaith/krakend-jose"
	muxjose "github.com/devopsfaith/krakend-jose/mux"
	muxkrakend "github.com/devopsfaith/krakend-lua/router/mux"
	metrics "github.com/devopsfaith/krakend-metrics/mux"
	opencensus "github.com/devopsfaith/krakend-opencensus/router/mux"
	muxlua "github.com/luraproject/lura/router/mux"

	//juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/luraproject/lura/logging"
	router "github.com/luraproject/lura/router/mux"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory, pe muxlua.ParamExtractor) router.HandlerFactory {
	//handlerFactory := juju.HandlerFactory
	handlerFactory := muxlua.EndpointHandler
	handlerFactory = muxkrakend.HandlerFactory(logger, handlerFactory, pe)
	handlerFactory = muxjose.HandlerFactory(handlerFactory, pe, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	//handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}

type handlerFactory struct{}

func (h handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory, pe muxlua.ParamExtractor) router.HandlerFactory {
	return NewHandlerFactory(l, m, r, pe)
}
