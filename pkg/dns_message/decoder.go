package message

import (
	"encoding/binary"
	"errors"
	"fmt"

	record "github.com/XxRoloxX/dns/pkg/dns_record"
)

type Decoder struct {
	buf []byte
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{
		buf: buf,
	}
}

func (d *Decoder) isIndexValid(index uint16) bool {
	return uint16(len(d.buf)) > index
}
func (d *Decoder) Decode(message *Message) error {

	header, err := d.decodeHeader()
	if err != nil {
		return err
	}

	body, err := d.decodeBody(header)
	if err != nil {
		return err
	}

	message.Header = *header
	message.Body = *body

	return nil

}

func (d *Decoder) decodeBody(header *Header) (*MessageBody, error) {

	var index uint16 = 12 // Message header always has 12 bytes
	queries := make([]Query, 0)
	answers := make([]Answer, 0)
	authorative := make([]Answer, 0)
	additonal := make([]Answer, 0)

	for _ = range header.NumberOfQuestions {
		query, read, err := d.decodeQuery(d.buf[index:])
		if err != nil {
			return nil, err
		}

		queries = append(queries, *query)
		index += read
	}

	for _ = range header.NumberOfAnswers {
		answer, newIndex, err := d.decodeAnswer(index)
		if err != nil {
			return nil, err
		}

		answers = append(answers, *answer)
		index = newIndex
	}

	for _ = range header.NumberOfAuthorityRR {
		answer, newIndex, err := d.decodeAnswer(index)
		if err != nil {
			return nil, err
		}

		authorative = append(authorative, *answer)
		index = newIndex
	}

	for _ = range header.NumberOfAdditionalRR {
		answer, newIndex, err := d.decodeAnswer(index)
		if err != nil {
			return nil, err
		}

		additonal = append(additonal, *answer)
		index = newIndex
	}

	return &MessageBody{
		Queries:     queries,
		Answers:     answers,
		Authorative: authorative,
		Additional:  additonal,
	}, nil
}

func (d *Decoder) decodeAnswer(index uint16) (*Answer, uint16, error) {

	name, index, err := d.decodeNameWithPointers(index)
	if err != nil {
		return nil, 0, err
	}

	t, err := record.NewResourceRecordType(binary.BigEndian.Uint16(d.buf[index : index+2]))
	class, err := record.NewResourceRecordClass(binary.BigEndian.Uint16(d.buf[index+2 : index+4]))
	ttl := binary.BigEndian.Uint32(d.buf[index+4 : index+8])
	rDataLength := binary.BigEndian.Uint16(d.buf[index+8 : index+10])

	rData := d.buf[index+10 : index+10+rDataLength]

	return &Answer{
		Name:                name,
		ResourceRecordType:  t,
		ResourceRecordClass: class,
		Ttl:                 ttl,
		RDataLength:         rDataLength,
		RData:               rData,
	}, index + 10 + rDataLength, nil

}

func (d *Decoder) decodeNameWithPointers(index uint16) ([]string, uint16, error) {

	if !d.isIndexValid(index) {
		return nil, 0, errors.New(fmt.Sprintf("Invalid index: %d", index))
	}

	groups := make([]string, 0)
	nameEndsAt := index
	wasPointerUsed := false

	// In the case of OPT record the name is set to single 0x00 byte (root domain)
	isRoot := d.isRootDomain(d.buf[index])
	if isRoot {
		return groups, nameEndsAt + 1, nil // there is no termination byte in this case
	}

	for {

		initialByte := uint8(d.buf[index])

		isTerminated := d.isNameTerminated(initialByte)
		if isTerminated {
			// slog.Info("NAME", "name", groups, "endsAt", nameEndsAt)
			return groups, nameEndsAt + 2, nil
		}

		isPointer := d.isPoinerToDomain(initialByte)
		if isPointer {
			pointer := d.pointerFrom(index)

			// Assign current index to pointer
			index = pointer
			wasPointerUsed = true
		}

		groupLength := uint8(d.buf[index])

		if !d.isIndexValid(uint16(groupLength) + index + 1) {
			return nil,
				0,
				errors.New(
					fmt.Sprintf(
						"Invalid group length: Expected < %d, got %d",
						len(d.buf),
						index+1+uint16(groupLength)))
		}

		// Get bytes as group after the group length byte
		group := d.buf[index+1 : index+1+uint16(groupLength)]
		groups = append(groups, string(group))

		index += uint16(groupLength) + 1

		if !wasPointerUsed {
			nameEndsAt = index
		}
	}
}

func (d *Decoder) isPoinerToDomain(b byte) bool {
	return (b&128 > 0) && (b&64 > 0) // 110000 -> marks a start of an pointer
}

func (d *Decoder) pointerFrom(index uint16) uint16 {

	return binary.BigEndian.Uint16([]byte{d.buf[index] & (63), d.buf[index+1]})
}

func (d *Decoder) isNameTerminated(b byte) bool {
	return b == 0 // 00000000 -> marks termination
}

func (d *Decoder) isRootDomain(b byte) bool {
	return b == 0 // 00000000 -> root domain -> . used for OPT records
}

func (d *Decoder) decodeQuery(buf []byte) (*Query, uint16, error) {

	var index uint16 = 0
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
		group := buf[index+1 : index+uint16(groupLength)+1]
		groups = append(groups, string(group))

		index += uint16(groupLength) + 1
	}

	t, err := record.NewResourceRecordType(binary.BigEndian.Uint16(buf[index : index+2]))

	if err != nil {
		return nil, index, err
	}

	class, err := record.NewResourceRecordClass(binary.BigEndian.Uint16(buf[index+2 : index+4]))
	if err != nil {
		return nil, index, err
	}

	return &Query{
		Name:                groups,
		ResourceRecordType:  t,
		ResourceRecordClass: class,
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
		TransactionId:        transactionId,
		Flags:                *flags,
		NumberOfAuthorityRR:  numberOfAuthorityRR,
		NumberOfAnswers:      numberOfAnswers,
		NumberOfAdditionalRR: numberOfAdditionalRR,
		NumberOfQuestions:    numberOfQuestions,
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
		Query:              query,
		OperationCode:      opcode,
		AuthorativeAnswer:  authorative,
		Truncation:         truncated,
		RecursionDesired:   recursionDesired,
		RecursionAvailable: recursionAvailable,
		ResponseCode:       responseCode,
	}

	return &flags, nil
}
