package http

import "net/http"

func StartServer(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
