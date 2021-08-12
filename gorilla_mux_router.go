/* Package gorilla provides some basic implementations for building routers based on gorilla/mux
 */
// SPDX-License-Identifier: Apache-2.0
package krakend

import (
	"net/http"
	"strings"
	"fmt"
	"os"

	gorilla "github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
	"github.com/luraproject/lura/router"
	"github.com/luraproject/lura/router/mux"
)

// DefaultFactory returns a net/http mux router factory with the injected proxy factory and logger
func DefaultFactory(pf proxy.Factory, logger logging.Logger) router.Factory {
	return mux.NewFactory(DefaultConfig(pf, logger))
}

// DefaultConfig returns the struct that collects the parts the router should be builded from
func DefaultConfig(pf proxy.Factory, logger logging.Logger) mux.Config {
	return mux.Config{
		Engine:         GorillaEngine{gorilla.NewRouter()},
		Middlewares:    []mux.HandlerMiddleware{},
		HandlerFactory: mux.CustomEndpointHandler(mux.NewRequestBuilder(GorillaParamsExtractor)),
		ProxyFactory:   pf,
		Logger:         logger,
		DebugPattern:   "/__debug/{params}",
		RunServer:      router.RunServer,
	}
}

func GorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	fmt.Println("Gorilla ParamsExtractor is called...........")
	for key, value := range gorilla.Vars(r) {
		params[strings.Title(key)] = value
		fmt.Println(strings.Title(key))
	}
	return params
}

func GorillaNewEngine(r *gorilla.Router) GorillaEngine {
	fmt.Println("Creating Gorilla Engine..............")
	r.Use(loggingMiddleware)
	return GorillaEngine{r}
}

type GorillaEngine struct {
	r *gorilla.Router
}

// Handle implements the mux.Engine interface from the lura router package
func (g GorillaEngine) Handle(pattern, method string, handler http.Handler) {
	g.r.Handle(pattern, handler).Methods(method)
}

// ServeHTTP implements the http:Handler interface from the stdlib
func (g GorillaEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	g.r.ServeHTTP(mux.NewHTTPErrorInterceptor(w), r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
