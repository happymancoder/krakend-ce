package krakend

import (
	"io"

	"github.com/devopsfaith/krakend/router/httptreemux"
	orighttptreemux "github.com/dimfeld/httptreemux"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) httptreemux.Engine {
	return httptreemux.NewEngine(orighttptreemux.NewContextMux())
}

type engineFactory struct{}

func (e engineFactory) NewEngine(cfg config.ServiceConfig, l logging.Logger, w io.Writer) httptreemux.Engine {
	return NewEngine(cfg, l, w)
}
