# mini-tcp-go

A small Redis-like TCP server written in Go to understand how backend servers, TCP connections, custom protocols, concurrency, timeouts, and shared state work underneath normal HTTP frameworks.

This project was built from scratch for learning backend fundamentals, without hiding everything behind a web framework.

## What this project is

`mini-tcp-go` is a custom TCP application protocol server.

It does not rebuild TCP itself. TCP is provided by the operating system. This project builds a simple protocol on top of TCP, similar in spirit to how HTTP, Redis, SMTP, PostgreSQL, and other systems define their own rules over a TCP byte stream.

The server listens on port `9000`, accepts multiple clients, handles each client in its own goroutine, reads line-based commands, processes them, and sends protocol responses.

## Why I built this

Most backend frameworks hide the lower-level machinery:

```text
open port
  ↓
accept connection
  ↓
read bytes
  ↓
parse protocol
  ↓
route command/request
  ↓
run logic
  ↓
write response
  ↓
close / keep alive / timeout
```

This project helped me understand what actually happens underneath libraries like `net/http`, Express, Axum, Redis clients, database drivers, and other backend abstractions.

## Concepts learned

### TCP and ports

The server opens a TCP port:

```go
net.Listen("tcp", ":9000")
```

This means the operating system opens a TCP door on port `9000`, and clients can connect to it.

### Listener vs connection

```text
listener = server door
conn     = one private TCP pipe with one client
```

The server keeps accepting connections:

```go
conn, err := listener.Accept()
```

Each accepted connection represents one client conversation.

### Goroutine per client

Each client is handled independently:

```go
go handleConn(conn)
```

This keeps the main accept loop free.

```text
main goroutine   = keeps accepting clients
client goroutine = handles one connected client
```

Without this, one slow or silent client could block other clients.

### TCP is a byte stream, not messages

TCP does not understand commands like `PING`, `GET`, or `SET`.

It only transports bytes.

So the protocol must define message boundaries. This project uses newline-based framing:

```text
one command = one line ending with \n
```

The server reads one full command line using:

```go
reader.ReadString('\n')
```

### Custom protocol layer

The protocol currently supports:

```text
PING              -> PONG
ECHO hello        -> hello
SET name Veerbal  -> OK
GET name          -> Veerbal
QUIT              -> BYE + close connection
unknown           -> ERR unknown command
```

This is the core idea of protocols:

```text
TCP transports bytes.
Protocol gives meaning to those bytes.
Application logic acts on that meaning.
```

### Timeout for silent clients

The server does not wait forever for silent clients.

```go
conn.SetReadDeadline(time.Now().Add(10 * time.Second))
```

If a client connects but sends nothing for 10 seconds, the server closes that connection.

This prevents goroutines from waiting forever on dead or idle clients.

### Shared state and locking

The server stores key-value data in memory:

```go
var store = map[string]string{}
```

Since many clients run in different goroutines, shared state must be protected:

```go
var storeMu sync.RWMutex
```

`SET` uses a write lock:

```go
storeMu.Lock()
store[key] = value
storeMu.Unlock()
```

`GET` uses a read lock:

```go
storeMu.RLock()
value, ok := store[key]
storeMu.RUnlock()
```

This prevents unsafe concurrent map access.

## How to run

Start the server:

```bash
go run main.go
```

Expected output:

```text
server listening on :9000
```

In another terminal, connect using `nc`:

```bash
nc localhost 9000
```

Now send commands.

Ping the server:

```text
PING
```

Response:

```text
PONG
```

Set a value:

```text
SET name Veerbal
```

Response:

```text
OK
```

Get a value:

```text
GET name
```

Response:

```text
Veerbal
```

Echo text:

```text
ECHO hello bro
```

Response:

```text
hello bro
```

Close the connection:

```text
QUIT
```

Response:

```text
BYE
```

## Example session

```text
PING
PONG

SET name Veerbal
OK

GET name
Veerbal

ECHO learning tcp in go
learning tcp in go

QUIT
BYE
```

## Current architecture

```text
main()
  ↓
net.Listen("tcp", ":9000")
  ↓
accept loop
  ↓
go handleConn(conn)
  ↓
read one newline-ended command
  ↓
handleCommand(msg)
  ↓
write response
  ↓
continue / close connection
```

## Code structure

### `main`

Starts the TCP server, listens on port `9000`, and accepts clients forever.

### `handleConn`

Handles one client connection.

Responsibilities:

```text
set read timeout
read command line
trim input
call handleCommand
write response
close connection when needed
```

### `handleCommand`

Handles protocol logic.

Responsibilities:

```text
parse command
execute command
read/write store
return response
tell connection whether to close
```

This keeps networking logic separate from protocol logic.

## What this taught me about real backend systems

This project made backend internals less mysterious.

Normal HTTP servers do the same kind of work, but with a standardized HTTP protocol:

```text
TCP connection
  ↓
HTTP parser
  ↓
method/path/headers/body
  ↓
router
  ↓
handler
  ↓
HTTP response
```

In this project, I manually built a smaller version:

```text
TCP connection
  ↓
line reader
  ↓
custom command parser
  ↓
command handler
  ↓
custom response
```

This helped me understand:

```text
TCP transport
ports
connections
client lifecycle
goroutines
timeouts
message framing
protocol parsing
shared state
mutexes
server architecture
```

## What this is not

This is not a production database.

It does not support:

```text
persistence
authentication
replication
clustering
binary protocol
advanced error types
memory limits
TLS
observability
graceful shutdown
```

The goal is learning backend fundamentals.

## Possible next improvements

```text
DEL key
EXISTS key
TTL / EXPIRE
case-insensitive commands
structured protocol responses
Go client
unit tests for handleCommand
graceful shutdown
basic metrics
Dockerfile
```

## Learning takeaway

The biggest lesson:

```text
Frameworks are not magic.

They are layers over:
TCP → protocol parsing → routing/command handling → response writing.
```

This project shows the lower layer directly.

## Repository description

```text
mini-tcp-go — A Redis-like TCP server in Go for learning backend internals, custom protocols, goroutines, timeouts, and shared state.
```
