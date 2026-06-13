package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("client disconnected:", err)
			return
		}

		msg := strings.TrimSpace(line)
		fmt.Println("received:", msg)

		response, shouldClose := handleCommand(msg)

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
		if shouldClose {
			return
		}
	}
}
