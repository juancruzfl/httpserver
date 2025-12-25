package main

import (
	"fmt"
	"log"
	"net"
	"io"
)
		
func getLinesChannel(f io.ReadCloser) <- chan string {
	channel := make(chan string, 1)
	
	go func(){
	
		defer f.Close()
		defer close(channel)

		var accumulator []byte

		for {

			buffer := make([]byte, 8)

			count, error := f.Read(buffer)

			if error != nil {
				break
			}

			for i := 0; i < count; i++ {
				if buffer[i] == '\n'{
					channel <- string(accumulator)
					accumulator = nil				
				} else {
					accumulator = append(accumulator, buffer[i])
				}
			}
		}
	}()

	return channel 
}	

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
	
		fmt.Print("Connection has been accepted")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println("Read in line: ", line)
		}

	}

}
