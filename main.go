package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var store = map[string]string{}
var storeMu sync.RWMutex

func handleCommand(msg string) (response string, shouldClose bool) {
	if msg == "PING" {
		return "PONG\n", false
	}
	if msg == "QUIT" {
		return "BYE\n", true
	}
	if strings.HasPrefix(msg, "ECHO ") {
		subStr := strings.TrimPrefix(msg, "ECHO ")
		return subStr + "\n", false
	}
	if strings.HasPrefix(msg, "SET ") {
		parts := strings.SplitN(msg, " ", 3)
		if len(parts) != 3 {
			return "ERR usage: SET key value\n", false
		}
		key := parts[1]
		value := parts[2]

		storeMu.Lock()
		store[key] = value
		storeMu.Unlock()

		return "OK\n", false
	}
	if strings.HasPrefix(msg, "GET ") {
		parts := strings.SplitN(msg, " ", 2)
		if len(parts) != 2 {
			return "ERR usage: GET key\n", false
		}
		key := parts[1]

		storeMu.RLock()
		value, ok := store[key]
		storeMu.RUnlock()

		if !ok {
			return "ERR key not found\n", false
		}
		return value + "\n", false
	}
	return "ERR unknown command\n", false
}

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
