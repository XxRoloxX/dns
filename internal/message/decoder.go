package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/stretchr/testify/assert"
)

type Decoder struct {
	// reader io.Reader
	buf []byte
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{
		buf: buf,
	}
}

func (d *Decoder) Decode(message *Message) error {

	header, err := d.decodeHeader()
	if err != nil {
		return err
	}

	queries, err := d.decodeBody(header)
	if err != nil {
		return err
	}

	message.header = *header
	message.queries = queries

	return nil

}

func (d *Decoder) decodeBody(header *Header) ([]Query, []Answer, error) {

	index := 0
	queries := make([]Query, 0)
	answers := make([]Answer, 0)

	buf := make([]byte, 500)

	_, err := d.reader.Read(buf)
	if err != nil {
		return nil, nil, err
	}

	for _ = range header.numberOfQuestions {
		query, read, err := d.decodeQuery(buf[index:])
		if err != nil {
			return nil, nil, err
		}

		queries = append(queries, *query)
		index += read
	}

	for _ = range header.numberOfAnswers {
		answer, read, err := d.decodeAnswer(buf[index:])
		if err != nil {
			return nil, nil, err
		}

		answers = append(answers, *answer)
		index += read
	}

	return queries, answers, nil
}

func (d *Decoder) decodeAnswer(index int) (*Answer, int, error) {
	isPointer := d.isPoinerToDomain(d.buf[index])

	if isPointer {
		pointer := d.pointerFrom([2]byte{d.buf[index], d.buf[index+1]})

	}
}

func (d *Decoder) decodeNameRecursively(index uint16, groups []string) ([]string, error) {

	//copy the groups array
	groups = make([]string, len(groups), 0)
	for _, group := range groups {
		groups = append(groups, group)
	}

	for {
		groupLength := uint8(d.buf[index])

		isTerminated := d.isNameTerminated(groupLength)
		if isTerminated {
			return groups, nil
		}

		isPointer := d.isPoinerToDomain(groupLength)

		if isPointer {
			pointer := d.pointerFrom([2]byte{d.buf[index], d.buf[index+1]})
			return d.decodeNameRecursively(pointer, groups)
		}

		if groupLength >= uint8(len(d.buf)) {
			return nil, errors.New(fmt.Sprintf("Invalid group length: Expected %d, got %d", d.buf, groupLength))
		}

		// Get bytes as group after the group length byte
		group := d.buf[index+1 : index+1+uint16(groupLength)]
		groups = append(groups, string(group))

		index += uint16(groupLength) + 1
	}

}

func (d *Decoder) isPoinerToDomain(b byte) bool {
	return (b&128 > 0) && (b&64 > 0) // 110000 -> marks a start of an pointer
}

func (d *Decoder) pointerFrom(b [2]byte) uint16 {

	return binary.BigEndian.Uint16([]byte{b[0] & (63), b[1]})
}

func (d *Decoder) isNameTerminated(b byte) bool {
	return b == 0 // 110000 -> marks a start of an pointer
}

func (d *Decoder) decodeQuery(buf []byte) (*Query, int, error) {

	index := 0
	groups := make([]string, 0)

	for {
		groupLength := uint8(buf[index])

		if groupLength == 0 {
			//Skip null termination byte
			index++
			break
		}

		if groupLength >= uint8(len(buf)) {
			return nil, index, errors.New(fmt.Sprintf("Invalid group length: Expected %d, got %d", buf, groupLength))
		}

		// Get bytes as group after the group length byte
		group := buf[index+1 : index+int(groupLength)+1]
		groups = append(groups, string(group))

		index += int(groupLength) + 1
	}

	t, err := NewResourceRecordType(binary.BigEndian.Uint16(buf[index : index+2]))

	if err != nil {
		return nil, index, err
	}

	class, err := NewResourceRecordClass(binary.BigEndian.Uint16(buf[index+2 : index+4]))
	if err != nil {
		return nil, index, err
	}

	return &Query{
		groups: groups,
		t:      t,
		class:  class,
	}, index + 4, nil

}

func (d *Decoder) decodeHeader() (*Header, error) {

	if len(d.buf) < 12 {
		return nil, errors.New("Failed to decode header, invalid length")
	}

	transactionId := binary.BigEndian.Uint16(d.buf[0:2])
	flags, err := d.decodeHeaderFlags(d.buf[2:4])
	if err != nil {
		return nil, err
	}
	numberOfQuestions := binary.BigEndian.Uint16(d.buf[4:6])
	numberOfAnswers := binary.BigEndian.Uint16(d.buf[6:8])
	numberOfAuthorityRR := binary.BigEndian.Uint16(d.buf[8:10])
	numberOfAdditionalRR := binary.BigEndian.Uint16(d.buf[10:12])

	return &Header{
		transactionId:        transactionId,
		flags:                *flags,
		numberOfAuthorityRR:  numberOfAuthorityRR,
		numberOfAnswers:      numberOfAnswers,
		numberOfAdditionalRR: numberOfAdditionalRR,
		numberOfQuestions:    numberOfQuestions,
	}, nil
}

func (d *Decoder) decodeHeaderFlags(buf []byte) (*HeaderFlags, error) {

	if len(buf) != 2 {
		return nil, errors.New("Failed to decode header flags, invalid length")
	}

	query := buf[0]&128 == 0 // 10000000<- first bit (if zero then it is a query)

	opcode, err := NewOperationCode(uint16(buf[0] & 120)) // 01111000 (4 bites)
	if err != nil {
		return nil, err
	}

	authorative := buf[0]&4 > 0 // 00000100 (1 bit)

	truncated := buf[0]&2 > 0 // 00000010 (1 bit)

	recursionDesired := buf[0]&1 > 0 // 00000010 (1 bit)

	recursionAvailable := buf[1]&128 > 0 // 10000000 (1 bit)

	responseCode, err := NewResponseCode(uint16(buf[1] & 15)) // 00001111 (4 bits)
	if err != nil {
		return nil, err
	}

	flags := HeaderFlags{
		query:              query,
		operationCode:      opcode,
		authorativeAnswer:  authorative,
		truncation:         truncated,
		recursionDesired:   recursionDesired,
		recursionAvailable: recursionAvailable,
		responseCode:       responseCode,
	}

	return &flags, nil
}
