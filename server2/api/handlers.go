package api

import "net/http"

func handleTest(w http.ResponseWriter, r *http.Request) {
	res := "Hello"
	w.Header().Add("content-Type", "text/html")
	w.Write([]byte(res))
}
