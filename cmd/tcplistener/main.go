package main

import (
	"fmt"
	"log"
	"net"
	"github.com/juancruzfl/httpserver/internal/request"
)	

func main(){ 
	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal("error", err)
	}

	defer listener.Close()

	fmt.Println("Listening on port 8000")

	for {
		conn, error := listener.Accept()
		
		if error != nil {
			fmt.Println("Error trying to accept connection", error)
			continue
		}
	
		fmt.Print("Connection has been accepted\n")
		
		req, err := request.RequestFromReader(conn)
		requestLine := req.RequestLine
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Printf("Request Line: \n")
		fmt.Printf("- Method: %s\n", requestLine.Method)
		fmt.Printf("- Target: %s\n", requestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", requestLine.HttpVersion)
	}
}
