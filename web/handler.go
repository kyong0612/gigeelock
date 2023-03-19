package web

import (
	"io"
	"io/fs"
	"net/http"
	"sync/atomic"

	"github.com/kyong0612/gigeelock/html"
)

type API struct {
	templateFS fs.FS
	templates  *atomic.Value
}

func NewAPI() *API {
	api := &API{
		templateFS: html.Embedded,
		templates:  new(atomic.Value),
	}

	api.templates.Store(api.parseTemplates())

	return api
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// simple health check
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}
