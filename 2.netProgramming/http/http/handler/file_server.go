package handler

import "net/http"

func StaticFileServer() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}
