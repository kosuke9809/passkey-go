package router

import "net/http"

type Router struct {
	http.ServeMux
}
