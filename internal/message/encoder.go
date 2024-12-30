package message

import (
	"encoding/binary"
	bin "github.com/XxRoloxX/dns/pkg/binary"
)

func test() {

}

type Encoder struct {
	buffer []byte
}

func NewEncoder() *Encoder {
	return &Encoder{
		buffer: make([]byte, 0),
	}
}

func (e *Encoder) Encode(message *Message) []byte {

	encodedHeader := e.encodeHeader(&message.Header)
	e.buffer = append(e.buffer, encodedHeader...)

	for _, query := range message.Queries {
		e.buffer = append(e.buffer, e.encodeQuery(query)...)
	}

	for _, answer := range message.Answers {
		e.buffer = append(e.buffer, e.encodeAnswer(answer)...)
	}

	return e.buffer
}

func (e *Encoder) encodeHeaderFlags(headerFlags *HeaderFlags) []byte {

	flags := make([]byte, 2)

	flagsFirstBitSetter := bin.NewBitSetter(flags[0])

	flagsFirstBitSetter.Set(0, !headerFlags.Query)
	flagsFirstBitSetter.SetRange(1, 4, uint8(headerFlags.OperationCode))
	flagsFirstBitSetter.Set(5, headerFlags.AuthorativeAnswer)
	flagsFirstBitSetter.Set(6, headerFlags.Truncation)
	flagsFirstBitSetter.Set(7, headerFlags.RecursionDesired)
	flags[0] = flagsFirstBitSetter.Byte()

	flagsSecondBitSetter := bin.NewBitSetter(flags[1])
	flagsSecondBitSetter.Set(0, headerFlags.RecursionAvailable)
	flagsFirstBitSetter.SetRange(4, 7, uint8(headerFlags.ResponseCode))
	flags[1] = flagsSecondBitSetter.Byte()

	return flags
}

func (e *Encoder) encodeHeader(header *Header) []byte {

	encodedHeader := make([]byte, 0, 12)

	transactionId := make([]byte, 2)
	binary.BigEndian.PutUint16(transactionId, header.TransactionId)

	flags := e.encodeHeaderFlags(&header.Flags)

	numberOfQuestions := make([]byte, 2)
	binary.BigEndian.PutUint16(numberOfQuestions, header.NumberOfQuestions)

	numberOfAnswers := make([]byte, 2)
	binary.BigEndian.PutUint16(numberOfAnswers, header.NumberOfAnswers)

	numberOfAuthorityRR := make([]byte, 2)
	binary.BigEndian.PutUint16(numberOfAuthorityRR, header.NumberOfAuthorityRR)

	numberOfAdditionalRR := make([]byte, 2)
	binary.BigEndian.PutUint16(numberOfAdditionalRR, header.NumberOfAdditionalRR)

	encodedHeader = append(encodedHeader, transactionId...)
	encodedHeader = append(encodedHeader, flags...)
	encodedHeader = append(encodedHeader, numberOfQuestions...)
	encodedHeader = append(encodedHeader, numberOfAnswers...)
	encodedHeader = append(encodedHeader, numberOfAuthorityRR...)
	encodedHeader = append(encodedHeader, numberOfAdditionalRR...)

	return encodedHeader
}

func (e *Encoder) encodeName(name []string) []byte {

	encodedName := make([]byte, 0)

	for _, group := range name {
		// Group length byte
		encodedName = append(encodedName, uint8(len(group)))

		encodedName = append(encodedName, []byte(group)...)
	}

	// Termination byte
	encodedName = append(encodedName, 0)

	return encodedName
}

func (e *Encoder) encodeQuery(query Query) []byte {

	encodedQuery := make([]byte, 0)

	encodedName := e.encodeName(query.Name)

	resourceRecordType := make([]byte, 2)
	binary.BigEndian.PutUint16(resourceRecordType, uint16(query.ResourceRecordType))

	resourceRecordClass := make([]byte, 2)
	binary.BigEndian.PutUint16(resourceRecordClass, uint16(query.ResourceRecordClass))

	encodedQuery = append(encodedQuery, encodedName...)
	encodedQuery = append(encodedQuery, resourceRecordType...)
	encodedQuery = append(encodedQuery, resourceRecordClass...)

	return encodedQuery
}

func (e *Encoder) encodeAnswer(answer Answer) []byte {

	encodedAnswer := make([]byte, 0)

	encodedName := e.encodeName(answer.Name)

	resourceRecordType := make([]byte, 2)
	binary.BigEndian.PutUint16(resourceRecordType, uint16(answer.ResourceRecordType))

	resourceRecordClass := make([]byte, 2)
	binary.BigEndian.PutUint16(resourceRecordClass, uint16(answer.ResourceRecordClass))

	ttl := make([]byte, 4)
	binary.BigEndian.PutUint32(ttl, uint32(answer.Ttl))

	rDataLength := make([]byte, 2)
	binary.BigEndian.PutUint16(rDataLength, uint16(answer.RDataLength))

	encodedAnswer = append(encodedAnswer, encodedName...)
	encodedAnswer = append(encodedAnswer, resourceRecordType...)
	encodedAnswer = append(encodedAnswer, resourceRecordClass...)
	encodedAnswer = append(encodedAnswer, ttl...)
	encodedAnswer = append(encodedAnswer, rDataLength...)
	encodedAnswer = append(encodedAnswer, answer.RData...)

	return encodedAnswer
}
