package main

import (
	"net/http"
	"time"
)

func main() {
	srv := &http.Server{
		Addr:        ":8080",
		Handler:     handler,
		ReadTimeout: 10 * time.SSecond,
	}
}
