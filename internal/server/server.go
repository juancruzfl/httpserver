package server

import (
	"log"
	"net"
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/response"
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

var MyDefaultMux = NewServerMux()

func (m *MyServerMux) Get(route string) (handler.Handler, bool) {
	handler, ok := m.routes[route]
	return handler, ok			
}

func (m *MyServerMux) ServeHttp(w response.ResponseWriter, r *request.Request) {
	handler, ok := m.Get(r.RequestLine.RequestTarget)
	if !ok {
		w.Write([]byte("404 Not Found"))
	} else {
		handler.ServeHttp(w, r)
	}
}

// sidenote: we are casting the function that is being pass to the HandlerFunc method of the server multiplexer to the HandlerFunc apdater of the handler interface. It in turn, gives
// us a valid handler, which is useful since we can reuse our Handle method.
func (m *MyServerMux) HandlerFunc(route string, f func(w response.ResponseWriter, r *request.Request)) {
	m.Handle(route, handler.HandlerFunc(f))
}

func (m *MyServerMux) Handle(route string, h handler.Handler) {
	if m.routes == nil {
		m.routes = map[string]handler.Handler{}
	}
	m.routes[route] = h
}

func serve(conn net.Conn, h handler.Handler) error {
	defer conn.Close()

	request, err := request.RequestFromReader(conn)
	if err != nil {
		return err
	}
	if h == nil {
		h = MyDefaultMux
	}
	writer := response.NewResponseWriter(conn)
	h.ServeHttp(writer, request)
	return nil 
}

func CustomListenAndServe(addr string, h handler.Handler) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// sidenote: I have decided to change how the request is read here. If we read the incoming requests and then wait for them to be parsed, we disallow multiple people from
		// connecting to our server since we are occupaying the main thread in our method to wait for the request operations to finsih. We instead use a go rountine in a serve function
		// to offshore that work into a background thread.
		go serve(conn, h)
	}
}
