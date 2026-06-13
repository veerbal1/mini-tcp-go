package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	fmt.Println("server listening on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		fmt.Println("client connected:", conn.RemoteAddr())

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("read error:", err)
			conn.Close()
			continue
		}
		fmt.Println("received: ", string(buffer[:n]))

		conn.Write([]byte("Ok\n"))
		conn.Close()
	}
}
