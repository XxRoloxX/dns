package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/XxRoloxX/dns/internal/message"
)

type Server struct {
	conn *net.UDPConn
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
		conn: conn,
	}
}

func (s *Server) HandleMessage(rawMessage []byte) {

	var decodedMessage message.Message
	err := message.NewDecoder(rawMessage).Decode(&decodedMessage)
	if err != nil {
		slog.Error("failed to decode message ",
			"err", err.Error(), "msg", fmt.Sprintf("%+v", rawMessage),
		)
		return
	}

	slog.Info("Got message", "msg", fmt.Sprintf("%+v", decodedMessage))
}

func (s *Server) Listen(chan message.Message) {

	for {
		buf := make([]byte, 512)
		n, _, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("failed to read message", "err", err.Error())
			continue
		}
		rawMessage := buf[:n]
		go s.HandleMessage(rawMessage)
	}
}

func (s *Server) Close() {

	err := s.conn.Close()
	if err != nil {
		slog.Error("failed to close server", "err", err.Error())
		println("Failed to close connection to server")
	}
}
