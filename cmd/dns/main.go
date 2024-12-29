package main

import (
	"github.com/XxRoloxX/dns/internal/message"
	"github.com/XxRoloxX/dns/internal/server"
)

func main() {

	srv := server.NewServer()

	srv.Listen(make(chan message.Message))
}
