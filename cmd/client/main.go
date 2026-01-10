package main

import (
	"fmt"
	"net"
	"log"
)

func main() {
	conn, error := net.Dial("tcp", "localhost:8000")
	
	if error != nil {
		log.Fatal("Error: ", error)
		return
	}

	defer conn.Close()

	buffer := []byte("Hello from our client \n")
	
	n, err := conn.Write(buffer)

	if err != nil {
		log.Fatal("Error in trying to write message to buffer: ", err)
	}

	fmt.Println(n)
}
