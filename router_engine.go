package krakend

import (
	"io"

	//botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	//httpsecure "github.com/devopsfaith/krakend-httpsecure/mux"
	//lua "github.com/devopsfaith/krakend-lua/router/mux"
	//"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	muxlua "github.com/luraproject/lura/router/mux"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *muxlua.engine {
	// if !cfg.Debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	e := muxlua.DefaultEngine()
	// engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: w}), gin.Recovery())

	// engine.RedirectTrailingSlash = true
	// engine.RedirectFixedPath = true
	// engine.HandleMethodNotAllowed = true

	// if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
	// 	logger.Warning(err)
	// }

	// lua.Register(logger, cfg.ExtraConfig, engine)

	// botdetector.Register(cfg, logger, engine)

	return e
}

type engineFactory struct{}

func (e engineFactory) NewEngine(cfg config.ServiceConfig, l logging.Logger, w io.Writer) *muxlua.engine {
	return NewEngine(cfg, l, w)
}
