package client

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/XxRoloxX/dns/internal/message"
	"github.com/XxRoloxX/dns/internal/record"
)

type Client struct {
	rootNameserver net.IP
}

func NewClient() *Client {
	return &Client{
		// rootNameserver: []byte("127.0.0.1:53"),
		rootNameserver: []byte("198.41.0.4:53"),
	}
}

func (c *Client) QueryATypeRecords(name string) (*message.Message, error) {

	rootNameserver, err := net.ResolveUDPAddr("udp4", string(c.rootNameserver))
	if err != nil {
		slog.Error("Failed to resolve root nameserver", "err", err)
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, rootNameserver)
	if err != nil {
		slog.Error("Failed to connect to root nameserver", "addr", c.rootNameserver)
		return nil, err
	}

	defer conn.Close()

	msg := message.Message{
		Header: message.Header{
			TransactionId:     2137,
			NumberOfQuestions: 1,
			Flags: message.HeaderFlags{
				Query:            true,
				RecursionDesired: true,
			},
		},
		Body: message.MessageBody{
			Queries: []message.Query{
				{
					Name:                strings.Split(name, "."),
					ResourceRecordType:  record.ResourceRecordType__A,
					ResourceRecordClass: record.ResourceRecordClass__In,
				},
			},
		},
	}

	encodedMessage := message.NewEncoder().Encode(&msg)

	n, err := conn.Write(encodedMessage)

	if err != nil {
		slog.Error("Failed to send request to root nameserver", "err", err)
		return nil, err
	}

	slog.Info("Sent message!", "bytes", n)

	buffer := make([]byte, 512)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if err != io.EOF {
				slog.Error("Failed to read request from root nameserver", "err", err)
				return nil, err
			}

			break
		}

		if n < len(buffer) {
			break
		}

		slog.Info(fmt.Sprintf("Read %s bytes from dns server", n))
	}

	var decodedMessage message.Message

	err = message.NewDecoder(buffer).Decode(&decodedMessage)
	if err != nil {
		slog.Error("Failed to decode message from root nameserver", "err", err)
		return nil, err
	}

	slog.Info("Message", "msg", decodedMessage)

	return &decodedMessage, nil
}
