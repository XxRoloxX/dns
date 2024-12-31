package message

import (
	"net"
	"testing"

	"github.com/XxRoloxX/dns/pkg/dns_record"
	"github.com/stretchr/testify/assert"
)

func TestEncoder_encodeHeader(t *testing.T) {
	testCases := []struct {
		name            string
		header          Header
		expectedEncoded []byte
	}{
		{
			name: "Encode header for a message for example.com with an answer without pointer",
			header: Header{
				TransactionId: 0x0000,
				Flags: HeaderFlags{
					Query:              false,
					OperationCode:      0,
					AuthorativeAnswer:  false,
					Truncation:         false,
					RecursionDesired:   true,
					RecursionAvailable: true,
					ResponseCode:       0,
				},
				NumberOfQuestions:    1,
				NumberOfAnswers:      1,
				NumberOfAuthorityRR:  0,
				NumberOfAdditionalRR: 0,
			},
			expectedEncoded: []byte{
				0b00000000, 0b00000000, // Transaction ID
				0b10000001, 0b10000000, // Flags (query = false, response = true, etc.)
				0b00000000, 0b00000001, // Number of questions
				0b00000000, 0b00000001, // Number of answers
				0b00000000, 0b00000000, // Number of authority RR
				0b00000000, 0b00000000, // Number of additional RR
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := NewEncoder()

			// Encode Header
			encodedHeader := encoder.encodeHeader(&tc.header)

			// Assert that the encoded header matches the expected output
			assert.Equal(t, tc.expectedEncoded, encodedHeader)
		})
	}
}

func TestEncoder_encodeBody(t *testing.T) {
	testCases := []struct {
		name            string
		queries         []Query
		answers         []Answer
		expectedEncoded []byte
	}{
		{
			name: "Encode body for a message for example.com with an answer without pointer",
			queries: []Query{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  record.ResourceRecordType__A,
					ResourceRecordClass: record.ResourceRecordClass__In,
				},
			},
			answers: []Answer{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  record.ResourceRecordType__A,
					ResourceRecordClass: record.ResourceRecordClass__In,
					Ttl:                 60,
					RDataLength:         4,
					RData:               net.IPv4(192, 168, 1, 1).To4(),
				},
			},
			expectedEncoded: []byte{
				// Question Section: "example.com" + Type: A + Class: IN
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Type: A (IPv4 address)
				0b00000000, 0b00000001, // Class: IN

				// Answer section
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Type: A (IPv4 address)
				0b00000000, 0b00000001, // Class: IN
				0x00, 0x00, 0x00, 0x3C, // Time to Live: 60 seconds
				0x00, 0x04, // RDATA Length: 4 bytes
				0xC0, 0xA8, 0x01, 0x01, // RDATA: IPv4 Address 192.168.1.1
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := NewEncoder()

			// Encode Queries
			var encodedQueries []byte
			for _, query := range tc.queries {
				encodedQueries = append(encodedQueries, encoder.encodeQuery(query)...)
			}

			// Encode Answers
			var encodedAnswers []byte
			for _, answer := range tc.answers {
				encodedAnswers = append(encodedAnswers, encoder.encodeAnswer(answer)...)
			}

			// Combine all encoded parts
			encodedBody := append(encodedQueries, encodedAnswers...)

			// Assert that the encoded body matches the expected output
			assert.Equal(t, tc.expectedEncoded, encodedBody)
		})
	}
}
