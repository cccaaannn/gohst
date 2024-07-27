package server

import (
	"fmt"
	"net"
)

func (server *server) ListenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening ", err.Error())
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting ", err.Error())
			return
		}
		go server.handleConnection(conn)
	}
}
