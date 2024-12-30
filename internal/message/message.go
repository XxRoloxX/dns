package message

import (
	"errors"
	"fmt"
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

type ResourceRecordType uint

const (
	ResourceRecordType__A    = 1
	ResourceRecordType__AAAA = 2
	ResourceRecordType__MX   = 3
	ResourceRecordType__TXT  = 4
)

func NewResourceRecordType(code uint16) (ResourceRecordType, error) {
	switch code {
	case 1:
		return ResourceRecordType__A, nil
	case 2:
		return ResourceRecordType__AAAA, nil
	case 3:
		return ResourceRecordType__MX, nil
	case 4:
		return ResourceRecordType__TXT, nil
	default:
		return 0, errors.New("Invalid resource record type code")
	}
}

type ResourceRecordClass uint

const (
	ResourceRecordClass__In     = 1
	ResourceRecordClass__Review = 256
)

func NewResourceRecordClass(code uint16) (ResourceRecordClass, error) {
	switch code {
	case 1:
		return ResourceRecordClass__In, nil
	case 256:
		return ResourceRecordClass__Review, nil
	default:
		return 0, errors.New(fmt.Sprintf("Invalid resource record class code: %d", code))
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
	ResourceRecordType  ResourceRecordType
	ResourceRecordClass ResourceRecordClass
}

type Answer struct {
	Name                []string
	ResourceRecordType  ResourceRecordType
	ResourceRecordClass ResourceRecordClass
	Ttl                 uint32
	RDataLength         uint16
	RData               []byte
}

type Message struct {
	Header  Header
	Queries []Query
	Answers []Answer
}
