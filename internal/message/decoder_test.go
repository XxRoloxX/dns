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
			name: "A message for example.com should be decoded into single query",
			rawQuery: []byte{
				0x00, 0x00, 0x00, 0x00, //Header filler
				0x00, 0x00, 0x00, 0x00, //Header filler
				0x00, 0x00, 0x00, 0x00, //Header filler

				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Query Type: AAAA
				0b00000000, 0b00000001, // Query Class: IN

				// Answer section:
				0b11000000, 0b00001100, // Name (Pointer to "example.com" in the question section)
				0b00000000, 0b00000001, // Type: A (IPv4 address)
				0b00000000, 0b00000001, // Class: IN
				0x00, 0x00, 0x00, 0x3C, // Time to Live: 60 seconds
				0x00, 0x04, // RDATA Length: 4 bytes
				0xC0, 0xA8, 0x01, 0x01, // RDATA: IPv4 Address 192.168.1.1
			},
			expectedQueries: []Query{
				{
					groups: []string{"example", "com"},
					t:      ResourceRecordType__AAAA,
					class:  ResourceRecordClass__In,
				},
			},
			expectedAnswers: []Answer{
				{
					groups:      []string{"example", "com"},
					t:           ResourceRecordType__AAAA,
					class:       ResourceRecordClass__In,
					ttl:         60,
					rDataLength: 4,
					rData:       net.IPv4(192, 168, 1, 1).To4(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := Header{
				numberOfQuestions: 1,
				numberOfAnswers:   1,
			}

			query, answers, err := NewDecoder(tc.rawQuery).decodeBody(&header)
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedQueries, query)
			assert.Equal(t, tc.expectedAnswers, answers)
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
				transactionId: 0b0001101000101011,
				flags: HeaderFlags{
					query:              true,
					operationCode:      OpCode__Query,
					authorativeAnswer:  false,
					truncation:         false,
					recursionDesired:   true,
					recursionAvailable: false,
				},
				numberOfQuestions:    1,
				numberOfAnswers:      0,
				numberOfAuthorityRR:  0,
				numberOfAdditionalRR: 0,
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
