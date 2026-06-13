package main

import "strings"

func handleCommand(msg string) (response string, shouldClose bool) {
	upper := strings.ToUpper(msg)

	switch {
	case upper == "PING":
		return "PONG\n", false
	case upper == "QUIT":
		return "BYE\n", true
	case strings.HasPrefix(upper, "ECHO "):
		return msg[5:] + "\n", false
	case strings.HasPrefix(upper, "SET "):
		parts := strings.SplitN(msg[4:], " ", 2)
		if len(parts) != 2 {
			return "ERR usage: SET key value\n", false
		}
		storeMu.Lock()
		store[parts[0]] = parts[1]
		storeMu.Unlock()
		return "OK\n", false
	case strings.HasPrefix(upper, "GET "):
		key := strings.TrimSpace(msg[4:])
		storeMu.RLock()
		value, ok := store[key]
		storeMu.RUnlock()
		if !ok {
			return "ERR key not found\n", false
		}
		return value + "\n", false
	default:
		return "ERR unknown command\n", false
	}
}
