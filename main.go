package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var store = map[string]string{}
var storeMu sync.RWMutex

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client connected:", conn.RemoteAddr())

	buffer := make([]byte, 1024)

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("client disconnected:", err)
			return
		}

		msg := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("received:", msg)

		if msg == "PING" {
			_, err = conn.Write([]byte("PONG\n"))
		} else if msg == "QUIT" {
			_, err = conn.Write([]byte("BYE\n"))
			return
		} else if strings.HasPrefix(msg, "ECHO ") {
			subStr := strings.TrimPrefix(msg, "ECHO ")
			_, err = conn.Write([]byte(subStr + "\n"))
		} else if strings.HasPrefix(msg, "SET ") {
			parts := strings.SplitN(msg, " ", 3)
			if len(parts) != 3 {
				_, err = conn.Write([]byte("ERR usage: SET key value\n"))
			} else {
				key := parts[1]
				value := parts[2]

				storeMu.Lock()
				store[key] = value
				storeMu.Unlock()

				_, err = conn.Write([]byte("OK\n"))
			}
		} else if strings.HasPrefix(msg, "GET ") {
			parts := strings.SplitN(msg, " ", 2)
			if len(parts) != 2 {
				_, err = conn.Write([]byte("ERR usage: GET key\n"))
			} else {
				key := parts[1]

				storeMu.RLock()
				value, ok := store[key]
				storeMu.RUnlock()

				if !ok {
					_, err = conn.Write([]byte("ERR key not found\n"))
				} else {
					_, err = conn.Write([]byte(value + "\n"))
				}
			}
		} else {
			_, err = conn.Write([]byte("ERR unknown command\n"))
		}

		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

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
		go handleConn(conn)
	}
}
