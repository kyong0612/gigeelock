package web

import (
	"net/http"
)

func (api *API) RecorderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := api.getTemplate(ctx, "web/record").Execute(w, nil); err != nil {
		panic(err)
	}
}
