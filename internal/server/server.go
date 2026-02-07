package server

import (
	"net"
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/response"
	"github.com/juancruzfl/httpserver/internal/handler"
)

// sidenote: this is what is considered a 'router' in most http frameworks. In simple terms it helps the http server keep track of routes and their associated handlers
type routeKey struct {
	method string
	path string
}

type MyServerMux struct {
	routes map[routeKey]handler.Handler
}

func NewServerMux() *MyServerMux {
	return &MyServerMux{
		routes: map[routeKey]handler.Handler{},
	}
}

var MyDefaultMux = NewServerMux()

func (m *MyServerMux) Get(method string, path string) (handler.Handler, bool) {
	key := routeKey{method: method, path: path}
	handler, ok := m.routes[key]
	return handler, ok			
}

func (m *MyServerMux) ServeHttp(w response.ResponseWriter, r *request.Request) {
    key := routeKey{
        method: r.RequestLine.Method,
        path:   r.RequestLine.RequestTarget,
    }
    handler, ok := m.routes[key]
    if !ok {
        w.CustomWriteHeader(404)
        w.Write([]byte("404 Not Found"))
        return
    }
    handler.ServeHttp(w, r)
}

// sidenote: we are casting the function that is being pass to the HandlerFunc method of the server multiplexer to the HandlerFunc apdater of the handler interface. It in turn, gives
// us a valid handler, which is useful since we can reuse our Handle method.
func (m *MyServerMux) HandleFunc(method string, path string, f func(response.ResponseWriter, *request.Request)) {
    m.Handle(method, path, handler.HandlerFunc(f))
}

func (m *MyServerMux) Handle(method string, path string, h handler.Handler) {
    if m.routes == nil {
        m.routes = make(map[routeKey]handler.Handler)
    }
    key := routeKey{method: method, path: path}
    m.routes[key] = h
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

func CustomListenAndServe(addr string, h handler.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		// sidenote: I have decided to change how the request is read here. If we read the incoming requests and then wait for them to be parsed, we disallow multiple people from
		// connecting to our server since we are occupaying the main thread in our method to wait for the request operations to finsih. We instead use a go rountine in a serve function
		// to offshore that work into a background thread.
		go func(c net.Conn) {
			err := serve(c, h)
			if err != nil {
				println("Error in trying to serve connection", err.Error())
			}
		}(conn)
	}
}
