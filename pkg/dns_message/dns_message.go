package message

import (
	"errors"
	"fmt"

	// message "github.com/XxRoloxX/dns/pkg/dns_message"
	"github.com/XxRoloxX/dns/pkg/dns_record"
)

type OpCode uint

const (
	OpCode__Query  OpCode = 0
	OpCode__IQuery OpCode = 1
	OpCode__Status OpCode = 2
)

func NewOperationCode(code uint16) (OpCode, error) {
	switch code {
	case 0:
		return OpCode__Query, nil
	case 1:
		return OpCode__IQuery, nil
	case 2:
		return OpCode__Status, nil
	default:
		return 0, errors.New(fmt.Sprintf("Invalid operation code: %d", code))
	}
}

type ResponseCode uint

const (
	ResponseCode__NoError  = 0
	ResponseCode__FormErr  = 1
	ResponseCode__ServFail = 2
	ResponseCode__NxDomain = 3
)

func NewResponseCode(code uint16) (ResponseCode, error) {
	switch code {
	case 0:
		return ResponseCode__NoError, nil
	case 1:
		return ResponseCode__FormErr, nil
	case 2:
		return ResponseCode__ServFail, nil
	case 3:
		return ResponseCode__NxDomain, nil
	default:
		return 0, errors.New("Invalid response code")
	}
}

type HeaderFlags struct {
	Query              bool
	OperationCode      OpCode
	AuthorativeAnswer  bool
	Truncation         bool
	RecursionDesired   bool
	RecursionAvailable bool
	ResponseCode       ResponseCode
}

type Header struct {
	TransactionId        uint16
	Flags                HeaderFlags
	NumberOfQuestions    uint16
	NumberOfAnswers      uint16
	NumberOfAuthorityRR  uint16
	NumberOfAdditionalRR uint16
}

type Query struct {
	Name                []string
	ResourceRecordType  record.ResourceRecordType
	ResourceRecordClass record.ResourceRecordClass
}

type Answer struct {
	Name                []string
	ResourceRecordType  record.ResourceRecordType
	ResourceRecordClass record.ResourceRecordClass
	Ttl                 uint32
	RDataLength         uint16
	RData               []byte
}

type MessageBody struct {
	Queries     []Query
	Answers     []Answer
	Authorative []Answer
	Additional  []Answer
}

type Message struct {
	Header Header
	Body   MessageBody
}

func (m *Message) SetAsResponse() {
	m.Header.Flags.Query = false
}

// Update headers with numbers of Answers, Authorative and Additional RRs
func (m *Message) UpdateRRNumbers() {
	m.Header.NumberOfAnswers = uint16(len(m.Body.Answers))
	m.Header.NumberOfAuthorityRR = uint16(len(m.Body.Authorative))
	m.Header.NumberOfAdditionalRR = uint16(len(m.Body.Additional))
}

func (m *Message) AddAnswer(rr record.ResourceRecord) {

	m.Body.Answers = append(m.Body.Answers, Answer{
		Name:                rr.Name(),
		ResourceRecordType:  rr.Type(),
		ResourceRecordClass: rr.Class(),
		Ttl:                 1080,
		RDataLength:         uint16(len(rr.Data())),
		RData:               rr.Data(),
	})
}

func (m *Message) AddQuery(q Query) {
	m.Body.Queries = append(m.Body.Queries, q)
	m.Header.NumberOfQuestions++
}

func (m *Message) AddAuthorative(rr record.ResourceRecord) {

	m.Body.Authorative = append(m.Body.Answers, Answer{
		Name:                rr.Name(),
		ResourceRecordType:  rr.Type(),
		ResourceRecordClass: rr.Class(),
		Ttl:                 1080,
		RDataLength:         uint16(len(rr.Data())),
		RData:               rr.Data(),
	})
}

func (m *Message) AddAdditional(rr record.ResourceRecord) {

	m.Body.Additional = append(m.Body.Answers, Answer{
		Name:                rr.Name(),
		ResourceRecordType:  rr.Type(),
		ResourceRecordClass: rr.Class(),
		Ttl:                 1080,
		RDataLength:         uint16(len(rr.Data())),
		RData:               rr.Data(),
	})
}
