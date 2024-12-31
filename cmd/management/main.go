package main

import managementserver "github.com/XxRoloxX/dns/internal/management_server"

func main() {

	server := managementserver.NewServer()
	server.Start()
}
