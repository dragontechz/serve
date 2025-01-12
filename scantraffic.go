package main

import (
	"fmt"
	"log"
	"net"
)

var buff []byte = make([]byte, 1024)
func main(){
		sock, err := net.Listen("tcp", ":7777")
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	
		fmt.Println("server listening on port 7777")
	
		for {
			conn, err := sock.Accept()
			if err != nil {
				log.Println("ERROR: ", err)
			}
			log.Println("new connection established by:", conn.RemoteAddr().String())
			
			go display_traffic((conn))
			continue
	
		}

}
func display_traffic(conn net.Conn){
	for {
		n , err := conn.Read(buff)
		if err != nil {}

		data := string(buff[:n])

		log.Println("data: ",data)

	}
}