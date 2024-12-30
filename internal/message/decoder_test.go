package message

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder_decodeBody(t *testing.T) {
	testCases := []struct {
		name            string
		rawQuery        []byte
		expectedQueries []Query
		expectedAnswers []Answer
		expectedErr     error
	}{
		{
			name: "A message for example.com with an answer using pointer should be decoded",
			rawQuery: []byte{
				//Header bytes filler
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				// Question section:
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Query Type: AAAA
				0b00000000, 0b00000001, // Query Class: IN

				// Answer section (uses pointer):
				0b11000000, 0b00001100, // Name (Pointer to "example.com" in the question section)
				0b00000000, 0b00000001, // Type: A (IPv4 address)
				0b00000000, 0b00000001, // Class: IN
				0x00, 0x00, 0x00, 0x3C, // Time to Live: 60 seconds
				0x00, 0x04, // RDATA Length: 4 bytes
				0xC0, 0xA8, 0x01, 0x01, // RDATA: IPv4 Address 192.168.1.1
			},
			expectedQueries: []Query{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  ResourceRecordType__A,
					ResourceRecordClass: ResourceRecordClass__In,
				},
			},
			expectedAnswers: []Answer{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  ResourceRecordType__A,
					ResourceRecordClass: ResourceRecordClass__In,
					Ttl:                 60,
					RDataLength:         4,
					RData:               net.IPv4(192, 168, 1, 1).To4(),
				},
			},
		},
		{
			name: "A message for example.com with an answer without pointer should be decoded",
			rawQuery: []byte{

				//Header bytes
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				// Question section:
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Query Type: A
				0b00000000, 0b00000001, // Query Class: IN

				// Answer section (does not use pointer):
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Type: A (IPv4 address)
				0b00000000, 0b00000001, // Class: IN
				0x00, 0x00, 0x00, 0x3C, // Time to Live: 60 seconds
				0x00, 0x04, // RDATA Length: 4 bytes
				0xC0, 0xA8, 0x01, 0x01, // RDATA: IPv4 Address 192.168.1.1
			},
			expectedQueries: []Query{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  ResourceRecordType__A,
					ResourceRecordClass: ResourceRecordClass__In,
				},
			},
			expectedAnswers: []Answer{
				{
					Name:                []string{"example", "com"},
					ResourceRecordType:  ResourceRecordType__A,
					ResourceRecordClass: ResourceRecordClass__In,
					Ttl:                 60,
					RDataLength:         4,
					RData:               net.IPv4(192, 168, 1, 1).To4(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := Header{
				NumberOfQuestions: 1,
				NumberOfAnswers:   1,
			}

			body, err := NewDecoder(tc.rawQuery).decodeBody(&header)
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedQueries, body.Queries)
			assert.Equal(t, tc.expectedAnswers, body.Answers)
		})
	}
}

func TestDecoder_decodeHeader(t *testing.T) {

	testCases := []struct {
		name           string
		rawQuery       []byte
		expectedHeader Header
		expectedErr    error
	}{
		{
			name: "A message for example.com should be decoded into single query",
			rawQuery: []byte{
				0b00011010, 0b00101011, // Transaction ID
				0b00000001, 0b00000000, // Flags
				0b00000000, 0b00000001, // Questions
				0b00000000, 0b00000000, // Answer RRs
				0b00000000, 0b00000000, // Authority RRs
				0b00000000, 0b00000000, // Additional RRs
			},
			expectedHeader: Header{
				TransactionId: 0b0001101000101011,
				Flags: HeaderFlags{
					Query:              true,
					OperationCode:      OpCode__Query,
					AuthorativeAnswer:  false,
					Truncation:         false,
					RecursionDesired:   true,
					RecursionAvailable: false,
				},
				NumberOfQuestions:    1,
				NumberOfAnswers:      0,
				NumberOfAuthorityRR:  0,
				NumberOfAdditionalRR: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			header, err := NewDecoder(tc.rawQuery).decodeHeader()
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedHeader, *header)
		})
	}
}
