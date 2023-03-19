package web

import "net/http"

func RecorderHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
