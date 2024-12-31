package server

import (
	"fmt"
	message "github.com/XxRoloxX/dns/pkg/dns_message"
	managementserver "github.com/XxRoloxX/dns/pkg/management_server"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	conn       *net.UDPConn
	repository managementserver.RecordsRepository
}

func NewServer() *Server {

	address, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		panic(fmt.Sprintf("failed to resolve address: %s", err.Error()))
	}

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(fmt.Sprintf("failed to listen on socket: %s", err.Error()))
	}

	slog.Info("Listening for DNS messages on :53")

	repository := managementserver.NewPostgresRecordsRepository()

	return &Server{
		conn:       conn,
		repository: repository,
	}
}

func (s *Server) HandleRequest(req *Request) {

	for _, query := range req.msg.Body.Queries {
		records, err := s.repository.GetRecordsByName(strings.Join(query.Name, "."))
		if err != nil {
			slog.Error("Failed to get records for", "query", query)
			s.HandleInternalError(req)
			return
		}

		if len(records) == 0 {
			s.HandleNoResourceError(req)
			return
		}

		for _, record := range records {
			rr, err := record.ConvertToResourceRecord()
			if err != nil {
				slog.Error("Failed to convert Managed Resource Record to canonical form", "err", err)
				s.HandleInternalError(req)
				return
			}
			req.msg.AddAnswer(rr)
		}
	}

	req.msg.SetAsResponse()
	req.msg.UpdateRRNumbers()

	req.Send()
}

func (s *Server) HandleInternalError(req *Request) {
	req.msg.SetAsResponse()
	req.msg.SetResponseCode(message.ResponseCode__ServFail)
	req.Send()
}

func (s *Server) HandleFormattingError(req *Request) {
	req.msg.SetAsResponse()
	req.msg.SetResponseCode(message.ResponseCode__FormErr)
	req.Send()
}

func (s *Server) HandleNoResourceError(req *Request) {
	req.msg.SetAsResponse()
	req.msg.SetResponseCode(message.ResponseCode__NxDomain)
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
		if err != nil {
			s.HandleFormattingError(&Request{
				msg:  &message.Message{},
				conn: s.conn,
				addr: addr,
			})
			continue
		}

		go s.HandleRequest(req)

	}
}

func (s *Server) Close() {

	err := s.conn.Close()
	if err != nil {
		slog.Error("failed to close server", "err", err.Error())
		panic("failed to close connection to server")
	}
}
