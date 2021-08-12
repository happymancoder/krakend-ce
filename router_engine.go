package krakend

import (
	"io"

	gorilla "github.com/gorilla/mux"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) GorillaEngine {
	return GorillaNewEngine(gorilla.NewRouter())
}
//.StrictSlash(true)

type engineFactory struct{}

func (e engineFactory) NewEngine(cfg config.ServiceConfig, l logging.Logger, w io.Writer) GorillaEngine {
	return NewEngine(cfg, l, w)
}
