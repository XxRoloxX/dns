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
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("failed to read message", "err", err.Error())
			continue
		}
		rawMessage := buf[:n]

		var msg message.Message

		err = message.NewDecoder(rawMessage).Decode(&msg)
		if err != nil {
			slog.Error("Failed to decode message", "msg", rawMessage)
			continue
		}

		msg.Header.Flags.Query = false
		msg.Header.NumberOfAnswers = 1
		msg.Header.NumberOfAuthorityRR = 1
		msg.Answers = []message.Answer{
			{
				Name:                msg.Queries[0].Name,
				ResourceRecordClass: msg.Queries[0].ResourceRecordClass,
				ResourceRecordType:  msg.Queries[0].ResourceRecordType,
				Ttl:                 60,
				RDataLength:         4,
				RData:               net.IPv4(192, 168, 1, 1).To4(),
			},
		}

		encodedMessage := message.NewEncoder().Encode(&msg)

		_, err = s.conn.WriteToUDP(encodedMessage, addr)
		if err != nil {
			slog.Error("Failed to send message", "msg", msg)
			continue
		}

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
