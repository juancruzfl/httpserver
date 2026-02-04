package handler

import (
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/response"
)

// as we can see here, a handlder is specified via this function signature. The response writer is what is being used to actually send a response back to the client. The request is the structure 
// that holds the already parsed information, which was done within our server. We are quite literally 'handling' a request.
type Handler interface {
	ServeHttp(response.ResponseWriter, *request.Request)
}

// sidenote: this is a function type that is meant to be the adpater to the handler (above) interface. We do this because we can't 
// attach interfaces to methods. We need some concrete type, so we use a function type to accomplish this.
type HandlerFunc func(response.ResponseWriter, *request.Request)

func (h HandlerFunc) ServeHttp(w response.ResponseWriter, r *request.Request) {
	h(w, r)
}

