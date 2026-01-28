package server

import (
	"net"
	"log"
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/handler"
)

// sidenote: this is what is considered a 'router' in most http frameworks. In simple terms it helps the http server keep track of routes and their associated handlers
type MyServerMux struct {
	routes map[string]handler.Handler
}

func NewServerMux() *MyServerMux {
	return &MyServerMux{
		routes: map[string]handler.Handler{},
	}
}

// sidenote: we are casting the function that is being pass to the HandleFunc method to the HandlerFunc apdater. It in turn, gives
// us a valid handler, which is useful since we can reuse our Handle method.
func (m *MyServerMux) HandleFunc(route string, f func(w response.CustomResponseWriter, r *request.Request) {
	m.Handle(route, handler.HanlderFunc(f))
}

func (m *MyServerMux) Handle(route string, handler handler.Handler) {
	if m.routes == nil {
		m.routes = map[string]handler.Handler{}
	}
	m.routes[route] = handler
}

func CustomListenAndServe(addr string, h handler.Handler) (
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("error", err)
	)

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Error in trying to accept connection", err)
		}

		// here, we are first parsing the request so that we can send the information stored in thr request strcut into the server http method
		// sidenote: request is a pointer. It's easy to forget and just felt like I should add this in too avoid any headaches. 
		request, err := request.RequestFromReader(conn)


	}
}
