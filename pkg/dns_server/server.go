package server

import (
	"fmt"
	"log/slog"
	"net"

	message "github.com/XxRoloxX/dns/pkg/dns_message"
	record "github.com/XxRoloxX/dns/pkg/dns_record"
)

type Server struct {
	conn  *net.UDPConn
	store *record.RRStore
}

func NewServer() *Server {

	address, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve address: %s", err.Error()))
	}

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on socket: %s", err.Error()))
	}

	slog.Info("starting to listen on :53")

	return &Server{
		store: record.NewRRStore(),
		conn:  conn,
	}
}

func (s *Server) HandleRequest(req *Request) {

	for _, query := range req.msg.Body.Queries {
		records := s.store.ResourceRecords(query.Name)
		for _, record := range records {
			req.msg.AddAnswer(record)
		}
	}

	req.msg.SetAsResponse()
	req.msg.UpdateRRNumbers()

	req.Send()
}

func (s *Server) Listen(chan message.Message) {

	for {
		buf := make([]byte, 512)
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("failed to read message", "err", err.Error())
			continue
		}

		req, err := NewRequest(s.conn, addr, buf[:n])

		go s.HandleRequest(req)

	}
}

func (s *Server) Close() {

	err := s.conn.Close()
	if err != nil {
		slog.Error("failed to close server", "err", err.Error())
		println("Failed to close connection to server")
	}
}
