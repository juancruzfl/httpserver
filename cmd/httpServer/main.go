package main

import (
	"fmt"
	"io"
	"github.com/juancruzfl/httpserver/internal/server"
	"github.com/juancruzfl/httpserver/internal/response"
	"github.com/juancruzfl/httpserver/internal/request"
)

func main() {
	errChan := make(chan error, 1)
	go func () {
		server.MyDefaultMux.HandleFunc("/", func(w response.ResponseWriter, r *request.Request) {
			fmt.Printf("got / request\n")
			// Write the response body.
			io.WriteString(w, "Hello, World!\n")
		})
		fmt.Printf("Server started running")
		errChan <- server.CustomListenAndServe(":8000", nil)
	}()
	err := <- errChan
	fmt.Printf("Server stopped: ", err)

}
