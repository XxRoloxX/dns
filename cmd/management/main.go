package main

import managementserver "github.com/XxRoloxX/dns/pkg/management_server"

func main() {

	server := managementserver.NewServer()
	server.Start()
}
