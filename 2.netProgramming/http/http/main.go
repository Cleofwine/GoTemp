package main

import (
	"http-case/http/handler"
	"net/http"
)

func main() {
	multiplexer := http.NewServeMux()
	multiplexer.Handle("/static/", handler.StaticFileServer())
	multiplexer.HandleFunc("/sum", handler.SumHandler)
	multiplexer.HandleFunc("/upload", handler.UploadHandler)
	http.ListenAndServe(":8899", multiplexer) // http 1.1
}
