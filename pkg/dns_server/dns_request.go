package server

import (
	"github.com/XxRoloxX/dns/pkg/dns_message"
	"log/slog"
	"net"
)

type Request struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	msg  *message.Message
}

func NewRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) (*Request, error) {

	var msg message.Message

	err := message.NewDecoder(buf).Decode(&msg)
	if err != nil {
		slog.Error("Failed to decode message", "msg", buf)
		return nil, err
	}

	slog.Info("Got message", "msg", msg)

	return &Request{
		conn: conn,
		addr: addr,
		msg:  &msg,
	}, nil
}

func (r *Request) Send() error {

	encodedMessage := message.NewEncoder().Encode(r.msg)

	_, err := r.conn.WriteToUDP(encodedMessage, r.addr)
	if err != nil {
		slog.Error("Failed to send message", "msg", r.msg)
		return err
	}

	slog.Info("Response", "msg", r.msg)

	return nil
}
