package main

import (
	// "github.com/XxRoloxX/dns/internal/client"
	"github.com/XxRoloxX/dns/internal/message"
	"github.com/XxRoloxX/dns/internal/server"
	// "log/slog"
)

func main() {

	srv := server.NewServer()

	srv.Listen(make(chan message.Message))

	// client := client.NewClient()
	//
	// _, err := client.QueryATypeRecords("wmsdev.pl")
	// if err != nil {
	// 	slog.Error("Failed to query A type record", "err", err)
	// 	return
	// }

}
