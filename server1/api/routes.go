package api

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", handleTest)

	return corsMiddleware(mux)
}
