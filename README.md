# TicTacToeTui Online

An infinite board online multiplayer tic tac toe game directly in your terminal.

# Client Game Setup

## Executable

(Coming soon)

## Manual 

Build (go 1.23.3 required):

```bash
cd cmd/client
go build -o client_game
```

Run:

```bash
./client_game
```

Or Run directly:

```bash
go run cmd/client/main.go
```

By default, the client will connect to the server at `tictactoetuionline.onrender.com`.

# Custom Server Setup

## Manual

Build (go 1.23.3 required):

```bash
cd cmd/server
go build -o server
```

Run:

```bash
./server
```

Or Run directly:

```bash
go run cmd/server/main.go
```

## Update Server Address

Edit the `serverAddr` constant in `cmd/client/main.go` without the protocol prefix.

Example:

```go
//...

const (
	serverAddr = "localhost:8080"
)

//...
```