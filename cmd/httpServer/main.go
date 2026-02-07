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
		server.MyDefaultMux.HandleFunc("GET", "/", func(w response.ResponseWriter, r *request.Request) {
			fmt.Printf("Handled GET / request\n")		
			w.Write([]byte("Hello, World!\n"))
		})
		server.MyDefaultMux.HandleFunc("POST", "/upload", func(w response.ResponseWriter, r *request.Request) {
			w.GetHeaders().Set("Content-Type", "application/json")
			w.CustomWriteHeader(201)
			w.Write([]byte(`{"status":"success"}`))
		})
		fmt.Printf("Server started running")
		errChan <- server.CustomListenAndServe(":8000", nil)
	}()
	err := <- errChan
	fmt.Printf("Server stopped: ", err)
}
