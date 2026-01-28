package handlder

import (
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/response"
)

type Handler interface {
	ServeHttp(response.CustomResponseWriter, *request.Request)
}

type HandlerFunc func(response.CustomResponseWriter, *request.Request)

func (h HandlerFunc) ServerHttp(w response.CustomResponseWriter, r *request.Request) {
	h(w, r)
}

