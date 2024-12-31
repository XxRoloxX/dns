package client

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"strings"

	message "github.com/XxRoloxX/dns/pkg/dns_message"
	record "github.com/XxRoloxX/dns/pkg/dns_record"
)

const (
	MESSAGE_BUFFER_SIZE = 512
	MESSAGE_MAX_SIZE    = 5048
)

type Client struct {
	rootNameserver net.UDPAddr
	header         message.Header

	// Each connection is identified by messages TransactionId
	conn map[uint16]*net.UDPConn
}

func NewClient(rootNameserver net.IP) (*Client, error) {

	rootNameserverAddr, err := net.ResolveUDPAddr("udp4", string(rootNameserver))
	if err != nil {
		slog.Error("Failed to resolve root nameserver", "err", err)
		return nil, err
	}

	return &Client{
		rootNameserver: *rootNameserverAddr,
		header: message.Header{
			Flags: message.HeaderFlags{
				Query:            true,
				RecursionDesired: true,
				OperationCode:    message.OpCode__Query,
			},
			NumberOfQuestions:    0,
			NumberOfAnswers:      0,
			NumberOfAuthorityRR:  0,
			NumberOfAdditionalRR: 0,
		},
	}, nil
}

func (c *Client) Query(queries []message.Query) (*message.Message, error) {
	msg := message.Message{
		Header: c.header,
		Body:   message.MessageBody{},
	}

	for _, query := range queries {
		msg.AddQuery(query)
	}

	id := c.NewTransactionId()
	msg.Header.TransactionId = id
	msg.UpdateRRNumbers()

	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c.conn[id] = conn

	response, err := c.sendAndAwaitResponse(&msg)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) createConnection() (*net.UDPConn, error) {

	conn, err := net.DialUDP("udp4", nil, &c.rootNameserver)
	if err != nil {
		slog.Error("Failed to connect to root nameserver", "addr", c.rootNameserver)
		return nil, err
	}

	return conn, nil
}

func (c *Client) getConnectionFromMessage(msg *message.Message) *net.UDPConn {
	return c.conn[msg.Header.TransactionId]
}

func (c *Client) readFromConnection(conn *net.UDPConn) ([]byte, error) {

	fullResponse := make([]byte, MESSAGE_MAX_SIZE)

	buffer := make([]byte, MESSAGE_BUFFER_SIZE)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if err != io.EOF {
				slog.Error("Failed to read request from root nameserver", "err", err)
				return nil, err
			}

			break
		}

		fullResponse = append(fullResponse, buffer[:n]...)

		if n < len(buffer) {
			break
		}

		slog.Info(fmt.Sprintf("Read %d bytes from the dns server", n))
	}

	return fullResponse, nil
}

func (c *Client) NewTransactionId() uint16 {
	return uint16(rand.Intn(1 << 16))
}

func (c *Client) sendAndAwaitResponse(msg *message.Message) (*message.Message, error) {

	encodedMessage := message.NewEncoder().Encode(msg)
	conn := c.getConnectionFromMessage(msg)

	n, err := conn.Write(encodedMessage)
	if err != nil {
		slog.Error("Failed to send request to root nameserver", "err", err)
		return nil, err
	}

	slog.Info("Message was sent", "bytes", n)

	response, err := c.readFromConnection(conn)
	if err != nil {
		slog.Error("Failed to read from connection", "err", err)
		return nil, err
	}

	var decodedMessage message.Message

	err = message.NewDecoder(response).Decode(&decodedMessage)
	if err != nil {
		slog.Error("Failed to decode message from root nameserver", "err", err)
		return nil, err
	}

	return &decodedMessage, nil
}

func (c *Client) QueryATypeRecords(name string) (*message.Message, error) {
	return c.Query([]message.Query{
		{
			Name:                strings.Split(name, "."),
			ResourceRecordClass: record.ResourceRecordClass__In,
			ResourceRecordType:  record.ResourceRecordType__A,
		},
	})
}
