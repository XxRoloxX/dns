package main

import (
	"github.com/XxRoloxX/dns/pkg/dns_message"
	"github.com/XxRoloxX/dns/pkg/dns_server"
)

func main() {

	srv := server.NewServer()

	srv.Listen(make(chan message.Message))
}
