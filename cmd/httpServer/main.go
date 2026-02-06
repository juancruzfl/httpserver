package main

import (
	"fmt"
	"github.com/juancruzfl/httpserver/internal/server"
	"github.com/juancruzfl/httpserver/internal/response"
	"github.com/juancruzfl/httpserver/internal/request"
)

func main() {
	errChan := make(chan error, 1)
	go func () {
		server.MyDefaultMux.HandleFunc("/upload", func(w response.ResponseWriter, r *request.Request) {
			fmt.Printf("Parsed Body: %s\n", string(r.Body))
			
			w.CustomWriteHeader(201)
			w.Write([]byte("I received your data: "))
			w.Write(r.Body)
		})
		server.MyDefaultMux.HandleFunc("/", func(w response.ResponseWriter, r *request.Request) {
			fmt.Printf("got / request\n")
			// Write the response body.
			w.CustomWriteHeader(201)
			w.Write([]byte("Hello World!\n"))
		})
		fmt.Printf("Server started running")
		errChan <- server.CustomListenAndServe(":8000", nil)
	}()
	err := <- errChan
	fmt.Printf("Server stopped: ", err)

}
